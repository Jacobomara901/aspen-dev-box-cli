package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"adb/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(PullCommand())
}

func PullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull Docker images",
		Long: `Pull all Docker images defined in docker-compose files.
This command scans the ASPEN_DOCKER directory for docker-compose files
and pulls all images defined in their services.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := config.GetProjectsDir()

			err := filepath.Walk(projectsDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("Error accessing path %q: %v\n", path, err)
					return err
				}

				if filepath.Ext(path) == ".yml" {
					// Read the Docker Compose file
					data, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Printf("Error reading Docker Compose file %s: %v\n", path, err)
						return err
					}

					// Parse the Docker Compose file
					var composeFile map[string]interface{}
					if err := yaml.Unmarshal(data, &composeFile); err != nil {
						fmt.Printf("Error parsing Docker Compose file %s: %v\n", path, err)
						return err
					}

					// Extract the services
					services, ok := composeFile["services"].(map[interface{}]interface{})
					if !ok {
						fmt.Printf("Warning: Docker Compose file %s does not define any services\n", path)
						return nil
					}

					// Pull the image for each service
					for _, service := range services {
						serviceMap, ok := service.(map[interface{}]interface{})
						if !ok {
							continue
						}

						image, ok := serviceMap["image"].(string)
						if !ok {
							continue
						}

						fmt.Printf("Pulling image: %s\n", image)
						pullCmd := exec.Command("docker", "pull", image)
						pullCmd.Stdout = os.Stdout
						pullCmd.Stderr = os.Stderr

						if err := pullCmd.Run(); err != nil {
							fmt.Printf("Error pulling image %s: %v\n", image, err)
							return err
						}
					}
				}
				return nil
			})

			if err != nil {
				fmt.Printf("Error scanning Docker Compose files: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
