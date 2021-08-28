package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/config"
	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewInitCommand inits the application
func NewInitCommand(activityRepo core.ActivityRepository) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "init",
		Short: "Inits the application",
		Long:  "Setup the application, by creating the necessary data.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := config.NewConfig()
			if config.AlreadySetup() {
				fmt.Println("Already configuration, you can now use the application")
			} else {
				err := config.Setup()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		},
	}
	return deleteCmd
}
