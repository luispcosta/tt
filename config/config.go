package config

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/utils"
)

// Config is a base struct with configuration options for the application.
type Config struct {
	UserDataLocation string
}

const ConfigFolder = ".gott"

// NewConfig returns a new app configuration.
func NewConfig() Config {
	config := Config{}
	initConfigWithDefaultValues(&config)
	return config
}

// DeleteConfig deletes the current config folder
func (config *Config) DeleteConfig() error {
	homeDir := utils.HomeDir()
	err := utils.DeleteDir(fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ConfigFolder))
	if err != nil {
		return err
	}
	return nil
}

// AlreadySetup returns true if the app has already been setup
func (config *Config) AlreadySetup() bool {
	exists, err := utils.PathExists(config.UserDataLocation)
	if err != nil || !exists {
		return false
	}

	return exists
}

// Setup tries to setup the app, otherwise returns an error
func (config *Config) Setup() error {
	if config.AlreadySetup() {
		return nil
	}

	return utils.CreateDir(config.UserDataLocation)
}

func initConfigWithDefaultValues(config *Config) {
	homeDir := utils.HomeDir()
	config.UserDataLocation = fmt.Sprintf("%s%s%s%s", homeDir, string(os.PathSeparator), ConfigFolder, string(os.PathSeparator))
}
