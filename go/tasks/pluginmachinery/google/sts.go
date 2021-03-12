package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/flyteorg/flytestdlib/logger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

const (
	requestedTokenType     = "urn:ietf:params:oauth:token-type:access_token"   // #nosec
	grantType              = "urn:ietf:params:oauth:grant-type:token-exchange" // #nosec
	subjectTokenType       = "urn:ietf:params:oauth:token-type:jwt"            // #nosec
	federatedTokenEndpoint = "https://securetoken.googleapis.com"              // #nosec
	contentType            = "application/json"
	httpTimeOutInSec       = 5
	maxRequestRetry        = 5
	lifetime               = "3000s"
	bearerTokenType        = "Bearer"
	defaultScope           = "https://www.googleapis.com/auth/cloud-platform"
)

type StsRequest struct {
	// This token is a an external credential issued by a workload identity pool provider
	SubjectToken string

	// Scope is OAuth 2.0 scopes to include on the resulting access token
	Scope []string

	// ServiceAccount defines Google service account name
	ServiceAccount string

	// IdentityNamespace Workload identity namespace, e.g. [project_id].svc.id.goog
	IdentityNamespace string

	federatedTokenEndpoint string // only for testing
	iamCredentialsEndpoint string // only for testing
}

// GenerateToken takes STS request and fetches token, returns Google credential
func ExchangeToken(ctx context.Context, request StsRequest) (*oauth2.Token, error) {
	if len(request.Scope) == 0 {
		request.Scope = []string{defaultScope}
	}

	if request.federatedTokenEndpoint == "" {
		request.federatedTokenEndpoint = federatedTokenEndpoint
	}

	// Exchange OIDC ID token for Google credentials. See https://cloud.google.com/iam/docs/access-resources-oidc
	// A lot of the GKE-specific code is inspired by https://pkg.go.dev/istio.io/istio/security/pkg/stsservice

	// 1. Obtain an OIDC ID token from your identity provider (passed in request.SubjectToken)

	// 2. Pass OIDC ID token to Security Token Service token() method to get a federated access token
	federatedToken, err := fetchFederatedToken(ctx, request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate federated token")
	}

	// 3. Call generateAccessToken() to exchange the federated token for a service account access token.
	// A limited number of Google Cloud APIs support federated tokens;
	// all Google Cloud APIs support service account access tokens.
	token, err := generateAccessToken(ctx, request, federatedToken)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate access token")
	}

	return token, nil
}

type federatedTokenResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       int64  `json:"expires_in"` // Expiration time in seconds
}

func constructFederatedTokenRequest(ctx context.Context, request StsRequest) (*http.Request, error) {
	iss, err := getIss(request.SubjectToken)

	if err != nil {
		return nil, errors.Wrapf(err, "can't get iss from subject token")
	}

	if iss == "" {
		return nil, errors.New("issued federated token doesn't have 'iss' in payload")
	}

	audience := fmt.Sprintf("identitynamespace:%s:%s", request.IdentityNamespace, iss)
	scope := strings.Join(request.Scope, ",")

	query := map[string]string{
		"audience":           audience,
		"grantType":          grantType,
		"requestedTokenType": requestedTokenType,
		"subjectTokenType":   subjectTokenType,
		"subjectToken":       request.SubjectToken,
		"scope":              scope,
	}

	redactedQuery := map[string]string{
		"audience":           scope,
		"grantType":          grantType,
		"requestedTokenType": requestedTokenType,
		"subjectTokenType":   subjectTokenType,
		"subjectToken":       "redacted",
		"scope":              audience,
	}

	if jsonRedactedQuery, err := json.Marshal(redactedQuery); err == nil {
		logger.Infof(ctx, "Prepared federated token request: %s", string(jsonRedactedQuery))
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query for get federated token request: %+v", err)
	}

	req, err := http.NewRequest("POST", request.federatedTokenEndpoint+"/v1/identitybindingtoken", bytes.NewBuffer(jsonQuery))
	if err != nil {
		return req, fmt.Errorf("failed to create get federated token request: %+v", err)
	}

	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// fetchFederatedToken exchanges a third-party issued Json Web Token for an OAuth2.0 access token
// which asserts a third-party identity within an identity namespace.
func fetchFederatedToken(ctx context.Context, parameters StsRequest) (*oauth2.Token, error) {
	respData := &federatedTokenResponse{}

	req, err := constructFederatedTokenRequest(ctx, parameters)
	if err != nil {
		logger.Errorf(ctx, "failed to create get federated token request: %+v", err)
		return nil, err
	}

	resp, err := sendRequestWithRetry(req)
	if err != nil {
		respCode := 0
		if resp != nil {
			respCode = resp.StatusCode
		}

		return nil, fmt.Errorf("failed to exchange federated token (HTTP status %d): %v", respCode,
			err)
	}
	// resp should not be nil.
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read federated token response body: %+v", err)
	}

	if err := json.Unmarshal(body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal federated token response data: %v", err)
	}

	if respData.AccessToken == "" {
		return nil, errors.Errorf("federated token response does not have access token. %v", body)
	}

	logger.Infof(ctx, "Federated token will expire in %d seconds", respData.ExpiresIn)

	token := oauth2.Token{
		AccessToken: respData.AccessToken,
		TokenType:   respData.TokenType,
		Expiry:      time.Now().Add(time.Duration(respData.ExpiresIn) * time.Second),
	}

	return &token, nil
}

// Send HTTP request every 0.01 seconds until successfully receive response or hit max retry numbers.
// If response code is 4xx, return immediately without retry.
func sendRequestWithRetry(req *http.Request) (resp *http.Response, err error) {
	httpClient := http.Client{Timeout: httpTimeOutInSec * time.Second}

	for i := 0; i < maxRequestRetry; i++ {
		resp, err = httpClient.Do(req)

		if resp != nil && resp.StatusCode == http.StatusOK {
			return resp, err
		}

		if resp != nil && resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	if resp != nil && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		return resp, fmt.Errorf("HTTP Status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, err
}

func generateAccessToken(ctx context.Context, stsRequest StsRequest, federatedAccessToken *oauth2.Token) (token *oauth2.Token, err error) {
	options := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(federatedAccessToken))}

	if stsRequest.iamCredentialsEndpoint != "" {
		options = append(options, option.WithEndpoint(stsRequest.iamCredentialsEndpoint))
	}

	service, err := iamcredentials.NewService(ctx, options...)

	if err != nil {
		return nil, errors.Wrap(err, "failed to construct iamcredentials service")
	}

	name := "projects/-/serviceAccounts/" + stsRequest.ServiceAccount
	request := iamcredentials.GenerateAccessTokenRequest{
		Scope:    stsRequest.Scope,
		Lifetime: lifetime,
	}

	logger.Infof(ctx, "Generating access token for [%v]", name)

	for i := 0; i < maxRequestRetry; i++ {
		respData, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, &request).Do()

		if err != nil {
			apiError, ok := err.(*googleapi.Error)

			if ok {
				logger.Errorf(ctx, "failed to generate access token: %v", apiError)

				if apiError.Code >= http.StatusBadRequest && apiError.Code < http.StatusInternalServerError {
					break
				}
			} else {
				break
			}
		} else {
			return &oauth2.Token{
				AccessToken: respData.AccessToken,
				TokenType:   bearerTokenType,
			}, nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil, errors.Wrapf(err, "failed to generate access token")
}
