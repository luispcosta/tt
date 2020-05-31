package cmd

import (
	"github.com/luispcosta/go-tt/persistence"
	"github.com/spf13/cobra"
)

// NewListCommand registers the list activity command
func NewListCommand(activityRepo persistence.ActivityRepository) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all activities",
		Long:  "Lists all the current registered activities in the system",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return listCmd
}
