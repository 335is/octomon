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
	Address string `yaml:"address" default:"https://demo.octopus.com"`
	APIKey  string `yaml:"apikey" default:"API-GUEST"`
}

// HealthCheck holds health check settings
type HealthCheck struct {
	Interval   time.Duration `yaml:"interval" default:"60s"`
	Version    *octopus.Version
	StuckTasks *octopus.StuckTasks
}

// New starts with a default config that works with a public demo Octopus Deploy server,
// then overrides settings from a YAML config file and env vars if they exist.
func New() *Config {
	c := Config{
		Octopus: &Octopus{},
		HealthCheck: &HealthCheck{
			Version:    &octopus.Version{},
			StuckTasks: &octopus.StuckTasks{},
		},
	}

	cfg.Load("", &c)

	return &c
}

// Dump returns the Octopus configuration in YAML string form
func (c *Config) Dump() string {
	s, _ := cfg.ToYaml(c)
	return s
}
