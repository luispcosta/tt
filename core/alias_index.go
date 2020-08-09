package core

import (
	"errors"
	"fmt"
	"strings"
)

// AliasIndexData represents the data contained within the index
type AliasIndexData map[string]string

// AliasIndex contains the mapping between an alias string and an activity name
type AliasIndex struct {
	Data AliasIndexData
}

// NewAliasIndex returns a new instance of an alias index
func NewAliasIndex() *AliasIndex {
	index := AliasIndex{}
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
func (index *AliasIndex) Update(aliasKey string, aliasValue string) error {
	if _, ok := index.Data[strings.ToLower(aliasKey)]; ok {
		return fmt.Errorf("Activity Alias %s is already being used", aliasKey)
	}

	if len(aliasValue) <= 0 {
		return errors.New("Alias name is not a valid string (it has no content)")
	}

	if len(aliasKey) <= 0 {
		return errors.New("Alias key is not a valid string (it has no content)")
	}

	index.Data[strings.ToLower(aliasKey)] = aliasValue
	return nil
}
