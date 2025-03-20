package landb

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	authURL = "https://auth.cern.ch/auth/realms/cern/api-access/token"
)

func Authenticate(clientID, clientSecret, audience string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"audience":      audience,
		}).
		SetResult(&AuthResponse{}).
		Post(authURL)

	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("authentication error: %s", resp.Status())
	}

	authResp := resp.Result().(*AuthResponse)
	return authResp.AccessToken, nil
}
