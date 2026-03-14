package ouractl

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
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
		link := lipgloss.NewStyle().Foreground(tui.Accent).Underline(true).
			Render("https://cloud.ouraring.com/personal-access-tokens")
		fmt.Printf("🔑 Generate a token at: %s\n", link)

		prompt := lipgloss.NewStyle().Bold(true).Render("Paste your access token:")
		fmt.Printf("%s ", prompt)

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

		success := lipgloss.NewStyle().
			Foreground(tui.Good).
			Bold(true).
			Render("✓ Logged in successfully.")
		fmt.Println(success)
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
		msg := lipgloss.NewStyle().
			Foreground(tui.Good).
			Bold(true).
			Render("✓ Successfully logged out.")
		fmt.Println(msg)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			msg := lipgloss.NewStyle().Foreground(tui.Warn).Render("✗ Not logged in.")
			fmt.Println(msg)
			return nil
		}

		client := api.NewClient(token)
		info, err := client.GetPersonalInfo()
		if err != nil {
			msg := lipgloss.NewStyle().
				Foreground(tui.Warn).
				Render("⚠ Logged in, but token is invalid or expired.")
			fmt.Println(msg)
			return nil
		}

		check := lipgloss.NewStyle().Foreground(tui.Good).Bold(true)
		if info.Email != nil && *info.Email != "" {
			fmt.Printf("%s %s\n", check.Render("✓ Logged in as"), *info.Email)
		} else {
			fmt.Println(check.Render("✓ Logged in."))
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd, logoutCmd, statusCmd)
	rootCmd.AddCommand(authCmd)
}
