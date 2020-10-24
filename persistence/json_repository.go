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
	Clock      utils.Clock
	AliasIndex *core.AliasIndex
}

// NewJSONActivityRepository initializes a new json repository
func NewJSONActivityRepository(config configuration.Config) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: "data",
		LogFolder:  "log",
		Clock:      utils.NewLiveClock(),
	}
	repo.AliasIndex = core.NewAliasIndex()
	return &repo
}

// NewCustomJSONActivityRepository initializes a new json repository where the data lives in a custom folder.
func NewCustomJSONActivityRepository(folder string, logFolder string, config configuration.Config, clock utils.Clock) *JSONActivityRepository {
	repo := JSONActivityRepository{
		Config:     config,
		DataFolder: folder,
		LogFolder:  logFolder,
		Clock:      clock,
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

	return nil
}

// Update updates the metadata of an activity
func (repo *JSONActivityRepository) Update(activity core.Activity) error {
	errValidate := activity.ValidateName()
	if errValidate != nil {
		return errValidate
	}

	if len(activity.Alias) > 0 {
		errUpdateAlias := repo.setActivityAlias(activity)

		if errUpdateAlias != nil {
			return errUpdateAlias
		}
	}

	bytes, err := json.Marshal(activity)
	if err != nil {
		// We couldn't marshall the activity but the index was already updated. Lets clear it.
		indexFilePath := fmt.Sprintf("%s%sindex.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator))

		repo.AliasIndex.Delete(activity.Alias)
		bytes, err := json.Marshal(repo.AliasIndex.Data)
		if err != nil {
			return err
		}

		errWrite := utils.WriteToFile(indexFilePath, bytes)
		if errWrite != nil {
			return errWrite
		}
	}

	path := fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activity.Name))
	errWrite := utils.WriteToFile(path, bytes)
	if errWrite != nil {
		// We couldn't write the activity data but the index was already updated. Lets clear it
		indexFilePath := fmt.Sprintf("%s%sindex.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator))

		repo.AliasIndex.Delete(activity.Alias)
		bytes, err := json.Marshal(repo.AliasIndex.Data)
		if err != nil {
			return err
		}

		errWrite := utils.WriteToFile(indexFilePath, bytes)
		if errWrite != nil {
			return errWrite
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
		if !info.IsDir() && filepath.Ext(path) == ".json" && filepath.Base(path) != "index.json" {
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
	if repo.activityExists(activityName) {
		activityFilePath, errGettingActivityPath := repo.activityFilePath(activityName)
		if errGettingActivityPath != nil {
			return errGettingActivityPath
		}
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

// Find returns an activity given its name or alias
func (repo *JSONActivityRepository) Find(activityNameOrAlias string) (*core.Activity, error) {
	if repo.activityExists(activityNameOrAlias) {
		activityFilePath, errGettingActivityPath := repo.activityFilePath(activityNameOrAlias)
		if errGettingActivityPath != nil {
			return nil, errGettingActivityPath
		}
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

	return nil, utils.NewNotFoundError(fmt.Sprintf("Activity with name and/or alias: %s not found", activityNameOrAlias))
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

	return nil, utils.NewNotFoundError(fmt.Sprintf("logs for day %s not found", day))
}

// Start sets the start time of an activity
func (repo *JSONActivityRepository) Start(activity core.Activity) error {
	instant := repo.Clock.Now()
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
		log.Start = strconv.FormatInt(repo.Clock.Now().Unix(), 10)

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
		log.Start = strconv.FormatInt(repo.Clock.Now().Unix(), 10)
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
	instant := repo.Clock.Now()
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
				lastEntry.End = strconv.FormatInt(repo.Clock.Now().Unix(), 10)
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

// Purge deletes all activity data
func (repo *JSONActivityRepository) Purge() error {
	userDataLocation := repo.Config.UserDataLocation
	userDataLocationExists, _ := utils.PathExists(userDataLocation)

	if userDataLocationExists {
		err := utils.DeleteDir(userDataLocation)
		if err != nil {
			errorStr := fmt.Sprintf("Could not delete data: %s", err.Error())
			return errors.New(errorStr)
		}
	} else {
		return errors.New("Nothing to delete. Have you run this command before?")
	}
	return nil
}

// Backup creates a zip file the data registered in the system
func (repo *JSONActivityRepository) Backup(destination string) (string, error) {
	userDataLocation := repo.Config.UserDataLocation
	err := utils.ZipIt(userDataLocation, destination)
	if err != nil {
		return "", err
	}

	return destination, nil
}

// Restore restores the repository data with a given backup file
func (repo *JSONActivityRepository) Restore(restoreFilePath string) error {
	errUnzip := utils.Unzip(restoreFilePath, "tmpUnzip")
	defer utils.DeleteDir("tmpUnzip")
	if errUnzip != nil {
		return errUnzip
	}
	errDelete := repo.Config.DeleteConfig()
	if errDelete != nil {
		return errDelete
	}

	errMove := utils.Move(fmt.Sprintf("tmpUnzip%s.gott", string(os.PathSeparator)), fmt.Sprintf("%s%s.gott", utils.HomeDir(), string(os.PathSeparator)))
	if errMove != nil {
		return errMove
	}

	return nil
}

func (repo *JSONActivityRepository) setActivityAlias(activity core.Activity) error {
	indexFilePath := fmt.Sprintf("%s%sindex.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator))

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
	return repo.activityFilePathOnDisk(activity.Name)
}

func (repo *JSONActivityRepository) dayLogFilePath(date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	logFilePath := fmt.Sprintf("%s%s%v%s%v%s%v.json", repo.Config.UserDataLocation+repo.LogFolder, string(os.PathSeparator), year, string(os.PathSeparator), int(month), string(os.PathSeparator), day)
	return logFilePath
}

func (repo *JSONActivityRepository) activityExists(activityName string) bool {
	path, err := repo.activityFilePath(activityName)
	if err != nil {
		return false
	}
	exists, _ := utils.PathExists(path)
	return exists
}

func (repo *JSONActivityRepository) activityFilePath(activityNameOrAlias string) (string, error) {
	if repo.AliasIndex.IsIndexed(activityNameOrAlias) {
		path, err := repo.AliasIndex.Get(activityNameOrAlias)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	return repo.activityFilePathOnDisk(activityNameOrAlias), nil
}

func (repo *JSONActivityRepository) activityFilePathOnDisk(activityName string) string {
	return fmt.Sprintf("%s%s%s.json", repo.Config.UserDataLocation+repo.DataFolder, string(os.PathSeparator), strings.ToLower(activityName))
}
