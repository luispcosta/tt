package cmd

import (
	"fmt"

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
			activities := activityRepo.List()
			for _, act := range activities {
				if act.Alias != "" {
					fmt.Printf("%s (%s)\n", act.Name, act.Alias)
				} else {
					fmt.Printf("%s\n", act.Name)
				}

			}
		},
	}
	return listCmd
}
