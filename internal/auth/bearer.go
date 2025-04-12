package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers["Authorization"]
	if len(bearerToken) == 0 {
		return "", errors.New("no authorization token found")
	}
	if len(bearerToken) > 1 {
		return "", errors.New("more than one auth tokens are provided")
	}

	TOKEN_STRING := strings.ReplaceAll(bearerToken[0], "Bearer ", "")
	return TOKEN_STRING, nil

}
