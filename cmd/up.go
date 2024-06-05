package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Bring up the Docker Compose project",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			commandArgs := []string{"compose"}

			composeFile := fmt.Sprintf("%s/docker-compose.yml", projectsDir)
			commandArgs = append(commandArgs, "-f", composeFile)

			if debugging {
				composeFile = fmt.Sprintf("%s/docker-compose.debug.yml", projectsDir)
				commandArgs = append(commandArgs, "-f", composeFile)
			}

			if dbgui {
				composeFile = fmt.Sprintf("%s/docker-compose.dbgui.yml", projectsDir)
				commandArgs = append(commandArgs, "-f", composeFile)
			}

			commandArgs = append(commandArgs, "up")

			if detached {
				commandArgs = append(commandArgs, "-d")
			}

			if pullUpdated {
				// Iterate over commandArgs to find all Docker Compose files
				for i, arg := range commandArgs {
					if arg == "-f" && i+1 < len(commandArgs) {
						composeFile := commandArgs[i+1]

						// Load and parse the docker-compose file
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

						// Extract the services from the loadedConfig
						services, ok := loadedConfig["services"].(map[string]interface{})
						if !ok {
							fmt.Println("No services found in docker-compose file")
							os.Exit(1)
						}

						// Iterate over the services and pull the images
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
							err := pullCmd.Run()
							if err != nil {
								fmt.Printf("Error pulling image %s: %v\n", imageName, err)
								os.Exit(1)
							}
						}
					}
				}
			}

			command := exec.Command("docker", commandArgs...)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error bringing up the project: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&detached, "detached", "d", false, "Run in detached mode")
	cmd.Flags().BoolVarP(&debugging, "debugging", "g", false, "Run with debugging compose file")
	cmd.Flags().BoolVarP(&dbgui, "dbgui", "b", false, "Run with dbgui compose file")
	cmd.Flags().BoolVarP(&pullUpdated, "pull", "p", false, "Pull the images for the project only if they have been updated")

	return cmd
}
