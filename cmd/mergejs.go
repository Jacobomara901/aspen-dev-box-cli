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
		Short: "Run the merge_javascript.php command",
		Run: func(cmd *cobra.Command, args []string) {
			jsDir := os.Getenv("ASPEN_CLONE") + "/code/web/interface/themes/responsive/js"
			phpFile := "merge_javascript.php"

			if jsDir == "" {
				fmt.Println("Error: ASPEN_CLONE environment variable not set.")
				os.Exit(1)
			}

			command := exec.Command("php", phpFile)
			command.Dir = jsDir

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
