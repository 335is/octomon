package octopus

import (
	"net/http"
	"strings"
)

// Client is an Octopus Deploy API client
type Client struct {
	Address    string
	APIKey     string
	httpClient *http.Client
}

// New creates a new Octopus Deploy API client
func New(address, apiKey string, httpClient *http.Client) *Client {
	client := Client{
		Address:    strings.TrimRight(address, "/"),
		APIKey:     apiKey,
		httpClient: httpClient,
	}

	return &client
}
