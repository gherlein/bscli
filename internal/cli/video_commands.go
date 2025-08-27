package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addVideoCommands() {
	videoCmd := &cobra.Command{
		Use:   "video",
		Short: "Video output management commands",
		Long:  "Commands for managing video outputs and settings",
	}

	// Output info command
	outputInfoCmd := &cobra.Command{
		Use:   "output-info [connector] [device]",
		Short: "Get video output information",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			info, err := client.Video.GetOutputInfo(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(info)
				return
			}

			fmt.Printf("Connector: %s\n", info.Connector)
			fmt.Printf("Device: %s\n", info.Device)
			fmt.Printf("Connected: %v\n", info.Connected)
			if info.Connected {
				fmt.Printf("Resolution: %dx%d @ %dHz\n", info.Width, info.Height, info.RefreshRate)
				if info.InterlaceMode != "" {
					fmt.Printf("Interlace Mode: %s\n", info.InterlaceMode)
				}
				if info.PreferredMode != "" {
					fmt.Printf("Preferred Mode: %s\n", info.PreferredMode)
				}
			}
		},
	}

	// EDID command
	edidCmd := &cobra.Command{
		Use:   "edid [connector] [device]",
		Short: "Get EDID information from connected display",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			edid, err := client.Video.GetEDID(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(edid)
				return
			}

			fmt.Printf("Manufacturer: %s\n", edid.Manufacturer)
			fmt.Printf("Product: %s\n", edid.Product)
			fmt.Printf("Serial Number: %s\n", edid.SerialNumber)
			fmt.Printf("Manufacturing: Week %d of %d\n", edid.WeekOfManufacture, edid.YearOfManufacture)
			fmt.Printf("EDID Version: %s\n", edid.Version)
			fmt.Printf("Digital: %v\n", edid.Digital)
			fmt.Printf("Display Size: %dx%d\n", edid.Width, edid.Height)
			
			if len(edid.SupportedModes) > 0 {
				fmt.Println("Supported Modes:")
				for _, mode := range edid.SupportedModes {
					fmt.Printf("  - %s\n", mode)
				}
			}
		},
	}

	// Power save commands
	powerSaveCmd := &cobra.Command{
		Use:   "power-save",
		Short: "Manage video output power save",
	}

	powerSaveGetCmd := &cobra.Command{
		Use:   "get [connector] [device]",
		Short: "Get power save status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			status, err := client.Video.GetPowerSaveStatus(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			if status.Enabled {
				fmt.Printf("Power save is enabled for %s/%s\n", args[0], args[1])
			} else {
				fmt.Printf("Power save is disabled for %s/%s\n", args[0], args[1])
			}
		},
	}

	powerSaveEnableCmd := &cobra.Command{
		Use:   "enable [connector] [device]",
		Short: "Enable power save",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Video.SetPowerSave(args[0], args[1], true)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Power save enabled for %s/%s\n", args[0], args[1])
		},
	}

	powerSaveDisableCmd := &cobra.Command{
		Use:   "disable [connector] [device]",
		Short: "Disable power save",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Video.SetPowerSave(args[0], args[1], false)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Power save disabled for %s/%s\n", args[0], args[1])
		},
	}

	powerSaveCmd.AddCommand(powerSaveGetCmd, powerSaveEnableCmd, powerSaveDisableCmd)

	// Video modes commands
	modesCmd := &cobra.Command{
		Use:   "modes",
		Short: "Manage video modes",
	}

	modesListCmd := &cobra.Command{
		Use:   "list [connector] [device]",
		Short: "List available video modes",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			modes, err := client.Video.GetAvailableModes(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Available video modes for %s/%s:\n", args[0], args[1])
			for _, mode := range modes {
				interlaced := ""
				if mode.Interlaced {
					interlaced = " (interlaced)"
				}
				preferred := ""
				if mode.PreferredMode {
					preferred = " [preferred]"
				}
				fmt.Printf("  %s: %dx%d @ %dHz%s%s\n", 
					mode.Mode, mode.Width, mode.Height, mode.RefreshRate, interlaced, preferred)
			}
		},
	}

	modesGetCmd := &cobra.Command{
		Use:   "current [connector] [device]",
		Short: "Get current video mode",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			mode, err := client.Video.GetCurrentMode(args[0], args[1])
			if err != nil {
				handleError(err)
			}

			interlaced := ""
			if mode.Interlaced {
				interlaced = " (interlaced)"
			}
			
			fmt.Printf("Current video mode for %s/%s:\n", args[0], args[1])
			fmt.Printf("  Mode: %s\n", mode.Mode)
			fmt.Printf("  Resolution: %dx%d @ %dHz%s\n", 
				mode.Width, mode.Height, mode.RefreshRate, interlaced)
			
			if mode.OverscanMode != "" {
				fmt.Printf("  Overscan Mode: %s\n", mode.OverscanMode)
			}
		},
	}

	modesSetCmd := &cobra.Command{
		Use:   "set [connector] [device] [mode]",
		Short: "Set video mode",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Video.SetVideoMode(args[0], args[1], args[2])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Video mode set to %s for %s/%s\n", args[2], args[0], args[1])
		},
	}

	modesCmd.AddCommand(modesListCmd, modesGetCmd, modesSetCmd)

	// CEC command
	cecCmd := &cobra.Command{
		Use:   "cec [hex-command]",
		Short: "Send CEC command (experimental)",
		Args:  cobra.ExactArgs(1),
		Long: `Send CEC payload out of HDMI-1 port.
The command should be a hex string (e.g., "40 04").

Note: This is an experimental feature.`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Video.SendCEC(args[0])
			if err != nil {
				handleError(err)
			}

			fmt.Printf("CEC command sent: %s\n", args[0])
		},
	}

	videoCmd.AddCommand(outputInfoCmd, edidCmd, powerSaveCmd, modesCmd, cecCmd)
	rootCmd.AddCommand(videoCmd)
}