package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

type backupCommand struct {
	baseCmd     *cobra.Command
	zipFilePath string
}

// NewBackupCommand backs up your data
func NewBackupCommand(activityRepo core.ActivityRepository) *cobra.Command {
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Backsup your data into a zip file",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			path := cmd.Flag("path").Value.String()
			fileBackupPath, err := activityRepo.Backup(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Println(fmt.Sprintf("Your data has been backed up at %s", fileBackupPath))
			}
		},
	}
	backup := backupCommand{}
	backupCmd.Flags().StringVarP(&backup.zipFilePath, "path", "p", "go-tt-backup.zip", "Backup file destination")
	backup.baseCmd = backupCmd
	return backupCmd
}
