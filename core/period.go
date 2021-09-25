package core

import (
	"strings"
	"time"

	"github.com/luispcosta/go-tt/utils"
)

// Period represents a time period
type Period struct {
	Sd time.Time
	Ed time.Time
}

// NumberOfDays returns the number of days in the period
func (period *Period) NumberOfDays() int {
	if dateEqual(period.Sd, period.Ed) {
		return 1
	}

	return int(period.Ed.Sub(period.Sd).Hours() / 24)
}

// PeriodFromDateStrings returns a new period struct from two date strings, if they are valid
func PeriodFromDateStrings(sd, ed string) (Period, error) {
	parsedStDate, err1 := parseSimpleDate(sd)
	if err1 != nil {
		return Period{}, err1
	}
	parsedEdDate, err2 := parseSimpleDate(ed)
	if err2 != nil {
		return Period{}, err2
	}

	if parsedStDate.After(*parsedEdDate) {
		return Period{Sd: *parsedEdDate, Ed: *parsedStDate}, nil
	}

	return Period{Sd: *parsedStDate, Ed: *parsedEdDate}, nil
}

const lastDayPeriod = "day"
const lastWeekPeriod = "week"
const lastMonthPeriod = "month"
const lastYearPeriod = "year"

// PeriodFromKeyWord returns a fixed period from a representation string, relative to the current date.
func PeriodFromKeyWord(keyword string) Period {
	now := time.Now()
	ed := now
	var sd time.Time
	switch strings.ToLower(keyword) {
	case lastDayPeriod:
		sd = now.AddDate(0, 0, 0)
	case lastWeekPeriod:
		sd = now.AddDate(0, 0, -7)
	case lastMonthPeriod:
		sd = now.AddDate(0, -1, 0)
	case lastYearPeriod:
		sd = now.AddDate(-1, 0, 0)
	default:
		sd = now.AddDate(0, 0, -1)
	}

	return Period{Sd: sd, Ed: ed}
}

func (period *Period) ForEachDay(fn func(time.Time) error) {
	if period.NumberOfDays() == 1 {
		fn(period.Sd)
	} else {
		i := 0
		for {
			nextDay := period.Sd.AddDate(0, 0, i)
			if nextDay == period.Ed.AddDate(0, 0, 1) {
				break
			}
			fn(nextDay)
			i += 1
		}
	}
}

func (period *Period) StartDateDay() string {
	return period.Sd.Format(utils.DateFormat)
}

func (period *Period) EndDateDay() string {
	return period.Ed.Format(utils.DateFormat)
}

// AllowedPeriodFixedTimeFrames returns an array of allowed period fixed time frames
func AllowedPeriodFixedTimeFrames() []string {
	return []string{lastDayPeriod, lastWeekPeriod, lastMonthPeriod, lastYearPeriod}
}

func parseSimpleDate(date string) (*time.Time, error) {
	parsedDate, err := time.Parse(utils.DateFormat, date)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
