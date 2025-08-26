package brightsign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInfoService_GetInfo(t *testing.T) {
	expectedInfo := DeviceInfo{
		Model:         "HD224",
		Serial:        "123456789",
		Family:        "HD2000",
		BootVersion:   "8.5.35",
		FWVersion:     "9.0.144",
		Uptime:        "2 days, 3:45:22",
		UptimeSeconds: 185722,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/info/" {
			t.Errorf("Expected path /api/v1/info/, got %s", r.URL.Path)
		}
		
		response := struct {
			Data struct {
				Result DeviceInfo `json:"result"`
			} `json:"data"`
		}{}
		response.Data.Result = expectedInfo

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	info, err := client.Info.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

	if info.Model != expectedInfo.Model {
		t.Errorf("Expected model %s, got %s", expectedInfo.Model, info.Model)
	}

	if info.Serial != expectedInfo.Serial {
		t.Errorf("Expected serial %s, got %s", expectedInfo.Serial, info.Serial)
	}

	if info.UptimeSeconds != expectedInfo.UptimeSeconds {
		t.Errorf("Expected uptime %d, got %d", expectedInfo.UptimeSeconds, info.UptimeSeconds)
	}
}

func TestInfoService_GetHealth(t *testing.T) {
	expectedHealth := HealthInfo{
		Status:     "running",
		StatusTime: time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health/" {
			t.Errorf("Expected path /api/v1/health/, got %s", r.URL.Path)
		}
		
		response := struct {
			Data struct {
				Result HealthInfo `json:"result"`
			} `json:"data"`
		}{}
		response.Data.Result = expectedHealth

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	health, err := client.Info.GetHealth()
	if err != nil {
		t.Fatalf("GetHealth failed: %v", err)
	}

	if health.Status != expectedHealth.Status {
		t.Errorf("Expected status %s, got %s", expectedHealth.Status, health.Status)
	}
}

func TestInfoService_SetTime(t *testing.T) {
	timeInfo := TimeInfo{
		Date: "2025-01-15",
		Time: "14:30:00",
		Timezone: "UTC",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/time/" {
			t.Errorf("Expected path /api/v1/time/, got %s", r.URL.Path)
		}

		var receivedTime TimeInfo
		if err := json.NewDecoder(r.Body).Decode(&receivedTime); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedTime.Date != timeInfo.Date {
			t.Errorf("Expected date %s, got %s", timeInfo.Date, receivedTime.Date)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	err := client.Info.SetTime(timeInfo)
	if err != nil {
		t.Fatalf("SetTime failed: %v", err)
	}
}

func TestInfoService_GetVideoMode(t *testing.T) {
	expectedMode := VideoMode{
		Resolution:    "1920x1080",
		FrameRate:     60,
		ScanMethod:    "progressive",
		PreferredMode: true,
		OverscanMode:  "none",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/video-mode/" {
			t.Errorf("Expected path /api/v1/video-mode/, got %s", r.URL.Path)
		}
		
		response := struct {
			Data struct {
				Result VideoMode `json:"result"`
			} `json:"data"`
		}{}
		response.Data.Result = expectedMode

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	mode, err := client.Info.GetVideoMode()
	if err != nil {
		t.Fatalf("GetVideoMode failed: %v", err)
	}

	if mode.Resolution != expectedMode.Resolution {
		t.Errorf("Expected resolution %s, got %s", expectedMode.Resolution, mode.Resolution)
	}

	if mode.FrameRate != expectedMode.FrameRate {
		t.Errorf("Expected frame rate %d, got %d", expectedMode.FrameRate, mode.FrameRate)
	}

	if mode.PreferredMode != expectedMode.PreferredMode {
		t.Errorf("Expected preferred mode %v, got %v", expectedMode.PreferredMode, mode.PreferredMode)
	}
}

func TestInfoService_ListAPIs(t *testing.T) {
	expectedAPIs := []string{
		"/info/",
		"/health/",
		"/time/",
		"/control/reboot/",
		"/files/sd/",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/" {
			t.Errorf("Expected path /api/v1/, got %s", r.URL.Path)
		}
		
		response := struct {
			Data struct {
				Result []string `json:"result"`
			} `json:"data"`
		}{}
		response.Data.Result = expectedAPIs

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		Host:     server.URL[7:],
		Username: "admin",
		Password: "password",
	}
	client := NewClient(config)
	client.baseURL = server.URL + "/api/v1"

	apis, err := client.Info.ListAPIs()
	if err != nil {
		t.Fatalf("ListAPIs failed: %v", err)
	}

	if len(apis) != len(expectedAPIs) {
		t.Errorf("Expected %d APIs, got %d", len(expectedAPIs), len(apis))
	}

	for i, expectedAPI := range expectedAPIs {
		if i >= len(apis) || apis[i] != expectedAPI {
			t.Errorf("Expected API %s at index %d, got %s", expectedAPI, i, apis[i])
		}
	}
}