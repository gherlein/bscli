package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addLogsCommands() {
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Log management commands",
		Long:  "Commands for retrieving and managing player logs",
	}

	// Get logs command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get player serial logs",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			logs, err := client.Logs.GetLogs()
			if err != nil {
				handleError(err)
			}

			fmt.Println(logs)
		},
	}

	// Supervisor logging level commands
	supervisorCmd := &cobra.Command{
		Use:   "supervisor",
		Short: "Manage supervisor logging level",
	}

	supervisorGetCmd := &cobra.Command{
		Use:   "get-level",
		Short: "Get supervisor logging level",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			level, err := client.Logs.GetSupervisorLoggingLevel()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Supervisor logging level: %s\n", level)
		},
	}

	supervisorSetCmd := &cobra.Command{
		Use:   "set-level [level]",
		Short: "Set supervisor logging level (0=error, 1=warn, 2=info, 3=trace)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var level int
			n, err := fmt.Sscanf(args[0], "%d", &level)
			if err != nil || n != 1 {
				handleError(fmt.Errorf("invalid level: must be 0-3"))
			}

			if level < 0 || level > 3 {
				handleError(fmt.Errorf("invalid level: must be 0-3"))
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Logs.SetSupervisorLoggingLevel(level)
			if err != nil {
				handleError(err)
			}

			levelNames := []string{"error", "warn", "info", "trace"}
			fmt.Printf("Supervisor logging level set to %d (%s)\n", level, levelNames[level])
		},
	}

	supervisorCmd.AddCommand(supervisorGetCmd, supervisorSetCmd)
	logsCmd.AddCommand(getCmd, supervisorCmd)
	rootCmd.AddCommand(logsCmd)
}