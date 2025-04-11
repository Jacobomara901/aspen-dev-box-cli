package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"adb/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(OAuthCommand())
}

func OAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oauth <client_id> <client_secret>",
		Short: "Update OAuth credentials",
		Long: `Update the OAuth client ID and secret for ILS logins.
This command updates the OAuth credentials in the database for the specified driver.
By default, it updates the Koha driver credentials.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			oAuthClientId := args[0]
			oAuthClientSecret := args[1]

			print, _ := cmd.Flags().GetBool("print")
			driver, _ := cmd.Flags().GetString("driver")
			if driver == "" {
				driver = "Koha"
			}

			// Build the SQL query
			sql := fmt.Sprintf(`
			SET @update_count = 0;
			UPDATE account_profiles 
			SET oAuthClientId='%s', 
			oAuthClientSecret='%s' 
			WHERE driver='%s';
			SET @update_count = ROW_COUNT();

			SELECT @update_count as Changed_Rows;
			`, oAuthClientId, oAuthClientSecret, driver)

			if print {
				sql += fmt.Sprintf(`
				SELECT * FROM account_profiles 
				WHERE driver='%s'
				`, driver)
			}

			// Execute the SQL query in the database container
			dockerCmd := exec.Command("docker", "exec", "-it",
				config.GetDBContainerName(),
				"/bin/bash", "-c",
				fmt.Sprintf("echo \"%s\" | mariadb %s", sql, config.GetDBConnectionString()))

			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr

			if err := dockerCmd.Run(); err != nil {
				fmt.Printf("Error updating OAuth credentials: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("driver", "d", "", "Specify the driver (default is 'Koha')")
	cmd.Flags().BoolP("print", "p", false, "Print the rows that match the driver")

	return cmd
}
