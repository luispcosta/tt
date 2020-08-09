package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	AliasIndex *core.AliasIndex
}

// NewJSONActivityRepository initializes a new json repository
func NewJSONActivityRepository(config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: "data",
		LogFolder:  "log",
	}
	repo.AliasIndex = core.NewAliasIndex()
	return &repo
}

// NewCustomJSONActivityRepository initializes a new json repository where the data lives in a custom folder.
func NewCustomJSONActivityRepository(folder string, logFolder string, config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: folder,
		LogFolder:  logFolder,
	}
	repo.AliasIndex = core.NewAliasIndex()
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
	errWrite := utils.WriteToFile(path, bytes)
	if errWrite != nil {
		return errWrite
	}

	if len(activity.Alias) > 0 {
		errUpdateAlias := repo.setActivityAlias(activity)

		if errUpdateAlias != nil {
			return errUpdateAlias
		}
	}

	return nil
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
		activity, errFind := repo.Find(activityName)
		if errFind != nil {
			return errFind
		}

		if activity.Alias != "" {
			indexFilePath := fmt.Sprintf("%s%sindex.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator))
			repo.AliasIndex.Delete(activity.Alias)
			aliasIndexData := repo.AliasIndex.Data

			bytes, err := json.Marshal(aliasIndexData)
			if err != nil {
				return nil
			}

			errWrite := utils.WriteToFile(indexFilePath, bytes)
			if errWrite != nil {
				return errWrite
			}
		}

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

// FindLogsForDay returns a day log file for the given date if it exists
func (repo *JSONActivityRepository) FindLogsForDay(day time.Time) (core.ActivityDayLog, error) {
	logFilePath := repo.dayLogFilePath(day)
	exists, _ := utils.PathExists(logFilePath)
	if exists {
		dayLog := core.ActivityDayLog{}
		fileData, errRead := ioutil.ReadFile(logFilePath)
		if errRead != nil {
			return nil, errRead
		}

		errUnmarshall := json.Unmarshal([]byte(fileData), &dayLog)

		if errUnmarshall != nil {
			return nil, errUnmarshall
		}

		return dayLog, nil
	}

	return nil, errors.New("Activity Log does not exist")
}

// Start sets the start time of an activity
func (repo *JSONActivityRepository) Start(activity core.Activity) error {
	instant := time.Now()
	year := instant.Year()
	month := instant.Month()
	logFilePath := repo.dayLogFilePath(instant)
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
		log.Start = strconv.FormatInt(time.Now().Unix(), 10)

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
		log.Start = strconv.FormatInt(time.Now().Unix(), 10)
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

// Stop sets the end time for an activity
func (repo *JSONActivityRepository) Stop(activity core.Activity) error {
	instant := time.Now()
	logFilePath := repo.dayLogFilePath(instant)
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

		logs, logExistsForActivity := dayLog[activity.Name]

		if logExistsForActivity {
			lastEntry := logs[len(logs)-1]
			if lastEntry.Start == "" {
				return errors.New("Last activity start time was not recorded") // TODO: We should do something about this
			} else if lastEntry.End != "" {
				return errors.New("Last activity has already been stoped. Please start a new one")
			} else {
				lastEntry.End = strconv.FormatInt(time.Now().Unix(), 10)
				duration, errCalcDuration := utils.CalcActivityLogDuration(lastEntry)
				if errCalcDuration != nil {
					return errCalcDuration
				}
				lastEntry.Duration = duration
				logs[len(logs)-1] = lastEntry
				dayLog[activity.Name] = logs
				bytes, err := json.Marshal(dayLog)
				if err != nil {
					return nil
				}

				errWrite := utils.WriteToFile(logFilePath, bytes)
				if errWrite != nil {
					return errWrite
				}
			}
		} else {
			return errors.New("Activity not yet started")
		}
	} else {
		return errors.New("No activity started yet today")
	}

	return nil
}

func (repo *JSONActivityRepository) setActivityAlias(activity core.Activity) error {
	indexFilePath := fmt.Sprintf("%s%sindex.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator))

	if exists, _ := utils.PathExists(indexFilePath); exists {
		aliasIndexData := core.AliasIndexData{}
		fileData, errRead := ioutil.ReadFile(indexFilePath)
		if errRead != nil {
			return errRead
		}

		errUnmarshall := json.Unmarshal([]byte(fileData), &aliasIndexData)

		if errUnmarshall != nil {
			return errUnmarshall
		}

		repo.AliasIndex.Load(aliasIndexData)
	}

	errUpdate := repo.AliasIndex.Update(activity.Alias, repo.activityAliasIndexValue(activity))
	if errUpdate != nil {
		return errUpdate
	}

	bytes, err := json.Marshal(repo.AliasIndex.Data)
	if err != nil {
		return nil
	}

	errWrite := utils.WriteToFile(indexFilePath, bytes)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func (repo *JSONActivityRepository) activityAliasIndexValue(activity core.Activity) string {
	return repo.activityFilePath(activity.Name)
}

func (repo *JSONActivityRepository) dayLogFilePath(date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	logFilePath := fmt.Sprintf("%s%s%v%s%v%s%v.json", repo.Config.UserDataLocation+repo.LogFolder, string(os.PathSeparator), year, string(os.PathSeparator), int(month), string(os.PathSeparator), day)
	return logFilePath
}

func (repo *JSONActivityRepository) activityExists(activityName string) bool {
	exists, _ := utils.PathExists(repo.activityFilePath(activityName))
	return exists
}

func (repo *JSONActivityRepository) activityFilePath(activityName string) string {
	return fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activityName))
}
