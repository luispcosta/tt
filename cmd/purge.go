package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

type purgeCmd struct {
	baseCmd *cobra.Command
	force   bool
}

// NewPurgeCommand deletes all activity data
func NewPurgeCommand(activityRepo core.ActivityRepository) *cobra.Command {
	purge := &cobra.Command{
		Use:   "purge",
		Short: "Deletes activity data",
		Long:  "Deletes all activities and related data from the system.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			force := cmd.Flag("force").Value.String()
			if force == "true" || userAllowedToContinue("Are you sure you want to proceed? This will erase all your data. Type 'y' to continue: ") {
				err := activityRepo.Purge()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				} else {
					fmt.Println("All data removed!")
				}
			} else {
				fmt.Println("Operation aborted")
				os.Exit(1)
			}
		},
	}

	cmd := purgeCmd{}
	purge.Flags().BoolVarP(&cmd.force, "force", "f", false, "Force the purge without confirmation")
	cmd.baseCmd = purge
	return purge
}
