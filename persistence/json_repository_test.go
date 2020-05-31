package persistence

import (
	"fmt"
	"os"
	"testing"

	"github.com/luispcosta/go-tt/core"

	"github.com/luispcosta/go-tt/utils"
)

func assertConfigurationFolderExists(t *testing.T) {
	expectedPath := fmt.Sprintf("%s%s.gott%sdata", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator))
	folderExists, err := utils.PathExists(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !folderExists {
		t.Fatal(fmt.Sprintf("Expected folder %s to exist", expectedPath))
	}
}

func assertActivityFileExists(activity core.Activity, t *testing.T) {
	assertConfigurationFolderExists(t)
	expectedPath := fmt.Sprintf("%s%s.gott%sdata%s%s.json", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator), activity.Name)
	folderExists, err := utils.PathExists(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !folderExists {
		t.Fatal(fmt.Sprintf("Expected activity file path %s to exist", expectedPath))
	}
}

func TestInitializeWhenNoConfigurationExists(t *testing.T) {
	repo := NewJSONActivityRepository()
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
}

func TestInitializeWhenConfigFolderExists(t *testing.T) {
	utils.CreateDir(fmt.Sprintf("%s%s.gott%s", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator)))

	repo := NewJSONActivityRepository()
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
}

func TestInitializeWhenDataFolderExists(t *testing.T) {
	utils.CreateDir(fmt.Sprintf("%s%s.gott%sdata", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator)))

	repo := NewJSONActivityRepository()
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
}

func TestUpdateActivityWhenActivityFileDoesNotExist(t *testing.T) {
	repo := NewJSONActivityRepository()
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	activity := core.Activity{}
	activity.Name = "Some_Name"
	errUpdating := repo.Update(activity)
	if errUpdating != nil {
		t.Fatal(errUpdating)
	}

	assertActivityFileExists(activity, t)
}

func TestUpdateActivityWhenActivityFileExists(t *testing.T) {
	utils.CreateDir(fmt.Sprintf("%s%s.gott%sdata", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator)))
	activity := core.Activity{}
	activity.Name = "a"

	homeDir := fmt.Sprintf("%s%s.gott%sdata", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator))

	filePath := fmt.Sprintf("%s%s%s.json", homeDir, string(os.PathSeparator), activity.Name)
	f, errCreate := os.Create(filePath)
	defer f.Close()

	if errCreate != nil {
		t.Fatal(errCreate)
	}

	repo := NewJSONActivityRepository()
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	errUpdating := repo.Update(activity)
	if errUpdating != nil {
		t.Fatal(errUpdating)
	}

	assertActivityFileExists(activity, t)
}
