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
	cmd := &cobra.Command{
		Use:   "oauth <client_id> <client_secret>",
		Short: "Update the OAuth client ID and secret for Koha ils logins for ADB",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dbContainerName := "aspen-db"
			oAuthClientId := args[0]
			oAuthClientSecret := args[1]

			print, _ := cmd.Flags().GetBool("print")
			driver, _ := cmd.Flags().GetString("driver")
			if driver == "" {
				driver = "Koha"
			}

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

			command := exec.Command("docker", "exec", "-it", dbContainerName, "/bin/bash", "-c", fmt.Sprintf("echo \"%s\" | mariadb -E -uroot -paspen aspen", sql))
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
	// Add this line outside the Run function, but inside the function where the command is defined.
	cmd.Flags().StringP("driver", "d", "", "Specify the driver (default is 'Koha')")
	cmd.Flags().BoolP("print", "p", false, "Print the rows that match the driver")

	return cmd
}
