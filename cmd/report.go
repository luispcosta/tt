package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/reporter"
	"github.com/luispcosta/go-tt/utils"
	"github.com/spf13/cobra"
)

type reportCmd struct {
	baseCmd *cobra.Command
	format  string
}

// NewReportCommand creates ativities reports
func NewReportCommand(activityRepo core.ActivityRepository) *cobra.Command {
	reportCommand := &cobra.Command{
		Use:   "report",
		Short: "Creates an activity report over a time period",
		Long:  "Generates a report that presents information about all the activities you performed in the required period.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			startDate := args[0]
			endDate := args[1]
			parsedStDate, err1 := utils.ParseSimpleDate(startDate)
			if err1 != nil {
				fmt.Println(err1)
				os.Exit(1)
			}
			parsedEdDate, err2 := utils.ParseSimpleDate(endDate)
			if err2 != nil {
				fmt.Println(err2)
				os.Exit(1)
			}

			format := strings.ToLower(cmd.Flag("format").Value.String())
			if !reporter.IsAllowedFormat(format) {
				fmt.Printf("%s is not an allowed report format, check this command's help to see the allowed formats", format)
				os.Exit(1)
			}

			reporter := reporter.Create(format)

			errInit := reporter.Initialize(activityRepo)

			if errInit != nil {
				fmt.Println(errInit)
				os.Exit(1)
			}

			err := reporter.ProduceReport(*parsedStDate, *parsedEdDate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	report := reportCmd{}
	reportCommand.Flags().StringVarP(&report.format, "format", "f", "cli", "Report format")
	report.baseCmd = reportCommand
	return reportCommand
}
