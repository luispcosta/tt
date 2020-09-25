package reporter

import (
	"errors"

	"github.com/luispcosta/go-tt/core"
)

// EmptyReporter represents a reporter not supported
type EmptyReporter struct{}

// NewEmptyReporter creates a new empty reporter
func NewEmptyReporter() *EmptyReporter {
	empty := EmptyReporter{}
	return &empty
}

// Initialize initializes a new empty reporter
func (reporter *EmptyReporter) Initialize(repo core.ActivityRepository) error {
	return errors.New("This reporter cannot be initialized")
}

// ProduceReport no-op
func (reporter *EmptyReporter) ProduceReport(period core.Period) error {
	return errors.New("Non implementation of a reporter")
}
