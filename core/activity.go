package core

import (
	"errors"
	"regexp"
	"strings"
)

// Activity represents an activity done by the user is some point in time
type Activity struct {
	Name  string
	Alias string
}

// ValidateName validates the correctness of the activity name.
func (activity *Activity) ValidateName() error {

	if activity.Name == "" {
		return errors.New("Activity name is not valid: it cannot be empty")
	}

	re := regexp.MustCompile(`^[0-9a-zA-Z_-]*$`)
	if !re.MatchString(activity.Name) {
		return errors.New("Activity name is not valid. It must only contain alpha numeric characters")
	}
	activity.Name = strings.ToLower(activity.Name)
	return nil
}
