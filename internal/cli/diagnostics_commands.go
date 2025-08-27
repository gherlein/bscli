package cli

import (
	"encoding/json"
	"fmt"

	"bscli/pkg/brightsign"
	"github.com/spf13/cobra"
)

func addDiagnosticsCommands() {
	diagCmd := &cobra.Command{
		Use:   "diagnostics",
		Aliases: []string{"diag"},
		Short: "Network and system diagnostics",
		Long:  "Commands for running network and system diagnostics",
	}

	// Run diagnostics command
	runDiagCmd := &cobra.Command{
		Use:   "run",
		Short: "Run all network diagnostics",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			results, err := client.Diagnostics.RunDiagnostics()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(results)
				return
			}

			fmt.Println("Diagnostic Results:")
			// Handle different possible result formats
			switch resultData := results.(type) {
			case []interface{}:
				for _, r := range resultData {
					if resultMap, ok := r.(map[string]interface{}); ok {
						status := "✓"
						if statusVal, ok := resultMap["status"]; ok && statusVal != "pass" {
							status = "✗"
						}
						test := ""
						if testVal, ok := resultMap["test"]; ok {
							test = fmt.Sprintf("%v", testVal)
						}
						message := ""
						if msgVal, ok := resultMap["message"]; ok {
							message = fmt.Sprintf("%v", msgVal)
						}
						fmt.Printf("%s %s: %s\n", status, test, message)
					}
				}
			case map[string]interface{}:
				// Single diagnostic result object
				status := "✓"
				if statusVal, ok := resultData["status"]; ok && statusVal != "pass" {
					status = "✗"
				}
				test := ""
				if testVal, ok := resultData["test"]; ok {
					test = fmt.Sprintf("%v", testVal)
				}
				message := ""
				if msgVal, ok := resultData["message"]; ok {
					message = fmt.Sprintf("%v", msgVal)
				}
				fmt.Printf("%s %s: %s\n", status, test, message)
			default:
				fmt.Printf("%v\n", results)
			}
		},
	}

	// Ping command
	pingCmd := &cobra.Command{
		Use:   "ping [ip-address]",
		Short: "Ping an IP address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			result, err := client.Diagnostics.Ping(args[0])
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(result)
				return
			}

			if result.Success {
				fmt.Printf("PING %s: %d/%d packets received\n", result.Address, result.PacketsRecv, result.PacketsSent)
				fmt.Printf("Packet Loss: %.1f%%\n", result.PacketLoss)
				fmt.Printf("RTT min/avg/max = %.2f/%.2f/%.2f ms\n", result.MinTime, result.AvgTime, result.MaxTime)
			} else {
				fmt.Printf("PING %s failed: %s\n", result.Address, result.ErrorMessage)
			}
		},
	}

	// DNS lookup command
	dnsCmd := &cobra.Command{
		Use:   "dns-lookup [hostname]",
		Short: "Perform DNS lookup",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resolveAddr, _ := cmd.Flags().GetBool("resolve")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			result, err := client.Diagnostics.DNSLookup(args[0], resolveAddr)
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(result)
				return
			}

			if result.Success {
				fmt.Printf("DNS lookup for %s:\n", result.Hostname)
				for _, addr := range result.Addresses {
					fmt.Printf("  %s\n", addr)
				}
			} else {
				fmt.Printf("DNS lookup failed: %s\n", result.Error)
			}
		},
	}
	dnsCmd.Flags().Bool("resolve", false, "Resolve addresses")

	// Traceroute command
	tracerouteCmd := &cobra.Command{
		Use:   "traceroute [address]",
		Short: "Run traceroute to address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resolveAddr, _ := cmd.Flags().GetBool("resolve")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			result, err := client.Diagnostics.TraceRoute(args[0], resolveAddr)
			if err != nil {
				handleError(err)
			}

			if result.Success {
				fmt.Printf("Traceroute to %s:\n", result.Target)
				for _, hop := range result.Hops {
					if hop.Hostname != "" {
						fmt.Printf("%2d  %s (%s)  %.2f ms\n", hop.Number, hop.Hostname, hop.Address, hop.RTT)
					} else {
						fmt.Printf("%2d  %s  %.2f ms\n", hop.Number, hop.Address, hop.RTT)
					}
				}
			} else {
				fmt.Printf("Traceroute failed: %s\n", result.Error)
			}
		},
	}
	tracerouteCmd.Flags().Bool("resolve", false, "Resolve addresses")

	// Network interfaces command
	interfacesCmd := &cobra.Command{
		Use:   "interfaces",
		Short: "List network interfaces",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			interfaces, err := client.Diagnostics.GetInterfaces()
			if err != nil {
				handleError(err)
			}

			if jsonOutput {
				outputJSON(interfaces)
				return
			}

			fmt.Println("Network interfaces:")
			for _, iface := range interfaces {
				fmt.Printf("  - %s\n", iface)
			}
		},
	}

	// Network configuration command
	netConfigCmd := &cobra.Command{
		Use:   "network-config [interface]",
		Short: "Get network configuration for interface",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config, err := client.Diagnostics.GetNetworkConfiguration(args[0])
			if err != nil {
				handleError(err)
			}

			data, _ := json.MarshalIndent(config, "", "  ")
			fmt.Println(string(data))
		},
	}

	// Packet capture commands
	pcapCmd := &cobra.Command{
		Use:   "packet-capture",
		Aliases: []string{"pcap"},
		Short: "Packet capture operations",
	}

	pcapStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get packet capture status",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			status, err := client.Diagnostics.GetPacketCaptureStatus()
			if err != nil {
				handleError(err)
			}

			if status.Running {
				fmt.Println("Packet capture is running")
				fmt.Printf("Interface: %s\n", status.Interface)
				fmt.Printf("Duration: %d seconds\n", status.Duration)
				fmt.Printf("Bytes captured: %d\n", status.BytesCaptured)
				if status.OutputFile != "" {
					fmt.Printf("Output file: %s\n", status.OutputFile)
				}
			} else {
				fmt.Println("Packet capture is not running")
			}
		},
	}

	pcapStartCmd := &cobra.Command{
		Use:   "start [interface]",
		Short: "Start packet capture",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			duration, _ := cmd.Flags().GetInt("duration")
			maxSize, _ := cmd.Flags().GetInt("max-size")
			filter, _ := cmd.Flags().GetString("filter")
			output, _ := cmd.Flags().GetString("output")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.PacketCaptureConfig{
				Interface:   args[0],
				Duration:    duration,
				MaxFileSize: maxSize,
				Filter:      filter,
				OutputFile:  output,
			}

			err = client.Diagnostics.StartPacketCapture(config)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Packet capture started")
		},
	}
	pcapStartCmd.Flags().Int("duration", 60, "Capture duration in seconds")
	pcapStartCmd.Flags().Int("max-size", 0, "Maximum file size in bytes")
	pcapStartCmd.Flags().String("filter", "", "Capture filter expression")
	pcapStartCmd.Flags().String("output", "", "Output file path")

	pcapStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop packet capture",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			err = client.Diagnostics.StopPacketCapture()
			if err != nil {
				handleError(err)
			}

			fmt.Println("Packet capture stopped")
		},
	}

	pcapCmd.AddCommand(pcapStatusCmd, pcapStartCmd, pcapStopCmd)

	// Telnet configuration
	telnetCmd := &cobra.Command{
		Use:   "telnet",
		Short: "Manage telnet configuration",
	}

	telnetStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get telnet configuration",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config, err := client.Diagnostics.GetTelnetConfig()
			if err != nil {
				handleError(err)
			}

			if config.Enabled {
				fmt.Printf("Telnet is enabled on port %d\n", config.PortNumber)
			} else {
				fmt.Println("Telnet is disabled")
			}
		},
	}

	telnetEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable telnet",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			reboot, _ := cmd.Flags().GetBool("reboot")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.TelnetConfig{
				Enabled:    true,
				PortNumber: port,
				Reboot:     reboot,
			}

			err = client.Diagnostics.SetTelnetConfig(config)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Telnet enabled")
			if reboot {
				fmt.Println("Player will reboot")
			}
		},
	}
	telnetEnableCmd.Flags().Int("port", 23, "Telnet port number")
	telnetEnableCmd.Flags().Bool("reboot", false, "Reboot after change")

	telnetDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable telnet",
		Run: func(cmd *cobra.Command, args []string) {
			reboot, _ := cmd.Flags().GetBool("reboot")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.TelnetConfig{
				Enabled: false,
				Reboot:  reboot,
			}

			err = client.Diagnostics.SetTelnetConfig(config)
			if err != nil {
				handleError(err)
			}

			fmt.Println("Telnet disabled")
			if reboot {
				fmt.Println("Player will reboot")
			}
		},
	}
	telnetDisableCmd.Flags().Bool("reboot", false, "Reboot after change")

	telnetCmd.AddCommand(telnetStatusCmd, telnetEnableCmd, telnetDisableCmd)

	// SSH configuration
	sshCmd := &cobra.Command{
		Use:   "ssh",
		Short: "Manage SSH configuration",
	}

	sshStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get SSH configuration",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config, err := client.Diagnostics.GetSSHConfig()
			if err != nil {
				handleError(err)
			}

			if config.Enabled {
				fmt.Printf("SSH is enabled on port %d\n", config.PortNumber)
			} else {
				fmt.Println("SSH is disabled")
			}
		},
	}

	sshEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable SSH",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			sshPassword, _ := cmd.Flags().GetString("ssh-password")
			reboot, _ := cmd.Flags().GetBool("reboot")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.SSHConfig{
				Enabled:    true,
				PortNumber: port,
				Password:   sshPassword,
				Reboot:     reboot,
			}

			err = client.Diagnostics.SetSSHConfig(config)
			if err != nil {
				handleError(err)
			}

			fmt.Println("SSH enabled")
			if reboot {
				fmt.Println("Player will reboot")
			}
		},
	}
	sshEnableCmd.Flags().Int("port", 22, "SSH port number")
	sshEnableCmd.Flags().String("ssh-password", "", "SSH password")
	sshEnableCmd.Flags().Bool("reboot", false, "Reboot after change")

	sshDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable SSH",
		Run: func(cmd *cobra.Command, args []string) {
			reboot, _ := cmd.Flags().GetBool("reboot")

			client, err := getClient()
			if err != nil {
				handleError(err)
			}

			config := brightsign.SSHConfig{
				Enabled: false,
				Reboot:  reboot,
			}

			err = client.Diagnostics.SetSSHConfig(config)
			if err != nil {
				handleError(err)
			}

			fmt.Println("SSH disabled")
			if reboot {
				fmt.Println("Player will reboot")
			}
		},
	}
	sshDisableCmd.Flags().Bool("reboot", false, "Reboot after change")

	sshCmd.AddCommand(sshStatusCmd, sshEnableCmd, sshDisableCmd)

	diagCmd.AddCommand(runDiagCmd, pingCmd, dnsCmd, tracerouteCmd, interfacesCmd, 
		netConfigCmd, pcapCmd, telnetCmd, sshCmd)
	rootCmd.AddCommand(diagCmd)
}