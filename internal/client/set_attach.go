// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package landb

import (
	"fmt"
	"time"
)

const setAttachmentURL = "beta/sets/%s/ip-addresses"

type SetAttachment struct {
    Name        string    `json:"name"`
    IPv4        string    `json:"ipv4"`
    IPv6        string    `json:"ipv6"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"_createdAt,omitempty"`
    UpdatedAt   time.Time `json:"_updatedAt,omitempty"`
}

func (c *Client) GetSetAttachments(setName string) ([]SetAttachment, error) {
    url := fmt.Sprintf("%s"+setAttachmentURL, landbURL, setName)

    var result []SetAttachment
    var apiErr APIError

    resp, err := c.HTTPClient.R().
        SetResult(&result).
        SetError(&apiErr).
        Get(url)
    if err != nil {
        return nil, err
    }
    if resp.IsError() {
        return nil, fmt.Errorf("list set attachments failed: %s", apiErr.Message)
    }
    return result, nil
}

func (c *Client) CreateSetAttachment(setName string, att SetAttachment) (SetAttachment, error) {
    url := fmt.Sprintf("%s"+setAttachmentURL, landbURL, setName)

    var result []SetAttachment
    var apiErr  APIError

    resp, err := c.HTTPClient.R().
        SetBody([]SetAttachment{att}).
        SetResult(&result).
        SetError(&apiErr).
        Post(url)
    if err != nil {
        return SetAttachment{}, err
    }
    if resp.IsError() {
        return SetAttachment{}, fmt.Errorf("create set attachment failed: %s", apiErr.Message)
    }
    return result[0], nil
}

func (c *Client) UpdateSetAttachment(setName, attachmentName string, att SetAttachment) (*SetAttachment, error) {
    url := fmt.Sprintf("%s"+setAttachmentURL+"/%s", landbURL, setName, attachmentName)

    var apiErr APIError
    resp, err := c.HTTPClient.R().
        SetBody(att).
        SetResult(&SetAttachment{}).
        SetError(&apiErr).
        Put(url)
    if err != nil {
        return nil, err
    }
    if resp.IsError() {
        return nil, fmt.Errorf("update set attachment failed: %s", apiErr.Message)
    }
    return resp.Result().(*SetAttachment), nil
}

func (c *Client) DeleteSetAttachment(setName, attachmentName string) error {
    url := fmt.Sprintf("%s"+setAttachmentURL+"/%s", landbURL, setName, attachmentName)

    var apiErr APIError
    resp, err := c.HTTPClient.R().
        SetError(&apiErr).
        Delete(url)
    if err != nil {
        return err
    }
    if resp.IsError() {
        return fmt.Errorf("delete set attachment failed: %s", apiErr.Message)
    }
    return nil
}
