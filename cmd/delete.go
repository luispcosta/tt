package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/persistence"
	"github.com/spf13/cobra"
)

// NewDeleteCommand deletes an activity registered from the system
func NewDeleteCommand(activityRepo persistence.ActivityRepository) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del",
		Short: "Deletes an activity",
		Long:  "Deletes an activity, if it exists, via its name.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activityName := args[0]
			errDelete := activityRepo.Delete(activityName)
			if errDelete != nil {
				fmt.Println(errDelete)
				os.Exit(1)
			} else {
				fmt.Printf("Activity %s deleted\n", activityName)
			}
		},
	}
	return deleteCmd
}
