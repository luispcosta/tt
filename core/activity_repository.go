package core

import "github.com/luispcosta/go-tt/config"

// ActivityRepository is the generic interface that exposes methods to create and read activities from a store.
type ActivityRepository interface {
	Initialize(config.Config) error
	Shutdown() error
	Add(Activity) error
	Delete(string) error
	List() ([]Activity, error)
	Update(string, UpdateActivity) error
	Find(string) (*Activity, error)
	Start(Activity) error
	LogsForPeriod(Period) (map[string][]ActivityDurationDayAggregation, error)
	Stop(Activity) error
	CurrentlyTrackedActivity() (*Activity, error)
}
