package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

// NewMigrateCommand inits the migrate command
func NewMigrateCommand(activityRepo core.ActivityRepository) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrates the schema of the repository",
		Long:  "Migrates the schema of the repository. NOTE: This should only be executed once!",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			direction := args[0]

			if direction != "up" && direction != "down" {
				fmt.Println("Direction must be either 'up' or 'down'")
				os.Exit(1)
				return
			}

			errMigrate := activityRepo.SchemaMigrate(direction)
			if errMigrate != nil {
				fmt.Println(errMigrate.Error())
				os.Exit(1)
			}
		},
	}
	return migrateCmd
}
