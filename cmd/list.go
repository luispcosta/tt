package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewListCommand registers the list activity command
func NewListCommand(activityRepo core.ActivityRepository) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all activities",
		Long:  "Lists all the current registered activities in the system",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ExitIfAppNotConfigured()
			activities, err := activityRepo.List()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for _, act := range activities {
				fmt.Println(act.ToPrintableString())
			}
		},
	}
	return listCmd
}
