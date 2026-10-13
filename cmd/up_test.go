package cmd

import (
	"os"
	"testing"
)

func TestProxiedHost(t *testing.T) {
	originalDir, originalStack := cfg.ProjectsDir, cfg.StackName
	t.Cleanup(func() {
		cfg.ProjectsDir, cfg.StackName = originalDir, originalStack
	})
	cfg.ProjectsDir = "/srv/aspen-dev-box-image"

	cases := []struct {
		name     string
		stack    string
		disabled bool
		host     string
		want     string
	}{
		{name: "no proxy keeps the env site name", stack: "aspen-dev-box-image", disabled: true, want: ""},
		{name: "explicit host wins", stack: "aspen-dev-box-image", host: "demo.localhost", want: "demo.localhost"},
		{name: "default instance", stack: "aspen-dev-box-image", want: "localhost"},
		{name: "named stack", stack: "holidays", want: "holidays.localhost"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg.StackName = tc.stack
			if got := proxiedHost(tc.disabled, tc.host); got != tc.want {
				t.Errorf("proxiedHost(%v, %q) = %q, want %q", tc.disabled, tc.host, got, tc.want)
			}
		})
	}
}

func TestSiteNameLookupOverridesOnlyWhenProxied(t *testing.T) {
	t.Setenv("SITE_NAME", "dev.localhost")
	t.Setenv("KOHA_STACK", "kohadev")

	proxied := siteNameLookup("localhost")
	if site, _ := proxied("SITE_NAME"); site != "localhost" {
		t.Errorf("proxied SITE_NAME = %q, want localhost", site)
	}
	if stack, _ := proxied("KOHA_STACK"); stack != "kohadev" {
		t.Errorf("other vars should come from the env, got %q", stack)
	}

	direct := siteNameLookup("")
	if site, _ := direct("SITE_NAME"); site != "dev.localhost" {
		t.Errorf("unproxied SITE_NAME = %q, want dev.localhost", site)
	}

	os.Unsetenv("SITE_NAME")
	if _, ok := direct("SITE_NAME"); ok {
		t.Error("an unset SITE_NAME must stay unset")
	}
}
