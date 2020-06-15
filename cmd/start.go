package cmd

import (
	"fmt"

	"github.com/luispcosta/go-tt/persistence"
	"github.com/spf13/cobra"
)

// NewStartCommand starts tracking the time for an activity
func NewStartCommand(activityRepo persistence.ActivityRepository) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts an activity",
		Long:  "Starts counting the time for an activity",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activityName := args[0]
			activity, err := activityRepo.Find(activityName)
			if err != nil {
				fmt.Printf("Could find activity with name %s - error: %s", activityName, err.Error())
			}
			errStart := activityRepo.Start(*activity)
			if errStart != nil {
				fmt.Printf("Could not start activity with name %s - error: %s", activityName, errStart.Error())
			}
		},
	}
	return deleteCmd
}
