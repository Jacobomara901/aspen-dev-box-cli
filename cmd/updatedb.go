package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(UpdateDBCommand())
}

func UpdateDBCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "updatedb",
		Short: "Runs any pending updates for the aspen database",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			mainContainerName := "containeraspen"
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			curlCommand := "curl -k http://localhost/API/SystemAPI?method=runPendingDatabaseUpdates"
			command := exec.Command("docker", "exec", "-it", mainContainerName, "/bin/bash", "-c", curlCommand)
			command.Dir = fmt.Sprintf(projectsDir)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error running db updates in the container: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
