package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(PullCommand())
}
func PullCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull all images for all docker compose files in the ASPEN_DOCKER directory",
		Run: func(cmd *cobra.Command, args []string) {
			projectsDir := os.Getenv("ASPEN_DOCKER")
			if projectsDir == "" {
				fmt.Println("Error: ASPEN_DOCKER environment variable not set.")
				os.Exit(1)
			}

			err := filepath.Walk(projectsDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("Error accessing a path %q: %v\n", path, err)
					return err
				}

				if filepath.Ext(path) == ".yml" {
					// Read the Docker Compose file
					data, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Printf("Error reading Docker Compose file %s: %v\n", path, err)
						os.Exit(1)
					}

					// Parse the Docker Compose file
					var composeFile map[string]interface{}
					err = yaml.Unmarshal(data, &composeFile)
					if err != nil {
						fmt.Printf("Error parsing Docker Compose file %s: %v\n", path, err)
						os.Exit(1)
					}

					// Extract the services
					services, ok := composeFile["services"].(map[interface{}]interface{})
					if !ok {
						fmt.Printf("Error: Docker Compose file %s does not define any services.\n", path)
						os.Exit(1)
					}

					// Pull the image for each service
					for _, service := range services {
						serviceMap, ok := service.(map[interface{}]interface{})
						if ok {
							image, ok := serviceMap["image"].(string)
							if ok {
								pullCmd := exec.Command("docker", "pull", image)
								pullCmd.Stdout = os.Stdout
								pullCmd.Stderr = os.Stderr

								err := pullCmd.Run()
								if err != nil {
									fmt.Printf("Error pulling Docker image %s: %v\n", image, err)
									os.Exit(1)
								}
							}
						}
					}
				}
				return nil
			})

			if err != nil {
				fmt.Printf("Error walking the path %v: %v\n", projectsDir, err)
				os.Exit(1)
			}
		},
	}
}
