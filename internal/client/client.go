package landb

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiEndpoint string
	token       string
	httpClient  *resty.Client
}

func NewClient(authEndpoint, apiEndpoint, clientID, clientSecret string) (*Client, error) {
	const audience = "production-microservice-landb-rest"

	token, err := authenticate(authEndpoint, clientID, clientSecret, audience)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &Client{
		apiEndpoint: apiEndpoint,
		token:       token,
		httpClient:  resty.New(),
	}, nil
}

// GetDevice retrieves a device by its ID.
func (c *Client) GetDevice(deviceID string) (*Device, error) {
	resp, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.token).
		SetResult(&Device{}).
		Get(fmt.Sprintf("%s/api/v1/devices/%s", c.apiEndpoint, deviceID))

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s - %s", resp.Status(), resp.String())
	}

	return resp.Result().(*Device), nil
}
