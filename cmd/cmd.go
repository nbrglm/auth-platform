package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nbrglm-auth-platform",
	Short: "NBRGLM Auth Platform CLI application",
	Long:  "A command line interface for the NBRGLM Auth Platform",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln("Please specify a command to run!")
		cmd.Help()
	},
}

// Exec starts the command application
// This is the entry point for the command line application.
// It is responsible for setting up the command line interface and executing the commands.
// Only supposed to be called once, when the application is started, by the main function.
func Exec() {
	initServeCommand()
	initKeygenCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		// Exit with a non-zero status code to indicate an error
		// This is important for CI/CD pipelines and other automated systems.
		os.Exit(1)
	}
}
