package cli

import (
	"testing"

	"bscli/pkg/brightsign"
)

func TestGetClient_ValidConfig(t *testing.T) {
	// Set global vars for testing
	host = "192.168.1.100"
	username = "admin"
	password = "testpass"
	debug = true

	client, err := getClient()
	if err != nil {
		t.Fatalf("getClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	// Verify client was configured correctly by checking that services exist
	if client.Info == nil {
		t.Error("Info service not initialized")
	}
	if client.Control == nil {
		t.Error("Control service not initialized")
	}
	if client.Storage == nil {
		t.Error("Storage service not initialized")
	}
}

func TestGetClient_MissingHost(t *testing.T) {
	// Reset global vars
	host = ""
	username = "admin"
	password = "testpass"

	_, err := getClient()
	if err == nil {
		t.Error("Expected error when host is missing, got nil")
	}

	expectedError := "host is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// Test helper functions
func TestFormatSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{5368709120, "5.0 GB"},
	}

	for _, test := range tests {
		result := formatSize(test.input)
		if result != test.expected {
			t.Errorf("formatSize(%d): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

// Mock test to verify brightsign client creation
func TestBrightSignClientCreation(t *testing.T) {
	config := brightsign.Config{
		Host:     "test.local",
		Username: "admin",
		Password: "password",
		Debug:    false,
	}

	client := brightsign.NewClient(config)
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// Verify all services are initialized
	services := map[string]interface{}{
		"Info":        client.Info,
		"Control":     client.Control,
		"Storage":     client.Storage,
		"Diagnostics": client.Diagnostics,
		"Display":     client.Display,
		"Registry":    client.Registry,
		"Logs":        client.Logs,
		"Video":       client.Video,
	}

	for name, service := range services {
		if service == nil {
			t.Errorf("%s service not initialized", name)
		}
	}
}