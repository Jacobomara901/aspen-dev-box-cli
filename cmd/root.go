/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adb",
	Short: "A CLI helper for aspen-dev-box",
	Long: `Aspen CLI

The Aspen CLI simplifies the management of your Aspen Dev Box, a containerized development environment using Docker Compose. 
This command-line tool offers essential features for an efficient development workflow.

Configuration:

The Aspen CLI assumes that the Aspen Dev Box project is located at a specific directory path. 
The project directory path is specified by the PROJECTS_DIR environment variable accorging to the aspen-dev-box repository readme.

Getting Started:

Ensure Docker and Docker Compose are installed.
Set PROJECTS_DIR to your Aspen Dev Box project location.
Run the Aspen CLI commands for efficient containerized development.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aspen-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

}
