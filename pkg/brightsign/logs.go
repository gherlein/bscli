package brightsign

// LogsService handles log retrieval
type LogsService struct {
	client *Client
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// GetLogs retrieves player serial logs
func (s *LogsService) GetLogs() (string, error) {
	resp, err := s.client.doRequest("GET", "/logs/", nil)
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

// GetSupervisorLoggingLevel returns current logging level
func (s *LogsService) GetSupervisorLoggingLevel() (string, error) {
	resp, err := s.client.doRequest("GET", "/system/supervisor/logging/", nil)
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

// SetSupervisorLoggingLevel sets logging level on player (0-3: error, warn, info, trace)
func (s *LogsService) SetSupervisorLoggingLevel(level int) error {
	if level < 0 || level > 3 {
		level = 2 // default to info
	}

	payload := map[string]int{"level": level}
	resp, err := s.client.doRequest("PUT", "/system/supervisor/logging/", payload)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}