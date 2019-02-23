package health

import (
	"context"
	"log"
	"sync"
	"time"
)

// Checker defines the health checker interface
type Checker interface {
	AddCheck(name string, exec Exec) *Check
	GetResult(name string) *Result
	Run(interval time.Duration, count int64)
	RunAsync(interval time.Duration, count int64)
	Stop()
	DoChecks()
}

type checker struct {
	// checks is the collection of health checks to execute
	checks    []*Check
	context   context.Context
	cancel    context.CancelFunc
	waitGroup sync.WaitGroup
}

// NewChecker creates a new health checker
func NewChecker() Checker {
	context, cancel := context.WithCancel(context.Background())
	r := checker{
		checks:  []*Check{},
		context: context,
		cancel:  cancel,
	}

	return &r
}

// AddCheck adds a health check to the collection
func (r *checker) AddCheck(name string, exec Exec) *Check {
	chk := NewCheck(name, exec)
	r.checks = append(r.checks, chk)

	return chk
}

// GetCheck returns a copy of a health check by name
func (r *checker) GetResult(name string) *Result {
	for _, check := range r.checks {
		if check.Name == name {
			r := check.GetResult()
			return r
		}
	}

	return nil
}

// Run is a blocking call that periodically executes all health checks
//		interval: time period between health checks
//		count: total number of health checks to make, negative number means infinite, zero skips it
func (r *checker) Run(interval time.Duration, count int64) {
	ticker := time.NewTicker(interval)
	for i := int64(0); count < 0 || i < count; i++ {
		r.DoChecks()
		<-ticker.C
	}
}

// Run is a go routine that periodically executes all health checks
//		interval: time period between health checks
//		count: total number of health checks to make, negative number means infinite
func (r *checker) RunAsync(interval time.Duration, count int64) {
	go func() {
		r.waitGroup.Add(1)
		ticker := time.NewTicker(interval)
		for i := int64(0); count < 0 || i < count; i++ {
			r.DoChecks()

			select {
			case <-r.context.Done():
				r.waitGroup.Done()
				return
			case <-ticker.C:
				continue
			}
		}
	}()
}

// Stop cancels the RunAsync go routine
func (r *checker) Stop() {
	r.cancel()
	r.waitGroup.Wait()
}

// DoChecks runs the set of health checks once and saves the result.
func (r *checker) DoChecks() {
	for _, check := range r.checks {
		check.result.LastCheck = time.Now()
		check.result.Count++
		check.result.LastStatus, check.result.LastMessage = check.Exec()

		log.Printf("%s health check: status=%s, message=%s", check.Name, check.result.LastStatus, check.result.LastMessage)
	}
}
