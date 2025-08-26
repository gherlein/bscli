package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func addDisplayCommands() {
	displayCmd := &cobra.Command{
		Use:   "display",
		Short: "Display control commands (Moka displays, BOS 9.0.189+)",
		Long:  "Commands for controlling Moka displays connected to the BrightSign player",
	}

	// Get all display settings
	getAllCmd := &cobra.Command{
		Use:   "get-all",
		Short: "Get all display settings",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			settings, err := client.Display.GetAll()
			if err != nil {
				handleError(err)
			}

			data, _ := json.MarshalIndent(settings, "", "  ")
			fmt.Println(string(data))
		},
	}

	// Display info
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get display information",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			info, err := client.Display.GetInfo()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Model: %s\n", info.Model)
			fmt.Printf("Serial: %s\n", info.SerialNumber)
			fmt.Printf("Version: %s\n", info.Version)
			fmt.Printf("Resolution: %dx%d\n", info.Width, info.Height)
		},
	}

	// Brightness commands
	brightnessCmd := &cobra.Command{
		Use:   "brightness",
		Short: "Manage display brightness",
	}

	brightnessGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get brightness setting",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			brightness, err := client.Display.GetBrightness()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Brightness: %d (min: %d, max: %d)\n", 
				brightness.Value, brightness.Min, brightness.Max)
		},
	}

	brightnessSetCmd := &cobra.Command{
		Use:   "set [value]",
		Short: "Set brightness value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var value int
			fmt.Sscanf(args[0], "%d", &value)

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Display.SetBrightness(value)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Brightness set to %d\n", value)
		},
	}

	brightnessCmd.AddCommand(brightnessGetCmd, brightnessSetCmd)

	// Contrast commands
	contrastCmd := &cobra.Command{
		Use:   "contrast",
		Short: "Manage display contrast",
	}

	contrastGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get contrast setting",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			contrast, err := client.Display.GetContrast()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Contrast: %d (min: %d, max: %d)\n", 
				contrast.Value, contrast.Min, contrast.Max)
		},
	}

	contrastSetCmd := &cobra.Command{
		Use:   "set [value]",
		Short: "Set contrast value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var value int
			fmt.Sscanf(args[0], "%d", &value)

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Display.SetContrast(value)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Contrast set to %d\n", value)
		},
	}

	contrastCmd.AddCommand(contrastGetCmd, contrastSetCmd)

	// Volume commands
	volumeCmd := &cobra.Command{
		Use:   "volume",
		Short: "Manage display volume",
	}

	volumeGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get volume setting",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			volume, err := client.Display.GetVolume()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Volume: %d (min: %d, max: %d)\n", 
				volume.Value, volume.Min, volume.Max)
		},
	}

	volumeSetCmd := &cobra.Command{
		Use:   "set [value]",
		Short: "Set volume value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var value int
			fmt.Sscanf(args[0], "%d", &value)

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Display.SetVolume(value)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Volume set to %d\n", value)
		},
	}

	volumeCmd.AddCommand(volumeGetCmd, volumeSetCmd)

	// Power commands
	powerCmd := &cobra.Command{
		Use:   "power",
		Short: "Manage display power",
	}

	powerGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get power state",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			power, err := client.Display.GetPowerSettings()
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Power state: %s\n", power.State)
		},
	}

	powerOnCmd := &cobra.Command{
		Use:   "on",
		Short: "Turn display on",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Display.SetPowerSettings("on")
			if err != nil {
				handleError(err)
			}

			fmt.Println("Display turned on")
		},
	}

	powerStandbyCmd := &cobra.Command{
		Use:   "standby",
		Short: "Put display in standby",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Display.SetPowerSettings("standby")
			if err != nil {
				handleError(err)
			}

			fmt.Println("Display in standby mode")
		},
	}

	powerCmd.AddCommand(powerGetCmd, powerOnCmd, powerStandbyCmd)

	// Firmware update
	firmwareUpdateCmd := &cobra.Command{
		Use:   "firmware-update [file-or-url]",
		Short: "Update display firmware",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Update display firmware? This may take several minutes. Continue? (y/N): ")
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

			err = client.Display.UpdateFirmware(args[0])
			if err != nil {
				handleError(err)
			}

			fmt.Println("Firmware update initiated")
		},
	}

	displayCmd.AddCommand(getAllCmd, infoCmd, brightnessCmd, contrastCmd, 
		volumeCmd, powerCmd, firmwareUpdateCmd)
	rootCmd.AddCommand(displayCmd)
}