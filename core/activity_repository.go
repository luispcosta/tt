package core

// ActivityRepository is the generic interface that exposes methods to create and read activities from a store.
type ActivityRepository interface {
	Initialize() error
	Shutdown() error
	Add(Activity) error
	Update(Activity) error
	List() ([]Activity, error)
	Delete(string) error
	Start(Activity) error
	Stop(Activity) error
	Find(string) (*Activity, error)
	Purge() error
	Backup(string) (string, error)
	Restore(string) error
	SchemaMigrate(string) error
}
