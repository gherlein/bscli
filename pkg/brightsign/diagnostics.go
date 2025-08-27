package brightsign

import (
	"fmt"
)

// DiagnosticsService handles diagnostic operations
type DiagnosticsService struct {
	client *Client
}

// DiagnosticResult represents a diagnostic test result
type DiagnosticResult struct {
	Test    string `json:"test"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// PingResult represents ping test results
type PingResult struct {
	Success      bool    `json:"success"`
	Address      string  `json:"address"`
	PacketsSent  int     `json:"packetsSent"`
	PacketsRecv  int     `json:"packetsReceived"`
	PacketLoss   float64 `json:"packetLoss"`
	MinTime      float64 `json:"minTime"`
	MaxTime      float64 `json:"maxTime"`
	AvgTime      float64 `json:"avgTime"`
	ErrorMessage string  `json:"error,omitempty"`
}

// DNSLookupResult represents DNS lookup results
type DNSLookupResult struct {
	Success   bool     `json:"success"`
	Hostname  string   `json:"hostname"`
	Addresses []string `json:"addresses"`
	Error     string   `json:"error,omitempty"`
}

// TraceRouteResult represents trace route results
type TraceRouteResult struct {
	Success bool        `json:"success"`
	Target  string      `json:"target"`
	Hops    []TraceHop  `json:"hops"`
	Error   string      `json:"error,omitempty"`
}

// TraceHop represents a single hop in trace route
type TraceHop struct {
	Number   int     `json:"number"`
	Address  string  `json:"address"`
	Hostname string  `json:"hostname,omitempty"`
	RTT      float64 `json:"rtt"`
}

// NetworkConfig represents network interface configuration
type NetworkConfig struct {
	Interface   string   `json:"interface"`
	DHCP        bool     `json:"dhcp"`
	IP          string   `json:"ip,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Gateway     string   `json:"gateway,omitempty"`
	DNS         []string `json:"dns,omitempty"`
	VLANID      int      `json:"vlanId,omitempty"`
}

// PacketCaptureConfig represents packet capture configuration
type PacketCaptureConfig struct {
	Interface    string `json:"interface"`
	Duration     int    `json:"duration"`
	MaxFileSize  int    `json:"maxFileSize,omitempty"`
	Filter       string `json:"filter,omitempty"`
	OutputFile   string `json:"outputFile,omitempty"`
}

// PacketCaptureStatus represents packet capture status
type PacketCaptureStatus struct {
	Running      bool   `json:"running"`
	Interface    string `json:"interface,omitempty"`
	Duration     int    `json:"duration,omitempty"`
	BytesCaptured int64  `json:"bytesCaptured,omitempty"`
	OutputFile   string `json:"outputFile,omitempty"`
}

// TelnetConfig represents telnet configuration
type TelnetConfig struct {
	Enabled    bool `json:"enabled"`
	PortNumber int  `json:"portNumber,omitempty"`
	Reboot     bool `json:"reboot,omitempty"`
}

// SSHConfig represents SSH configuration
type SSHConfig struct {
	Enabled    bool   `json:"enabled"`
	PortNumber int    `json:"portNumber,omitempty"`
	Password   string `json:"password,omitempty"`
	Reboot     bool   `json:"reboot,omitempty"`
}

// RunDiagnostics runs network diagnostics
func (s *DiagnosticsService) RunDiagnostics() (interface{}, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result interface{} `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Result, nil
}

// DNSLookup performs DNS lookup
func (s *DiagnosticsService) DNSLookup(address string, resolveAddress bool) (*DNSLookupResult, error) {
	path := fmt.Sprintf("/diagnostics/dns-lookup/%s", address)
	if resolveAddress {
		path += "?resolveAddress=true"
	}

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result DNSLookupResult `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// Ping performs ping test
func (s *DiagnosticsService) Ping(ipAddress string) (*PingResult, error) {
	path := fmt.Sprintf("/diagnostics/ping/%s", ipAddress)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result PingResult `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// TraceRoute performs trace route
func (s *DiagnosticsService) TraceRoute(address string, resolveAddress bool) (*TraceRouteResult, error) {
	path := fmt.Sprintf("/diagnostics/trace-route/%s", address)
	if resolveAddress {
		path += "?resolveAddress=true"
	}

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result TraceRouteResult `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// GetNetworkNeighborhood retrieves network neighborhood information
func (s *DiagnosticsService) GetNetworkNeighborhood() (map[string]interface{}, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/network-neighborhood/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result map[string]interface{} `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Result, nil
}

// GetNetworkConfiguration gets network configuration for interface
func (s *DiagnosticsService) GetNetworkConfiguration(interfaceName string) (*NetworkConfig, error) {
	path := fmt.Sprintf("/diagnostics/network-configuration/%s/", interfaceName)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result NetworkConfig `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetNetworkConfiguration applies test network configuration
func (s *DiagnosticsService) SetNetworkConfiguration(interfaceName string, config NetworkConfig) error {
	path := fmt.Sprintf("/diagnostics/network-configuration/%s/", interfaceName)

	resp, err := s.client.doRequest("PUT", path, config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set network configuration: status %d", resp.StatusCode)
	}

	return nil
}

// GetInterfaces returns list of applied network interfaces
func (s *DiagnosticsService) GetInterfaces() ([]string, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/interfaces/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result []string `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Result, nil
}

// GetPacketCaptureStatus returns packet capture operation status
func (s *DiagnosticsService) GetPacketCaptureStatus() (*PacketCaptureStatus, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/packet-capture/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result PacketCaptureStatus `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// StartPacketCapture starts packet capture operation
func (s *DiagnosticsService) StartPacketCapture(config PacketCaptureConfig) error {
	resp, err := s.client.doRequest("POST", "/diagnostics/packet-capture/", config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to start packet capture: status %d", resp.StatusCode)
	}

	return nil
}

// StopPacketCapture stops packet capture operation
func (s *DiagnosticsService) StopPacketCapture() error {
	resp, err := s.client.doRequest("DELETE", "/diagnostics/packet-capture/", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to stop packet capture: status %d", resp.StatusCode)
	}

	return nil
}

// GetTelnetConfig returns telnet configuration
func (s *DiagnosticsService) GetTelnetConfig() (*TelnetConfig, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/telnet/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result TelnetConfig `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetTelnetConfig configures telnet settings
func (s *DiagnosticsService) SetTelnetConfig(config TelnetConfig) error {
	resp, err := s.client.doRequest("PUT", "/diagnostics/telnet/", config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set telnet configuration: status %d", resp.StatusCode)
	}

	return nil
}

// GetSSHConfig returns SSH configuration
func (s *DiagnosticsService) GetSSHConfig() (*SSHConfig, error) {
	resp, err := s.client.doRequest("GET", "/diagnostics/ssh/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result SSHConfig `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetSSHConfig configures SSH settings
func (s *DiagnosticsService) SetSSHConfig(config SSHConfig) error {
	resp, err := s.client.doRequest("PUT", "/diagnostics/ssh/", config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set SSH configuration: status %d", resp.StatusCode)
	}

	return nil
}