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
	identity := Identity{
		K8sNamespace:      "flytesnacks-development",
		K8sServiceAccount: "default",
	}

	t.Run("has no cached token", func(t *testing.T) {
		ts := gkeTokenSource{}

		_, ok := ts.getCachedToken(ctx, identity)

		assert.False(t, ok)
	})

	t.Run("has cached token", func(t *testing.T) {
		token := oauth2.Token{
			AccessToken: "secret",
			Expiry:      time.Now().Add(time.Hour),
		}
		ts := gkeTokenSource{}

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
		ts := gkeTokenSource{}

		ts.tokens.Store(identity, &token)

		_, ok := ts.getCachedToken(ctx, identity)
		assert.False(t, ok)

		_, ok = ts.tokens.Load(identity)
		assert.False(t, ok)
	})
}
