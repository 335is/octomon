package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/335is/log"
	"github.com/335is/octomon/internal/config"
	"github.com/335is/octomon/internal/health"
	"github.com/335is/octopus"
)

const (
	// AppName defines the prefix for any configuration environment variables, as in OCTOMON_OCTOPUS_ADDRESS
	appName    = "octomon"
	appVersion = "0.0.1"
)

var (
	appInstance         string
	lastCanceledTime    time.Time
	lastInterruptedTime time.Time
)

func init() {
	appInstance = uuid.NewV4().String()
}

func main() {
	log.Infof("Starting %s %s %s LOG_LEVEL=%s", appName, appVersion, appInstance, log.GetLevel().String())

	cfg := config.New()
	log.Debugf("Config settings: " + cfg.Dump())

	octo := octopus.New(cfg.Octopus.Address, cfg.Octopus.APIKey, &http.Client{})
	log.Infof("Monitoring %s", cfg.Octopus.Address)

	checker := health.NewChecker()
	checker.AddCheck("Version", version_check(octo, cfg.HealthCheck.Version))
	checker.AddCheck("StuckTasks", stucktasks_check(octo, cfg.HealthCheck.StuckTasks))
	checker.RunAsync(cfg.HealthCheck.Interval, 10)

	waitForExit()

	log.Infof("Stopping %s %s %s", appName, appVersion, appInstance)
	checker.Stop()

	log.Infof("Shutting down")
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	sig := <-sigs
	log.Infof("Received signal %s, exiting...", sig)
}

// version_check health checks the root API and gets the Octopus Deploy version
func version_check(o *octopus.Client, c *config.Version) func() (health.Status, string) {
	return func() (health.Status, string) {
		s, err := o.Version()
		if err != nil {
			return health.Failure, s
		}

		return health.Ok, s
	}
}

// stucktasks_check health checks the currently executing deploy tasks looking for any that are stuck
func stucktasks_check(o *octopus.Client, c *config.StuckTasks) func() (health.Status, string) {
	return func() (health.Status, string) {
		// check for at least one task stuck in canceling state beyond a warning or failure time duration threshold
		qp := map[string]string{
			"active": "true",
			"states": "Cancelling",
		}
		cTasks, err := o.GetFiltered(qp)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request get cancelling tasks to %s failed: %v", o.Address, err)
		}

		if len(cTasks) > 0 {
			if lastCanceledTime.IsZero() {
				lastCanceledTime = time.Now()
			}

			if time.Since(lastCanceledTime) > c.FailureCancellingDuration {
				return health.Failure, fmt.Sprintf("Task in canceled state beyond failure threshold (%s)", c.FailureCancellingDuration)
			}
			if time.Since(lastCanceledTime) > c.WarningCancellingDuration {
				return health.Warning, fmt.Sprintf("Task in canceled state beyond warning threshold (%s)", c.WarningCancellingDuration)
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
		iTasks, err := o.GetFiltered(qi)
		if err != nil {
			return health.Failure, fmt.Sprintf("HTTP request get interrupted tasks to %s failed: %v", o.Address, err)
		}

		if len(iTasks) > 0 {
			if lastInterruptedTime.IsZero() {
				lastInterruptedTime = time.Now()
			}

			if time.Since(lastInterruptedTime) > c.FailureInterruptedDuration {
				return health.Failure, fmt.Sprintf("Task in interrupted state beyond failure threshold (%s)", c.FailureInterruptedDuration)
			}
			if time.Since(lastInterruptedTime) > c.WarningInterruptedDuration {
				return health.Warning, fmt.Sprintf("Task in interrupted state beyond warning threshold (%s)", c.WarningInterruptedDuration)
			}
		} else {
			// reset last canceled time
			lastInterruptedTime = time.Time{}
		}

		return health.Ok, "No stuck tasks"
	}
}
