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

	if len(args) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	source := args[0]
	destination := args[1]

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

	if !strings.HasPrefix(remotePath, "/") {
		return fmt.Errorf("remote path must be absolute")
	}

	if !fileExists(source) {
		return fmt.Errorf("source file does not exist: %s", source)
	}

	password, err := promptPassword(host)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	client := dws.NewClient(host, password)

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
	fmt.Printf("Usage: bscp <source> <host:destination>\n")
	fmt.Printf("\nCopy files to BrightSign player using DWS API\n")
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  bscp file.txt 192.168.1.100:/my/file.txt\n")
	fmt.Printf("  bscp video.mp4 player.local:/videos/video.mp4\n")
	fmt.Printf("\nFiles are always copied to the SD card (/storage/sd) on the player.\n")
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