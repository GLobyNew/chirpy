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

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers["Authorization"]
	if len(apiKey) == 0 {
		return "", errors.New("no api key found")
	}
	if len(apiKey) > 1 {
		return "", errors.New("more than one api key provided")
	}

	API_KEY := strings.ReplaceAll(apiKey[0], "ApiKey ", "")
	API_KEY = strings.TrimSpace(API_KEY)
	return API_KEY, nil
}
