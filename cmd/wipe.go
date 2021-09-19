package cmd

import (
	"fmt"
	"os"

	"github.com/luispcosta/go-tt/core"
	"github.com/spf13/cobra"
)

type wipeCommand struct {
	activity string
	baseCmd  *cobra.Command
}

// NewWipeCommand deletes log data for a given period
func NewWipeCommand(activityRepo core.ActivityRepository) *cobra.Command {
	wipeCmd := &cobra.Command{
		Use:   "wipe",
		Short: "Deletes logs for a given period, and for a sepcific activity (optional)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ExitIfAppNotConfigured()
			alias := cmd.Flag("activity").Value.String()
			period, errPeriod := core.PeriodFromDateStrings(args[0], args[1])
			if errPeriod != nil {
				fmt.Println(errPeriod)
				os.Exit(1)
				return
			}
			if AllowedToContinue() {
				if alias != "" {
					activity, err := activityRepo.Find(alias)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
						return
					}
					errWipe := activityRepo.WipeLogsPeriodAndActivity(period, activity)
					if errWipe != nil {
						fmt.Println(err.Error())
						os.Exit(1)
						return
					}
					fmt.Println("Done")
					os.Exit(1)
					return
				}

				errWipe := activityRepo.WipeLogsPeriod(period)
				if errWipe != nil {
					fmt.Println(errWipe.Error())
					os.Exit(1)
				}
				fmt.Println("Done")
			}
		},
	}
	wipe := wipeCommand{}
	wipeCmd.Flags().StringVarP(&wipe.activity, "activity", "a", "", "Activity name or alias")
	wipe.baseCmd = wipeCmd
	return wipeCmd
}
