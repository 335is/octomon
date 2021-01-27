package config

import (
	"time"

	cfg "github.com/335is/config"
	"github.com/335is/octomon/internal/octopus"
)

// Config holds our configuration settings
//	OCTOMON_OCTOPUS_ADDRESS
//	OCTOMON_OCTOPUS_APIKEY
//	OCTOMON_HEALTHCHECK_INTERVAL
//	OCTOMON_HEALTHCHECK_STUCKTASKS_WARNINGCANCELLINGDURATION
type Config struct {
	Octopus     *Octopus     `yaml:"octopus"`
	HealthCheck *HealthCheck `yaml:"healthcheck"`
}

// Octopus holds our octopus settings
type Octopus struct {
	Address string `yaml:"address"`
	APIKey  string `yaml:"apikey"`
}

// HealthCheck holds health check settings
type HealthCheck struct {
	Interval   time.Duration `yaml:"interval"`
	Version    *octopus.Version
	StuckTasks *octopus.StuckTasks
}

// New starts with a default config that works with a public demo Octopus Deploy server,
// then overrides settings from a YAML config file and env vars if they exist.
func New(appName string) *Config {
	c := Config{
		Octopus:     &Octopus{},
		HealthCheck: &HealthCheck{},
	}

	cfg.Load(appName, "", &c)

	return &c
}
