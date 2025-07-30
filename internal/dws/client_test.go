package dws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("localhost", "password123", false)
	
	if client.host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", client.host)
	}
	
	if client.password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", client.password)
	}
	
	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestUploadFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	testContent := "Hello, BrightSign!"
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/files/sd/" {
			t.Errorf("Expected path '/api/v1/files/sd/', got '%s'", r.URL.Path)
		}
		
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got '%s'", r.Method)
		}
		
		// For testing, just accept any request (skip auth validation)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}
		
		response := UploadResponse{Success: true, Message: "Upload successful"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := &Client{
		host:     server.URL[7:], // Remove http://
		password: "testpass",
		client:   server.Client(),
		debug:    false,
	}
	
	err = client.UploadFile(tempFile.Name(), "/storage/sd/test.txt")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestUploadFileError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For testing, just return error (skip auth validation)
		w.WriteHeader(http.StatusInternalServerError)
		response := UploadResponse{Success: false, Message: "Server error"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tempFile, err := os.CreateTemp("", "test-upload-error-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	client := &Client{
		host:     server.URL[7:],
		password: "testpass",
		client:   server.Client(),
		debug:    false,
	}
	
	err = client.UploadFile(tempFile.Name(), "/storage/sd/test.txt")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestListFiles(t *testing.T) {
	expectedFiles := []FileInfo{
		{Name: "file1.txt", Path: "/storage/sd/file1.txt", Size: 1024},
		{Name: "file2.mp4", Path: "/storage/sd/file2.mp4", Size: 2048},
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/files/sd/" {
			t.Errorf("Expected path '/api/v1/files/sd/', got '%s'", r.URL.Path)
		}
		
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got '%s'", r.Method)
		}
		
		// For testing, just accept any request (skip auth validation)
		response := ListResponse{Files: expectedFiles}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := &Client{
		host:     server.URL[7:],
		password: "testpass",
		client:   server.Client(),
		debug:    false,
	}
	
	files, err := client.ListFiles("/storage/sd")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}
	
	for i, file := range files {
		if file.Name != expectedFiles[i].Name {
			t.Errorf("Expected file name '%s', got '%s'", expectedFiles[i].Name, file.Name)
		}
	}
}

func TestNewClientWithDebug(t *testing.T) {
	client := NewClient("localhost", "password123", true)
	
	if client.host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", client.host)
	}
	
	if client.password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", client.password)
	}
	
	if !client.debug {
		t.Error("Expected debug mode to be enabled")
	}
	
	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestVerifyFileExists(t *testing.T) {
	files := []FileInfo{
		{Name: "existing.txt", Path: "/storage/sd/existing.txt", Size: 1024},
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ListResponse{Files: files}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := &Client{
		host:     server.URL[7:],
		password: "testpass",
		client:   server.Client(),
		debug:    false,
	}
	
	exists, err := client.VerifyFileExists("/storage/sd/existing.txt")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !exists {
		t.Error("Expected file to exist")
	}
	
	exists, err = client.VerifyFileExists("/storage/sd/nonexistent.txt")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if exists {
		t.Error("Expected file to not exist")
	}
}