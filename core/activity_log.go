package core

import "go.mongodb.org/mongo-driver/bson/primitive"

// ActivityLog represents one run of an activity in a given point in time
type ActivityLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ActivityID primitive.ObjectID `bson:"ActivityID,omitempty"`
	Date       string
	Start      int64
	End        int64
}
