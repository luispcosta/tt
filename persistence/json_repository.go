package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

var config = configuration.NewConfig()

const dataFolder = "data"

// JSONActivityRepository is an implementation of the activity repository that stores activities as json data.
type JSONActivityRepository struct {
	DataLocation string
}

// NewJSONActivityRepository initializes a new json repository
func NewJSONActivityRepository() *JSONActivityRepository {
	repo := JSONActivityRepository{
		DataLocation: fmt.Sprintf(dataFolder),
	}
	return &repo
}

// Initialize initializes the json repository by creating the data location folder if it doesn't exist.
func (repo *JSONActivityRepository) Initialize() error {
	_, err := os.Stat(config.UserDataLocation + repo.DataLocation)
	if os.IsNotExist(err) {
		errCreating := utils.CreateDir(config.UserDataLocation + repo.DataLocation)

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

	path := fmt.Sprintf("%s%s%s.json", config.UserDataLocation+repo.DataLocation, string(os.PathSeparator), activity.Name)
	return utils.WriteToFile(path, bytes)
}
