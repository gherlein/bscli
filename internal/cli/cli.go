package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	jsonOutput bool
	insecure bool

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
	// Check environment variables for default values
	debugDefault := os.Getenv("BSCLI_TEST_DEBUG") == "true"
	insecureDefault := os.Getenv("BSCLI_TEST_INSECURE") == "true"
	
	// Global flags (no longer need host flag)
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "admin", "Username for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for authentication")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", debugDefault, "Enable debug output")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output raw JSON (for scripts)")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "local", "l", insecureDefault, "Accept locally signed certificates (use HTTPS with insecure TLS)")

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
		Insecure: insecure,
	}

	return brightsign.NewClient(config), nil
}

// handleError prints an error message and exits
func handleError(err error) {
	errMsg := err.Error()
	
	// Check for TLS certificate errors and provide helpful suggestions
	if isTLSError(errMsg) {
		helpfulMsg := errMsg + "\n\nThis appears to be a TLS certificate error. The player may be using a self-signed certificate.\nTry one of the following:\n  1. Use the --local or -l flag to accept locally signed certificates\n  2. Set environment variable: export BSCLI_TEST_INSECURE=true"
		if jsonOutput {
			// For JSON mode, include the helpful message in JSON
			errorObj := map[string]string{
				"error": errMsg,
				"suggestion": "This appears to be a TLS certificate error. Try using --local or -l flag, or set BSCLI_TEST_INSECURE=true",
			}
			json.NewEncoder(os.Stdout).Encode(errorObj)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", helpfulMsg)
		}
	} else {
		// Regular error handling
		if jsonOutput {
			// For JSON mode, output error as JSON to stdout (not stderr for proper JSON parsing)
			errorObj := map[string]string{"error": errMsg}
			json.NewEncoder(os.Stdout).Encode(errorObj)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
	os.Exit(1)
}

// isTLSError checks if an error message indicates a TLS certificate problem
func isTLSError(errMsg string) bool {
	tlsIndicators := []string{
		"x509:",
		"certificate",
		"tls:",
		"TLS",
		"self-signed",
		"verify certificate",
		"certificate is not standards compliant",
	}
	
	for _, indicator := range tlsIndicators {
		if strings.Contains(errMsg, indicator) {
			return true
		}
	}
	return false
}

// outputJSON outputs data as JSON when --json flag is used
func outputJSON(data interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
		handleError(fmt.Errorf("failed to encode JSON: %w", err))
	}
}