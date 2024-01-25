package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ShellCommand())
}

func ShellCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Open a shell inside the main container",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			mainContainerName := "containeraspen"
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			command := exec.Command("docker", "exec", "-itw", "/usr/local/aspen-discovery", mainContainerName, "/bin/bash")
			command.Dir = fmt.Sprintf(projectsDir)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error opening a shell in the container: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
