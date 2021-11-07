package octopus

import (
	"fmt"
	"time"

	"github.com/335is/octomon/internal/health"
)

// StuckTasks is the configuration for the stuck tasks health check
type StuckTasks struct {
	WarningCancellingDuration  time.Duration `yaml:"warning_cancelling_duration" default:"5m"`
	FailureCancellingDuration  time.Duration `yaml:"failure_cancelling_duration" default:"30m"`
	WarningInterruptedDuration time.Duration `yaml:"warning_interrupted_duration" default:"1h"`
	FailureInterruptedDuration time.Duration `yaml:"failure_interrupted_duration" default:"6h"`
}

var (
	lastCanceledTime    time.Time
	lastInterruptedTime time.Time
)

// StuckTasks health checks the currently executing deploy tasks looking for any that are stuck
func (c *Client) StuckTasks(cfg *StuckTasks) func() (health.Status, string) {
	return func() (health.Status, string) {
		// check for at least one task stuck in canceling state beyond a warning or failure time duration threshold
		qp := map[string]string{
			"active": "true",
			"states": "Cancelling",
		}
		cTasks, err := c.GetFiltered(qp)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request get cancelling tasks to %s failed: %v", c.Address, err)
		}

		if len(cTasks) > 0 {
			if lastCanceledTime.IsZero() {
				lastCanceledTime = time.Now()
			}

			if time.Since(lastCanceledTime) > cfg.FailureCancellingDuration {
				return health.Failure, fmt.Sprintf("Task in canceled state beyond failure threshold (%s)", cfg.FailureCancellingDuration)
			}
			if time.Since(lastCanceledTime) > cfg.WarningCancellingDuration {
				return health.Warning, fmt.Sprintf("Task in canceled state beyond warning threshold (%s)", cfg.WarningCancellingDuration)
			}
		} else {
			// reset last canceled time
			lastCanceledTime = time.Time{}
		}

		// check for at least one task stuck in interrupted state beyond a warning or failure time duration threshold
		qi := map[string]string{
			"active":                  "true",
			"hasPendingInterruptions": "true",
		}
		iTasks, err := c.GetFiltered(qi)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request get interrupted tasks to %s failed: %v", c.Address, err)
		}

		if len(iTasks) > 0 {
			if lastInterruptedTime.IsZero() {
				lastInterruptedTime = time.Now()
			}

			if time.Since(lastInterruptedTime) > cfg.FailureInterruptedDuration {
				return health.Failure, fmt.Sprintf("Task in interrupted state beyond failure threshold (%s)", cfg.FailureInterruptedDuration)
			}
			if time.Since(lastInterruptedTime) > cfg.WarningInterruptedDuration {
				return health.Warning, fmt.Sprintf("Task in interrupted state beyond warning threshold (%s)", cfg.WarningInterruptedDuration)
			}
		} else {
			// reset last canceled time
			lastInterruptedTime = time.Time{}
		}

		return health.Ok, "No stuck tasks"
	}
}
