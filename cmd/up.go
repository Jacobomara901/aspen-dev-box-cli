package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"adb/pkg/config"
	"github.com/compose-spec/compose-go/loader"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(UpCommand())
}

func UpCommand() *cobra.Command {
	var detached bool
	var debugging bool
	var dbgui bool
	var pullUpdated bool
	var kohaStack string
	var ils string

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Bring up the Docker Compose project",
		Long: `Bring up the Docker Compose project with optional configurations.
You can run in detached mode, with debugging enabled, or with the database GUI.
You can also select which ILS to use (koha or evergreen).`,
		Run: func(cmd *cobra.Command, args []string) {
			commandArgs := []string{"compose", "-f", config.GetDefaultComposeFile()}

			if debugging {
				commandArgs = append(commandArgs, "-f", config.GetDebugComposeFile())
			}

			if dbgui {
				commandArgs = append(commandArgs, "-f", config.GetDBGUIComposeFile())
			}

			// Get ASPEN_DOCKER directory
			aspenDocker := os.Getenv("ASPEN_DOCKER")
			if aspenDocker == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable is not set")
				os.Exit(1)
			}

			// Add ILS-specific compose file
			switch ils {
			case "koha":
				if kohaStack != "" {
					os.Setenv("KOHA_STACK", kohaStack)
				}
				kohaOverride := filepath.Join(aspenDocker, "docker-compose.koha.yml")
				if _, err := os.Stat(kohaOverride); err != nil {
					fmt.Printf("Error: Koha override file not found at %s\n", kohaOverride)
					os.Exit(1)
				}
				commandArgs = append(commandArgs, "-f", kohaOverride)
			case "evergreen":
				evergreenOverride := filepath.Join(aspenDocker, "docker-compose.evergreen.yml")
				if _, err := os.Stat(evergreenOverride); err != nil {
					fmt.Printf("Error: Evergreen override file not found at %s\n", evergreenOverride)
					os.Exit(1)
				}
				commandArgs = append(commandArgs, "-f", evergreenOverride)
			default:
				fmt.Printf("Error: Unsupported ILS '%s'. Supported values: koha, evergreen\n", ils)
				os.Exit(1)
			}

			// Add local.yml file, if it exists
			localComposeFile := filepath.Join(aspenDocker, "local.yml")
			if _, err := os.Stat(localComposeFile); err == nil {
				commandArgs = append(commandArgs, "-f", localComposeFile)
			}

			commandArgs = append(commandArgs, "up")

			if detached {
				commandArgs = append(commandArgs, "-d")
			}

			if pullUpdated {
				pullImages(commandArgs)
			}

			command := exec.Command("docker", commandArgs...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			if err := command.Run(); err != nil {
				fmt.Printf("Error bringing up the project: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&detached, "detached", "d", false, "Run in detached mode")
	cmd.Flags().BoolVarP(&debugging, "debugging", "g", false, "Run with debugging compose file")
	cmd.Flags().BoolVarP(&dbgui, "dbgui", "b", false, "Run with dbgui compose file")
	cmd.Flags().BoolVarP(&pullUpdated, "pull", "p", false, "Pull the images for the project only if they have been updated")
	cmd.Flags().StringVarP(&kohaStack, "koha-stack", "k", "", "Specify the Koha stack to connect to (default: kohadev)")
	cmd.Flags().StringVarP(&ils, "ils", "i", "koha", "Select ILS to use (koha|evergreen)")

	return cmd
}

func pullImages(commandArgs []string) {
	for i, arg := range commandArgs {
		if arg == "-f" && i+1 < len(commandArgs) {
			composeFile := commandArgs[i+1]
			composeFileContent, err := os.ReadFile(composeFile)
			if err != nil {
				fmt.Printf("Error reading docker-compose file: %v\n", err)
				os.Exit(1)
			}

			loadedConfig, err := loader.ParseYAML(composeFileContent)
			if err != nil {
				fmt.Printf("Error parsing docker-compose file: %v\n", err)
				os.Exit(1)
			}

			services, ok := loadedConfig["services"].(map[string]interface{})
			if !ok {
				fmt.Println("No services found in docker-compose file")
				os.Exit(1)
			}

			for _, service := range services {
				serviceMap, ok := service.(map[string]interface{})
				if !ok {
					fmt.Println("Invalid service format in docker-compose file")
					os.Exit(1)
				}

				imageName, ok := serviceMap["image"].(string)
				if !ok {
					fmt.Println("No image name found for service in docker-compose file")
					os.Exit(1)
				}

				pullCmd := exec.Command("docker", "pull", imageName)
				pullCmd.Stdout = os.Stdout
				pullCmd.Stderr = os.Stderr
				if err := pullCmd.Run(); err != nil {
					fmt.Printf("Error pulling image %s: %v\n", imageName, err)
					os.Exit(1)
				}
			}
		}
	}
}
