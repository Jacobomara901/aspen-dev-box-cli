package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(OAuthCommand())
}

func OAuthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "oauth <client_id> <client_secret>",
		Short: "Update the OAuth client ID and secret for Koha ils logins for ADB",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dbContainerName := "aspen-db"
			oAuthClientId := args[0]
			oAuthClientSecret := args[1]

			sql := fmt.Sprintf("UPDATE account_profiles SET oAuthClientId='%s', oAuthClientSecret='%s' WHERE driver='Koha'", oAuthClientId, oAuthClientSecret)

			command := exec.Command("docker", "exec", "-it", dbContainerName, "/bin/bash", "-c", fmt.Sprintf("echo \"%s\" | mariadb -uroot -paspen aspen", sql))
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Error updating OAuth credentials in the database: %v\n", err)
				os.Exit(1)
			}
		},
	}
}
