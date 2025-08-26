package brightsign

import (
	"fmt"
)

// DisplayService handles display control endpoints (Moka displays, BOS 9.0.189+)
type DisplayService struct {
	client *Client
}

// DisplaySettings represents all display control settings
type DisplaySettings struct {
	Brightness      *BrightnessSettings      `json:"brightness,omitempty"`
	Contrast        *ContrastSettings        `json:"contrast,omitempty"`
	AlwaysConnected *AlwaysConnectedSettings `json:"alwaysConnected,omitempty"`
	PowerSettings   *PowerSettings           `json:"powerSettings,omitempty"`
	StandbyTimeout  *StandbyTimeoutSettings  `json:"standbyTimeout,omitempty"`
	SDConnection    *SDConnectionSettings    `json:"sdConnection,omitempty"`
	VideoOutput     *VideoOutputSettings     `json:"videoOutput,omitempty"`
	Volume          *VolumeSettings          `json:"volume,omitempty"`
	WhiteBalance    *WhiteBalanceSettings    `json:"whiteBalance,omitempty"`
}

// BrightnessSettings represents brightness settings
type BrightnessSettings struct {
	Value int `json:"value"`
	Min   int `json:"min,omitempty"`
	Max   int `json:"max,omitempty"`
}

// ContrastSettings represents contrast settings
type ContrastSettings struct {
	Value int `json:"value"`
	Min   int `json:"min,omitempty"`
	Max   int `json:"max,omitempty"`
}

// AlwaysConnectedSettings represents connection settings
type AlwaysConnectedSettings struct {
	Enabled bool `json:"enabled"`
}

// PowerSettings represents power settings
type PowerSettings struct {
	State string `json:"state"` // "on" or "standby"
}

// StandbyTimeoutSettings represents standby timeout settings
type StandbyTimeoutSettings struct {
	Seconds int `json:"seconds"`
	Min     int `json:"min,omitempty"`
	Max     int `json:"max,omitempty"`
}

// SDConnectionSettings represents SD connection settings
type SDConnectionSettings struct {
	Target string `json:"target"` // "display" or "brightsign"
}

// VideoOutputSettings represents video output settings
type VideoOutputSettings struct {
	Output string `json:"output"` // "HDMI1" or "HDMI2"
}

// VolumeSettings represents volume settings
type VolumeSettings struct {
	Value int `json:"value"`
	Min   int `json:"min,omitempty"`
	Max   int `json:"max,omitempty"`
}

// WhiteBalanceSettings represents white balance settings
type WhiteBalanceSettings struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Min   int `json:"min,omitempty"`
	Max   int `json:"max,omitempty"`
}

// DisplayInfo represents display information
type DisplayInfo struct {
	Model        string `json:"model"`
	SerialNumber string `json:"serialNumber"`
	Version      string `json:"version"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// GetAll returns all control settings for connected display
func (s *DisplayService) GetAll() (*DisplaySettings, error) {
	resp, err := s.client.doRequest("GET", "/display-control/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result DisplaySettings `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// GetBrightness returns brightness settings
func (s *DisplayService) GetBrightness() (*BrightnessSettings, error) {
	resp, err := s.client.doRequest("GET", "/display-control/brightness/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result BrightnessSettings `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetBrightness changes brightness setting
func (s *DisplayService) SetBrightness(value int) error {
	payload := BrightnessSettings{Value: value}
	resp, err := s.client.doRequest("PUT", "/display-control/brightness/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set brightness: status %d", resp.StatusCode)
	}

	return nil
}

// GetContrast returns contrast settings
func (s *DisplayService) GetContrast() (*ContrastSettings, error) {
	resp, err := s.client.doRequest("GET", "/display-control/contrast/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result ContrastSettings `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetContrast changes contrast setting
func (s *DisplayService) SetContrast(value int) error {
	payload := ContrastSettings{Value: value}
	resp, err := s.client.doRequest("PUT", "/display-control/contrast/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set contrast: status %d", resp.StatusCode)
	}

	return nil
}

// GetVolume returns volume settings
func (s *DisplayService) GetVolume() (*VolumeSettings, error) {
	resp, err := s.client.doRequest("GET", "/display-control/volume/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result VolumeSettings `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetVolume changes volume level
func (s *DisplayService) SetVolume(value int) error {
	payload := VolumeSettings{Value: value}
	resp, err := s.client.doRequest("PUT", "/display-control/volume/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set volume: status %d", resp.StatusCode)
	}

	return nil
}

// GetPowerSettings returns power settings
func (s *DisplayService) GetPowerSettings() (*PowerSettings, error) {
	resp, err := s.client.doRequest("GET", "/display-control/power-settings/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result PowerSettings `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// SetPowerSettings changes power setting
func (s *DisplayService) SetPowerSettings(state string) error {
	payload := PowerSettings{State: state}
	resp, err := s.client.doRequest("PUT", "/display-control/power-settings/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set power settings: status %d", resp.StatusCode)
	}

	return nil
}

// GetInfo returns display information
func (s *DisplayService) GetInfo() (*DisplayInfo, error) {
	resp, err := s.client.doRequest("GET", "/display-control/info/", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Result DisplayInfo `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data.Result, nil
}

// UpdateFirmware updates display firmware
func (s *DisplayService) UpdateFirmware(filepathOrURL string) error {
	payload := map[string]string{"source": filepathOrURL}
	resp, err := s.client.doRequest("PUT", "/display-control/firmware/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to update firmware: status %d", resp.StatusCode)
	}

	return nil
}