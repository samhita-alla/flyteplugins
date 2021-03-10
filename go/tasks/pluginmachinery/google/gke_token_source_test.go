package google

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestGetCachedToken(t *testing.T) {
	ctx := context.TODO()

	t.Run("has no cached token", func(t *testing.T) {
		ts := GkeTokenSource{}

		_, ok := ts.getCachedToken(ctx, "flytesnacks-development", "default")

		assert.False(t, ok)
	})

	t.Run("has cached token", func(t *testing.T) {
		token := oauth2.Token{
			AccessToken: "secret",
			Expiry:      time.Now().Add(time.Hour),
		}
		ts := GkeTokenSource{}

		ts.tokens.Store(tokensKey{
			k8sNamespace:      "flytesnacks-development",
			k8sServiceAccount: "default",
		}, &token)

		cached, ok := ts.getCachedToken(ctx, "flytesnacks-development", "default")

		assert.True(t, ok)
		assert.NotNil(t, cached)
		assert.Equal(t, token, *cached)
	})

	t.Run("has expired token", func(t *testing.T) {
		token := oauth2.Token{
			AccessToken: "secret",
			Expiry:      time.Now(),
		}
		ts := GkeTokenSource{}
		cacheKey := tokensKey{
			k8sNamespace:      "flytesnacks-development",
			k8sServiceAccount: "default",
		}

		ts.tokens.Store(cacheKey, &token)

		_, ok := ts.getCachedToken(ctx, "flytesnacks-development", "default")
		assert.False(t, ok)

		_, ok = ts.tokens.Load(cacheKey)
		assert.False(t, ok)
	})
}
