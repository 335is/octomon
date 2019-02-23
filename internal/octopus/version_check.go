package octopus

import (
	"encoding/json"
	"fmt"

	"github.com/335is/octomon/internal/health"
)

// Root describes the JSON response from the root API call
type Root struct {
	Application    string
	Version        string
	APIVersion     string `json:"ApiVersion"`
	InstallationID string `json:"InstallationId"`
	Links          map[string]string
}

// Version health checks the root API and gets the Octopus Deploy version
func (c *Client) Version() func() (health.Status, string) {
	return func() (health.Status, string) {
		qp := map[string]string{}
		response, err := c.DoGetRequest(qp)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request to %s failed: %v", c.Address, err)
		}

		root := Root{}
		err = json.NewDecoder(response.Body).Decode(&root)
		if err != nil {
			return health.Failure, fmt.Sprintf("failed to decode JSON response from %s: %v", c.Address, err)
		}

		status := fmt.Sprintf("%s %s %s %s", root.Application, root.Version, root.APIVersion, root.InstallationID)

		return health.Ok, status
	}
}
