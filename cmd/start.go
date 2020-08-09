package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewStartCommand starts tracking the time for an activity
func NewStartCommand(activityRepo core.ActivityRepository) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts an activity",
		Long:  "Starts counting the time for an activity",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activityNameOrAlias := args[0]
			activity, err := activityRepo.Find(activityNameOrAlias)
			if err != nil {
				fmt.Printf("Could not find activity with name and or alias %s - error: %s\n", activityNameOrAlias, err.Error())
				os.Exit(1)
			}
			errStart := activityRepo.Start(*activity)
			if errStart != nil {
				fmt.Printf("Could not start activity with name and or alias %s - error: %s\n", activityNameOrAlias, errStart.Error())
				os.Exit(1)
			}
		},
	}
	return deleteCmd
}
