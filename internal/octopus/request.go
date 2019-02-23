package octopus

import (
	"fmt"
	"net/http"
)

// DoGetRequest makes a Get API call to Octopus
func (c *Client) DoGetRequest(queryParams map[string]string) (*http.Response, error) {
	url := c.Address + "/api/"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// HTTP request header required to make Octopus Deploy API calls
	request.Header.Add("X-Octopus-ApiKey", c.APIKey)

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP status %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp, nil
}
