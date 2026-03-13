package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hagelstam/oura-cli/internal/api"
	"github.com/hagelstam/oura-cli/internal/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with an Oura access token",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Generate a token at: https://cloud.ouraring.com/personal-access-tokens")
		fmt.Print("Paste your access token: ")

		raw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		token := strings.TrimSpace(string(raw))
		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		// Validate the token with a test API call.
		client := api.NewClient(token)
		if _, err := client.Get("/v2/usercollection/personal_info", nil); err != nil {
			return fmt.Errorf("token validation failed: %w", err)
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Logged in successfully.")
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and remove stored token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.RemoveToken(); err != nil {
			return fmt.Errorf("failed to remove token: %w", err)
		}
		fmt.Println("Logged out.")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			fmt.Println("Not logged in.")
			return nil
		}

		// Verify the token still works.
		client := api.NewClient(token)
		if _, err := client.Get("/v2/usercollection/personal_info", nil); err != nil {
			fmt.Println("Logged in, but token is invalid or expired.")
			return nil
		}

		fmt.Println("Logged in.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd, logoutCmd, statusCmd)
	rootCmd.AddCommand(authCmd)
}
