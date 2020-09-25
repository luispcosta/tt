package core

// Reporter is an object responsible for producing activities reports
type Reporter interface {
	Initialize(ActivityRepository) error
	ProduceReport(Period) error
}
