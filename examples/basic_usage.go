package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"bscli/pkg/brightsign"
)

func main() {
	// Define command-line flags
	var (
		localFlag = flag.Bool("local", false, "Accept locally signed certificates (use HTTPS with insecure TLS)")
		lFlag     = flag.Bool("l", false, "Accept locally signed certificates (short form)")
		debugFlag = flag.Bool("debug", false, "Enable debug output")
		dFlag     = flag.Bool("d", false, "Enable debug output (short form)")
	)
	
	// Parse command-line flags
	flag.Parse()
	
	// Get configuration from environment variables (same as integration tests)
	host := os.Getenv("BSCLI_TEST_HOST")
	if host == "" {
		log.Fatal("Error: BSCLI_TEST_HOST environment variable is required\nExample: export BSCLI_TEST_HOST=192.168.1.100")
	}

	password := os.Getenv("BSCLI_TEST_PASSWORD")
	if password == "" {
		log.Fatal("Error: BSCLI_TEST_PASSWORD environment variable is required\nExample: export BSCLI_TEST_PASSWORD=yourpassword")
	}

	username := os.Getenv("BSCLI_TEST_USERNAME")
	if username == "" {
		username = "admin"
	}

	// Check both environment variable and flags for debug and insecure settings
	debug := os.Getenv("BSCLI_TEST_DEBUG") == "true" || *debugFlag || *dFlag
	insecure := os.Getenv("BSCLI_TEST_INSECURE") == "true" || *localFlag || *lFlag

	// Create a new BrightSign client with environment configuration
	protocol := "HTTP"
	if insecure {
		protocol = "HTTPS (insecure)"
	}
	fmt.Printf("Connecting to BrightSign player at %s as user '%s' using %s...\n\n", host, username, protocol)
	
	client := brightsign.NewClient(brightsign.Config{
		Host:     host,
		Username: username,
		Password: password,
		Debug:    debug,
		Insecure: insecure,
		Timeout:  30 * time.Second,
	})

	// Example 1: Get device information
	fmt.Println("=== Device Information ===")
	info, err := client.Info.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get device info: %v", err)
	}
	fmt.Printf("Model: %s\n", info.Model)
	fmt.Printf("Serial: %s\n", info.Serial)
	fmt.Printf("Family: %s\n", info.Family)
	fmt.Printf("Firmware: %s\n", info.FWVersion)
	fmt.Printf("Boot Version: %s\n", info.BootVersion)
	fmt.Printf("Uptime: %s (%d seconds)\n", info.Uptime, info.UptimeSeconds)
	
	// Show network interfaces
	if info.Network.Interfaces != nil {
		fmt.Println("\nNetwork Interfaces:")
		for _, iface := range info.Network.Interfaces {
			fmt.Printf("  %s (%s/%s): %s\n", iface.Name, iface.Type, iface.Proto, iface.IP)
		}
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

	// Example 3: Get current time configuration
	fmt.Println("\n=== Time Configuration ===")
	timeInfo, err := client.Info.GetTime()
	if err != nil {
		log.Printf("Failed to get time: %v", err)
	} else {
		fmt.Printf("Date: %v\n", timeInfo.Date)
		fmt.Printf("Time: %s\n", timeInfo.Time)
		if timeInfo.Timezone != "" {
			fmt.Printf("Timezone: %s\n", timeInfo.Timezone)
		}
	}

	// Example 4: List files on SD card (if command arg provided)
	args := flag.Args()
	if len(args) > 0 && args[0] == "list-files" {
		fmt.Println("\n=== Files on SD Card ===")
		files, err := client.Storage.ListFiles("/storage/sd/", nil)
		if err != nil {
			log.Printf("Failed to list files: %v", err)
		} else {
			if len(files) == 0 {
				fmt.Println("No files found")
			} else {
				for _, file := range files {
					if file.Type == "file" {
						fmt.Printf("  📄 %s (%d bytes)\n", file.Name, file.Size)
					} else {
						fmt.Printf("  📁 %s/\n", file.Name)
					}
				}
			}
		}
	}

	// Example 5: Run diagnostics (if command arg provided)
	if len(args) > 0 && args[0] == "diagnostics" {
		fmt.Println("\n=== Running Network Diagnostics ===")
		
		// Ping test
		pingResult, err := client.Diagnostics.Ping("8.8.8.8")
		if err != nil {
			log.Printf("Failed to ping: %v", err)
		} else {
			if pingResult.Success {
				fmt.Printf("✓ Ping to %s: %d/%d packets, %.2fms avg\n",
					pingResult.Address, pingResult.PacketsRecv,
					pingResult.PacketsSent, pingResult.AvgTime)
			} else {
				fmt.Printf("✗ Ping failed: %s\n", pingResult.ErrorMessage)
			}
		}

		// DNS lookup
		dnsResult, err := client.Diagnostics.DNSLookup("google.com", false)
		if err != nil {
			log.Printf("Failed DNS lookup: %v", err)
		} else {
			if dnsResult.Success {
				fmt.Printf("✓ DNS lookup for %s:\n", dnsResult.Hostname)
				for _, addr := range dnsResult.Addresses {
					fmt.Printf("    %s\n", addr)
				}
			} else {
				fmt.Printf("✗ DNS lookup failed: %s\n", dnsResult.Error)
			}
		}

		// List network interfaces
		interfaces, err := client.Diagnostics.GetInterfaces()
		if err != nil {
			log.Printf("Failed to get interfaces: %v", err)
		} else {
			fmt.Println("✓ Network interfaces:")
			for _, iface := range interfaces {
				fmt.Printf("    %s\n", iface)
			}
		}
	}

	// Example 6: Registry operations (if command arg provided)
	if len(args) > 0 && args[0] == "registry" {
		fmt.Println("\n=== Registry Operations ===")
		
		// Try to get hostname from registry
		hostname, err := client.Registry.GetValue("networking", "hostname")
		if err != nil {
			// Try alternative section
			hostname, err = client.Registry.GetValue("system", "hostname")
			if err != nil {
				log.Printf("Failed to get hostname from registry: %v", err)
			} else {
				fmt.Printf("System hostname: %s\n", hostname)
			}
		} else {
			fmt.Printf("Network hostname: %s\n", hostname)
		}

		// Set a test value (safe, non-critical)
		testKey := "bscli_example_test"
		testValue := fmt.Sprintf("test_%d", time.Now().Unix())
		err = client.Registry.SetValue("networking", testKey, testValue)
		if err != nil {
			log.Printf("Failed to set test value: %v", err)
		} else {
			fmt.Printf("Set test value: networking/%s = %s\n", testKey, testValue)
			
			// Read it back
			readValue, err := client.Registry.GetValue("networking", testKey)
			if err != nil {
				log.Printf("Failed to read back test value: %v", err)
			} else {
				fmt.Printf("Read back value: %s\n", readValue)
			}
			
			// Clean up - delete test value
			err = client.Registry.DeleteValue("networking", testKey)
			if err != nil {
				log.Printf("Failed to delete test value: %v", err)
			} else {
				fmt.Printf("Cleaned up test value\n")
			}
		}
	}

	// Example 7: Video output information (if command arg provided)
	if len(args) > 0 && args[0] == "video" {
		fmt.Println("\n=== Video Output Information ===")
		
		// Try common video output combinations
		connectors := []string{"hdmi", "HDMI"}
		devices := []string{"0", "1"}
		
		found := false
		for _, connector := range connectors {
			for _, device := range devices {
				info, err := client.Video.GetOutputInfo(connector, device)
				if err == nil {
					found = true
					fmt.Printf("Output: %s/%s\n", info.Connector, info.Device)
					fmt.Printf("Connected: %v\n", info.Connected)
					if info.Connected {
						fmt.Printf("Resolution: %dx%d @ %dHz\n", info.Width, info.Height, info.RefreshRate)
						if info.PreferredMode != "" {
							fmt.Printf("Preferred Mode: %s\n", info.PreferredMode)
						}
					}
					break
				}
			}
			if found {
				break
			}
		}
		
		if !found {
			fmt.Println("No video outputs found (this may be normal for some configurations)")
		}
	}

	// Show available commands if none provided
	if len(args) == 0 {
		fmt.Println("\n=== Available Commands ===")
		fmt.Println("Run with additional arguments to see more examples:")
		fmt.Println("  ./examples/basic_usage list-files   # List files on SD card")
		fmt.Println("  ./examples/basic_usage diagnostics  # Run network diagnostics")
		fmt.Println("  ./examples/basic_usage registry     # Test registry operations")
		fmt.Println("  ./examples/basic_usage video        # Get video output info")
		fmt.Println("\nFlags:")
		fmt.Println("  -l, --local    Accept locally signed certificates (HTTPS with insecure TLS)")
		fmt.Println("  -d, --debug    Enable debug output")
		fmt.Println("\nEnvironment variables:")
		fmt.Println("  BSCLI_TEST_DEBUG=true     Enable debug output")
		fmt.Println("  BSCLI_TEST_INSECURE=true  Accept locally signed certificates")
	}

	fmt.Println("\n=== Done ===")
}