package cmd

import (
	"context"
	"fmt"
	"os"

	"adb/pkg/docker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(SeedCommand())
}

func SeedCommand() *cobra.Command {
	return &cobra.Command{
		Use:                "seed [-- args...]",
		Short:              "Seed dev data into aspen (libraries, etc.)",
		Long:               "Run the aspen seeder. Examples:\n  adb seed list\n  adb seed build library 500",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			runner, err := docker.NewRunner()
			if err != nil {
				return fmt.Errorf("initialize docker: %w", err)
			}
			defer runner.Close()
			return seedExec(cmd.Context(), runner, args)
		},
	}
}

func seedExec(ctx context.Context, runner *docker.SDKRunner, args []string) error {
	cmd := append([]string{"php", "/seeder/seed.php"}, args...)
	result, err := runner.Exec(ctx, docker.ExecConfig{
		Container: cfg.MainContainerName(),
		Cmd:       cmd,
		User:      "www-data",
	})
	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	fmt.Print(result.Stdout)
	if result.Stderr != "" {
		fmt.Fprint(os.Stderr, result.Stderr)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("seed exited with %d", result.ExitCode)
	}
	return nil
}
