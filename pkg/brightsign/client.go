// Package brightsign provides a Go client library for interacting with BrightSign players
// via their Diagnostic Web Server (DWS) API.
package brightsign

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Client is the main client for interacting with a BrightSign player's DWS API
type Client struct {
	host     string
	username string
	password string
	client   *http.Client
	debug    bool
	baseURL  string

	// Services
	Info        *InfoService
	Control     *ControlService
	Storage     *StorageService
	Diagnostics *DiagnosticsService
	Display     *DisplayService
	Registry    *RegistryService
	Logs        *LogsService
	Video       *VideoService
}

// Config contains configuration options for the client
type Config struct {
	Host     string
	Username string // Default is "admin"
	Password string
	Debug    bool
	Timeout  time.Duration
}

// Response is the standard API response wrapper
type Response struct {
	Data struct {
		Result interface{} `json:"result"`
	} `json:"data"`
}

// NewClient creates a new BrightSign DWS API client
func NewClient(config Config) *Client {
	if config.Username == "" {
		config.Username = "admin"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	c := &Client{
		host:     config.Host,
		username: config.Username,
		password: config.Password,
		client:   httpClient,
		debug:    config.Debug,
		baseURL:  fmt.Sprintf("http://%s/api/v1", config.Host),
	}

	// Initialize services
	c.Info = &InfoService{client: c}
	c.Control = &ControlService{client: c}
	c.Storage = &StorageService{client: c}
	c.Diagnostics = &DiagnosticsService{client: c}
	c.Display = &DisplayService{client: c}
	c.Registry = &RegistryService{client: c}
	c.Logs = &LogsService{client: c}
	c.Video = &VideoService{client: c}

	return c
}

// doRequest performs an HTTP request with digest authentication if needed
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	return c.doRequestWithBody(method, url, bodyReader, "application/json")
}

// doRequestWithBody performs an HTTP request with a pre-formatted body
func (c *Client) doRequestWithBody(method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if contentType != "" && body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if c.debug {
		fmt.Printf("DEBUG: %s %s\n", method, url)
	}

	// First attempt without authentication
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// If we get 401, handle digest authentication
	if resp.StatusCode == http.StatusUnauthorized {
		wwwAuth := resp.Header.Get("WWW-Authenticate")
		resp.Body.Close()

		if !strings.HasPrefix(wwwAuth, "Digest") {
			return nil, fmt.Errorf("server requires digest authentication but sent: %s", wwwAuth)
		}

		// Parse digest challenge
		authParams := parseDigestAuth(wwwAuth)

		// Create new request with same body
		var newBody io.Reader
		if body != nil {
			// Need to re-read the body
			if seeker, ok := body.(io.Seeker); ok {
				seeker.Seek(0, io.SeekStart)
				newBody = body
			} else if bodyReader, ok := body.(*bytes.Reader); ok {
				bodyReader.Seek(0, io.SeekStart)
				newBody = bodyReader
			} else {
				return nil, fmt.Errorf("cannot retry request with non-seekable body")
			}
		}

		req, err = http.NewRequest(method, url, newBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticated request: %w", err)
		}

		if contentType != "" && newBody != nil {
			req.Header.Set("Content-Type", contentType)
		}

		// Create digest authorization header
		authHeader := createDigestAuthHeader(c.username, c.password, method, req.URL.RequestURI(), authParams)
		req.Header.Set("Authorization", authHeader)

		// Retry with authentication
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("authenticated request failed: %w", err)
		}
	}

	return resp, nil
}

// parseJSON parses the JSON response body
func parseJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if target == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// parseDigestAuth parses digest authentication parameters from WWW-Authenticate header
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

// createDigestAuthHeader creates a digest authentication header
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

// md5Hash returns MD5 hash of input string
func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}