package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"adb/pkg/docker"

	"github.com/spf13/cobra"
)

type aspenJob struct {
	Dir  string
	Bin  []string
}

var aspenJobs = map[string]aspenJob{
	"koha-export":  {Dir: "/usr/local/aspen-discovery/code/koha_export", Bin: []string{"java", "-jar", "koha_export.jar"}},
	"reindexer":    {Dir: "/usr/local/aspen-discovery/code/reindexer", Bin: []string{"java", "-jar", "reindexer.jar"}},
	"oai-indexer":  {Dir: "/usr/local/aspen-discovery/code/oai_indexer", Bin: []string{"java", "-jar", "oai_indexer.jar"}},
	"sideload":     {Dir: "/usr/local/aspen-discovery/code/sideload_processing", Bin: []string{"java", "-jar", "sideload_processing.jar"}},
	"user-lists":   {Dir: "/usr/local/aspen-discovery/code/user_list_indexer", Bin: []string{"java", "-jar", "user_list_indexer.jar"}},
	"cron":         {Dir: "/usr/local/aspen-discovery/docker/files/cron", Bin: []string{"php", "checkBackgroundProcessesDocker.php"}},
	"sitemaps":     {Dir: "/usr/local/aspen-discovery/code/web/cron", Bin: []string{"php", "createSitemaps.php"}},
}

func init() {
	rootCmd.AddCommand(RunCommand())
}

func RunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run <job> [extra args...]",
		Short: "Run an aspen background job (reindexer, koha-export, etc.)",
		Long:  buildRunLong(),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runner, err := docker.NewRunner()
			if err != nil {
				return fmt.Errorf("initialize docker: %w", err)
			}
			defer runner.Close()
			return runJob(cmd.Context(), runner, args[0], args[1:])
		},
	}
}

func buildRunLong() string {
	names := make([]string, 0, len(aspenJobs))
	for k := range aspenJobs {
		names = append(names, k)
	}
	sort.Strings(names)
	return "Available jobs: " + strings.Join(names, ", ") +
		"\n\nExamples:\n  adb run koha-export\n  adb run reindexer nightly\n  adb run cron"
}

func runJob(ctx context.Context, runner *docker.SDKRunner, name string, extra []string) error {
	job, ok := aspenJobs[name]
	if !ok {
		return fmt.Errorf("unknown job %q (try one of: %s)", name, sortedJobNames())
	}
	parts := append([]string{}, job.Bin...)
	parts = append(parts, "${SITE_NAME:-dev.localhost}")
	parts = append(parts, extra...)
	script := fmt.Sprintf("cd %s && %s", job.Dir, strings.Join(parts, " "))
	return runner.ExecInteractive(ctx, docker.ExecConfig{
		Container:  cfg.MainContainerName(),
		User:       "www-data",
		Cmd:        []string{"sh", "-c", script},
	})
}

func sortedJobNames() string {
	names := make([]string, 0, len(aspenJobs))
	for k := range aspenJobs {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
