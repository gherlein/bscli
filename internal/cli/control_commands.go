package cli

import (
	"fmt"
	"strings"

	"bscli/pkg/brightsign"
	"github.com/spf13/cobra"
)

func addControlCommands() {
	controlCmd := &cobra.Command{
		Use:   "control",
		Short: "Player control commands",
		Long:  "Commands for controlling the BrightSign player",
	}

	// Reboot command
	rebootCmd := &cobra.Command{
		Use:   "reboot",
		Short: "Reboot the player",
		Run: func(cmd *cobra.Command, args []string) {
			crashReport, _ := cmd.Flags().GetBool("crash-report")
			factoryReset, _ := cmd.Flags().GetBool("factory-reset")
			disableAutorun, _ := cmd.Flags().GetBool("disable-autorun")

			// Confirm dangerous operations
			if factoryReset {
				fmt.Print("WARNING: Factory reset will erase all settings. Continue? (y/N): ")
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

			options := &brightsign.RebootOptions{
				CrashReport:    crashReport,
				FactoryReset:   factoryReset,
				DisableAutorun: disableAutorun,
			}

			err = client.Control.Reboot(options)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Reboot initiated")
		},
	}
	rebootCmd.Flags().Bool("crash-report", false, "Generate crash report")
	rebootCmd.Flags().Bool("factory-reset", false, "Perform factory reset")
	rebootCmd.Flags().Bool("disable-autorun", false, "Disable autorun after reboot")

	// Snapshot command
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Take a snapshot of current display",
		Run: func(cmd *cobra.Command, args []string) {
			width, _ := cmd.Flags().GetInt("width")
			height, _ := cmd.Flags().GetInt("height")
			fullRes, _ := cmd.Flags().GetBool("full-resolution")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			options := &brightsign.SnapshotOptions{
				Width:                      width,
				Height:                     height,
				ShouldCaptureFullResolution: fullRes,
			}

			filename, err := client.Control.TakeSnapshot(options)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Snapshot saved: %s\n", filename)
		},
	}
	snapshotCmd.Flags().Int("width", 0, "Width of snapshot")
	snapshotCmd.Flags().Int("height", 0, "Height of snapshot")
	snapshotCmd.Flags().Bool("full-resolution", false, "Capture at full resolution")

	// DWS password commands
	dwsPasswordCmd := &cobra.Command{
		Use:   "dws-password",
		Short: "Manage DWS password",
	}

	dwsPasswordGetCmd := &cobra.Command{
		Use:   "status",
		Short: "Check if DWS password is set",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			info, err := client.Control.GetDWSPassword()
			if err != nil {
				handleError(err)
			}

			if info.IsSet {
				fmt.Println("DWS password is set")
			} else {
				fmt.Println("DWS password is not set")
			}
		},
	}

	dwsPasswordSetCmd := &cobra.Command{
		Use:   "set [password]",
		Short: "Set DWS password",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			reset, _ := cmd.Flags().GetBool("reset")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.DWSPassword{
				Reset: reset,
			}

			if !reset && len(args) > 0 {
				config.Password = args[0]
			}

			err = client.Control.SetDWSPassword(config)
			if err != nil {
				handleError(err)
			}

			if reset {
				fmt.Println("DWS password reset to default")
			} else {
				fmt.Println("DWS password set")
			}
		},
	}
	dwsPasswordSetCmd.Flags().Bool("reset", false, "Reset password to default")

	dwsPasswordCmd.AddCommand(dwsPasswordGetCmd, dwsPasswordSetCmd)

	// Local DWS commands
	localDWSCmd := &cobra.Command{
		Use:   "local-dws",
		Short: "Manage local DWS",
	}

	localDWSStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check if local DWS is enabled",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config, err := client.Control.GetLocalDWS()
			if err != nil {
				handleError(err)
			}

			if config.Enabled {
				fmt.Println("Local DWS is enabled")
			} else {
				fmt.Println("Local DWS is disabled")
			}
		},
	}

	localDWSEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable local DWS",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Control.SetLocalDWS(true)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Local DWS enabled")
		},
	}

	localDWSDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable local DWS",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Control.SetLocalDWS(false)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Local DWS disabled")
		},
	}

	localDWSCmd.AddCommand(localDWSStatusCmd, localDWSEnableCmd, localDWSDisableCmd)

	// Download firmware command
	downloadFirmwareCmd := &cobra.Command{
		Use:   "download-firmware [url]",
		Short: "Download and install firmware from URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]

			// Validate URL
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				handleError(fmt.Errorf("invalid URL: must start with http:// or https://"))
			}

			fmt.Printf("WARNING: This will download and install firmware from %s\n", url)
			fmt.Print("The player will reboot automatically. Continue? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return
			}

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Control.DownloadFirmware(url)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Firmware download initiated, player will reboot")
		},
	}

	controlCmd.AddCommand(rebootCmd, snapshotCmd, dwsPasswordCmd, localDWSCmd, downloadFirmwareCmd)
	rootCmd.AddCommand(controlCmd)
}