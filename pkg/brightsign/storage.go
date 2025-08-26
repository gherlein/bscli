package brightsign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StorageService handles file and storage operations
type StorageService struct {
	client *Client
}

// FileInfo represents information about a file or directory
type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Modified string `json:"lastModified,omitempty"`
}

// ListOptions contains options for listing files
type ListOptions struct {
	Raw bool // If true, returns raw directory listing
}

// ListFiles lists files and directories in the specified path
func (s *StorageService) ListFiles(path string, options *ListOptions) ([]FileInfo, error) {
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Convert path like "/storage/sd/" to API path "/files/sd/"
	apiPath := strings.Replace(path, "/storage/", "/files/", 1)
	
	if options != nil && options.Raw {
		apiPath += "?raw"
	}

	resp, err := s.client.doRequest("GET", apiPath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read raw response to understand structure
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if s.client.debug {
		fmt.Printf("DEBUG: ListFiles API response: %s\n", string(bodyBytes))
	}

	// Try to parse as array first (directory listing)
	var arrayResult struct {
		Data struct {
			Result []FileInfo `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &arrayResult); err == nil && len(arrayResult.Data.Result) > 0 {
		return arrayResult.Data.Result, nil
	}

	// Try to parse as single object (single file info)
	var singleResult struct {
		Data struct {
			Result FileInfo `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &singleResult); err == nil {
		return []FileInfo{singleResult.Data.Result}, nil
	}

	// Try to parse as object with files property
	var objectResult struct {
		Data struct {
			Result struct {
				Files []FileInfo `json:"files"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &objectResult); err == nil {
		return objectResult.Data.Result.Files, nil
	}

	// If none of the above worked, return the parsing error
	return nil, fmt.Errorf("failed to parse response as known format: %s", string(bodyBytes))
}

// UploadFile uploads a file to the specified path on the player
func (s *StorageService) UploadFile(localPath, remotePath string) error {
	// Open the local file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field
	filename := filepath.Base(remotePath)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Convert path like "/storage/sd/file.txt" to API path "/files/sd/"
	dir := filepath.Dir(remotePath)
	apiPath := strings.Replace(dir, "/storage/", "/files/", 1) + "/"

	// Make request
	url := s.client.baseURL + apiPath
	resp, err := s.client.doRequestWithBody("PUT", url, bytes.NewReader(body.Bytes()), contentType)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if s.client.debug {
		fmt.Printf("DEBUG: Uploaded %s (%d bytes) to %s\n", localPath, fileInfo.Size(), remotePath)
	}

	return nil
}

// DownloadFile downloads a file from the player to local path
func (s *StorageService) DownloadFile(remotePath, localPath string) error {
	// Convert path like "/storage/sd/file.txt" to API path "/files/sd/file.txt?contents&stream"
	apiPath := strings.Replace(remotePath, "/storage/", "/files/", 1) + "?contents&stream"

	resp, err := s.client.doRequest("GET", apiPath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create local file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	// Copy content
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if s.client.debug {
		fmt.Printf("DEBUG: Downloaded %s (%d bytes) to %s\n", remotePath, written, localPath)
	}

	return nil
}

// DeleteFile deletes a file or directory
func (s *StorageService) DeleteFile(path string) error {
	// Convert path like "/storage/sd/file.txt" to API path "/files/sd/file.txt"
	apiPath := strings.Replace(path, "/storage/", "/files/", 1)

	resp, err := s.client.doRequest("DELETE", apiPath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// RenameFile renames a file
func (s *StorageService) RenameFile(oldPath, newName string) error {
	// Convert path like "/storage/sd/file.txt" to API path "/files/sd/"
	dir := filepath.Dir(oldPath)
	apiPath := strings.Replace(dir, "/storage/", "/files/", 1) + "/"

	payload := map[string]string{
		"oldName": filepath.Base(oldPath),
		"newName": newName,
	}

	resp, err := s.client.doRequest("POST", apiPath, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("rename failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// CreateDirectory creates a new directory
func (s *StorageService) CreateDirectory(path string) error {
	// Convert path like "/storage/sd/newdir" to API path "/files/sd/"
	dir := filepath.Dir(path)
	dirName := filepath.Base(path)
	apiPath := strings.Replace(dir, "/storage/", "/files/", 1) + "/"

	// Create form data for directory creation
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	
	if err := writer.WriteField("directory", dirName); err != nil {
		return fmt.Errorf("failed to write directory field: %w", err)
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	url := s.client.baseURL + apiPath
	resp, err := s.client.doRequestWithBody("PUT", url, bytes.NewReader(body.Bytes()), contentType)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create directory failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// FormatStorage formats a storage device
func (s *StorageService) FormatStorage(device string) error {
	// device should be like "sd", "usb1", etc.
	apiPath := fmt.Sprintf("/storage/%s/", device)

	resp, err := s.client.doRequest("DELETE", apiPath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("format failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}