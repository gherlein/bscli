package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func addRegistryCommands() {
	registryCmd := &cobra.Command{
		Use:     "registry",
		Aliases: []string{"reg"},
		Short:   "Registry management commands",
		Long:    "Commands for managing player registry settings",
	}

	// Get all registry
	getAllCmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get entire registry dump",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			registry, err := client.Registry.GetAll()
			if err != nil {
				handleError(err)
			}

			data, _ := json.MarshalIndent(registry, "", "  ")
			fmt.Println(string(data))
		},
	}

	// Get specific value
	getCmd := &cobra.Command{
		Use:   "get [section] [key]",
		Short: "Get specific registry value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			value, err := client.Registry.GetValue(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("%s/%s = %s\n", args[0], args[1], value)
		},
	}

	// Set value
	setCmd := &cobra.Command{
		Use:   "set [section] [key] [value]",
		Short: "Set registry value",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Registry.SetValue(args[0], args[1], args[2])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Set %s/%s = %s\n", args[0], args[1], args[2])
		},
	}

	// Delete value
	deleteCmd := &cobra.Command{
		Use:   "delete [section] [key]",
		Short: "Delete registry value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				fmt.Printf("Delete %s/%s? (y/N): ", args[0], args[1])
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Cancelled")
					return
				}
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Registry.DeleteValue(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Deleted %s/%s\n", args[0], args[1])
		},
	}
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Delete section
	deleteSectionCmd := &cobra.Command{
		Use:   "delete-section [section]",
		Short: "Delete entire registry section",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				fmt.Printf("WARNING: Delete entire section %s? This will remove all keys. (y/N): ", args[0])
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Cancelled")
					return
				}
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Registry.DeleteSection(args[0])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Deleted section %s\n", args[0])
		},
	}
	deleteSectionCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Recovery URL commands
	recoveryURLCmd := &cobra.Command{
		Use:   "recovery-url",
		Short: "Manage recovery URL",
	}

	recoveryURLGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get recovery URL",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			url, err := client.Registry.GetRecoveryURL()
			if err != nil {
				handleError(err)
			}

			if url != "" {
				fmt.Printf("Recovery URL: %s\n", url)
			} else {
				fmt.Println("No recovery URL set")
			}
		},
	}

	recoveryURLSetCmd := &cobra.Command{
		Use:   "set [url]",
		Short: "Set recovery URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]

			// Basic URL validation
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				handleError(fmt.Errorf("invalid URL: must start with http:// or https://"))
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Registry.SetRecoveryURL(url)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Recovery URL set to: %s\n", url)
		},
	}

	recoveryURLCmd.AddCommand(recoveryURLGetCmd, recoveryURLSetCmd)

	// Flush command
	flushCmd := &cobra.Command{
		Use:   "flush",
		Short: "Flush registry to persistent storage (BOS 9.0.107+)",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Registry.Flush()
			if err != nil {
				handleError(err)
			}

			fmt.Println("Registry flushed to persistent storage")
		},
	}

	// Search command
	searchCmd := &cobra.Command{
		Use:   "search [term]",
		Short: "Search registry keys and values",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := strings.ToLower(args[0])
			ignoreCase := true

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			registry, err := client.Registry.GetAll()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Search results for '%s':\n", args[0])
			found := false

			for section, keys := range registry {
				sectionLower := strings.ToLower(section)
				
				for key, value := range keys {
					keyLower := strings.ToLower(key)
					valueLower := strings.ToLower(value)

					if (ignoreCase && (strings.Contains(sectionLower, searchTerm) ||
						strings.Contains(keyLower, searchTerm) ||
						strings.Contains(valueLower, searchTerm))) ||
						(!ignoreCase && (strings.Contains(section, args[0]) ||
						strings.Contains(key, args[0]) ||
						strings.Contains(value, args[0]))) {
						
						fmt.Printf("  %s/%s = %s\n", section, key, value)
						found = true
					}
				}
			}

			if !found {
				fmt.Println("  No matches found")
			}
		},
	}

	registryCmd.AddCommand(getAllCmd, getCmd, setCmd, deleteCmd, deleteSectionCmd, 
		recoveryURLCmd, flushCmd, searchCmd)
	rootCmd.AddCommand(registryCmd)
}