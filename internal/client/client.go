// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package landb

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	HTTPClient   *resty.Client
	clientID     string
	clientSecret string
	audience     string
}

type APIError struct {
	Code      string `json:"code"`
	ErrorType string `json:"error"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

var ErrDeleteNotSupported = errors.New("delete operation not supported by API")

func NewClient(apiURL, clientID, clientSecret, audience string) (*Client, error) {
	client := &Client{
		HTTPClient:   resty.New(),
		clientID:     clientID,
		clientSecret: clientSecret,
		audience:     audience,
	}

	client.HTTPClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		authResp, err := Authenticate(client.clientID, client.clientSecret, client.audience)
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		r.SetAuthToken(authResp.AccessToken)

		return nil
	})

	client.HTTPClient.SetDebug(true)

	return client, nil
}
