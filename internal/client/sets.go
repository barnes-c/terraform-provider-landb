// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package landb

import (
	"fmt"
)

const setsURL = "beta/sets/"

type Set struct {
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	NetworkDomain        string  `json:"networkDomain"`
	Responsible          Contact `json:"responsible"`
	Description          string  `json:"description"`
	ProjectURL           string  `json:"projectUrl"`
	ReceiveNotifications bool    `json:"receiveNotifications"`
	Version              int     `json:"version"`
}

func (c *Client) CreateSet(set Set) (Set, error) {
	url := fmt.Sprintf("%s%s", landbURL, setsURL)

	var result []Set
	var apiErr APIError

	resp, err := c.HTTPClient.R().
		SetBody([]Set{set}).
		SetResult(&result).
		SetError(&apiErr).
		Post(url)
	if err != nil {
		return Set{}, err
	}
	if resp.IsError() {
		return Set{}, fmt.Errorf("create set failed: %s", apiErr.Message)
	}
	return result[0], nil
}

func (c *Client) GetSet(name string) (*Set, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetResult(&Set{}).
		SetError(&apiErr).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get set failed: %s", apiErr.Message)
	}
	return resp.Result().(*Set), nil
}

func (c *Client) UpdateSet(name string, set Set) (*Set, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetBody(set).
		SetResult(&Set{}).
		SetError(&apiErr).
		Put(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("update set failed: %s", apiErr.Message)
	}
	return resp.Result().(*Set), nil
}

func (c *Client) DeleteSet(name string, version int) error {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetQueryParam("version", fmt.Sprintf("%d", version)).
		SetError(&apiErr).
		Delete(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("delete set failed: %s", apiErr.Message)
	}
	return nil
}
