package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityLog represents one run of an activity in a given point in time
type ActivityLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ActivityID primitive.ObjectID `bson:"ActivityID,omitempty"`
	Date       string
	Start      int64
	End        int64
}

func (log *ActivityLog) IsDone() bool {
	return log.Start > 0 && log.End > 0
}

// Duration returns the duration of the activity in hours
func (log *ActivityLog) Duration() float64 {
	startTs := time.Unix(log.Start, 0)
	endTs := time.Unix(log.End, 0)

	diff := endTs.Sub(startTs)
	return diff.Hours()
}
