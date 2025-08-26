package main

import (
	"fmt"
	"log"
	"time"

	"bscli/pkg/brightsign"
)

func main() {
	// Create a new BrightSign client
	client := brightsign.NewClient(brightsign.Config{
		Host:     "192.168.1.100", // Replace with your player's IP
		Username: "admin",
		Password: "your_password", // Replace with actual password
		Debug:    true,           // Enable debug output
		Timeout:  30 * time.Second,
	})

	// Example 1: Get device information
	fmt.Println("=== Device Information ===")
	info, err := client.Info.GetInfo()
	if err != nil {
		log.Printf("Failed to get device info: %v", err)
	} else {
		fmt.Printf("Model: %s\n", info.Model)
		fmt.Printf("Serial: %s\n", info.Serial)
		fmt.Printf("Firmware: %s\n", info.FWVersion)
		fmt.Printf("Uptime: %s\n", info.Uptime)
	}

	// Example 2: Check player health
	fmt.Println("\n=== Player Health ===")
	health, err := client.Info.GetHealth()
	if err != nil {
		log.Printf("Failed to get health: %v", err)
	} else {
		fmt.Printf("Status: %s\n", health.Status)
		fmt.Printf("Status Time: %s\n", health.StatusTime)
	}

	// Example 3: List files on SD card
	fmt.Println("\n=== Files on SD Card ===")
	files, err := client.Storage.ListFiles("/storage/sd/", nil)
	if err != nil {
		log.Printf("Failed to list files: %v", err)
	} else {
		if len(files) == 0 {
			fmt.Println("No files found")
		} else {
			for _, file := range files {
				fmt.Printf("%s: %s (%d bytes)\n", file.Type, file.Name, file.Size)
			}
		}
	}

	// Example 4: Run a ping test
	fmt.Println("\n=== Network Diagnostics ===")
	pingResult, err := client.Diagnostics.Ping("8.8.8.8")
	if err != nil {
		log.Printf("Failed to ping: %v", err)
	} else {
		if pingResult.Success {
			fmt.Printf("Ping to %s: %d/%d packets, %.2fms avg\n",
				pingResult.Address, pingResult.PacketsRecv, 
				pingResult.PacketsSent, pingResult.AvgTime)
		} else {
			fmt.Printf("Ping failed: %s\n", pingResult.ErrorMessage)
		}
	}

	// Example 5: Get registry value
	fmt.Println("\n=== Registry Example ===")
	hostname, err := client.Registry.GetValue("networking", "hostname")
	if err != nil {
		log.Printf("Failed to get hostname: %v", err)
	} else {
		fmt.Printf("Player hostname: %s\n", hostname)
	}

	// Example 6: Take a snapshot (optional - uncomment to test)
	/*
	fmt.Println("\n=== Taking Snapshot ===")
	filename, err := client.Control.TakeSnapshot(&brightsign.SnapshotOptions{
		Width:  1920,
		Height: 1080,
	})
	if err != nil {
		log.Printf("Failed to take snapshot: %v", err)
	} else {
		fmt.Printf("Snapshot saved: %s\n", filename)
	}
	*/

	fmt.Println("\n=== Done ===")
}