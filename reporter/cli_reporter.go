package reporter

import (
	"errors"

	"github.com/luispcosta/go-tt/core"
)

// CliReporter is an activity reporter that presents activity information in the standard out
type CliReporter struct {
	activityRepo core.ActivityRepository
}

// NewCliReporter creates a new CLI reporter
func NewCliReporter() *CliReporter {
	cliReporter := CliReporter{}
	return &cliReporter
}

// Initialize initializes a new CLI reporter
func (reporter *CliReporter) Initialize(repo core.ActivityRepository) error {
	reporter.activityRepo = repo
	return nil
}

// ProduceReport creates a new cli report in the given period
func (reporter *CliReporter) ProduceReport(period core.Period) error {
	return errors.New("TODO")
}
