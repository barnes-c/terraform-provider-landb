// Copyright (c) Christopher Barnes <christopher@barnes.biz>
// SPDX-License-Identifier: MPL-2.0

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
	_, err := c.HTTPClient.R().
		SetBody([]Set{set}).
		SetResult(&result).
		Post(url)
	if err != nil {
		return Set{}, err
	}

	if len(result) != 1 {
		return Set{}, fmt.Errorf("unexpected number of sets returned after create: got %d, want 1", len(result))
	}

	return result[0], nil
}

func (c *Client) GetSet(name string) (*Set, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	resp, err := c.HTTPClient.R().
		SetResult(&Set{}).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Set), nil
}

func (c *Client) UpdateSet(name string, device Set) (*Set, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	resp, err := c.HTTPClient.R().
		SetBody(device).
		SetResult(&Set{}).
		Put(url)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Set), nil
}

func (c *Client) DeleteSet(name string, version int) error {
	url := fmt.Sprintf("%s%s%s", landbURL, setsURL, name)

	_, err := c.HTTPClient.R().
		SetQueryParam("version", fmt.Sprintf("%d", version)).
		Delete(url)
	if err != nil {
		return err
	}

	return nil
}
