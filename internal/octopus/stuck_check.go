package octopus

import (
	"encoding/json"
	"fmt"

	"github.com/335is/octomon/internal/health"
)

// Tasks holds metrics on the deployment tasks
type Tasks struct {
	ItemType       string
	TotalResults   int
	ItemsPerPage   int
	NumberOfPages  int
	LastPageNumber int
	TotalCounts    TotalCounts
}

// TotalCounts holds metrics on currently running deployment tasks
type TotalCounts struct {
	Canceled    int
	Cancelling  int
	Executing   int
	Failed      int
	Queued      int
	Success     int
	TimedOut    int
	Interrupted int
}

// StuckTasks health checks the currently executing deploy tasks looking for any that are stuck
func (c *Client) StuckTasks() func() (health.Status, string) {
	return func() (health.Status, string) {
		qp := map[string]string{}
		response, err := c.DoGetRequest(qp)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request to %s failed: %v", c.Address, err)
		}

		tasks := Tasks{}
		err = json.NewDecoder(response.Body).Decode(&tasks)
		if err != nil {
			return health.Failure, fmt.Sprintf("failed to decode JSON response from %s: %v", c.Address, err)
		}

		// ToDo: add actual checks for stuck deploy tasks

		return health.Ok, ""
	}
}
