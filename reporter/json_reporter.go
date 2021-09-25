package reporter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// JsonReporter is an activity reporter that exports activity information to a json file.
type JsonReporter struct {
	Period         core.Period
	Repo           core.ActivityRepository
	DurationFormat core.DurationFormat
	Clock          utils.Clock
}

// NewJsonReporter creates a new JSON reporter
func NewJsonReporter() *JsonReporter {
	jsonReporter := JsonReporter{
		DurationFormat: core.HumanDurationFormat{},
		Clock:          utils.NewLiveClock(),
	}
	return &jsonReporter
}

// Initialize initializes a new JSON reporter
func (reporter *JsonReporter) Initialize(repo core.ActivityRepository, period core.Period) error {
	reporter.Repo = repo
	reporter.Period = period
	return nil
}

// SetDurationFormat sets the duration formatter
func (reporter *JsonReporter) SetDurationFormat(f core.DurationFormat) {
	reporter.DurationFormat = f
}

// Struct example:
/*
	{
		'2020-10-10: {
			'act1': "123",
			'act2': "12",
			'act3': "121",
			'act4': "23",
			'act5': "10",
		},
		...
	}
*/
type jsonData map[string]map[string]string

// ProduceReport creates a new json report in the given period
func (reporter *JsonReporter) ProduceReport() error {
	logs, err := reporter.Repo.LogsForPeriod(reporter.Period)
	if err != nil {
		return err
	}

	data := make(jsonData)

	reporter.Period.ForEachDay(func(d time.Time) error {
		date := d.Format("2006-01-02")
		activityLogs := logs[date]

		if len(activityLogs) != 0 {
			actData := make(map[string]string)
			for i := range activityLogs {
				entry := activityLogs[i]
				actData[entry.Activity.Name] = reporter.DurationFormat.Format(entry.Duration)
			}
			data[date] = actData
		}

		return nil
	})

	fileName := fmt.Sprintf("report_%s_%s_%v.json", reporter.Period.Sd.Format("2006_01_02"), reporter.Period.Ed.Format("2006_01_02"), reporter.Clock.Now().Unix())
	fileData, _ := json.MarshalIndent(data, "", " ")
	utils.WriteToFile(fileName, fileData)

	return nil
}
