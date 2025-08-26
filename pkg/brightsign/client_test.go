package brightsign

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := Config{
		Host:     "192.168.1.100",
		Username: "admin",
		Password: "password123",
		Debug:    true,
		Timeout:  10 * time.Second,
	}

	client := NewClient(config)

	if client.host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, client.host)
	}

	if client.username != config.Username {
		t.Errorf("Expected username %s, got %s", config.Username, client.username)
	}

	if client.password != config.Password {
		t.Errorf("Expected password %s, got %s", config.Password, client.password)
	}

	if client.debug != config.Debug {
		t.Errorf("Expected debug %v, got %v", config.Debug, client.debug)
	}

	expectedBaseURL := "http://192.168.1.100/api/v1"
	if client.baseURL != expectedBaseURL {
		t.Errorf("Expected baseURL %s, got %s", expectedBaseURL, client.baseURL)
	}

	// Check that services are initialized
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

func TestNewClientDefaults(t *testing.T) {
	config := Config{
		Host:     "test.local",
		Password: "test",
	}

	client := NewClient(config)

	if client.username != "admin" {
		t.Errorf("Expected default username 'admin', got %s", client.username)
	}

	if client.client.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", client.client.Timeout)
	}
}

func TestParseDigestAuth(t *testing.T) {
	wwwAuth := `Digest realm="BrightSign", nonce="abc123", qop="auth", opaque="xyz789"`
	
	params := parseDigestAuth(wwwAuth)
	
	expected := map[string]string{
		"realm":  "BrightSign",
		"nonce":  "abc123",
		"qop":    "auth",
		"opaque": "xyz789",
	}

	for key, expectedValue := range expected {
		if value, exists := params[key]; !exists || value != expectedValue {
			t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, value)
		}
	}
}

func TestDoRequestSuccess(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"result":"success"}}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := Config{
		Host:     server.URL[7:], // Remove http:// prefix
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	// Test request
	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDoRequestWithDigestAuth(t *testing.T) {
	authAttempts := 0
	
	// Create a test server that requires digest auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			authAttempts++
			w.Header().Set("WWW-Authenticate", `Digest realm="BrightSign", nonce="abc123", qop="auth"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Check if digest auth header is present
		if authHeader[:6] != "Digest" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"result":"authenticated"}}`))
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:], // Remove http:// prefix
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	// Test request
	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if authAttempts != 1 {
		t.Errorf("Expected 1 auth attempt, got %d", authAttempts)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMd5Hash(t *testing.T) {
	input := "test"
	expected := "098f6bcd4621d373cade4e832627b4f6"
	
	result := md5Hash(input)
	
	if result != expected {
		t.Errorf("Expected MD5 hash %s, got %s", expected, result)
	}
}