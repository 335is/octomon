package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"

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
	log.Printf("Starting %s %s %s", appName, appVersion, appInstance)

	cfg := config.New(appName)
	octo := octopus.New(cfg.Octopus.Address, cfg.Octopus.APIKey, &http.Client{})
	log.Printf("Monitoring %s", cfg.Octopus.Address)

	checker := health.NewChecker()
	checker.AddCheck("Version", octo.Version(cfg.HealthCheck.Version))
	checker.AddCheck("StuckTasks", octo.StuckTasks(cfg.HealthCheck.StuckTasks))
	checker.RunAsync(cfg.HealthCheck.Interval, 10)

	waitForExit()

	log.Printf("Stopping %s %s %s", appName, appVersion, appInstance)
	checker.Stop()

	log.Printf("Shutting down")
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Received signal %s, exiting...", sig)
}
