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
	"reindexer":               jarJob("reindexer"),
	"oai-indexer":             jarJob("oai_indexer"),
	"sideload":                jarJob("sideload_processing"),
	"user-lists":              jarJob("user_list_indexer"),
	"course-reserves":         jarJob("course_reserves_indexer"),
	"events-indexer":          jarJob("events_indexer"),
	"series-indexer":          jarJob("series_indexer"),
	"web-indexer":             jarJob("web_indexer"),
	"marc-merge":              jarJob("marcMergeUtility"),
	"cron-jar":                jarJob("cron"),
	"koha-export":             jarJob("koha_export"),
	"evergreen-export":        jarJob("evergreen_export"),
	"polaris-export":          jarJob("polaris_export"),
	"sierra-export":           jarJob("sierra_export_api"),
	"carlx-export":            jarJob("carlx_export"),
	"symphony-export":         jarJob("symphony_export"),
	"evolve-export":           jarJob("evolve_export"),
	"axis-360-export":         jarJob("axis_360_export"),
	"hoopla-export":           jarJob("hoopla_export"),
	"overdrive-export":        jarJob("overdrive_extract"),
	"cloud-library-export":    jarJob("cloud_library_export"),
	"palace-project-export":   jarJob("palace_project_export"),
	"cron":                    {Dir: "/usr/local/aspen-discovery/docker/files/cron", Bin: []string{"php", "checkBackgroundProcessesDocker.php"}},
	"sitemaps":                {Dir: "/usr/local/aspen-discovery/code/web/cron", Bin: []string{"php", "createSitemaps.php"}},
}

func jarJob(name string) aspenJob {
	return aspenJob{
		Dir: "/usr/local/aspen-discovery/code/" + name,
		Bin: []string{"java", "-jar", name + ".jar"},
	}
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
