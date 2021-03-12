package google

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iamcredentials/v1"
)

func TestExchangeToken(t *testing.T) {
	server := newFakeServer()
	defer server.Close()

	ctx := context.TODO()

	header := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
	payload := "eyJpc3MiOiJkb2dlIiwiaWF0IjpudWxsLCJleHAiOm51bGwsImF1ZCI6IiIsInN1YiI6IiJ9"
	signature := "zRDLWGQa25HqLesVLgrIbG3pVFTiD7WbjTg-2f6v5FI"

	t.Run("exchange token", func(t *testing.T) {
		token, err := ExchangeToken(ctx, StsRequest{
			SubjectToken:           header + "." + payload + "." + signature,
			ServiceAccount:         "service-account",
			IdentityNamespace:      "flyte.svc.id.doge",
			iamCredentialsEndpoint: server.URL,
			federatedTokenEndpoint: server.URL,
		})

		assert.NoError(t, err)
		assert.Equal(t, "google-access-token", token.AccessToken)
	})
}

func newFakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v1/identitybindingtoken" && request.Method == "POST" {
			writer.WriteHeader(200)
			response := federatedTokenResponse{
				AccessToken: "federated-access-token",
			}
			bytes, _ := json.Marshal(response)
			_, _ = writer.Write(bytes)
			return
		}

		if request.URL.Path == "/v1/projects/-/serviceAccounts/service-account:generateAccessToken" &&
			request.Method == "POST" {

			auth := request.Header.Get("Authorization")

			if auth != "Bearer federated-access-token" {
				writer.WriteHeader(401)
				return
			}

			writer.WriteHeader(200)
			response := iamcredentials.GenerateAccessTokenResponse{AccessToken: "google-access-token"}
			bytes, _ := json.Marshal(response)
			_, _ = writer.Write(bytes)
			return
		}

		writer.WriteHeader(500)
	}))
}
