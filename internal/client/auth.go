// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package landb

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const authURL = "https://auth.cern.ch/auth/realms/cern/api-access/token"


type AuthResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func Authenticate(clientID, clientSecret, audience string) (*AuthResponse, error) {
	client := resty.New()

	var authResp AuthResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"audience":      audience,
		}).
		SetResult(&authResp).
		Post(authURL)

	if err != nil {
		return nil, fmt.Errorf("authentication request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("authentication error [%d]: %s", resp.StatusCode(), resp.String())
	}

	return &authResp, nil
}
