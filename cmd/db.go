package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(DBCommand())
}

func DBCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "Opens the database shell",
		Long: `Opens an interactive MariaDB shell connected to the Aspen database.
This command provides direct access to the database for running SQL queries and managing data.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command("docker", "exec", "-it", config.GetDBContainerName(), "/bin/bash", "-c", "mariadb "+config.GetDBConnectionString())
			command.Dir = config.GetProjectsDir()
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Printf("Error opening database shell: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
