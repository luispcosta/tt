package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// JSONActivityRepository is an implementation of the activity repository that stores activities as json data.
type JSONActivityRepository struct {
	Config     configuration.Config
	DataFolder string
	LogFolder  string
}

// NewJSONActivityRepository initializes a new json repository
func NewJSONActivityRepository(config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: "data",
		LogFolder:  "log",
	}
	return &repo
}

// NewCustomJSONActivityRepository initializes a new json repository where the data lives in a custom folder.
func NewCustomJSONActivityRepository(folder string, logFolder string, config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: folder,
		LogFolder:  logFolder,
	}
	return &repo
}

// Initialize initializes the json repository by creating the data location folder if it doesn't exist.
func (repo *JSONActivityRepository) Initialize() error {
	_, err := os.Stat(repo.Config.UserDataLocation + repo.DataFolder)
	if os.IsNotExist(err) {
		errCreating := utils.CreateDir(repo.Config.UserDataLocation + repo.DataFolder)

		if errCreating != nil {
			return errCreating
		}
	}

	_, errLogFolder := os.Stat(repo.Config.UserDataLocation + repo.LogFolder)
	if os.IsNotExist(errLogFolder) {
		errCreatingLogFolder := utils.CreateDir(repo.Config.UserDataLocation + repo.LogFolder)

		if errCreatingLogFolder != nil {
			return errCreatingLogFolder
		}
	}

	return nil
}

// Update updates the metadata of an activity
func (repo *JSONActivityRepository) Update(activity core.Activity) error {
	errValidate := activity.ValidateName()
	if errValidate != nil {
		return errValidate
	}

	bytes, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activity.Name))
	return utils.WriteToFile(path, bytes)
}

// List returns all activities currently registered in the system.
func (repo *JSONActivityRepository) List() []core.Activity {
	var activities []core.Activity
	var readingErrors []string

	folder := fmt.Sprintf("%s", repo.Config.UserDataLocation+repo.DataFolder)

	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			shouldIgnoreFile := false
			fileData, errRead := ioutil.ReadFile(path)
			if errRead != nil {
				readingErrors = append(readingErrors, fmt.Sprintf("Could not read data from file %s - error: %s", path, errRead))
				shouldIgnoreFile = true
			}
			activity := core.Activity{}
			errUnmarshall := json.Unmarshal([]byte(fileData), &activity)
			if errUnmarshall != nil {
				readingErrors = append(readingErrors, fmt.Sprintf("Could not load json data from file %s - error: %s", path, errUnmarshall))
				shouldIgnoreFile = true
			}

			if !shouldIgnoreFile {
				activities = append(activities, activity)
			}
		}
		return nil
	})

	if len(readingErrors) > 0 {
		for _, err := range readingErrors {
			log.Printf("WARN: %s", err)
		}
	}

	return activities
}

// Delete deletes an activity via a name (that needs to match a file name). If the files was deleted with success, then no error is returned.
func (repo *JSONActivityRepository) Delete(activityName string) error {
	activityFilePath := repo.activityFilePath(activityName)
	if repo.activityExists(activityName) {
		errDelete := utils.DeleteAtPath(activityFilePath)
		return errDelete
	}

	return errors.New("Activity not registered")
}

// Find returns an activity given its name
func (repo *JSONActivityRepository) Find(activityName string) (*core.Activity, error) {
	activityFilePath := repo.activityFilePath(activityName)
	if repo.activityExists(activityName) {
		fileData, errRead := ioutil.ReadFile(activityFilePath)
		if errRead != nil {
			return nil, errRead
		}

		activity := core.Activity{}
		errUnmarshall := json.Unmarshal([]byte(fileData), &activity)
		if errUnmarshall != nil {
			return nil, errUnmarshall
		}

		return &activity, nil
	}

	return nil, errors.New("Activity does not exist")
}

// Start sets the start time of an activity
func (repo *JSONActivityRepository) Start(activity core.Activity) error {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	logFilePath := fmt.Sprintf("%s%s%v%s%v%s%v.json", repo.Config.UserDataLocation+repo.LogFolder, string(os.PathSeparator), year, string(os.PathSeparator), int(month), string(os.PathSeparator), day)
	exists, _ := utils.PathExists(logFilePath)
	if exists {
		dayLog := core.ActivityDayLog{}
		fileData, errRead := ioutil.ReadFile(logFilePath)
		if errRead != nil {
			return errRead
		}

		errUnmarshall := json.Unmarshal([]byte(fileData), &dayLog)

		if errUnmarshall != nil {
			return errUnmarshall
		}

		logs := dayLog[activity.Name]
		log := core.ActivityLog{}
		log.Start = time.Now().String()

		logs = append(logs, log)
		dayLog[activity.Name] = logs

		bytes, err := json.Marshal(dayLog)

		if err != nil {
			return nil
		}

		errWrite := utils.WriteToFile(logFilePath, bytes)
		if errWrite != nil {
			return errWrite
		}
	} else {
		log := core.ActivityLog{}
		log.Start = time.Now().String()
		dayLog := make(core.ActivityDayLog)
		dayLog[activity.Name] = []core.ActivityLog{log}
		bytes, err := json.Marshal(dayLog)
		if err != nil {
			return err
		}
		yearFolderPath := fmt.Sprintf("%s%s%v", repo.Config.UserDataLocation+repo.LogFolder, string(os.PathSeparator), year)
		existsYearFolder, _ := utils.PathExists(yearFolderPath)
		if !existsYearFolder {
			errCreateYearFolder := utils.CreateDir(yearFolderPath)
			if errCreateYearFolder != nil {
				return errCreateYearFolder
			}
		}
		monthFolderPath := fmt.Sprintf("%s%s%v%s%v", repo.Config.UserDataLocation+repo.LogFolder, string(os.PathSeparator), year, string(os.PathSeparator), int(month))
		existsMonthFolder, _ := utils.PathExists(monthFolderPath)
		if !existsMonthFolder {
			errCreateMonthFolder := utils.CreateDir(monthFolderPath)
			if errCreateMonthFolder != nil {
				return errCreateMonthFolder
			}
		}

		errWrite := utils.WriteToFile(logFilePath, bytes)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

func (repo *JSONActivityRepository) activityExists(activityName string) bool {
	exists, _ := utils.PathExists(repo.activityFilePath(activityName))
	return exists
}

func (repo *JSONActivityRepository) activityFilePath(activityName string) string {
	return fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activityName))
}
