package brightsign

import (
	"fmt"
)

// VideoService handles video output management
type VideoService struct {
	client *Client
}

// VideoOutputInfo represents video output information
type VideoOutputInfo struct {
	Connector    string `json:"connector"`
	Device       string `json:"device"`
	Connected    bool   `json:"connected"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	RefreshRate  int    `json:"refreshRate"`
	InterlaceMode string `json:"interlaceMode,omitempty"`
	PreferredMode string `json:"preferredMode,omitempty"`
}

// EDIDInfo represents EDID information from connected display
type EDIDInfo struct {
	Manufacturer  string   `json:"manufacturer"`
	Product       string   `json:"product"`
	SerialNumber  string   `json:"serialNumber"`
	WeekOfManufacture int    `json:"weekOfManufacture"`
	YearOfManufacture int    `json:"yearOfManufacture"`
	Version       string   `json:"version"`
	Digital       bool     `json:"digital"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	SupportedModes []string `json:"supportedModes"`
}

// PowerSaveStatus represents power save status
type PowerSaveStatus struct {
	Enabled bool `json:"enabled"`
}

// VideoModeInfo represents a video mode
type VideoModeInfo struct {
	Mode          string `json:"mode"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	RefreshRate   int    `json:"refreshRate"`
	Interlaced    bool   `json:"interlaced"`
	PreferredMode bool   `json:"preferredMode,omitempty"`
	OverscanMode  string `json:"overscanMode,omitempty"`
}

// GetOutputInfo retrieves video output information
func (s *VideoService) GetOutputInfo(connector, device string) (*VideoOutputInfo, error) {
	path := fmt.Sprintf("/video/%s/output/%s/", connector, device)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result VideoOutputInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// GetEDID gets EDID information from connected display
func (s *VideoService) GetEDID(connector, device string) (*EDIDInfo, error) {
	path := fmt.Sprintf("/video/%s/output/%s/edid/", connector, device)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result EDIDInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// GetPowerSaveStatus returns power save status
func (s *VideoService) GetPowerSaveStatus(connector, device string) (*PowerSaveStatus, error) {
	path := fmt.Sprintf("/video/%s/output/%s/power-save/", connector, device)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result PowerSaveStatus `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetPowerSave changes power save setting
func (s *VideoService) SetPowerSave(connector, device string, enabled bool) error {
	path := fmt.Sprintf("/video/%s/output/%s/power-save/", connector, device)
	payload := PowerSaveStatus{Enabled: enabled}

	resp, err := s.client.doRequest("PUT", path, payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set power save: status %d", resp.StatusCode)
	}

	return nil
}

// GetAvailableModes gets available video modes
func (s *VideoService) GetAvailableModes(connector, device string) ([]VideoModeInfo, error) {
	path := fmt.Sprintf("/video/%s/output/%s/modes/", connector, device)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result []VideoModeInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Result, nil
}

// GetCurrentMode returns current video mode
func (s *VideoService) GetCurrentMode(connector, device string) (*VideoModeInfo, error) {
	path := fmt.Sprintf("/video/%s/output/%s/mode/", connector, device)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result VideoModeInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetVideoMode sets video mode
func (s *VideoService) SetVideoMode(connector, device, mode string) error {
	path := fmt.Sprintf("/video/%s/output/%s/mode/", connector, device)
	payload := map[string]string{"mode": mode}

	resp, err := s.client.doRequest("PUT", path, payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set video mode: status %d", resp.StatusCode)
	}

	return nil
}

// SendCEC sends CEC payload out of HDMI port (experimental)
func (s *VideoService) SendCEC(hexCommand string) error {
	payload := map[string]string{"hexCommand": hexCommand}

	resp, err := s.client.doRequest("POST", "/sendCecX/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send CEC command: status %d", resp.StatusCode)
	}

	return nil
}