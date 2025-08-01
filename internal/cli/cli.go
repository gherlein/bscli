package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"bscp/internal/dws"
	"golang.org/x/term"
)

const (
	sdCardPath = "/storage/sd"
)

func Run(args []string) error {
	if len(args) == 0 {
		return showUsage()
	}

	if args[0] == "-h" || args[0] == "--help" {
		return showUsage()
	}

	// Parse debug flag
	debug := false
	var filteredArgs []string
	
	for _, arg := range args {
		if arg == "-debug" || arg == "-d" {
			debug = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	if len(filteredArgs) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	source := filteredArgs[0]
	destination := filteredArgs[1]

	if !strings.Contains(destination, ":") {
		return fmt.Errorf("destination must be in format host:path")
	}

	// Split by the last colon to handle host:port:path or host:path
	lastColon := strings.LastIndex(destination, ":")
	if lastColon == -1 {
		return fmt.Errorf("invalid destination format, expected host:path")
	}

	host := destination[:lastColon]
	remotePath := destination[lastColon+1:]

	// If path doesn't start with /, assume it's relative to /storage/sd
	if !strings.HasPrefix(remotePath, "/") {
		remotePath = "/storage/sd/" + remotePath
	}

	if !fileExists(source) {
		return fmt.Errorf("source file does not exist: %s", source)
	}

	password, err := promptPassword(host)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	client := dws.NewClient(host, password, debug)

	// If remotePath ends with '/', it's a directory - append source filename
	var targetPath string
	if strings.HasSuffix(remotePath, "/") {
		targetPath = filepath.Join(remotePath, filepath.Base(source))
	} else {
		targetPath = remotePath
	}
	
	// Ensure the path starts with /storage/sd
	if !strings.HasPrefix(targetPath, sdCardPath) {
		return fmt.Errorf("remote path must be under %s, got: %s", sdCardPath, targetPath)
	}
	
	// Check if target path tries to create subdirectories (only allow root level files)
	pathParts := strings.Split(strings.Trim(targetPath, "/"), "/")
	if len(pathParts) > 3 {
		return fmt.Errorf("subdirectories not supported to prevent directory creation, use only /storage/sd/filename")
	}

	fmt.Printf("Uploading %s to %s:%s...\n", source, host, targetPath)

	err = client.UploadFile(source, targetPath)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("Verifying file exists at destination...\n")

	exists, err := client.VerifyFileExists(targetPath)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	if !exists {
		return fmt.Errorf("file upload succeeded but file not found at destination")
	}

	fmt.Printf("Successfully copied %s to %s:%s\n", source, host, targetPath)
	return nil
}

func showUsage() error {
	fmt.Printf("Usage: bscp [-debug|-d] <source> <host:destination>\n")
	fmt.Printf("\nCopy files to BrightSign player using DWS API\n")
	fmt.Printf("\nOptions:\n")
	fmt.Printf("  -debug, -d    Enable debug output\n")
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  bscp file.txt 192.168.1.100:/storage/sd/file.txt\n")
	fmt.Printf("  bscp -debug video.mp4 player.local:/storage/sd/video.mp4\n")
	fmt.Printf("  bscp -d video.mp4 player.local:/storage/sd/\n")
	fmt.Printf("\nFiles are copied to the root of the SD card (/storage/sd/) only.\n")
	fmt.Printf("Subdirectories are not supported to prevent directory creation errors.\n")
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func promptPassword(host string) (string, error) {
	fmt.Printf("Password for %s: ", host)
	
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	
	fmt.Println()
	return string(bytePassword), nil
}