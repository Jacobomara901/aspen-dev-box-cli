package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(LogsCommand())
}

func LogsCommand() *cobra.Command {
	var includeIndexingLogs bool
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Tails the logs inside the main container",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			mainContainerName := "containeraspen"
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			logPath := "/var/log/aspen-discovery/test.localhostaspen/"

			logsToTail := "./* "
			if includeIndexingLogs {
				logsToTail += "./logs/* "
			}

			tailOption := ""
			if follow {
				tailOption = "-f "
			}

			command := exec.Command("docker", "exec", "-it", mainContainerName, "/bin/bash", "-c", "(cd "+logPath+"; tail "+tailOption+logsToTail+")")
			command.Dir = fmt.Sprintf(projectsDir)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Run()
			if err != nil {
				fmt.Printf("Error tailing logs in the container: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&includeIndexingLogs, "include-indexing", "i", false, "Include indexing logs")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the logs")

	return cmd
}
