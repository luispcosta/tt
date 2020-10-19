package reporter

import (
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// EmptyReporter represents a reporter not supported
type EmptyReporter struct{}

// NewEmptyReporter creates a new empty reporter
func NewEmptyReporter() *EmptyReporter {
	empty := EmptyReporter{}
	return &empty
}

// Initialize initializes a new empty reporter
func (reporter *EmptyReporter) Initialize(repo core.ActivityRepository, period core.Period) error {
	return utils.NewReportNotImplementedError()
}

// ProduceReport no-op
func (reporter *EmptyReporter) ProduceReport() error {
	return utils.NewReportNotImplementedError()
}
