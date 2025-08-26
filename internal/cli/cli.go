package cli

import (
	"fmt"
	"os"
	"syscall"

	"bscli/pkg/brightsign"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// Global flags
	host     string
	username string
	password string
	debug    bool

	// Root command
	rootCmd = &cobra.Command{
		Use:   "bscli [host] [command]",
		Short: "BrightSign CLI for controlling players via DWS API",
		Long: `bscli is a command-line interface for managing BrightSign players
through their Diagnostic Web Server (DWS) API.

Usage: bscli [host] [command] [args...]

Examples:
  bscli 192.168.1.100 info device
  bscli player.local file list /storage/sd/
  bscli 10.0.0.50 control reboot

It provides commands for:
  - Device information and status
  - File management (upload, download, list)
  - System control (reboot, snapshot, etc.)
  - Network diagnostics
  - Registry management
  - Display control
  - And more...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

// Execute runs the CLI
func Execute() error {
	// Parse host from command line arguments manually
	args := os.Args[1:] // Skip program name
	
	if len(args) == 0 {
		return rootCmd.Help()
	}
	
	// First argument should be the host
	host = args[0]
	
	// Set remaining arguments for cobra to parse
	rootCmd.SetArgs(args[1:])
	
	return rootCmd.Execute()
}

func init() {
	// Global flags (no longer need host flag)
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "admin", "Username for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for authentication")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

	// Add command groups
	addInfoCommands()
	addControlCommands()
	addFileCommands()
	addDiagnosticsCommands()
	addDisplayCommands()
	addRegistryCommands()
	addLogsCommands()
	addVideoCommands()
}

// getClient creates a BrightSign client with authentication
func getClient() (*brightsign.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host is required")
	}

	// Prompt for password if not provided
	if password == "" {
		fmt.Printf("Password for %s@%s: ", username, host)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()
		password = string(bytePassword)
	}

	config := brightsign.Config{
		Host:     host,
		Username: username,
		Password: password,
		Debug:    debug,
	}

	return brightsign.NewClient(config), nil
}

// handleError prints an error message and exits
func handleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}