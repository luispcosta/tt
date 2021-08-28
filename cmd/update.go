package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"

	"github.com/spf13/cobra"
)

type updateCommand struct {
	name        string
	description string
	baseCmd     *cobra.Command
}

// NewUpdateCommand builds the "update" command
func NewUpdateCommand(activityRepo core.ActivityRepository) *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates the metadata of an activity. Only the name and description can be updated.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ExitIfAppNotConfigured()
			name := cmd.Flag("name")
			description := cmd.Flag("desc")
			var updateOp core.UpdateActivity
			if name.Changed && description.Changed {
				updateOp = core.UpdateActivityNameAndDescription{Name: name.Value.String(), Desc: description.Value.String()}
			} else if !name.Changed && description.Changed {
				updateOp = core.UpdateActivityDescription{Desc: description.Value.String()}
			} else if name.Changed && !description.Changed {
				updateOp = core.UpdateActivityName{Name: name.Value.String()}
			} else {
				updateOp = core.NoActivityUpdate{}
			}

			errUpdate := activityRepo.Update(args[0], updateOp)
			if errUpdate != nil {
				fmt.Println(errUpdate)
				os.Exit(1)
			}
			fmt.Println("Activity updated")
		},
	}
	upd := updateCommand{}
	updateCmd.Flags().StringVarP(&upd.name, "name", "n", "", "Activity name")
	updateCmd.Flags().StringVarP(&upd.description, "desc", "d", "xx", "Activity description")
	upd.baseCmd = updateCmd
	return updateCmd
}
