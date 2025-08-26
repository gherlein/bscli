package cli

import (
	"encoding/json"
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

			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
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

			fmt.Printf("Status: %s\n", health.Status)
			fmt.Printf("Status Time: %s\n", health.StatusTime)
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

			time, err := client.Info.GetTime()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Date: %s\n", time.Date)
			fmt.Printf("Time: %s\n", time.Time)
			if time.Timezone != "" {
				fmt.Printf("Timezone: %s\n", time.Timezone)
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

			fmt.Println("Time set successfully")
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

			fmt.Printf("Resolution: %s\n", mode.Resolution)
			fmt.Printf("Frame Rate: %d\n", mode.FrameRate)
			fmt.Printf("Scan Method: %s\n", mode.ScanMethod)
			fmt.Printf("Preferred Mode: %v\n", mode.PreferredMode)
			if mode.OverscanMode != "" {
				fmt.Printf("Overscan Mode: %s\n", mode.OverscanMode)
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

			fmt.Println("Available APIs:")
			for _, api := range apis {
				fmt.Printf("  - %s\n", api)
			}
		},
	}

	infoCmd.AddCommand(deviceInfoCmd, healthCmd, timeCmd, setTimeCmd, videoModeCmd, listAPIsCmd)
	rootCmd.AddCommand(infoCmd)
}