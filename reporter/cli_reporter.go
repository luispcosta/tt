package reporter

import (
	"fmt"

	"github.com/luispcosta/go-tt/core"
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
	/*reporter.Period.ForEachDay(func(d time.Time) error {
		header := fmt.Sprintf("Day %s: \n", utils.TimeToStandardFormat(d))
		activityLogs, err := reporter.Repo.LogsForDay(d)
		if err != nil {
			return err
		}

		var content string
		if len(activityLogs) == 0 {
			content = "\t\tNo activities found for this day"
		} else {
			for i := range activityLogs {
				content += fmt.Sprintf("\tActivity %s", act.Name)
				for _, log := range logs {
					content += fmt.Sprintf(" %vh", log.Duration())
				}
				content += "\n"
			}
		}
		reporter.Printer(header)
		reporter.Printer(content)
		reporter.Printer("\n\n")

		return nil
	})

	return nil*/
}
