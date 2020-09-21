package reporter

import (
	"errors"
	"time"

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
func (reporter *EmptyReporter) ProduceReport(startDate time.Time, endDate time.Time) error {
	return errors.New("Non implementation of a reporter")
}
