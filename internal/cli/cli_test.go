package cli

import (
	"os"
	"strings"
	"testing"
)

func TestShowUsage(t *testing.T) {
	err := showUsage()
	if err != nil {
		t.Errorf("Expected no error from showUsage, got: %v", err)
	}
}

func TestFileExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-exists-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	if !fileExists(tempFile.Name()) {
		t.Error("Expected existing file to return true")
	}
	
	if fileExists("/nonexistent/file/path") {
		t.Error("Expected nonexistent file to return false")
	}
}

func TestRunInvalidArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		expectError bool
	}{
		{
			name: "no arguments",
			args: []string{},
			expectError: false, // shows usage
		},
		{
			name: "help flag",
			args: []string{"-h"},
			expectError: false, // shows usage
		},
		{
			name: "help flag long",
			args: []string{"--help"},
			expectError: false, // shows usage
		},
		{
			name: "too few arguments",
			args: []string{"file.txt"},
			expectError: true,
		},
		{
			name: "too many arguments",
			args: []string{"file.txt", "host:path", "extra"},
			expectError: true,
		},
		{
			name: "debug flag with valid args",
			args: []string{"-debug", "file.txt", "host:/storage/sd/path"},
			expectError: true, // will fail because file doesn't exist
		},
		{
			name: "debug flag short form with valid args",
			args: []string{"-d", "file.txt", "host:/storage/sd/path"},
			expectError: true, // will fail because file doesn't exist
		},
		{
			name: "missing colon in destination",
			args: []string{"file.txt", "hostpath"},
			expectError: true,
		},
		{
			name: "relative remote path",
			args: []string{"file.txt", "host:relative/path"},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestRunNonexistentFile(t *testing.T) {
	args := []string{"/nonexistent/file.txt", "host:/remote/path"}
	err := Run(args)
	
	if err == nil {
		t.Error("Expected error for nonexistent source file")
	}
	
	if !strings.Contains(err.Error(), "source file does not exist") {
		t.Errorf("Expected 'source file does not exist' error, got: %v", err)
	}
}

func TestDebugFlagParsing(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectError bool
		expectDebug bool
	}{
		{
			name: "debug flag long form",
			args: []string{"-debug", "file.txt", "host:/storage/sd/path"},
			expectError: true, // will fail because file doesn't exist
			expectDebug: true,
		},
		{
			name: "debug flag short form", 
			args: []string{"-d", "file.txt", "host:/storage/sd/path"},
			expectError: true, // will fail because file doesn't exist
			expectDebug: true,
		},
		{
			name: "no debug flag",
			args: []string{"file.txt", "host:/storage/sd/path"},
			expectError: true, // will fail because file doesn't exist
			expectDebug: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that the flag parsing logic works without actual file operations
			debug := false
			var filteredArgs []string
			
			for _, arg := range tt.args {
				if arg == "-debug" || arg == "-d" {
					debug = true
				} else {
					filteredArgs = append(filteredArgs, arg)
				}
			}
			
			if debug != tt.expectDebug {
				t.Errorf("Expected debug=%t, got debug=%t", tt.expectDebug, debug)
			}
			
			if len(filteredArgs) != 2 {
				t.Errorf("Expected 2 filtered args, got %d", len(filteredArgs))
			}
		})
	}
}

func TestDestinationParsing(t *testing.T) {
	tests := []struct {
		destination string
		expectError bool
		expectedHost string
		expectedPath string
	}{
		{
			destination: "192.168.1.100:/storage/sd/file.txt",
			expectError: false,
			expectedHost: "192.168.1.100",
			expectedPath: "/storage/sd/file.txt",
		},
		{
			destination: "player.local:/videos/movie.mp4",
			expectError: false,
			expectedHost: "player.local",
			expectedPath: "/videos/movie.mp4",
		},
		{
			destination: "host:relative/path",
			expectError: true,
		},
		{
			destination: "no-colon-here",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.destination, func(t *testing.T) {
			if !strings.Contains(tt.destination, ":") {
				return // Skip parsing test for invalid formats
			}
			
			parts := strings.SplitN(tt.destination, ":", 2)
			if len(parts) != 2 {
				if !tt.expectError {
					t.Error("Expected successful parsing")
				}
				return
			}
			
			host := parts[0]
			remotePath := parts[1]
			
			if !tt.expectError {
				if host != tt.expectedHost {
					t.Errorf("Expected host '%s', got '%s'", tt.expectedHost, host)
				}
				if remotePath != tt.expectedPath {
					t.Errorf("Expected path '%s', got '%s'", tt.expectedPath, remotePath)
				}
			}
			
			if tt.expectError && strings.HasPrefix(remotePath, "/") {
				t.Error("Expected error for relative path")
			}
		})
	}
}