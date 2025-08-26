package cli

import (
	"fmt"

	"bscli/pkg/brightsign"
	"github.com/spf13/cobra"
)

func addInfoCommands() {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get player information",
		Long:  "Commands for retrieving various player information",
	}

	// Device info command
	deviceInfoCmd := &cobra.Command{
		Use:   "device",
		Short: "Get device information",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			info, err := client.Info.GetInfo()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(info)
			} else {
				fmt.Printf("Model: %s\n", info.Model)
				fmt.Printf("Serial: %s\n", info.Serial)
				fmt.Printf("Family: %s\n", info.Family)
				fmt.Printf("Boot Version: %s\n", info.BootVersion)
				fmt.Printf("Firmware Version: %s\n", info.FWVersion)
				fmt.Printf("Uptime: %s (%d seconds)\n", info.Uptime, info.UptimeSeconds)
				
				if len(info.Network.Interfaces) > 0 {
					fmt.Printf("\nNetwork Interfaces:\n")
					for _, iface := range info.Network.Interfaces {
						fmt.Printf("  %s (%s): %s\n", iface.Name, iface.Type, iface.IP)
					}
				}
				
				if info.Network.Hostname != "" {
					fmt.Printf("Hostname: %s\n", info.Network.Hostname)
				}
				
				if len(info.Extensions) > 0 {
					fmt.Printf("\nExtensions:\n")
					for key, value := range info.Extensions {
						fmt.Printf("  %s: %s\n", key, value)
					}
				}
			}
		},
	}

	// Health command
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Get player health status",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			health, err := client.Info.GetHealth()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(health)
			} else {
				fmt.Printf("Status: %s\n", health.Status)
				fmt.Printf("Status Time: %s\n", health.StatusTime.Format("2006-01-02 15:04:05"))
			}
		},
	}

	// Time command
	timeCmd := &cobra.Command{
		Use:   "time",
		Short: "Get current time configuration",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			timeInfo, err := client.Info.GetTime()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(timeInfo)
			} else {
				fmt.Printf("Date: %s\n", timeInfo.Date)
				fmt.Printf("Time: %s\n", timeInfo.Time)
				if timeInfo.Timezone != "" {
					fmt.Printf("Timezone: %s\n", timeInfo.Timezone)
				}
			}
		},
	}

	// Set time command
	setTimeCmd := &cobra.Command{
		Use:   "set-time [date] [time]",
		Short: "Set player time",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			timezone, _ := cmd.Flags().GetString("timezone")
			
			err = client.Info.SetTime(brightsign.TimeInfo{
				Date:     args[0],
				Time:     args[1],
				Timezone: timezone,
			})
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(map[string]bool{"success": true})
			} else {
				fmt.Println("Time set successfully")
			}
		},
	}
	setTimeCmd.Flags().String("timezone", "", "Timezone to apply")

	// Video mode command
	videoModeCmd := &cobra.Command{
		Use:   "video-mode",
		Short: "Get current video mode",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			mode, err := client.Info.GetVideoMode()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(mode)
			} else {
				fmt.Printf("Resolution: %s\n", mode.Resolution)
				fmt.Printf("Frame Rate: %d Hz\n", mode.FrameRate)
				fmt.Printf("Scan Method: %s\n", mode.ScanMethod)
				fmt.Printf("Preferred Mode: %v\n", mode.PreferredMode)
				if mode.OverscanMode != "" {
					fmt.Printf("Overscan Mode: %s\n", mode.OverscanMode)
				}
			}
		},
	}

	// List APIs command
	listAPIsCmd := &cobra.Command{
		Use:   "apis",
		Short: "List all available APIs",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			apis, err := client.Info.ListAPIs()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(apis)
			} else {
				fmt.Println("Available APIs:")
				for _, api := range apis {
					fmt.Printf("  - %s\n", api)
				}
			}
		},
	}

	infoCmd.AddCommand(deviceInfoCmd, healthCmd, timeCmd, setTimeCmd, videoModeCmd, listAPIsCmd)
	rootCmd.AddCommand(infoCmd)
}