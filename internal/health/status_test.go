package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	assert.Equal(t, "Ok", Ok.String())
	assert.Equal(t, "Warning", Warning.String())
	assert.Equal(t, "Failure", Failure.String())
	assert.Equal(t, "Unknown", Status(42).String())
}
