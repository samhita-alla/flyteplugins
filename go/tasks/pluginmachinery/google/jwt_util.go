package google

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// getIss returns the claim `iss` from the JWT token
func getIss(token string) (string, error) {
	claims, err := parseJwtClaims(token)
	if err != nil {
		return "", err
	}

	rawIss := claims["iss"]
	if rawIss == nil {
		return "", fmt.Errorf("no aud in the token claims")
	}

	data, err := json.Marshal(rawIss)
	if err != nil {
		return "", err
	}

	var iss string
	err = json.Unmarshal(data, &iss)

	return iss, err
}

func parseJwtClaims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token contains an invalid number of segments: %d, expected: 3", len(parts))
	}

	// Decode the second part.
	claimBytes, err := decodeSegment(parts[1])
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))

	claims := make(map[string]interface{})
	if err := dec.Decode(&claims); err != nil {
		return nil, fmt.Errorf("failed to decode the JWT claims")
	}
	return claims, nil
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
