package health

import (
	"time"
)

// Check defines a health check
type Check struct {
	Name   string
	Exec   Exec
	result *Result
}

// Result stores the latest health check result
type Result struct {
	LastStatus  Status
	LastMessage string
	LastCheck   time.Time
	Count       int64
}

// Exec defines the health check interface
type Exec func() (Status, string)

// NewCheck creates a fresh Check
func NewCheck(name string, exec Exec) *Check {
	return &Check{
		Name:   name,
		Exec:   exec,
		result: &Result{},
	}
}

// GetResult returns a copy of the health check's current Result
func (c *Check) GetResult() *Result {
	cpy := (*c.result)
	return &cpy
}
