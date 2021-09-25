package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewDeleteCommand deletes an activity registered from the system
func NewDeleteCommand(activityRepo core.ActivityRepository) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del",
		Short: "Deletes an activity",
		Long:  "Deletes an activity, if it exists. The argument can be either the activity name or alias",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ExitIfAppNotConfigured()
			activityNameOrAlias := args[0]
			errDelete := activityRepo.Delete(activityNameOrAlias)
			if errDelete != nil {
				fmt.Println(errDelete)
				os.Exit(1)
			} else {
				fmt.Printf("Activity with name or alias %s deleted\n", strings.ToLower(activityNameOrAlias))
			}
		},
	}
	return deleteCmd
}
