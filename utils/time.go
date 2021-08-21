package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/luispcosta/go-tt/core"
)

// CalcActivityLogDuration calculates the duration of an activity log
func CalcActivityLogDuration(log core.ActivityLog) (float64, error) {
	startTime := log.StartedAt.Unix()
	endTime := log.StoppedAt.Unix()

	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	duration := end.Sub(start).Seconds()
	return duration, nil
}

// SecondsToHours returns a string with the representation of a number that is in seconds, to hours
func SecondsToHours(seconds float64) string {
	return fmt.Sprintf("%.2f", seconds/3600)
}

func TimeToStandardFormat(time time.Time) string {
	year, month, day := time.Date()
	return fmt.Sprintf("%s-%s-%s", strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
}
