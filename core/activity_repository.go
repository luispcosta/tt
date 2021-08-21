package core

import "time"

// ActivityRepository is the generic interface that exposes methods to create and read activities from a store.
type ActivityRepository interface {
	Initialize() error
	Shutdown() error
	Add(Activity) error
	Delete(string) error
	List() ([]Activity, error)
	Update(string, UpdateActivity) error
	Find(string) (*Activity, error)
	Start(Activity) error
	LogsForDay(time.Time) ([]ActivityLog, error)
	Stop(Activity) error
}
