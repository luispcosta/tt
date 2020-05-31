package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"

	"github.com/luispcosta/go-tt/persistence"

	"github.com/spf13/cobra"
)

type addCommand struct {
	alias   string
	baseCmd *cobra.Command
}

// NewAddCommand builds the "add" command
func NewAddCommand(activityRepo persistence.ActivityRepository) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a new activity",
		Long:  "Registers a new activity to be tracked. You can also add an alias to the activity.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activity := core.Activity{Name: args[0], Alias: cmd.Flag("alias").Value.String()}
			errUpdate := activityRepo.Update(activity)
			if errUpdate != nil {
				fmt.Println(errUpdate)
				os.Exit(1)
			}
		},
	}
	add := addCommand{}
	addCmd.Flags().StringVarP(&add.alias, "alias", "a", "", "Activity alias")
	add.baseCmd = addCmd
	return addCmd
}
