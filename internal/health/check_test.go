package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCheck(t *testing.T) {
	name := "Test"
	exec := func() (Status, string) {
		return Ok, ""
	}
	check := NewCheck(name, exec)
	assert.NotNil(t, check)
	assert.Equal(t, name, check.Name)
	assert.Equal(t, Ok, check.result.LastStatus)
	assert.Equal(t, "", check.result.LastMessage)
	assert.True(t, check.result.LastCheck.IsZero())
	assert.Zero(t, check.result.Count)
}
