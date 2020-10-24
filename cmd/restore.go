package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"

	"github.com/spf13/cobra"
)

type restoreCmd struct {
	baseCmd *cobra.Command
	force   bool
}

// NewRestoreCommand builds the "restore" command
func NewRestoreCommand(activityRepo core.ActivityRepository) *cobra.Command {
	resCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores your activity database",
		Long:  "Restores your activity database with the given dataset. NOTE!: This is destructive. It will delete your current data folder!",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			backupFilePath := args[0]
			force := cmd.Flag("force").Value.String()
			_, err := utils.PathExists(backupFilePath)
			if err != nil {
				fmt.Println("Restore file does not exist or could not be loaded")
				os.Exit(1)
			}

			if utils.IsExtension(backupFilePath, "zip") {
				if force == "true" || userAllowedToContinue() {
					err := activityRepo.Restore(backupFilePath)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
					fmt.Println("Restore complete")
				} else {
					os.Exit(1)
				}
			} else {
				fmt.Println("Restore file must be a zip file")
				os.Exit(1)
			}

		},
	}

	cmd := restoreCmd{}
	resCmd.Flags().BoolVarP(&cmd.force, "force", "f", false, "Force the restore without prompting for confirmation")
	cmd.baseCmd = resCmd
	return resCmd
}

func userAllowedToContinue() bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Are you sure you want to proceed? This will erase any data you already have locally. Type 'y' to continue: ")
	scanner.Scan()
	text := scanner.Text()
	return text == "y"
}
