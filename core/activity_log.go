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
