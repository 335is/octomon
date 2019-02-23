package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/335is/octomon/internal/config"
	"github.com/335is/octomon/internal/health"
	"github.com/335is/octomon/internal/octopus"
)

const (
	// AppName defines the prefix for any configuration environment variables, as in OCTOMON_OCTOPUS_ADDRESS
	appName    = "OCTOMON"
	appVersion = "0.0.1"
)

func main() {
	log.Printf("Starting %s %s", appName, appVersion)

	cfg, err := config.FromEnvironment(appName)
	if err != nil {
		log.Fatalf("Failed to load configuration from environment variables: %v", err)
	}

	octo := octopus.New(cfg.Octopus.Address, cfg.Octopus.APIKey, &http.Client{})

	checker := health.NewChecker()
	checker.AddCheck("Version", octo.Version())
	checker.AddCheck("Stuck Tasks", octo.StuckTasks())
	checker.RunAsync(cfg.HealthCheck.Interval, 10)

	waitForExit()

	checker.Stop()
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Received signal %s, exiting...", sig)
}
