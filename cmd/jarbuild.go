package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(JarBuildCommand())
}

func runDockerBuild(jarFile string) error {
	workDir := fmt.Sprintf("/app/code/%s", jarFile)

	fmt.Printf("\n\033[1;34mRecompiling JAR file: %s\033[0m\n", jarFile)

	command := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", config.GetAspenCloneDir()),
		"-w", workDir,
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		config.GetJavaBuildImage(), "bash", "-c", `
			mkdir -p bin && \
			javac -cp "$(find /app -name '*.jar' | tr '\n' ':')" -d bin $(find src -name '*.java') $(find `+config.GetJavaSharedLibrariesPath()+` -name '*.java') && \
			jar cfm $(basename $(pwd)).jar META-INF/MANIFEST.MF -C bin . && \
			rm -rf bin
		`,
	)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func JarBuildCommand() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "jarbuild",
		Short: "Build Java JAR files",
		Long: `Build Java JAR files from source code.
This command can build either a single JAR file selected via fzf or all JAR files at once.`,
		Run: func(cmd *cobra.Command, args []string) {
			if all {
				buildAllJars()
			} else {
				buildSingleJar()
			}
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Build all JAR files")
	return cmd
}

func buildAllJars() {
	findCmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", config.GetAspenCloneDir()),
		"-w", "/app",
		config.GetAlpineImage(), "sh", "-c", fmt.Sprintf(`
			apk add --no-cache findutils > /dev/null && \
			find /app/code -mindepth 2 -maxdepth 2 -name '*.jar' | grep -v "%s" | xargs -n 1 basename | sed 's/\.jar$//'
		`, strings.ReplaceAll(config.GetExcludedJarPatterns(), " ", "\\|")),
	)

	findOutput, err := findCmd.Output()
	if err != nil {
		fmt.Printf("Error finding JAR files: %v\n", err)
		os.Exit(1)
	}

	jarFiles := strings.Split(strings.TrimSpace(string(findOutput)), "\n")
	for _, jarFile := range jarFiles {
		if jarFile == "" {
			continue
		}

		if err := runDockerBuild(jarFile); err != nil {
			fmt.Printf("Error building JAR file %s: %v\n", jarFile, err)
			os.Exit(1)
		}
	}
}

func buildSingleJar() {
	tmpFile, err := os.CreateTemp(config.GetAspenCloneDir(), "fzf-output")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())

	tmpFileName := filepath.Join("/app", filepath.Base(tmpFile.Name()))

	fzfCmd := exec.Command("docker", "run", "--rm", "-it",
		"-v", fmt.Sprintf("%s:/app", config.GetAspenCloneDir()),
		"-w", "/app",
		config.GetAlpineImage(), "sh", "-c", fmt.Sprintf(`
			apk add --no-cache fzf findutils > /dev/null && \
			find /app/code -mindepth 2 -maxdepth 2 -name '*.jar' | grep -v "%s" | xargs -n 1 basename | sed 's/\.jar$//' | fzf > %s
		`, strings.ReplaceAll(config.GetExcludedJarPatterns(), " ", "\\|"), tmpFileName),
	)

	fzfCmd.Stdin = os.Stdin
	fzfCmd.Stdout = os.Stdout
	fzfCmd.Stderr = os.Stderr

	if err := fzfCmd.Run(); err != nil {
		fmt.Printf("Error selecting JAR file with fzf: %v\n", err)
		os.Exit(1)
	}

	fzfOutput, err := os.ReadFile(filepath.Join(config.GetAspenCloneDir(), filepath.Base(tmpFile.Name())))
	if err != nil {
		fmt.Printf("Error reading fzf output: %v\n", err)
		os.Exit(1)
	}

	selectedJar := strings.TrimSpace(string(fzfOutput))
	if selectedJar == "" {
		fmt.Println("No JAR file selected.")
		os.Exit(1)
	}

	if err := runDockerBuild(selectedJar); err != nil {
		fmt.Printf("Error building JAR file: %v\n", err)
		os.Exit(1)
	}
}