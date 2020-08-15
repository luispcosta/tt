package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewPurgeCommand deletes all activity data
func NewPurgeCommand(activityRepo core.ActivityRepository) *cobra.Command {
	purgeCmd := &cobra.Command{
		Use:   "purge",
		Short: "Deletes activity data",
		Long:  "Deletes all activities and related data from the system.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := activityRepo.Purge()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Println("All data removed!")
			}
		},
	}
	return purgeCmd
}
