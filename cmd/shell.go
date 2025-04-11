package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ShellCommand())
}

func ShellCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Open a shell inside the main container",
		Long: `Open an interactive shell inside the main container.
This command opens a bash shell in the main container with the working directory set to the Aspen Discovery installation.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command("docker", "exec", "-itw", config.GetMainContainerWorkDir(), config.GetMainContainerName(), "/bin/bash")
			command.Dir = config.GetProjectsDir()
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Printf("Error opening a shell in the container: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
