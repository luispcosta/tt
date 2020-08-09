package core

import (
	"testing"
)

func TestDeleteWhenMapIsEmpty(t *testing.T) {
	aliasIndex := NewAliasIndex()
	aliasIndex.Delete("someKey")
	data := aliasIndex.Data
	if len(data) != 0 {
		t.Error("Index should remain empty")
	}
}

func TestDeleteWhenMapDoesNotContainGivenKey(t *testing.T) {
	aliasIndex := NewAliasIndex()
	data := AliasIndexData{
		"key1": "value1",
		"key2": "Value2",
	}
	aliasIndex.Load(data)
	aliasIndex.Delete("key3")
	if len(data) != 2 {
		t.Error("Index should have remained with 2 keys")
	}
}

func TestDeleteWhenMapDoesContainKeyButGivenBadCase(t *testing.T) {
	aliasIndex := NewAliasIndex()
	data := AliasIndexData{
		"key1": "value1",
		"key2": "Value2",
	}
	aliasIndex.Load(data)
	aliasIndex.Delete("kEy2")
	if len(data) != 2 {
		t.Error("Index should have remained with 2 keys")
	}
}

func TestDeleteWhenMapDoesContainKey(t *testing.T) {
	aliasIndex := NewAliasIndex()
	data := AliasIndexData{
		"key1": "value1",
		"key2": "Value2",
	}
	aliasIndex.Load(data)
	aliasIndex.Delete("key2")
	if len(data) != 1 {
		t.Error("One entry should have been deleted from the data")
	}
}

func TestUpdateWhenMapIsEmpty(t *testing.T) {
	aliasIndex := NewAliasIndex()
	err := aliasIndex.Update("key1", "val1")
	if err != nil {
		t.Error("Should not have failed on a correct update")
	}
	data := aliasIndex.Data
	if len(data) != 1 {
		t.Error("Alias index data should have increased by 1 after correct update")
	}
}

func TestUpdateWhenKeyAlreadyExists(t *testing.T) {
	aliasIndex := NewAliasIndex()
	errUpdate1 := aliasIndex.Update("key1", "val1")
	if errUpdate1 != nil {
		t.Error("Should not have failed on a correct update")
	}
	errUpdate2 := aliasIndex.Update("kEy1", "val2")
	if errUpdate2 == nil {
		t.Error("Should have failed on second wrong update with an existent key but different case")
	}
}

func TestUpdateWhenKeyIsEmpty(t *testing.T) {
	aliasIndex := NewAliasIndex()
	err := aliasIndex.Update("", "val1")
	if err == nil {
		t.Error("Should have failed when using a empty key alias key")
	}
}

func TestUpdateWhenKeyValueIsEmpty(t *testing.T) {
	aliasIndex := NewAliasIndex()
	err := aliasIndex.Update("xx", "")
	if err == nil {
		t.Error("Should have failed when using a empty key value")
	}
}

func TestUpdateWhenKeyAndValueIsEmpty(t *testing.T) {
	aliasIndex := NewAliasIndex()
	err := aliasIndex.Update("", "")
	if err == nil {
		t.Error("Should have failed when using a empty key and empty value")
	}
}
