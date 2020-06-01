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

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

// JSONActivityRepository is an implementation of the activity repository that stores activities as json data.
type JSONActivityRepository struct {
	Config     configuration.Config
	DataFolder string
}

// NewJSONActivityRepository initializes a new json repository
func NewJSONActivityRepository(config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: "data",
	}
	return &repo
}

// NewCustomJSONActivityRepository initializes a new json repository where the data lives in a custom folder.
func NewCustomJSONActivityRepository(folder string, config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: folder,
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
	activityFilePath := fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activityName))
	exists, _ := utils.PathExists(activityFilePath)
	if exists {
		errDelete := utils.DeleteAtPath(activityFilePath)
		return errDelete
	}

	return errors.New("Activity not registered")
}
