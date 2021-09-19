package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luispcosta/go-tt/config"
	"github.com/luispcosta/go-tt/persistence"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "go-tt is a simple CLI app time tracker",
	Long: `
		With go-tt you can easily track the time you spend on different activities throughout your day.
		go-tt provides several reports outlining the time you have spent in all your registered activities.
		The goal of this small app is to help you fight procrastination, by making you aware of where you chose to spend your time.
	`,
}

func ExitIfAppNotConfigured() {
	config := config.NewConfig()
	if !config.AlreadySetup() {
		fmt.Println("Application not yet configured. Please configure with `tt init`")
		os.Exit(1)
	}
}

func AllowedToContinue() bool {
	var input string
	fmt.Printf("Do you want to continue with this operation? [y|n]: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		panic(err)
	}
	input = strings.ToLower(input)
	if input == "y" {
		return true
	}
	return false
}

// Execute executes the root commmand.
func Execute() {
	configuration := config.NewConfig()
	repo, err := persistence.NewSqliteRepository()

	if err != nil {
		fmt.Printf("Could not connect to mongo with error: %s", err.Error())
		os.Exit(1)
	}

	errorInitRepo := repo.Initialize(configuration)

	if errorInitRepo != nil {
		fmt.Printf("Error initializing mongo with error: %s", errorInitRepo.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(NewInitCommand(repo))
	rootCmd.AddCommand(NewAddCommand(repo))
	rootCmd.AddCommand(NewListCommand(repo))
	rootCmd.AddCommand(NewDeleteCommand(repo))
	rootCmd.AddCommand(NewStartCommand(repo))
	rootCmd.AddCommand(NewStopCommand(repo))
	rootCmd.AddCommand(NewReportCommand(repo))
	rootCmd.AddCommand(NewUpdateCommand((repo)))
	rootCmd.AddCommand(NewCurrentCommand((repo)))
	rootCmd.AddCommand(NewWipeCommand((repo)))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
