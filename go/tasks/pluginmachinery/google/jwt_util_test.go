package google

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIss(t *testing.T) {
	t.Run("has iss", func(t *testing.T) {
		payload := "eyJpc3MiOiJkb2dlIiwiaWF0IjpudWxsLCJleHAiOm51bGwsImF1ZCI6Ind3d" +
			"y5leGFtcGxlLmNvbSIsInN1YiI6Impyb2NrZXRAZXhhbXBsZS5jb20ifQ"

		iss, err := getIss("header." + payload + ".signature")

		assert.NoError(t, err)
		assert.Equal(t, "doge", iss)
	})
}
