package reporter

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// CsvReporter is an activity reporter that exports activity information to a csv file.
type CsvReporter struct {
	Period         core.Period
	Repo           core.ActivityRepository
	DurationFormat core.DurationFormat
	Clock          utils.Clock
}

// NewCsvReporter creates a new CSV reporter
func NewCsvReporter() *CsvReporter {
	csvReporter := CsvReporter{
		DurationFormat: core.HumanDurationFormat{},
		Clock:          utils.NewLiveClock(),
	}
	return &csvReporter
}

// Initialize initializes a new CSV reporter
func (reporter *CsvReporter) Initialize(repo core.ActivityRepository, period core.Period) error {
	reporter.Repo = repo
	reporter.Period = period
	return nil
}

// SetDurationFormat sets the duration formatter
func (reporter *CsvReporter) SetDurationFormat(f core.DurationFormat) {
	reporter.DurationFormat = f
}

// ProduceReport creates a new CSV report in the given period
func (reporter *CsvReporter) ProduceReport() error {
	logs, err := reporter.Repo.LogsForPeriod(reporter.Period)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("report_%s_%s_%v.csv", reporter.Period.Sd.Format("2006_01_02"), reporter.Period.Ed.Format("2006_01_02"), reporter.Clock.Now().Unix())
	file, err := utils.NewFile(fileName)
	defer func() error {
		err := file.Close()
		if err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	reporter.Period.ForEachDay(func(d time.Time) error {
		date := d.Format("2006-01-02")
		activityLogs := logs[date]

		if len(activityLogs) != 0 {
			for i := range activityLogs {
				entry := activityLogs[i]
				row := []string{date, entry.Activity.Name, reporter.DurationFormat.Format(entry.Duration)}
				err := w.Write(row)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return nil
}
