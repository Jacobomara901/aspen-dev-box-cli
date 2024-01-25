package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(UpCommand())
}

func UpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Bring up the Docker Compose project",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			composeFile := fmt.Sprintf("%s/docker-compose.yml", projectsDir)
			command := exec.Command("docker", "compose", "-f", composeFile, "up", "-d")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error bringing up the project: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
