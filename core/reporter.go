package core

import "time"

// Reporter is an object responsible for producing activities reports
type Reporter interface {
	Initialize(ActivityRepository) error
	ProduceReport(time.Time, time.Time) error
}
