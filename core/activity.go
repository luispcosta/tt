package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Activity represents an activity done by the user is some point in time
type Activity struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name""`
	Alias       string             `json:"alias""`
	Description string             `json:"description"`
}

// HasAlias returns true if the activity has an alias defined on it
func (activity *Activity) HasAlias() bool {
	return activity.Alias != ""
}

// HasDescription returns true if the activity has a description
func (activity *Activity) HasDescription() bool {
	return activity.Description != ""
}

// ValidateName validates the correctness of the activity name.
func (activity *Activity) ValidateName() error {

	if activity.Name == "" {
		return errors.New("Activity name is not valid: it cannot be empty")
	}

	if strings.ToLower(activity.Name) == "index" {
		return errors.New("Activity name is not valid: 'index' is a reserved keyword and cannot be used as name")
	}

	re := regexp.MustCompile(`^[0-9a-zA-Z_-]*$`)
	if !re.MatchString(activity.Name) {
		return errors.New("Activity name is not valid. It must only contain alpha numeric characters")
	}
	activity.Name = strings.ToLower(activity.Name)
	return nil
}

// ToPrintableString returns a pretty string with the activity data
func (activity *Activity) ToPrintableString() string {
	format := "Name: %s"
	var res string
	res = fmt.Sprintf(format, activity.Name)
	if activity.HasAlias() {
		format = "Name: %s\n Alias: %s"
		res = fmt.Sprintf(format, activity.Name, activity.Alias)
	}

	if activity.HasDescription() {
		res = fmt.Sprintf("%s\n Description: %s", res, activity.Description)
	}

	return res
}
