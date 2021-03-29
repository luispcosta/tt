package utils

import (
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/core"
)

// CalcActivityLogDuration calculates the duration of an activity log
func CalcActivityLogDuration(log core.ActivityLog) (float64, error) {
	startTime := log.Start
	endTime := log.End

	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	duration := end.Sub(start).Seconds()
	return duration, nil
}

// SecondsToHours returns a string with the representation of a number that is in seconds, to hours
func SecondsToHours(seconds float64) string {
	return fmt.Sprintf("%.2f", seconds/3600)
}
