package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewChecker(t *testing.T) {
	checker := NewChecker()

	assert.NotNil(t, checker)
}

func TestAddCheck(t *testing.T) {
	name := "check1"
	exec := func() (Status, string) { return Ok, "" }
	checker := NewChecker()
	check := checker.AddCheck(name, exec)

	assert.NotNil(t, check)
	assert.Equal(t, name, check.Name)
	assert.Zero(t, check.result.LastStatus)
	assert.Equal(t, "", check.result.LastMessage)
	assert.True(t, check.result.LastCheck.IsZero())
	assert.Zero(t, check.result.Count)
}

func TestGetResultBadCheckName(t *testing.T) {
	exec := func() (Status, string) { return Ok, "" }
	checker := NewChecker()
	checker.AddCheck("name", exec)
	checker.DoChecks()
	result := checker.GetResult("not name")

	assert.Nil(t, result)
}

func TestDoChecks(t *testing.T) {
	name := "Good Test"
	msg := "My test is good!"
	exec := func() (Status, string) { return Ok, msg }
	checker := NewChecker()
	checker.AddCheck(name, exec)
	checker.DoChecks()
	result := checker.GetResult(name)

	assert.NotNil(t, result)
	assert.Equal(t, Ok, result.LastStatus)
	assert.Equal(t, msg, result.LastMessage)
	assert.False(t, result.LastCheck.IsZero())
	assert.Equal(t, int64(1), result.Count)
}

func TestRunZero(t *testing.T) {
	name := "No Test"
	msg := "None shall pass!!!"
	var count int64
	exec := func() (Status, string) { return Failure, msg }
	checker := NewChecker()
	checker.AddCheck(name, exec)
	checker.Run(time.Millisecond, count)
	result := checker.GetResult(name)

	assert.NotNil(t, result)
	assert.Equal(t, Ok, result.LastStatus)
	assert.Equal(t, "", result.LastMessage)
	assert.True(t, result.LastCheck.IsZero())
	assert.Equal(t, count, result.Count)
}

func TestRunFive(t *testing.T) {
	name := "Warning Test"
	msg := "My test has a warning"
	var count int64 = 5
	exec := func() (Status, string) { return Warning, msg }
	checker := NewChecker()
	checker.AddCheck(name, exec)
	checker.Run(time.Millisecond, count)
	checker.Stop()

	result := checker.GetResult(name)
	assert.NotNil(t, result)
	assert.Equal(t, Warning, result.LastStatus)
	assert.Equal(t, msg, result.LastMessage)
	assert.False(t, result.LastCheck.IsZero())
	assert.Equal(t, count, result.Count)
}

func TestStop(t *testing.T) {
	name := "Test42"
	msg := "Somewhere out there..."
	interval := time.Millisecond
	exec := func() (Status, string) { return Ok, msg }
	checker := NewChecker()
	checker.AddCheck(name, exec)
	checker.RunAsync(interval, -1)
	time.Sleep(5 * interval)
	checker.Stop()

	result := checker.GetResult(name)
	assert.NotNil(t, result)
	assert.Equal(t, Ok, result.LastStatus)
	assert.Equal(t, msg, result.LastMessage)
	assert.False(t, result.LastCheck.IsZero())
	assert.NotZero(t, result.Count)
}

func TestStopNoAsync(t *testing.T) {
	var count int64 = 1
	exec := func() (Status, string) { return Ok, "" }
	checker := NewChecker()
	checker.AddCheck("", exec)
	checker.Run(time.Millisecond, 1)
	checker.Stop()

	result := checker.GetResult("")
	assert.NotNil(t, result)
	assert.Equal(t, Ok, result.LastStatus)
	assert.Equal(t, "", result.LastMessage)
	assert.False(t, result.LastCheck.IsZero())
	assert.Equal(t, count, result.Count)
}
