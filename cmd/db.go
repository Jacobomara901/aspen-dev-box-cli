package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(DBCommand())
}

func DBCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "Opens the database shell",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			dbContainerName := "aspen-db"
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			command := exec.Command("docker", "exec", "-it", dbContainerName, "/bin/bash", "-c", "mysql -uroot -paspen aspen")
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
}
