package core

// ActivityDurationDayAggregation represents the total duration of an activity for a given date.
type ActivityDurationDayAggregation struct {
	Activity Activity
	Date     string
	Duration int
}
