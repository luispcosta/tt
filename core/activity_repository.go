package core

import (
	"time"
)

// ActivityRepository is the generic interface that exposes methods to create and read activities from a store.
type ActivityRepository interface {
	Initialize() error
	Update(Activity) error
	List() []Activity
	Delete(string) error
	Start(Activity) error
	Stop(Activity) error
	Find(string) (*Activity, error)
	FindLogsForDay(time.Time) (ActivityDayLog, error)
	Purge() error
	Backup(string) (string, error)
}
