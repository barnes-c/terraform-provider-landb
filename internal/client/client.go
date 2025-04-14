// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

package landb

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiEndpoint string
	token       string
	HTTPClient  *resty.Client
}

func NewClient(apiEndpoint, clientID, clientSecret, audience string) (*Client, error) {
	authResp, err := Authenticate(clientID, clientSecret, audience)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	client := resty.New().
		SetBaseURL(apiEndpoint).
		SetAuthToken(authResp.AccessToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	return &Client{
		apiEndpoint: apiEndpoint,
		token:       authResp.AccessToken,
		HTTPClient:  client,
	}, nil
}
