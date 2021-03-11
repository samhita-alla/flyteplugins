package google

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIss(t *testing.T) {
	t.Run("has iss", func(t *testing.T) {
		header := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
		payload := "eyJpc3MiOiJkb2dlIiwiaWF0IjpudWxsLCJleHAiOm51bGwsImF1ZCI6IiIsInN1YiI6IiJ9"
		signature := "zRDLWGQa25HqLesVLgrIbG3pVFTiD7WbjTg-2f6v5FI"

		iss, err := getIss(header + "." + payload + "." + signature)

		assert.NoError(t, err)
		assert.Equal(t, "doge", iss)
	})
}
