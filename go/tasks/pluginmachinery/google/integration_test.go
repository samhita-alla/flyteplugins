package google

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authentication "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iamcredentials/v1"
)

func TestEndToEnd(t *testing.T) {
	server := newFakeServer(t)
	defer server.Close()

	header := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
	payload := "eyJpc3MiOiJkb2dlIiwiaWF0IjpudWxsLCJleHAiOm51bGwsImF1ZCI6IiIsInN1YiI6IiJ9"
	signature := "zRDLWGQa25HqLesVLgrIbG3pVFTiD7WbjTg-2f6v5FI"

	kubeClient := fake.NewSimpleClientset()

	kubeClient.PrependReactor("*", "serviceaccounts", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		if action.GetSubresource() == "token" {
			return true, &authentication.TokenRequest{
				Status: authentication.TokenRequestStatus{
					Token: header + "." + payload + "." + signature,
				},
			}, nil
		}

		return true, &corev1.ServiceAccount{
			ObjectMeta: v1.ObjectMeta{
				Name:      "name",
				Namespace: "namespace",
				Annotations: map[string]string{
					"owner":                          "abc",
					"iam.gke.io/gcp-service-account": "gcp-service-account",
				},
			},
		}, nil
	})

	ctx := context.TODO()

	gkeTokenSource := newGKETokenSource(kubeClient, TokenSourceFactoryConfig{})
	gkeTokenSource.federatedTokenEndpoint = server.URL
	gkeTokenSource.iamCredentialsEndpoint = server.URL

	tokenSource, err := gkeTokenSource.GetTokenSource(ctx, Identity{
		K8sNamespace:      "namespace",
		K8sServiceAccount: "name",
	})

	assert.NoError(t, err)

	token, err := tokenSource.Token()

	assert.NoError(t, err)
	assert.Equal(t, "google-access-token", token.AccessToken)
}

func newFakeServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v1/identitybindingtoken" && request.Method == "POST" {
			writer.WriteHeader(200)
			response := federatedTokenResponse{
				AccessToken: "federated-access-token",
			}

			bytes, err := json.Marshal(response)
			assert.NoError(t, err)

			_, err = writer.Write(bytes)
			assert.NoError(t, err)

			return
		}

		if request.URL.Path == "/v1/projects/-/serviceAccounts/gcp-service-account:generateAccessToken" &&
			request.Method == "POST" {

			auth := request.Header.Get("Authorization")

			if auth != "Bearer federated-access-token" {
				writer.WriteHeader(401)
				return
			}

			writer.WriteHeader(200)
			response := iamcredentials.GenerateAccessTokenResponse{AccessToken: "google-access-token"}

			bytes, err := json.Marshal(response)
			assert.NoError(t, err)

			_, err = writer.Write(bytes)
			assert.NoError(t, err)

			return
		}

		writer.WriteHeader(500)
	}))
}
