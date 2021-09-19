package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewCurrentCommand displays the current running activity
func NewCurrentCommand(activityRepo core.ActivityRepository) *cobra.Command {
	currentCmd := &cobra.Command{
		Use:   "current",
		Short: "Displays the current running activity, if any",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ExitIfAppNotConfigured()
			activity, err := activityRepo.CurrentlyTrackedActivity()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if activity != nil {
				fmt.Printf("%s (%s)\n", activity.Name, activity.Alias)
			} else {
				fmt.Println("Not currently tracking any activity")
			}
		},
	}
	return currentCmd
}
