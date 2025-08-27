package brightsign

import (
	"fmt"
	"io"
	"os"
)

// InfoService handles player information endpoints
type InfoService struct {
	client *Client
}

// DeviceInfo represents basic device information
type DeviceInfo struct {
	Model           string            `json:"model"`
	Serial          string            `json:"serial"`
	Family          string            `json:"family"`
	BootVersion     string            `json:"bootVersion"`
	FWVersion       string            `json:"fwVersion"`
	Network         NetworkInfo       `json:"network"`
	Uptime          string            `json:"uptime"`
	UptimeSeconds   int64             `json:"uptimeSeconds"`
	Extensions      interface{} `json:"extensions"`
}

// NetworkInfo represents network information
type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	Hostname   string             `json:"hostname"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Proto     string `json:"proto"`
	IP        string `json:"ip"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	DNS       string `json:"dns"`
	MAC       string `json:"mac"`
	Metric    int    `json:"metric"`
}

// HealthInfo represents player health status
type HealthInfo struct {
	Status     string `json:"status"`
	StatusTime string `json:"statusTime"` // Changed to string to handle various date formats
}

// TimeInfo represents time configuration
type TimeInfo struct {
	Date     interface{} `json:"date"` // Can be string or number
	Time     string      `json:"time"`
	Timezone string      `json:"timezone,omitempty"`
}

// VideoMode represents video output mode
type VideoMode struct {
	Resolution       string `json:"resolution"`
	FrameRate        int    `json:"frameRate"`
	ScanMethod       string `json:"scanMethod"`
	PreferredMode    bool   `json:"preferredMode"`
	OverscanMode     string `json:"overscanMode,omitempty"`
}

// GetInfo retrieves basic player information
func (s *InfoService) GetInfo() (*DeviceInfo, error) {
	resp, err := s.client.doRequest("GET", "/info/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result DeviceInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		// For better debugging, show what we got
		if s.client.debug {
			resp.Body.Close()
			// Re-read the response body for debugging
			resp2, _ := s.client.doRequest("GET", "/info/", nil)
			if resp2 != nil {
				body, _ := io.ReadAll(resp2.Body)
				fmt.Fprintf(os.Stderr, "DEBUG: Failed to parse GetInfo response: %s\n", string(body))
				resp2.Body.Close()
			}
		}
		return nil, fmt.Errorf("failed to parse device info response: %w", err)
	}

	return &result.Data.Result, nil
}

// GetHealth retrieves player health status
func (s *InfoService) GetHealth() (*HealthInfo, error) {
	resp, err := s.client.doRequest("GET", "/health/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result HealthInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// GetTime retrieves current time configuration
func (s *InfoService) GetTime() (*TimeInfo, error) {
	resp, err := s.client.doRequest("GET", "/time/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result TimeInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetTime sets the time on the player
func (s *InfoService) SetTime(info TimeInfo) error {
	resp, err := s.client.doRequest("PUT", "/time/", info)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set time: status %d", resp.StatusCode)
	}

	return nil
}

// GetVideoMode retrieves current video mode
func (s *InfoService) GetVideoMode() (*VideoMode, error) {
	resp, err := s.client.doRequest("GET", "/video-mode/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result VideoMode `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// ListAPIs returns list of all available APIs
func (s *InfoService) ListAPIs() (interface{}, error) {
	resp, err := s.client.doRequest("GET", "/", nil)
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