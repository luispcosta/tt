package cmd

import (
	"fmt"

	"github.com/luispcosta/go-tt/persistence"
	"github.com/spf13/cobra"
)

// NewStopCommand stops tracking an activity
func NewStopCommand(activityRepo persistence.ActivityRepository) *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an activity",
		Long:  "Stops counting the time for an activity",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activityName := args[0]
			activity, err := activityRepo.Find(activityName)
			if err != nil {
				fmt.Printf("Could not find activity with name %s", activityName)
			}
			errStop := activityRepo.Stop(*activity)
			if errStop != nil {
				fmt.Printf("Could not stop activity with name %s - error: %s", activityName, errStop.Error())
			}
		},
	}
	return stopCmd
}
