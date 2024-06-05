package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(MergeJSCommand())
}

func MergeJSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mergejs",
		Short: "Run the merge_javascript.php command inside the container",
		Run: func(cmd *cobra.Command, args []string) {
			containerName := "containeraspen"
			if containerName == "" {
				fmt.Println("Error: Container name not set.")
				os.Exit(1)
			}

			phpFile := "merge_javascript.php"
			workDir := "/usr/local/aspen-discovery/code/web/interface/themes/responsive/js"

			command := exec.Command("docker", "exec", "-w", workDir, containerName, "php", phpFile)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error running the merge_javascript.php command: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
