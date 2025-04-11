package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"adb/pkg/config"
)

func init() {
	rootCmd.AddCommand(CSSCommand())
}

// CSSCommand returns the command for compiling CSS
func CSSCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compilecss",
		Short: "Compile CSS files",
		Long: `Compile CSS files using LESS.
This command compiles the main.less file into main.css in the specified directory.
Use the --rtl flag to compile RTL (right-to-left) CSS files.`,
		Run: func(cmd *cobra.Command, args []string) {
			rtl, _ := cmd.Flags().GetBool("rtl")
			cssDir := config.GetCSSDir(rtl)

			// Check if the CSS directory exists
			if _, err := os.Stat(cssDir); os.IsNotExist(err) {
				fmt.Printf("Error: CSS directory does not exist: %s\n", cssDir)
				os.Exit(1)
			}

			// Run the LESS compilation
			dockerCmd := exec.Command("docker", "run", "--rm",
				"-v", fmt.Sprintf("%s:/src", cssDir),
				config.GetLessImage(),
				config.GetLessInputFile(), config.GetLessOutputFile())

			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				fmt.Printf("Error compiling CSS: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Successfully compiled CSS in %s\n", cssDir)
		},
	}

	cmd.Flags().BoolP("rtl", "r", false, "Compile RTL (right-to-left) CSS files")

	return cmd
}
