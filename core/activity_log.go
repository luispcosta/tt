package core

import (
	"time"
)

// ActivityLog represents one run of an activity in a given point in time
type ActivityLog struct {
	Id        int
	Date      string
	StartedAt *time.Time
	StoppedAt *time.Time
	Activity  Activity
}

func (log *ActivityLog) IsDone() bool {
	return log.StartedAt.Unix() > 0 && log.StoppedAt.Unix() > 0
}

// Duration returns the duration of the activity in hours
func (log *ActivityLog) Duration() float64 {
	startTs := time.Unix(log.StartedAt.Unix(), 0)
	endTs := time.Unix(log.StoppedAt.Unix(), 0)

	diff := endTs.Sub(startTs)
	return diff.Hours()
}
