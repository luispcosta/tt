package persistence

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/luispcosta/go-tt/configuration"

	"github.com/luispcosta/go-tt/core"

	"github.com/luispcosta/go-tt/utils"
)

const testDataFolder = "test"
const logTestFolder = "logTest"

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
	expectedPath := fmt.Sprintf("%s%s.gott%s%s%s%s.json", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), testDataFolder, string(os.PathSeparator), strings.ToLower(activity.Name))
	folderExists, err := utils.PathExists(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !folderExists {
		t.Fatal(fmt.Sprintf("Expected activity file path %s to exist", expectedPath))
	}
}

func assertActivityLogFileExistsCurrentDay(t *testing.T) {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	expectedPath := fmt.Sprintf("%s%s.gott%s%s%s%v%s%v%s%v.json", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), logTestFolder, string(os.PathSeparator), year, string(os.PathSeparator), int(month), string(os.PathSeparator), day)
	folderExists, err := utils.PathExists(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	if !folderExists {
		t.Fatal(fmt.Sprintf("Expected activity log file path %s to exist", expectedPath))
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

func clearLogTestFolder() {
	path := fmt.Sprintf("%s%s.gott%s%s", utils.HomeDir(), string(os.PathSeparator), string(os.PathSeparator), logTestFolder)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Could not delete log test folder %s with error %s", path, err.Error())
		os.Exit(-1)
	}
}

func TestInitializeWhenNoConfigurationExists(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	assertConfigurationFolderExists(t)
	defer clearTestFolder()
}

func TestUpdateActivityWhenActivityFileDoesNotExist(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
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

func TestDeleteActivitiesWhenActivityExists(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "act1", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	errDelete := repo.Delete(activity1.Name)

	if errDelete != nil {
		t.Fatalf("Should not have failed deleting an activity that is registered. Failed with %s", errDelete.Error())
	}
	defer clearTestFolder()
}

func TestDeleteActivitiesWhenActivityDoesNotExist(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "act1", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	errDelete := repo.Delete("xxx")

	if errDelete == nil {
		t.Fatalf("Should have failed deleting an activity that does not exist")
	}
	defer clearTestFolder()
}

func TestDeleteActivitiesWhenActivityNameDoesntMatchCase(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	errDelete := repo.Delete("aCt")

	if errDelete != nil {
		t.Fatalf("Should not have failed deleting an activity ignoring case")
	}
	defer clearTestFolder()
}

func TestStartActivity(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}
	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	repo.Start(activity1)
	assertActivityLogFileExistsCurrentDay(t)

	_, errFind := repo.Find(activity1.Name)
	if errFind != nil {
		t.Fatal("Should not have failed to find activity that was created before")
	}

	activityDayLog, errDayLog := repo.FindLogsForDay(time.Now())
	if errDayLog != nil {
		t.Fatal("Should have created one activity day log after starting an activity")
	}

	if activityDayLog[activity1.Name] == nil {
		t.Fatal("Should have created a log entry for activity after starting it")
	}

	defer clearTestFolder()
	defer clearLogTestFolder()
}

func TestStopActivityWhenActivityDoesNotExist(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}
	errStop := repo.Stop(activity1)
	if errStop == nil {
		t.Fatal("Should have failed stopping an activity that does not exist")
	}

	defer clearTestFolder()
	defer clearLogTestFolder()
}

func TestStopActivityWhenNoActivityHasNotBeenStartedYet(t *testing.T) {
	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}

	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	errStop := repo.Stop(activity1)
	if errStop == nil {
		t.Fatal("Should have failed stopping an activity that hasn't been starrted yet")
	}

	if errStop.Error() != "No activity started yet today" {
		t.Fatal("Error should have been: no activity started yet today")
	}

	defer clearTestFolder()
	defer clearLogTestFolder()
}

func TestStopActivityWhenActivityHasNotBeenStartedYet(t *testing.T) {
	defer clearTestFolder()
	defer clearLogTestFolder()

	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}
	activity2 := core.Activity{Name: "ACT2", Alias: ""}

	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	errActivity2 := repo.Update(activity2)
	if errActivity2 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	repo.Start(activity2)
	repo.Stop(activity2)

	errStop := repo.Stop(activity1)
	if errStop == nil {
		t.Fatal("Should have failed stopping an activity that hasn't been starrted yet")
	}

	if errStop.Error() != "Activity not yet started" {
		t.Fatal("Error should have been: Activity not yet started")
	}
}

func TestStopActivityWhenActivityHasAlreadyBeenStopped(t *testing.T) {
	defer clearTestFolder()
	defer clearLogTestFolder()

	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}

	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	repo.Start(activity1)
	errInitialStop := repo.Stop(activity1)

	if errInitialStop != nil {
		t.Fatal("Should not have failed stopping activity")
	}

	errStop := repo.Stop(activity1)
	if errStop == nil {
		t.Fatal("Should have failed stopping an activity that hasn't been starrted yet")
	}

	if errStop.Error() != "Last activity has already been stoped. Please start a new one" {
		t.Fatal("Error should have been: Last activity has already been stoped. Please start a new one")
	}

}

func TestStopActivity(t *testing.T) {
	defer clearTestFolder()
	defer clearLogTestFolder()

	config := configuration.NewConfig()
	repo := NewCustomJSONActivityRepository(testDataFolder, logTestFolder, *config)
	err := repo.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	activity1 := core.Activity{Name: "ACT", Alias: ""}

	errActivity1 := repo.Update(activity1)
	if errActivity1 != nil {
		t.Fatal("Should not have failed creating a valid activity")
	}

	repo.Start(activity1)
	time.Sleep(1 * time.Second)
	errInitialStop := repo.Stop(activity1)

	if errInitialStop != nil {
		t.Fatal("Should not have failed stopping activity")
	}

	activityDayLog, errDayLog := repo.FindLogsForDay(time.Now())
	if errDayLog != nil {
		t.Fatal("Should have created one activity day log after starting an activity")
	}

	if activityDayLog[activity1.Name] == nil {
		t.Fatal("Should have created a log entry for activity after starting it")
	}

	if len(activityDayLog[activity1.Name]) != 1 {
		t.Fatal("Should have only one activity log for activity started only once")
	}

	log := activityDayLog[activity1.Name][0]

	if log.Start == "" {
		t.Fatalf("Activity log Start field should have been filled after starting and stopping activity")
	}

	if log.End == "" {
		t.Fatalf("Activity log End field should have been filled after starting and stopping activity")
	}

	if log.Duration == 0 {
		t.Fatalf("Activity log Duration field should have been filled after starting and stopping activity")
	}
}
