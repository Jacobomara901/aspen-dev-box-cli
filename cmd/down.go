package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(DownCommand())
}

func DownCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Bring down the Docker Compose project",
		Long: `Bring down the Docker Compose project and remove orphaned containers.
This command stops and removes all containers defined in the docker-compose file.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command("docker", "compose", "-f", config.GetDefaultComposeFile(), "down", "--remove-orphans")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Printf("Error bringing down the project: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
