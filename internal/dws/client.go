package dws

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	host     string
	password string
	client   *http.Client
	debug    bool
}

type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type ListResponse struct {
	Files []FileInfo `json:"files"`
}

type NestedListResponse struct {
	Data struct {
		Result struct {
			Files []FileInfo `json:"files"`
		} `json:"result"`
	} `json:"data"`
}

type UploadResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Results []string `json:"results"`
}

type NestedUploadResponse struct {
	Data struct {
		Result UploadResponse `json:"result"`
	} `json:"data"`
}

func NewClient(host, password string, debug bool) *Client {
	return &Client{
		host:     host,
		password: password,
		debug:    debug,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// debugf prints debug messages only if debug mode is enabled
func (c *Client) debugf(format string, args ...interface{}) {
	if c.debug {
		fmt.Printf("DEBUG: "+format, args...)
	}
}

func (c *Client) UploadFile(localPath, remotePath string) error {
	// Extract storage device (e.g., "sd") from path like "/storage/sd/file.txt"
	pathParts := strings.Split(strings.Trim(remotePath, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "storage" {
		return fmt.Errorf("invalid remote path format, expected /storage/{device}/...")
	}
	
	// Only allow root level files to prevent directory creation
	if len(pathParts) != 3 {
		return fmt.Errorf("only root level files are supported (no subdirectories), got: %s", remotePath)
	}
	
	storageDevice := pathParts[1]
	url := fmt.Sprintf("http://%s/api/v1/files/%s/", c.host, storageDevice)
	
	// Use only the filename - no paths that could trigger directory creation
	targetFilename := pathParts[2]
	
	c.debugf("Upload URL: %s", url)
	c.debugf("Target filename in form: %s", targetFilename)

	// Create a function that builds the multipart body and returns bytes
	createBodyBytes := func() ([]byte, string, error) {
		file, err := os.Open(localPath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		
		part, err := writer.CreateFormFile("file", targetFilename)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create form file: %w", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", fmt.Errorf("failed to copy file content: %w", err)
		}

		contentType := writer.FormDataContentType()
		err = writer.Close()
		if err != nil {
			return nil, "", fmt.Errorf("failed to close writer: %w", err)
		}

		return body.Bytes(), contentType, nil
	}

	// Get body bytes and content type
	bodyBytes, contentType, err := createBodyBytes()
	if err != nil {
		return err
	}
	
	// First attempt without authentication
	req, err := http.NewRequest("PUT", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.ContentLength = int64(len(bodyBytes))
	
	// Use a simple HTTP client without digest transport
	client := &http.Client{Timeout: 30 * time.Second}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	
	// If we get 401, handle digest authentication manually
	if resp.StatusCode == 401 {
		wwwAuth := resp.Header.Get("WWW-Authenticate")
		resp.Body.Close()
		
		if !strings.HasPrefix(wwwAuth, "Digest") {
			return fmt.Errorf("server requires digest authentication but sent: %s", wwwAuth)
		}
		
		// Parse digest challenge
		authParams := parseDigestAuth(wwwAuth)
		
		// Create fresh request with digest auth
		req, err := http.NewRequest("PUT", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return fmt.Errorf("failed to create authenticated request: %w", err)
		}
		
		req.Header.Set("Content-Type", contentType)
		req.ContentLength = int64(len(bodyBytes))
		
		// Create digest authorization header
		authHeader := createDigestAuthHeader("admin", c.password, "PUT", req.URL.RequestURI(), authParams)
		req.Header.Set("Authorization", authHeader)
		
		// Retry with authentication
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to upload file with auth: %w", err)
		}
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to parse as nested response first
	var nestedResp NestedUploadResponse
	if err := json.Unmarshal(respBody, &nestedResp); err == nil {
		// Check if this is actually a nested response by seeing if Data has content
		if nestedResp.Data.Result.Success || nestedResp.Data.Result.Message != "" || len(nestedResp.Data.Result.Results) > 0 {
			uploadResp := nestedResp.Data.Result
			if !uploadResp.Success {
				return fmt.Errorf("upload failed: %s (full response: %s)", uploadResp.Message, string(respBody))
			}
			// Success!
			return nil
		}
	}

	// Fall back to direct response format
	var uploadResp UploadResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		// If JSON parsing fails, show the raw response
		return fmt.Errorf("failed to decode response as JSON (got: %s): %w", string(respBody), err)
	}

	if !uploadResp.Success {
		// Show both the message and the full response for debugging
		return fmt.Errorf("upload failed: %s (full response: %s)", uploadResp.Message, string(respBody))
	}

	return nil
}

func (c *Client) ListFiles(path string) ([]FileInfo, error) {
	// Convert path like "/storage/sd" to "sd"
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "storage" {
		return nil, fmt.Errorf("invalid path format, expected /storage/{device}/...")
	}
	
	storageDevice := pathParts[1]
	url := fmt.Sprintf("http://%s/api/v1/files/%s/", c.host, storageDevice)
	
	// First attempt without authentication
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use a simple HTTP client without digest transport
	client := &http.Client{Timeout: 30 * time.Second}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	
	// If we get 401, handle digest authentication manually
	if resp.StatusCode == 401 {
		wwwAuth := resp.Header.Get("WWW-Authenticate")
		resp.Body.Close()
		
		if !strings.HasPrefix(wwwAuth, "Digest") {
			return nil, fmt.Errorf("server requires digest authentication but sent: %s", wwwAuth)
		}
		
		// Parse digest challenge
		authParams := parseDigestAuth(wwwAuth)
		
		// Create fresh request with digest auth
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticated request: %w", err)
		}
		
		// Create digest authorization header
		authHeader := createDigestAuthHeader("admin", c.password, "GET", req.URL.RequestURI(), authParams)
		req.Header.Set("Authorization", authHeader)
		
		// Retry with authentication
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to list files with auth: %w", err)
		}
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug: show what we got
	c.debugf("List API response: %s\n", string(respBody))

	// Try to parse as nested response first
	var nestedResp NestedListResponse
	if err := json.Unmarshal(respBody, &nestedResp); err == nil {
		// Check if this is actually a nested response by seeing if Data has content
		if len(nestedResp.Data.Result.Files) > 0 {
			c.debugf("Found %d files (nested): %+v\n", len(nestedResp.Data.Result.Files), nestedResp.Data.Result.Files)
			return nestedResp.Data.Result.Files, nil
		}
	}

	// Fall back to direct response format
	var listResp ListResponse
	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response as JSON (got: %s): %w", string(respBody), err)
	}

	// Debug: show parsed files
	c.debugf("Found %d files (direct): %+v\n", len(listResp.Files), listResp.Files)

	return listResp.Files, nil
}

func (c *Client) VerifyFileExists(path string) (bool, error) {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	
	c.debugf("Looking for file '%s' in directory '%s'\n", filename, dir)

	files, err := c.ListFiles(dir)
	if err != nil {
		return false, err
	}

	c.debugf("Comparing filename '%s' against:\n", filename)
	for _, file := range files {
		c.debugf("  - '%s' (match: %t)\n", file.Name, strings.EqualFold(file.Name, filename))
		if strings.EqualFold(file.Name, filename) {
			return true, nil
		}
	}

	return false, nil
}

func parseDigestAuth(wwwAuth string) map[string]string {
	params := make(map[string]string)
	
	// Remove "Digest " prefix
	auth := strings.TrimPrefix(wwwAuth, "Digest ")
	
	// Split by comma and parse key=value pairs
	parts := strings.Split(auth, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx != -1 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			// Remove quotes
			value = strings.Trim(value, `"`)
			params[key] = value
		}
	}
	
	return params
}

func createDigestAuthHeader(username, password, method, uri string, params map[string]string) string {
	realm := params["realm"]
	nonce := params["nonce"]
	qop := params["qop"]
	opaque := params["opaque"]
	
	// Generate cnonce
	rand.Seed(time.Now().UnixNano())
	cnonce := fmt.Sprintf("%08x", rand.Uint32())
	nc := "00000001"
	
	// Calculate response hash
	ha1 := md5Hash(fmt.Sprintf("%s:%s:%s", username, realm, password))
	ha2 := md5Hash(fmt.Sprintf("%s:%s", method, uri))
	
	var response string
	if qop == "auth" || qop == "auth-int" {
		response = md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, nonce, nc, cnonce, qop, ha2))
	} else {
		response = md5Hash(fmt.Sprintf("%s:%s:%s", ha1, nonce, ha2))
	}
	
	// Build authorization header
	authHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s"`,
		username, realm, nonce, uri, response)
	
	if qop != "" {
		authHeader += fmt.Sprintf(`, qop=%s, nc=%s, cnonce="%s"`, qop, nc, cnonce)
	}
	
	if opaque != "" {
		authHeader += fmt.Sprintf(`, opaque="%s"`, opaque)
	}
	
	return authHeader
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}