package cmd

import (
	"testing"

	"adb/pkg/docker"
)

func TestILSStatusReportsTheSetupJob(t *testing.T) {
	cases := []struct {
		name  string
		stack docker.StackSummary
		want  string
	}{
		{
			name:  "koha ok",
			stack: docker.StackSummary{ILS: "koha", Services: map[string]string{kohaOAuthService: "exited(0)"}},
			want:  "koha (oauth ok)",
		},
		{
			name:  "koha failed",
			stack: docker.StackSummary{ILS: "koha", Services: map[string]string{kohaOAuthService: "exited(1)"}},
			want:  "koha (oauth failed, see adb logs --service koha-oauth-init)",
		},
		{
			name:  "koha pending",
			stack: docker.StackSummary{ILS: "koha", Services: map[string]string{kohaOAuthService: "running"}},
			want:  "koha (oauth pending)",
		},
		{
			name:  "evergreen ok",
			stack: docker.StackSummary{ILS: "evergreen", Services: map[string]string{evergreenImportService: "exited(0)"}},
			want:  "evergreen (import ok)",
		},
		{
			name:  "evergreen failed",
			stack: docker.StackSummary{ILS: "evergreen", Services: map[string]string{evergreenImportService: "exited(2)"}},
			want:  "evergreen (import failed, see adb logs --service evergreen-import-init)",
		},
		{
			name:  "evergreen pending",
			stack: docker.StackSummary{ILS: "evergreen", Services: map[string]string{evergreenImportService: "created"}},
			want:  "evergreen (import pending)",
		},
		{
			name:  "no setup job",
			stack: docker.StackSummary{ILS: "evergreen", Services: map[string]string{"aspen-dev-box": "healthy"}},
			want:  "evergreen",
		},
		{
			name:  "no ils",
			stack: docker.StackSummary{ILS: "none", Services: map[string]string{kohaOAuthService: "exited(0)"}},
			want:  "none",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ilsStatus(tc.stack); got != tc.want {
				t.Errorf("ilsStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}
