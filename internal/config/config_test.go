package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Octopus)
	assert.Equal(t, "https://demo.octopusdeploy.com", cfg.Octopus.Address)
	assert.Equal(t, "API-GUEST", cfg.Octopus.APIKey)
	assert.Equal(t, time.Minute, cfg.HealthCheck.Interval)
}

var yml = `
---
octopus:
   address: http://example.com/
   apikey: API-DEADBEEFBAADFOODFAADDEEDED
healthcheck:
   interval: 5m30s
`
var badYml = "This is really NOT YAML."

// FromYaml - good YAML
func TestFromYaml(t *testing.T) {
	cfg, err := FromYaml([]byte(yml))
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Octopus)
	assert.Equal(t, "http://example.com/", cfg.Octopus.Address)
	assert.Equal(t, "API-DEADBEEFBAADFOODFAADDEEDED", cfg.Octopus.APIKey)
	assert.Equal(t, (5*time.Minute)+(30*time.Second), cfg.HealthCheck.Interval)
}

// FromYaml - bad YAML
func TestFromYamlBad(t *testing.T) {
	cfg, err := FromYaml([]byte(badYml))
	assert.NotNil(t, err, "Expected an error become of bad YAML")
	assert.Nil(t, cfg)
}

// FromYamlFile - good YAML file
func TestFromYamlFile(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml_test")
	assert.Nil(t, err, "Got error trying to create temporary YAML file")
	assert.NotNil(t, file, "Failed to create temporary YAML file")
	defer os.Remove(file.Name())

	l, err := file.Write([]byte(yml))
	file.Close()
	assert.Nil(t, err, "Got error trying to write contents to temporary YAML file")
	assert.Equal(t, len(yml), l, "Mismatched bytes written to temporary YAML file")
	assert.FileExists(t, file.Name(), "Temporary YAML file doesn't exist")

	cfg, err := FromYamlFile(file.Name())
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Octopus)
	assert.Equal(t, "http://example.com/", cfg.Octopus.Address)
	assert.Equal(t, "API-DEADBEEFBAADFOODFAADDEEDED", cfg.Octopus.APIKey)
	assert.Equal(t, (5*time.Minute)+(30*time.Second), cfg.HealthCheck.Interval)
}

// FromYamlFile - bad YAML file
func TestFromYamlFileBad(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml_test_bad")
	assert.Nil(t, err, "Got error trying to create temporary YAML file")
	assert.NotNil(t, file, "Failed to create temporary YAML file")
	defer os.Remove(file.Name())

	l, err := file.Write([]byte(badYml))
	file.Close()
	assert.Nil(t, err, "Got error trying to write contents to temporary YAML file")
	assert.Equal(t, len(badYml), l, "Mismatched bytes written to temporary YAML file")
	assert.FileExists(t, file.Name(), "Temporary YAML file doesn't exist")

	cfg, err := FromYamlFile(file.Name())
	assert.NotNil(t, err, "Expected an error")
	assert.Nil(t, cfg)
}

// FromYamlFile - missing YAML file
func TestFromYamlFileNoFile(t *testing.T) {
	cfg, err := FromYamlFile("bogus_file_name")
	assert.NotNil(t, err, "Expected an error")
	assert.Nil(t, cfg)
}

// FromEnvironment - env vars exist
func TestFromEnvironment(t *testing.T) {
	os.Setenv("OCTOMON_OCTOPUS_ADDRESS", "http://example.com/")
	os.Setenv("OCTOMON_OCTOPUS_APIKEY", "API-DEADBEEFBAADFOODFAADDEEDED")
	os.Setenv("OCTOMON_HEALTHCHECK_INTERVAL", "10m")

	cfg, err := FromEnvironment("OCTOMON")
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Octopus)
	assert.Equal(t, "http://example.com/", cfg.Octopus.Address)
	assert.Equal(t, "API-DEADBEEFBAADFOODFAADDEEDED", cfg.Octopus.APIKey)
	assert.Equal(t, 10*time.Minute, cfg.HealthCheck.Interval)
}

// FromEnvironment - missing env vars
func TestFromEnvironmentBad(t *testing.T) {
	// explicitly make sure these don't exist
	os.Unsetenv("OCTOMON_OCTOPUS_ADDRESS")
	os.Unsetenv("OCTOMON_OCTOPUS_APIKEY")
	os.Unsetenv("OCTOMON_HEALTHCHECK_INTERVAL")

	cfg, err := FromEnvironment("OCTOMON")
	assert.NotNil(t, err, "Expected an error")
	assert.Nil(t, cfg)
}
