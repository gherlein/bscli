package brightsign

import (
	"fmt"
)

// ControlService handles player control endpoints
type ControlService struct {
	client *Client
}

// RebootOptions contains options for rebooting the player
type RebootOptions struct {
	CrashReport    bool `json:"crash_report,omitempty"`
	FactoryReset   bool `json:"factory_reset,omitempty"`
	DisableAutorun bool `json:"disable_autorun,omitempty"`
}

// DWSPassword represents DWS password configuration
type DWSPassword struct {
	Password string `json:"password,omitempty"`
	Reset    bool   `json:"reset,omitempty"`
}

// DWSPasswordInfo represents DWS password information
type DWSPasswordInfo struct {
	IsSet bool `json:"isSet"`
}

// LocalDWSConfig represents local DWS configuration
type LocalDWSConfig struct {
	Enabled bool `json:"enabled"`
}

// SnapshotOptions contains options for taking a snapshot
type SnapshotOptions struct {
	Width                      int  `json:"width,omitempty"`
	Height                     int  `json:"height,omitempty"`
	ShouldCaptureFullResolution bool `json:"shouldCaptureFullResolution,omitempty"`
}

// Reboot reboots the player with optional parameters
func (s *ControlService) Reboot(options *RebootOptions) error {
	if options == nil {
		options = &RebootOptions{}
	}

	resp, err := s.client.doRequest("PUT", "/control/reboot/", options)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to reboot: status %d", resp.StatusCode)
	}

	return nil
}

// GetDWSPassword retrieves DWS password information (not the actual password)
func (s *ControlService) GetDWSPassword() (*DWSPasswordInfo, error) {
	resp, err := s.client.doRequest("GET", "/control/dws-password/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result DWSPasswordInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetDWSPassword sets or resets the DWS password
func (s *ControlService) SetDWSPassword(config DWSPassword) error {
	resp, err := s.client.doRequest("PUT", "/control/dws-password/", config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set DWS password: status %d", resp.StatusCode)
	}

	return nil
}

// GetLocalDWS retrieves local DWS status
func (s *ControlService) GetLocalDWS() (*LocalDWSConfig, error) {
	resp, err := s.client.doRequest("GET", "/control/local-dws/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result LocalDWSConfig `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetLocalDWS enables or disables local DWS
func (s *ControlService) SetLocalDWS(enabled bool) error {
	config := LocalDWSConfig{Enabled: enabled}
	resp, err := s.client.doRequest("PUT", "/control/local-dws/", config)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set local DWS: status %d", resp.StatusCode)
	}

	return nil
}

// TakeSnapshot captures a snapshot of the currently playing content
func (s *ControlService) TakeSnapshot(options *SnapshotOptions) (string, error) {
	if options == nil {
		options = &SnapshotOptions{}
	}

	resp, err := s.client.doRequest("POST", "/snapshot/", options)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Result string `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return "", err
	}

	return result.Data.Result, nil
}

// DownloadFirmware downloads OS from remote URL and reboots player
func (s *ControlService) DownloadFirmware(url string) error {
	path := fmt.Sprintf("/download-firmware/?url=%s", url)
	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to download firmware: status %d", resp.StatusCode)
	}

	return nil
}