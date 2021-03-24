package cmd

import (
	"bufio"
	"fmt"
	"os"

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

// Execute executes the root commmand.
func Execute() {
	repo, err := persistence.NewMongoRepository()

	if err != nil {
		fmt.Printf("Could not connect to mongo with error: %s", err.Error())
		os.Exit(1)
	}

	errorInitRepo := repo.Initialize()

	if errorInitRepo != nil {
		fmt.Printf("Error initializing mongo with error: %s", errorInitRepo.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(NewAddCommand(repo))
	rootCmd.AddCommand(NewListCommand(repo))
	rootCmd.AddCommand(NewDeleteCommand(repo))
	rootCmd.AddCommand(NewStartCommand(repo))
	rootCmd.AddCommand(NewStopCommand(repo))
	rootCmd.AddCommand(NewPurgeCommand(repo))
	rootCmd.AddCommand(NewBackupCommand(repo))
	rootCmd.AddCommand(NewReportCommand(repo))
	rootCmd.AddCommand(NewRestoreCommand(repo))
	rootCmd.AddCommand(NewMigrateCommand(repo))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func userAllowedToContinue(confirmationMsg string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(confirmationMsg)
	scanner.Scan()
	text := scanner.Text()
	return text == "y"
}
