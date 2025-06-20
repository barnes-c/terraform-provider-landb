// SPDX-FileCopyrightText: 2025 CERN
//
// SPDX-License-Identifier: GPL-3.0-or-later

package landb

import (
	"fmt"
)

const devicesURL = "beta/devices/"

type Device struct {
	Name                 string          `json:"name"`
	SerialNumber         string          `json:"serialNumber"`
	InventoryNumber      string          `json:"inventoryNumber"`
	Tag                  string          `json:"tag"`
	Description          string          `json:"description"`
	Zone                 string          `json:"zone"`
	DHCPResponse         string          `json:"dhcpResponse"`
	IPv4InDNSAndFirewall bool            `json:"ipv4InDnsAndFirewall"`
	IPv6InDNSAndFirewall bool            `json:"ipv6InDnsAndFirewall"`
	ManagerLock          string          `json:"managerLock"`
	Ownership            string          `json:"ownership"`
	Location             Location        `json:"location"`
	Parent               string          `json:"parent"`
	Type                 string          `json:"type"`
	Manufacturer         string          `json:"manufacturer"`
	Model                string          `json:"model"`
	OperatingSystem      OperatingSystem `json:"operatingSystem"`
	Manager              Contact         `json:"manager"`
	Responsible          Contact         `json:"responsible"`
	User                 Contact         `json:"user"`
	Version              int             `json:"version"`
}

func (c *Client) CreateDevice(device Device) (Device, error) {
	url := fmt.Sprintf("%s%s", landbURL, devicesURL)

	var result []Device
	var apiErr APIError

	resp, err := c.HTTPClient.R().
		SetBody([]Device{device}).
		SetResult(&result).
		SetError(&apiErr).
		Post(url)
	if err != nil {
		return Device{}, err
	}

	if resp.IsError() {
		return Device{}, fmt.Errorf("create device failed: %s", apiErr.Message)
	}

	return result[0], nil
}

func (c *Client) GetDevice(name string) (*Device, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, devicesURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetResult(&Device{}).
		SetError(&apiErr).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("get device failed: %s", apiErr.Message)
	}

	return resp.Result().(*Device), nil
}

func (c *Client) UpdateDevice(name string, device Device) (*Device, error) {
	url := fmt.Sprintf("%s%s%s", landbURL, devicesURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetBody(device).
		SetResult(&Device{}).
		SetError(&apiErr).
		Put(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("update device failed: %s", apiErr.Message)
	}

	return resp.Result().(*Device), nil
}

func (c *Client) DeleteDevice(name string, version int) error {
	url := fmt.Sprintf("%s%s%s", landbURL, devicesURL, name)

	var apiErr APIError
	resp, err := c.HTTPClient.R().
		SetQueryParam("version", fmt.Sprintf("%d", version)).
		SetError(&apiErr).
		Delete(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("delete device failed: %s", apiErr.Message)
	}

	return nil
}
