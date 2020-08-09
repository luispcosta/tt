package core

import (
	"fmt"
)

// AliasIndexData represents the data contained within the index
type AliasIndexData map[string]string

// AliasIndex contains the mapping between an alias string and an activity name
type AliasIndex struct {
	Data       AliasIndexData
	Repository ActivityRepository
}

// NewAliasIndex returns a new instance of an alias index
func NewAliasIndex(repository ActivityRepository) *AliasIndex {
	index := AliasIndex{}
	index.Repository = repository
	index.Data = AliasIndexData{}
	return &index
}

// Load loads data into the index
func (index *AliasIndex) Load(data AliasIndexData) {
	index.Data = data
}

// Delete deletes an entry from the index
func (index *AliasIndex) Delete(indexKey string) {
	delete(index.Data, indexKey)
}

// Update updates the index value for an activity
func (index *AliasIndex) Update(activity Activity) error {
	if _, ok := index.Data[activity.Alias]; ok {
		return fmt.Errorf("Activity Alias %s is already being used", activity.Alias)
	}

	index.Data[activity.Alias] = index.Repository.IndexKey(activity)
	return nil
}
