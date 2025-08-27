# BrightSign Go Library Usage

The `pkg/brightsign` package provides a clean, well-structured Go API for interacting with BrightSign players via their Diagnostic Web Server (DWS) API.

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"

    "bscli/pkg/brightsign"
)

func main() {
    // Create client
    client := brightsign.NewClient(brightsign.Config{
        Host:     "192.168.1.100",
        Username: "admin",
        Password: "mypassword",
        Debug:    false,
        Insecure: false,  // Set to true for self-signed certificates
        Timeout:  30 * time.Second,
    })

    // Get device information
    info, err := client.Info.GetInfo()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Model: %s\n", info.Model)
    fmt.Printf("Serial: %s\n", info.Serial)
    fmt.Printf("Firmware: %s\n", info.FWVersion)

    // Upload a file
    err = client.Storage.UploadFile("local.mp4", "/storage/sd/video.mp4")
    if err != nil {
        log.Fatal(err)
    }

    // List files
    files, err := client.Storage.ListFiles("/storage/sd/", nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Printf("%s: %d bytes\n", file.Name, file.Size)
    }

    // Run ping diagnostic
    result, err := client.Diagnostics.Ping("8.8.8.8")
    if err != nil {
        log.Fatal(err)
    }

    if result.Success {
        fmt.Printf("Ping successful: %d/%d packets, %.2fms avg\n",
            result.PacketsRecv, result.PacketsSent, result.AvgTime)
    }

    // Take a snapshot
    filename, err := client.Control.TakeSnapshot(nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Snapshot saved: %s\n", filename)
}
```

## Library Services

The library is organized into services that correspond to DWS API categories:

- **Info**: Device information, health, time, video modes
- **Control**: Player control, snapshots, DWS configuration
- **Storage**: File operations, storage management
- **Diagnostics**: Network diagnostics, SSH/telnet configuration
- **Display**: Display control for Moka displays
- **Registry**: Player registry management
- **Logs**: Log retrieval and logging configuration
- **Video**: Video output management, EDID, CEC

## Configuration

### Client Configuration

```go
client := brightsign.NewClient(brightsign.Config{
    Host:     "192.168.1.100",  // Player IP address or hostname
    Username: "admin",          // Username (usually "admin")
    Password: "mypassword",     // Player password
    Debug:    false,           // Enable debug HTTP logging
    Timeout:  30 * time.Second, // HTTP timeout
    Insecure: false,           // Skip TLS certificate verification for local certificates
})
```

### Authentication

The library automatically handles digest authentication. You only need to provide the username and password in the configuration.

### Debug Mode

Enable debug mode to see HTTP requests and responses:

```go
client := brightsign.NewClient(brightsign.Config{
    Debug: true,
    // ... other config
})
```

### TLS Configuration

For players using locally signed certificates (common in newer firmware), enable insecure mode:

```go
client := brightsign.NewClient(brightsign.Config{
    Host:     "192.168.1.100",
    Username: "admin",
    Password: "mypassword",
    Insecure: true,  // Use HTTPS with insecure TLS (skip certificate verification)
    // ... other config
})
```

**Note:** When `Insecure: true` is set, the client automatically uses HTTPS instead of HTTP and skips TLS certificate verification. This is necessary for BrightSign players that use self-signed certificates.

## Error Handling

All methods return an error as the second value. Always check for errors:

```go
info, err := client.Info.GetInfo()
if err != nil {
    log.Fatalf("Failed to get device info: %v", err)
}
```

## Service Examples

### Info Service

```go
// Get device information
info, err := client.Info.GetInfo()

// Get player health
health, err := client.Info.GetHealth()

// Get time information
timeInfo, err := client.Info.GetTime()

// Set time
err = client.Info.SetTime(time.Now())

// Get video mode
videoMode, err := client.Info.GetVideoMode()

// List available APIs
apis, err := client.Info.GetAPIs()
```

### Control Service

```go
// Reboot the player
err := client.Control.Reboot()

// Take a snapshot
filename, err := client.Control.TakeSnapshot(nil)

// Get DWS password status
status, err := client.Control.GetDWSPasswordStatus()

// Enable local DWS
err = client.Control.EnableLocalDWS(true)
```

### Storage Service

```go
// List files
files, err := client.Storage.ListFiles("/storage/sd/", nil)

// Upload a file
err = client.Storage.UploadFile("local.mp4", "/storage/sd/video.mp4")

// Download a file
err = client.Storage.DownloadFile("/storage/sd/video.mp4", "local.mp4")

// Delete a file
err = client.Storage.DeleteFile("/storage/sd/video.mp4")

// Create directory
err = client.Storage.CreateDirectory("/storage/sd/videos/")

// Rename file
err = client.Storage.RenameFile("/storage/sd/old.mp4", "new.mp4")
```

### Diagnostics Service

```go
// Ping test
result, err := client.Diagnostics.Ping("8.8.8.8")

// DNS lookup
result, err := client.Diagnostics.DNSLookup("google.com", false)

// Traceroute
result, err := client.Diagnostics.Traceroute("8.8.8.8")

// Get network interfaces
interfaces, err := client.Diagnostics.GetInterfaces()

// Network diagnostics
diagnostics, err := client.Diagnostics.GetDiagnostics()
```

### Registry Service

```go
// Get registry value
value, err := client.Registry.GetValue("networking", "hostname")

// Set registry value
err = client.Registry.SetValue("networking", "hostname", "myplayer")

// Delete registry value
err = client.Registry.DeleteValue("networking", "hostname")

// Search registry
results, err := client.Registry.Search("hostname")

// Get full registry dump
registry, err := client.Registry.GetRegistry()
```

### Display Service (Moka displays only)

```go
// Get display settings
settings, err := client.Display.GetSettings()

// Set brightness
err = client.Display.SetBrightness(80)

// Set contrast
err = client.Display.SetContrast(50)

// Set volume
err = client.Display.SetVolume(75)

// Power on/off
err = client.Display.SetPower(true)
```

### Logs Service

```go
// Get logs
logs, err := client.Logs.GetLogs(nil)

// Get supervisor logging level
level, err := client.Logs.GetSupervisorLogging()

// Set supervisor logging level
err = client.Logs.SetSupervisorLogging("debug")
```

### Video Service

```go
// Get video output info
output, err := client.Video.GetOutputInfo("hdmi", "0")

// Get EDID information
edid, err := client.Video.GetEDID("hdmi", "0")

// Send CEC command
err = client.Video.SendCEC("hdmi", "0", "power_on")
```

## API Coverage

This implementation covers all Local DWS API endpoints as documented in the BrightSign API documentation:

### Info Endpoints
- `GET /info/` - Basic player information
- `GET /health/` - Player health status
- `GET /time/` - Current time configuration
- `PUT /time/` - Set time
- `GET /video-mode/` - Current video mode
- `GET /` - List all available APIs

### Control Endpoints
- `PUT /control/reboot/` - Reboot player
- `GET /control/dws-password/` - DWS password status
- `PUT /control/dws-password/` - Set/reset DWS password
- `GET /control/local-dws/` - Local DWS status
- `PUT /control/local-dws/` - Enable/disable local DWS
- `POST /snapshot/` - Take screenshot
- `GET /download-firmware/` - Download firmware

### Storage Endpoints
- `GET /files/:path/` - List files/directories
- `POST /files/:path/` - Rename files
- `PUT /files/:path/` - Upload files/create directories
- `DELETE /files/:path/` - Delete files/directories
- `DELETE /storage/:device/` - Format storage

### Diagnostics Endpoints
- `GET /diagnostics/` - Run network diagnostics
- `GET /diagnostics/dns-lookup/:address/` - DNS lookup
- `GET /diagnostics/ping/:ipAddress/` - Ping test
- `GET /diagnostics/trace-route/:address/` - Traceroute
- `GET /diagnostics/network-neighborhood/` - Network neighborhood
- `GET /diagnostics/network-configuration/:interface/` - Network config
- `PUT /diagnostics/network-configuration/:interface/` - Set network config
- `GET /diagnostics/interfaces/` - List interfaces
- Packet capture operations
- SSH/Telnet configuration

### Display Control Endpoints (Moka displays, BOS 9.0.189+)
- `GET /display-control/` - All display settings
- Brightness, contrast, volume controls
- Power management
- Firmware updates
- Display information

### Registry Endpoints
- `GET /registry/` - Full registry dump
- `GET /registry/:section/:key/` - Get registry value
- `PUT /registry/:section/:key/` - Set registry value
- `DELETE /registry/:section/:key/` - Delete registry value
- `DELETE /registry/:section/` - Delete registry section
- Recovery URL management
- Registry flush

### Logs Endpoints
- `GET /logs/` - Get player logs
- `GET /system/supervisor/logging/` - Supervisor logging level
- `PUT /system/supervisor/logging/` - Set logging level

### Video Endpoints
- `GET /video/:connector/output/:device/` - Video output info
- `GET /video/:connector/output/:device/edid/` - EDID information
- Power save management
- Video mode operations
- `POST /sendCecX/` - CEC commands

## Requirements

- Go 1.21 or higher
- BrightSign player with DWS enabled
- Network connectivity to the player

## Safety Notes

- The library provides full access to the DWS API
- Be careful with destructive operations like formatting storage or rebooting
- Always test operations on non-production players first
- Consider implementing your own safety checks for critical operations