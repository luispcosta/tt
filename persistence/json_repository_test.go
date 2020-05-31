package persistence

import (
	"fmt"
	"os"
	"testing"

	"github.com/luispcosta/go-tt/configuration"

	"github.com/luispcosta/go-tt/core"

	"github.com/luispcosta/go-tt/utils"
)

const testDataFolder = "test"

func assertConfigurationFolderExists(t *testing.T) {
	expectedPath := fmt.Sprintf("%s%s.gott%s%s", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder)
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
	expectedPath := fmt.Sprintf("%s%s.gott%s%s%s%s.json", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder, string(os.PathSeparator), activity.Name)
	folderExists, err := utils.PathExists(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !folderExists {
		t.Fatal(fmt.Sprintf("Expected activity file path %s to exist", expectedPath))
	}
}

func clearTestFolder() {
	path := fmt.Sprintf("%s%s.gott%s%s", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Could not delete test folder %s with error %s", path, err.Error())
		os.Exit(-1)
	}
}

func TestInitializeWhenNoConfigurationExists(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
	defer clearTestFolder()
}

func TestInitializeWhenConfigFolderExists(t *testing.T) {
	utils.CreateDir(fmt.Sprintf("%s%s.gott%s", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator)))
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	assertConfigurationFolderExists(t)
	defer clearTestFolder()
}

func TestInitializeWhenDataFolderExists(t *testing.T) {
	utils.CreateDir(fmt.Sprintf("%s%s.gott%sdata", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator)))
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
	defer clearTestFolder()
}

func TestUpdateActivityWhenActivityFileDoesNotExist(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
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
	defer clearTestFolder()
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

	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	errUpdating := repo.Update(activity)
	if errUpdating != nil {
		t.Fatal(errUpdating)
	}

	assertActivityFileExists(activity, t)
	defer clearTestFolder()
}

func TestListActivitiesWhenFolderContainsActivities(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "act1", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	activity2 := core.Activity{Name: "act2", Alias: ""}
	errActivity2 := repo.Update(activity2)
	if errActivity2 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	activities := repo.List()

	if len(activities) != 2 {
		t.Fatalf("Wrong number of activities listed: got %v and expected %v", len(activities), 2)
	}
	defer clearTestFolder()
}

func TestListActivitiesWhenFolderIsEmpty(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activities := repo.List()

	if len(activities) != 0 {
		t.Fatalf("Wrong number of activities listed: got %v and expected %v", len(activities), 0)
	}
	defer clearTestFolder()
}

func TestListActivitiesWhenFolderContainsAnUnexpectedJsonFile(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "act1", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	p := fmt.Sprintf("%s%s.gott%s%s%s%s.json", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder, string(os.PathSeparator), "some_file")
	s1 := make([]byte, 4)
	utils.WriteToFile(p, s1)

	activities := repo.List()

	if len(activities) != 1 {
		t.Fatalf("Wrong number of activities listed: got %v and expected %v", len(activities), 10)
	}
	defer clearTestFolder()
}

func TestListActivitiesWhenFolderContainsAnUnexpectedFileType(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "act1", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	p := fmt.Sprintf("%s%s.gott%s%s%s%s.xx", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder, string(os.PathSeparator), "some_file")
	s1 := make([]byte, 4)
	utils.WriteToFile(p, s1)

	activities := repo.List()

	if len(activities) != 1 {
		t.Fatalf("Wrong number of activities listed: got %v and expected %v", len(activities), 10)
	}
	defer clearTestFolder()
}
