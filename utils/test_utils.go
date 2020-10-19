package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/luispcosta/go-tt/core"
)

// TestDataFolder is the data folder used in tests
const TestDataFolder = "test"

// LogTestFolder is the folder with the activity logs used in tests
const LogTestFolder = "logTest"

// ClearTestFolder eliminates the the testdatafolder from the fs (and its contents)
func ClearTestFolder() {
	path := fmt.Sprintf("%s%s.gott%s%s", HomeDir(), string(os.PathSeparator), string(os.PathSeparator), TestDataFolder)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Could not delete test folder %s with error %s", path, err.Error())
		os.Exit(-1)
	}
}

// ClearLogTestFolder eliminates the the log data folder from the fs (and its contents)
func ClearLogTestFolder() {
	path := fmt.Sprintf("%s%s.gott%s%s", HomeDir(), string(os.PathSeparator), string(os.PathSeparator), LogTestFolder)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Could not delete log test folder %s with error %s", path, err.Error())
		os.Exit(-1)
	}
}

// ActivityPath returns the path of an activity in the fs in a testing environment
func ActivityPath(activity core.Activity) string {
	return fmt.Sprintf("%s%s.gott%s%s%s%s.json", HomeDir(), string(os.PathSeparator), string(os.PathSeparator), TestDataFolder, string(os.PathSeparator), strings.ToLower(activity.Name))
}
