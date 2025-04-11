package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(MergeJSCommand())
}

func MergeJSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mergejs",
		Short: "Merge JavaScript files",
		Long: `Merge JavaScript files using the merge_javascript.php script.
This command runs the merge script inside the main container to combine and minify JavaScript files.`,
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command("docker", "exec", "-w", config.GetJSWorkDir(), config.GetMainContainerName(), "php", config.GetMergeJSScript())
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Printf("Error running the merge_javascript.php command: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
