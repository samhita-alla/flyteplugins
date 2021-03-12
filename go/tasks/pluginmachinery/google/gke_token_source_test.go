package google

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/flyteorg/flytestdlib/atomic"
	"github.com/golang/groupcache/singleflight"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetCachedToken(t *testing.T) {
	ctx := context.TODO()
	identity := Identity{
		K8sNamespace:      "flytesnacks-development",
		K8sServiceAccount: "default",
	}

	t.Run("has no cached token", func(t *testing.T) {
		ts := newGkeTokenSource()

		_, ok := ts.getCachedToken(ctx, identity)

		assert.False(t, ok)
	})

	t.Run("has cached token", func(t *testing.T) {
		token := oauth2.Token{
			AccessToken: "secret",
			Expiry:      time.Now().Add(time.Hour),
		}
		ts := newGkeTokenSource()

		ts.tokens.Store(identity, &token)

		cached, ok := ts.getCachedToken(ctx, identity)

		assert.True(t, ok)
		assert.NotNil(t, cached)
		assert.Equal(t, token, *cached)
	})

	t.Run("has expired token", func(t *testing.T) {
		token := oauth2.Token{
			AccessToken: "secret",
			Expiry:      time.Now(),
		}
		ts := newGkeTokenSource()

		ts.tokens.Store(identity, &token)

		_, ok := ts.getCachedToken(ctx, identity)
		assert.False(t, ok)

		_, ok = ts.tokens.Load(identity)
		assert.False(t, ok)
	})
}

func TestGetGcpServiceAccount(t *testing.T) {
	ctx := context.TODO()

	t.Run("get GCP service account", func(t *testing.T) {
		ts := newGkeTokenSource()
		ts.kubeClient = fake.NewSimpleClientset(&corev1.ServiceAccount{
			ObjectMeta: v1.ObjectMeta{
				Name:      "name",
				Namespace: "namespace",
				Annotations: map[string]string{
					"owner":                          "abc",
					"iam.gke.io/gcp-service-account": "gcp-service-account",
				},
			}})

		gcpServiceAccount, err := ts.getGcpServiceAccount(ctx, Identity{
			K8sNamespace:      "namespace",
			K8sServiceAccount: "name",
		})

		assert.NoError(t, err)
		assert.Equal(t, "gcp-service-account", gcpServiceAccount)
	})
}

func newGkeTokenSource() gkeTokenSource {
	deletion := atomic.NewBool(false)

	return gkeTokenSource{
		tokens:       &sync.Map{},
		singleflight: &singleflight.Group{},
		deletion:     &deletion,
	}
}
