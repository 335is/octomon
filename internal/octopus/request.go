package octopus

import (
	"fmt"
	"net/http"
)

// Links are part of each JSON response, used for pagination
type Links map[string]string

// DoGetRequest makes a Get API call to Octopus
func (c *Client) DoGetRequest(api string, queryParams map[string]string) (*http.Response, error) {
	// build up the URL
	url := c.Address + "/api/" + api
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set the query parameters
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// HTTP request header required to make Octopus Deploy API calls
	req.Header.Add("X-Octopus-ApiKey", c.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP status %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp, nil
}
