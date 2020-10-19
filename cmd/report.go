package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/reporter"
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
		Long: fmt.Sprintf(`
			Generates a report that presents information about all the activities you performed in the required period.
			This command accepts 1 or 2 arguments.

			If only 1 argument is provided, than it is assumed that the user wants a report in a fixed time frame. This argument represents
			that time frame. Accepted values are: %v. All time frames are relative to the current day.

			If two arguments are passed, they are used to construct a specific time frame period. 
			For example: $ go-tt report '2020-10-10' '2020-10-20'

			You can also provide an additional flag (-f <FORMAT> or --format <FORMAT>) to specify the report format. The default format is printing
			the report to STDOUT. Allowed values are: %v
		`, core.AllowedPeriodFixedTimeFrames(), reporter.AllowedFormatsCollection()),
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var period core.Period
			if len(args) == 2 {
				parsedPeriod, errPeriod := core.PeriodFromDateStrings(args[0], args[1])

				if errPeriod != nil {
					fmt.Println(errPeriod)
					os.Exit(1)
				}
				period = parsedPeriod
			} else {
				period = core.PeriodFromKeyWord(args[0])
			}

			format := strings.ToLower(cmd.Flag("format").Value.String())
			if !reporter.IsAllowedFormat(format) {
				fmt.Printf("%s is not an allowed report format, check this command's help to see the allowed formats", format)
				os.Exit(1)
			}

			reporter := reporter.CreateReporter(format)

			errInit := reporter.Initialize(activityRepo, period)

			if errInit != nil {
				fmt.Println(errInit)
				os.Exit(1)
			}

			err := reporter.ProduceReport()
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
