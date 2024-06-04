package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(CSSCommand())
}

func CSSCommand() *cobra.Command {
	var rtl bool

	cmd := &cobra.Command{
		Use:   "css",
		Short: "Run the CSS command",
		Run: func(cmd *cobra.Command, args []string) {
			cssDir := os.Getenv("ASPEN_CLONE") + "/code/web/interface/themes/responsive/css"
			if rtl {
				cssDir += "-rtl"
			}

			if cssDir == "" {
				fmt.Println("Error: ASPEN_CLONE environment variable not set.")
				os.Exit(1)
			}

			command := exec.Command("docker", "run", "--rm",
				"-u", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
				"-v", fmt.Sprintf("%s:%s", cssDir, cssDir),
				"-w", cssDir,
				"ghcr.io/sndsgd/less", "main.less", "main.css", "--source-map",
			)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error running the CSS command: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&rtl, "rtl", "r", false, "Use RTL CSS directory")

	return cmd
}
