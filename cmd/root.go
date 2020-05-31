package cmd

import (
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
	repo := persistence.NewJSONActivityRepository()
	errorInitRepo := repo.Initialize()

	if errorInitRepo != nil {
		fmt.Println("Error initialize activity repository")
		os.Exit(1)
	}

	rootCmd.AddCommand(NewAddCommand(repo))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
