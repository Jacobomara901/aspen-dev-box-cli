package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(LogsCommand())
}

func LogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View container logs",
		Long: `View logs from the main container.
This command allows you to view and follow logs in real-time.
You can optionally include indexing logs using the --include-indexing flag.`,
		Run: func(cmd *cobra.Command, args []string) {
			includeIndexingLogs, _ := cmd.Flags().GetBool("include-indexing")
			follow, _ := cmd.Flags().GetBool("follow")

			// Build the logs path pattern
			logsPattern := "./*"
			if includeIndexingLogs {
				logsPattern += " ./logs/*"
			}

			// Build the tail command
			tailCmd := "tail"
			if follow {
				tailCmd += " -f"
			}

			// Execute the command in the container
			dockerCmd := exec.Command("docker", "exec", "-it",
				config.GetMainContainerName(),
				"/bin/bash", "-c",
				fmt.Sprintf("(cd %s; %s %s)", config.GetLogPath(), tailCmd, logsPattern))

			dockerCmd.Dir = config.GetProjectsDir()
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				fmt.Printf("Error viewing logs: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolP("include-indexing", "i", false, "Include indexing logs")
	cmd.Flags().BoolP("follow", "f", false, "Follow logs in real-time")

	return cmd
}
