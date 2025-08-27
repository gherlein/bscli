# BrightSign CLI Examples

This directory contains example programs demonstrating how to use the BrightSign Go library (`pkg/brightsign`) directly in your own applications.

## Example Programs

### 1. basic_usage.go
A comprehensive example that uses environment variables for configuration and provides command-line options:
- Reads configuration from environment variables (same as integration tests)
- Provides multiple command modes
- Demonstrates comprehensive API usage including:
  - Getting device information
  - Checking player health
  - Listing files
  - Running network diagnostics
  - Reading registry values
  - Getting video output information
- Includes proper error handling

## Building the Examples

### Using Make (if working)
```bash
# Build the example program
make example

# Build and run the example
make run-example BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=yourpassword

# Run with specific command
make run-example BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=yourpassword ARGS=diagnostics
```

### Manual Build (Recommended)
```bash
# Build the example program
CGO_ENABLED=0 go build -o basic_usage basic_usage.go
```

## Running the Example

### Basic Usage Example

Set environment variables (same as integration tests):
```bash
export BSCLI_TEST_HOST=192.168.1.100
export BSCLI_TEST_PASSWORD=yourpassword
export BSCLI_TEST_USERNAME=admin       # Optional, defaults to "admin"
export BSCLI_TEST_DEBUG=true           # Optional, enables debug output
export BSCLI_TEST_INSECURE=true        # Optional, use HTTPS with insecure TLS for local certificates
```

Run basic information gathering:
```bash
./basic_usage
```

Run with specific commands:
```bash
./basic_usage list-files    # List files on SD card
./basic_usage diagnostics   # Run network diagnostics
./basic_usage registry      # Test registry operations
./basic_usage video         # Get video output information
```

## Example Output

```
Connecting to BrightSign player at 192.168.1.100 as user 'admin'...

=== Device Information ===
Model: XT1144
Serial: XTC35T000155
Family: Sebring
Firmware: 9.1.66
Boot Version: 8.5.47
Uptime: 3 days, 4:23:45 (276225 seconds)

Network Interfaces:
  eth0 (ethernet/dhcp): 192.168.1.100

=== Player Health ===
Status: healthy
Status Time: 2025-08-26 16:37:37 PST

=== Time Configuration ===
Date: 1724707200
Time: 12:34:56
Timezone: America/Los_Angeles
```

## Using the Library in Your Own Code

To use the BrightSign library in your own Go project:

1. Import the package:
```go
import "bscli/pkg/brightsign"
```

2. Create a client:
```go
client := brightsign.NewClient(brightsign.Config{
    Host:     "192.168.1.100",
    Username: "admin",
    Password: "yourpassword",
    Debug:    false,
    Timeout:  30 * time.Second,
})
```

3. Use the services:
```go
// Get device info
info, err := client.Info.GetInfo()

// List files
files, err := client.Storage.ListFiles("/storage/sd/", nil)

// Run ping test
result, err := client.Diagnostics.Ping("8.8.8.8")

// Get registry value
value, err := client.Registry.GetValue("networking", "hostname")

// Take screenshot
filename, err := client.Control.TakeSnapshot(nil)
```

## Available Services

The library provides the following services:

- **Info**: Device information, health, time, video modes
- **Control**: Reboot, snapshots, DWS configuration
- **Storage**: File operations (list, upload, download, delete)
- **Diagnostics**: Network diagnostics, ping, DNS, traceroute
- **Registry**: Configuration management
- **Logs**: Log retrieval and management
- **Video**: Video output configuration
- **Display**: Display control (for Moka displays)

## Error Handling

All methods return an error as the second value. Always check for errors:

```go
info, err := client.Info.GetInfo()
if err != nil {
    log.Fatalf("Failed to get device info: %v", err)
}
```

## Authentication

The library automatically handles digest authentication. You only need to provide the username and password in the configuration.

## Debug Mode

Enable debug mode to see HTTP requests and responses:

```go
client := brightsign.NewClient(brightsign.Config{
    Debug: true,
    // ... other config
})
```

Or set the environment variables for the example:
```bash
export BSCLI_TEST_DEBUG=true
export BSCLI_TEST_INSECURE=true   # For HTTPS with locally signed certificates
```

## TLS Support

The example supports both HTTP and HTTPS modes:

- **HTTP mode** (default): Used when `BSCLI_TEST_INSECURE` is not set or false
- **HTTPS mode** (insecure): Used when `BSCLI_TEST_INSECURE=true` for players with locally signed certificates

The example will automatically display which protocol it's using when connecting.

## Safety Notes

- The examples avoid destructive operations by default
- Registry operations use safe test keys that are cleaned up
- File operations work in designated test locations
- No reboot or factory reset commands are executed

## Troubleshooting

### Connection Failed
- Verify the player's IP address is correct
- Ensure the player is on the same network
- Check that DWS (Diagnostic Web Server) is enabled on the player

### Authentication Failed
- Verify the password is correct
- Default password is often the serial number
- Username is typically "admin"

### Commands Not Working
- Some commands require specific player models or firmware versions
- Check the player's compatibility with the DWS API version
- Enable debug mode to see the actual HTTP requests and responses