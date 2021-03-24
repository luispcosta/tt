package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"

	"github.com/spf13/cobra"
)

type addCommand struct {
	alias       string
	description string
	baseCmd     *cobra.Command
}

// NewAddCommand builds the "add" command
func NewAddCommand(activityRepo core.ActivityRepository) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a new activity",
		Long:  "Registers a new activity to be tracked. You can also add an alias to the activity. Case is ignoring for the activity name.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			alias := cmd.Flag("alias").Value.String()
			description := cmd.Flag("desc").Value.String()
			activity := core.Activity{Name: args[0], Alias: alias, Description: description}
			errAdd := activityRepo.Add(activity)
			if errAdd != nil {
				fmt.Println(errAdd)
				os.Exit(1)
			}
		},
	}
	add := addCommand{}
	addCmd.Flags().StringVarP(&add.alias, "alias", "a", "", "Activity alias")
	addCmd.Flags().StringVarP(&add.description, "desc", "d", "", "Activity description")
	add.baseCmd = addCmd
	return addCmd
}
