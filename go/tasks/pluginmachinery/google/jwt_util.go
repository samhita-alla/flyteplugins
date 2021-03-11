package google

import (
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/jwt"
)

// getIss returns the claim `iss` from the JWT token
func getIss(token string) (string, error) {
	parsed, err := jwt.ParseSigned(token)

	if err != nil {
		return "", errors.Wrapf(err, "failed to parse JWT token")
	}

	claims := jwt.Claims{}
	err = parsed.UnsafeClaimsWithoutVerification(&claims)

	if err != nil {
		return "", errors.Wrapf(err, "failed to get JWT token claims")
	}

	return claims.Issuer, nil
}
