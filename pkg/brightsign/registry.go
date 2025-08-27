package brightsign

import (
	"fmt"
)

// RegistryService handles registry operations
type RegistryService struct {
	client *Client
}

// RegistryValue represents a registry key-value pair
type RegistryValue struct {
	Value string `json:"value"`
}

// GetAll returns entire registry dump (excludes hidden sections)
func (s *RegistryService) GetAll() (interface{}, error) {
	resp, err := s.client.doRequest("GET", "/registry/", nil)
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

// GetValue returns specific registry key value
func (s *RegistryService) GetValue(section, key string) (string, error) {
	path := fmt.Sprintf("/registry/%s/%s/", section, key)

	resp, err := s.client.doRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Result RegistryValue `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return "", err
	}

	return result.Data.Result.Value, nil
}

// SetValue creates or updates registry value
func (s *RegistryService) SetValue(section, key, value string) error {
	path := fmt.Sprintf("/registry/%s/%s/", section, key)

	payload := RegistryValue{Value: value}
	resp, err := s.client.doRequest("PUT", path, payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set registry value: status %d", resp.StatusCode)
	}

	return nil
}

// DeleteValue removes specific registry value
func (s *RegistryService) DeleteValue(section, key string) error {
	path := fmt.Sprintf("/registry/%s/%s/", section, key)

	resp, err := s.client.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to delete registry value: status %d", resp.StatusCode)
	}

	return nil
}

// DeleteSection deletes entire registry section
func (s *RegistryService) DeleteSection(section string) error {
	path := fmt.Sprintf("/registry/%s/", section)

	resp, err := s.client.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to delete registry section: status %d", resp.StatusCode)
	}

	return nil
}

// GetRecoveryURL retrieves recovery URL from player registry
func (s *RegistryService) GetRecoveryURL() (string, error) {
	resp, err := s.client.doRequest("GET", "/registry/recovery_url/", nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Result struct {
				URL string `json:"url"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := parseJSON(resp, &result); err != nil {
		return "", err
	}

	return result.Data.Result.URL, nil
}

// SetRecoveryURL updates recovery URL in player registry
func (s *RegistryService) SetRecoveryURL(url string) error {
	payload := map[string]string{"url": url}

	resp, err := s.client.doRequest("PUT", "/registry/recovery_url/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to set recovery URL: status %d", resp.StatusCode)
	}

	return nil
}

// Flush flushes registry contents to persistent storage (BOS 9.0.107+)
func (s *RegistryService) Flush() error {
	resp, err := s.client.doRequest("PUT", "/registry/flush/", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to flush registry: status %d", resp.StatusCode)
	}

	return nil
}