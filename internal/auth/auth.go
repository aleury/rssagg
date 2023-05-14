package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const tokenPrefix = "ApiKey"

// GetAPIKey extracts an API key from the headers of an http request.
// Example:
// Authorization: ApiKey {insert api key here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("missing authorization header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization header")
	}
	if vals[0] != tokenPrefix {
		return "", fmt.Errorf("invalid token prefix: %s", vals[0])
	}
	return vals[1], nil
}
