package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"adb/pkg/config"
)

func TestImagesFromComposeFilesSubstitutesAndDedupes(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.yml")
	overlay := filepath.Join(dir, "overlay.yml")
	writeCompose(t, base, `
services:
  aspen-dev-box:
    image: ${ADB_TEST_IMAGE:-aspendiscovery/aspen:latest}
  aspen-db:
    image: "mariadb:10.11"
  helper:
    build: .
`)
	writeCompose(t, overlay, `
services:
  unit-tests:
    image: ${ADB_TEST_IMAGE:-aspendiscovery/aspen:latest}
  proxy:
    image: traefik:v3
`)
	t.Setenv("ADB_TEST_IMAGE", "")
	os.Unsetenv("ADB_TEST_IMAGE")

	images, err := imagesFromComposeFiles([]string{base, overlay})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"aspendiscovery/aspen:latest", "mariadb:10.11", "traefik:v3"}
	if len(images) != len(want) {
		t.Fatalf("images = %v, want %v", images, want)
	}
	for i := range want {
		if images[i] != want[i] {
			t.Errorf("images[%d] = %q, want %q", i, images[i], want[i])
		}
	}

	t.Setenv("ADB_TEST_IMAGE", "example/aspen:pinned")
	images, err = imagesFromComposeFiles([]string{base})
	if err != nil {
		t.Fatal(err)
	}
	if images[0] != "example/aspen:pinned" {
		t.Errorf("env override not applied: %v", images)
	}
}

func TestPullableComposeFilesTakesEveryOverlayAndTheProxy(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "compose"))
	for _, name := range []string{config.DefaultComposeFile, config.KohaComposeFile, config.EvergreenComposeFile} {
		mustWrite(t, filepath.Join(dir, "compose", name))
	}

	original := cfg.ProjectsDir
	cfg.ProjectsDir = dir
	t.Cleanup(func() { cfg.ProjectsDir = original })

	files, err := pullableComposeFiles()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(dir, "compose", config.EvergreenComposeFile),
		filepath.Join(dir, "compose", config.KohaComposeFile),
		filepath.Join(dir, "compose", config.DefaultComposeFile),
		cfg.ProxyComposeFilePath(),
	}
	if len(files) != len(want) {
		t.Fatalf("files = %v, want %v", files, want)
	}
	for i := range want {
		if files[i] != want[i] {
			t.Errorf("files[%d] = %q, want %q", i, files[i], want[i])
		}
	}
}

func writeCompose(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
