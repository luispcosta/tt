package reporter

import (
	"errors"
	"fmt"
	"sort"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// CliReporter is an activity reporter that presents activity information in the standard out
type CliReporter struct {
	Period  core.Period
	Repo    core.ActivityRepository
	Printer func(...interface{}) (int, error)
}

type activityDataCache struct {
	activity      *core.Activity
	totalDuration float64
}

type cacheEntryWithFormat struct {
	data activityDataCache
	line string
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
	eachPeriodDayReports, cache, overallDurationInPeriod, err := reporter.findLogsInPeriod()

	if err != nil {
		return err
	}

	overallPeriodReportLines := reporter.createOverallPeriodActivitiesReportSortedByDuration(overallDurationInPeriod, cache)

	for _, entry := range overallPeriodReportLines {
		reporter.Printer(entry.line)
	}

	reporter.Printer("\n")

	for _, line := range eachPeriodDayReports {
		reporter.Printer(line)
	}

	return nil
}

func (reporter *CliReporter) createOverallPeriodActivitiesReportSortedByDuration(totalPeriodDuration float64, periodDataCache map[string]activityDataCache) []cacheEntryWithFormat {
	values := []cacheEntryWithFormat{}

	headerLine1 := fmt.Sprintf("Number of activities found: %v\n", len(periodDataCache))
	reporter.Printer(headerLine1)
	reporter.Printer("Total time spent\n")

	for _, data := range periodDataCache {
		durationInHours := utils.SecondsToHours(data.totalDuration)
		totalDurationPercentage := utils.Percentage(data.totalDuration, totalPeriodDuration)
		line := fmt.Sprintf("  %s (%s) - %vs/%vh - (%v)\n", data.activity.Name, data.activity.Description, data.totalDuration, durationInHours, totalDurationPercentage)
		values = append(values, cacheEntryWithFormat{data: data, line: line})
	}

	f := func(n, n1 int) bool {
		return values[n].data.totalDuration > values[n1].data.totalDuration
	}

	sort.Slice(values, f)

	return values
}

type dayActivity struct {
	activity *core.Activity
	duration float64
}

func (reporter *CliReporter) findLogsInPeriod() ([]string, map[string]activityDataCache, float64, error) {
	numberOfDays := reporter.Period.NumberOfDays()
	sd := reporter.Period.Sd
	cache := make(map[string]activityDataCache)
	var overallDurationInPeriod float64
	eachPeriodDayReport := []string{}
	numberDaysWithoutReports := 0

	for i := 0; i <= numberOfDays; i++ {
		date := sd.AddDate(0, 0, i)
		log, err := reporter.Repo.FindLogsForDay(date)
		if err != nil {
			if _, ok := err.(*utils.NotFoundError); ok {
				numberDaysWithoutReports++
				continue
			} else {
				return eachPeriodDayReport, cache, overallDurationInPeriod, err
			}
		}

		activitiesForThisDay := []dayActivity{}

		for k, v := range log {
			var totalDuration float64
			for _, dayLog := range v {
				totalDuration += dayLog.Duration
			}

			overallDurationInPeriod += totalDuration

			if val, has := cache[k]; !has {
				activity, errFind := reporter.Repo.Find(k)
				if errFind != nil {
					if _, ok := errFind.(*utils.NotFoundError); ok {
						fmt.Printf("Did not found activity metadata for activity log with activity name: %s\n", k)
						continue
					} else {
						return eachPeriodDayReport, cache, overallDurationInPeriod, errFind
					}
				}
				cache[k] = activityDataCache{activity: activity, totalDuration: totalDuration}
			} else {
				summedDuration := val.totalDuration + totalDuration
				val.totalDuration = summedDuration
				cache[k] = val
			}
			activitiesForThisDay = append(activitiesForThisDay, dayActivity{activity: cache[k].activity, duration: totalDuration})
		}

		if len(activitiesForThisDay) > 0 {
			f := func(n, n1 int) bool {
				return activitiesForThisDay[n].duration > activitiesForThisDay[n1].duration
			}

			sort.Slice(activitiesForThisDay, f)

			eachPeriodDayReport = append(eachPeriodDayReport, fmt.Sprintf("Date: %s\n", date.Format("2006-01-02")))
			for _, activity := range activitiesForThisDay {
				duration := activity.duration
				eachPeriodDayReport = append(eachPeriodDayReport, fmt.Sprintf("  Activity: %s\n", activity.activity.Name))
				eachPeriodDayReport = append(eachPeriodDayReport, fmt.Sprintf("  Total Duration: %vs/%vh\n", duration, utils.SecondsToHours(duration)))
			}
		}
	}

	if numberDaysWithoutReports-1 == numberOfDays {
		return eachPeriodDayReport, cache, overallDurationInPeriod, errors.New("No activity data found in period")
	}

	return eachPeriodDayReport, cache, overallDurationInPeriod, nil
}
