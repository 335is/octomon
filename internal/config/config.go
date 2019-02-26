package config

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v2"
)

// Config holds our configuration settings
type Config struct {
	Octopus     *Octopus     `yaml:"octopus"`
	HealthCheck *HealthCheck `yaml:"healthcheck"`
}

// Octopus holds our octopus settings
type Octopus struct {
	Address string `yaml:"address" required:"true"`
	APIKey  string `yaml:"apikey" required:"true"`
}

// HealthCheck holds health check settings
type HealthCheck struct {
	Interval time.Duration `yaml:"interval" default:"1m"`
}

// Default returns settings that work with a public demo Octopus Deploy server
func Default() *Config {
	cfg := Config{
		Octopus: &Octopus{
			Address: "https://demo.octopusdeploy.com",
			APIKey:  "API-GUEST",
		},
		HealthCheck: &HealthCheck{
			Interval: time.Minute,
		},
	}

	return &cfg
}

// FromYaml extracts settings from a YAML string
func FromYaml(yml []byte) (*Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(yml, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// FromYamlFile extracts settings from a YAML file
func FromYamlFile(path string) (*Config, error) {
	// read YAML text file into a string
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal from string to struct
	return FromYaml(yml)
}

// FromEnvironment extracts settings from environment variables.
// We expect them to be named like following:
//	OCTOMON_OCTOPUS_ADDRESS
//	OCTOMON_OCTOPUS_APIKEY
//	OCTOMON_HEALTHCHECK_INTERVAL
func FromEnvironment(appName string) (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(strings.ToUpper(appName), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
