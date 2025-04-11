package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(UpdateDBCommand())
}

func UpdateDBCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updatedb",
		Short: "Run database updates",
		Long: `Run any pending database updates for Aspen Discovery.
This command triggers the database update process by calling the SystemAPI endpoint.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Build the curl command to trigger database updates
			curlCmd := "curl -k http://localhost/API/SystemAPI?method=runPendingDatabaseUpdates"

			// Execute the command in the main container
			dockerCmd := exec.Command("docker", "exec", "-it",
				config.GetMainContainerName(),
				"/bin/bash", "-c", curlCmd)

			dockerCmd.Dir = config.GetProjectsDir()
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				fmt.Printf("Error running database updates: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Database updates completed successfully")
		},
	}

	return cmd
}
