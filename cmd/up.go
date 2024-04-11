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
	var detached bool
	var debugging bool
	var dbgui bool

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Bring up the Docker Compose project",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			commandArgs := []string{"compose"}

			composeFile := fmt.Sprintf("%s/docker-compose.yml", projectsDir)
			commandArgs = append(commandArgs, "-f", composeFile)

			if debugging {
				composeFile = fmt.Sprintf("%s/docker-compose.debug.yml", projectsDir)
				commandArgs = append(commandArgs, "-f", composeFile)
			}

			if dbgui {
				composeFile = fmt.Sprintf("%s/docker-compose.dbgui.yml", projectsDir)
				commandArgs = append(commandArgs, "-f", composeFile)
			}

			commandArgs = append(commandArgs, "up")

			if detached {
				commandArgs = append(commandArgs, "-d")
			}

			command := exec.Command("docker", commandArgs...)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error bringing up the project: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&detached, "detached", "d", false, "Run in detached mode")
	cmd.Flags().BoolVarP(&debugging, "debugging", "g", false, "Run with debugging compose file")
	cmd.Flags().BoolVarP(&dbgui, "dbgui", "b", false, "Run with dbgui compose file")

	return cmd
}
