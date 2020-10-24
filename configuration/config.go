package configuration

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/utils"
)

// Config is a base struct with configuration options for the application.
type Config struct {
	UserDataLocation string
}

// NewConfig returns a new app configuration.
func NewConfig() *Config {
	config := Config{}
	initConfigWithDefaultValues(&config)
	return &config
}

// DeleteConfig deletes the current config folder
func (config *Config) DeleteConfig() error {
	homeDir := utils.HomeDir()
	err := utils.DeleteDir(fmt.Sprintf("%s%s.gott", homeDir, string(os.PathSeparator)))
	if err != nil {
		return err
	}
	return nil
}

func initConfigWithDefaultValues(config *Config) {
	homeDir := utils.HomeDir()
	config.UserDataLocation = fmt.Sprintf("%s%s.gott%s", homeDir, string(os.PathSeparator), string(os.PathSeparator))
}
