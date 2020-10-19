package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/luispcosta/go-tt/core"
)

// CalcActivityLogDuration calculates the duration of an activity log
func CalcActivityLogDuration(log core.ActivityLog) (float64, error) {
	startTime, errConvStartTime := strconv.ParseInt(log.Start, 10, 64)
	if errConvStartTime != nil {
		return 0, errConvStartTime
	}

	endTime, errConvEndTime := strconv.ParseInt(log.End, 10, 64)

	if errConvEndTime != nil {
		return 0, errConvEndTime
	}

	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	duration := end.Sub(start).Seconds()
	return duration, nil
}

// SecondsToHours returns a string with the representation of a number that is in seconds, to hours
func SecondsToHours(seconds float64) string {
	return fmt.Sprintf("%.2f", seconds/3600)
}
