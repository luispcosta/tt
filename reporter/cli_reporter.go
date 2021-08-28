package reporter

import (
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// CliReporter is an activity reporter that presents activity information in the standard out
type CliReporter struct {
	Period  core.Period
	Repo    core.ActivityRepository
	Printer func(...interface{}) (int, error)
}

// NewCliReporter creates a new CLI reporter
func NewCliReporter() *CliReporter {
	cliReporter := CliReporter{Printer: fmt.Print}
	return &cliReporter
}

// NewCustomCLIReporter creates a new CLI reporter
func NewCustomCLIReporter(printer func(...interface{}) (int, error)) *CliReporter {
	cliReporter := CliReporter{
		Printer: printer,
	}
	return &cliReporter
}

// Initialize initializes a new CLI reporter
func (reporter *CliReporter) Initialize(repo core.ActivityRepository, period core.Period) error {
	reporter.Repo = repo
	reporter.Period = period
	return nil
}

// ProduceReport creates a new cli report in the given period
func (reporter *CliReporter) ProduceReport() error {
	logs, err := reporter.Repo.LogsForPeriod(reporter.Period)
	if err != nil {
		return err
	}

	reporter.Period.ForEachDay(func(d time.Time) error {
		date := d.Format("2006-01-02")
		header := fmt.Sprintf("Day %s: \n", date)
		activityLogs := logs[date]

		var content string
		if len(activityLogs) == 0 {
			content = "  No activities found for this day"
		} else {
			for i := range activityLogs {
				entry := activityLogs[i]

				content += fmt.Sprintf("  Activity %s", entry.Activity.Name)
				content += fmt.Sprintf(" %v", utils.SecondsToHuman(entry.Duration))
				content += "\n"
			}
		}
		reporter.Printer(header)
		reporter.Printer(content)
		reporter.Printer("\n\n")

		return nil
	})

	return nil
}
