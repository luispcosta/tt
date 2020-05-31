package persistence

import "github.com/luispcosta/go-tt/core"

// ActivityRepository is the generic interface that exposes methods to create and read activities from a store.
type ActivityRepository interface {
	Initialize() error
	Update(core.Activity) error
	List() []core.Activity
}
