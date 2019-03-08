package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"

	"github.com/335is/log"
	"github.com/335is/octomon/internal/config"
	"github.com/335is/octomon/internal/health"
	"github.com/335is/octomon/internal/octopus"
)

const (
	// AppName defines the prefix for any configuration environment variables, as in OCTOMON_OCTOPUS_ADDRESS
	appName    = "octomon"
	appVersion = "0.0.1"
)

var (
	appInstance string
)

func init() {
	appInstance = fmt.Sprintf("%s", uuid.NewV4())
}

func main() {
	log.Infof("Starting %s %s %s", appName, appVersion, appInstance)

	cfg := config.New(appName)
	octo := octopus.New(cfg.Octopus.Address, cfg.Octopus.APIKey, &http.Client{})
	log.Infof("Monitoring %s", cfg.Octopus.Address)

	checker := health.NewChecker()
	checker.AddCheck("Version", octo.Version(cfg.HealthCheck.Version))
	checker.AddCheck("StuckTasks", octo.StuckTasks(cfg.HealthCheck.StuckTasks))
	checker.RunAsync(cfg.HealthCheck.Interval, 10)

	waitForExit()

	log.Infof("Stopping %s %s %s", appName, appVersion, appInstance)
	checker.Stop()

	log.Infof("Shutting down")
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	sig := <-sigs
	log.Infof("Received signal %s, exiting...", sig)
}
