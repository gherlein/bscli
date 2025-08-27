# BrightSign CLI (bscli)

A comprehensive command-line interface and Go library for managing BrightSign players via their Diagnostic Web Server (DWS) API.

## Features

- **Complete DWS API Coverage**: Implements all available Local DWS API endpoints
- **Go Library**: Clean, well-structured Go package for programmatic use
- **CLI Tool**: User-friendly command-line interface
- **Authentication**: Built-in digest authentication support
- **Comprehensive Testing**: Unit tests for all components

## Installation

### Download Binary

Download the latest binary from the releases page, or build from source:

### Build from Source

```bash
git clone <repository-url>
cd bscli
make build
```

### Install to System

```bash
make install
```

## Usage

### CLI Usage

The CLI takes the host as the first argument and will prompt for authentication:

```bash
# Get device information
bscli 192.168.1.100 info device

# List files on SD card
bscli player.local file list /storage/sd/

# Upload a file
bscli 192.168.1.100 file upload local.mp4 /storage/sd/video.mp4

# Reboot the player
bscli 192.168.1.100 control reboot

# Run network diagnostics
bscli 192.168.1.100 diagnostics ping 8.8.8.8
```

### Available Commands

- **info**: Get player information (device, health, time, video-mode, APIs)
- **control**: Player control (reboot, snapshot, DWS settings, firmware)
- **file**: File management (list, upload, download, delete, rename, mkdir, format)
- **diagnostics**: Network diagnostics (ping, DNS, traceroute, interfaces, SSH, telnet)
- **display**: Display control (brightness, contrast, volume, power - Moka displays)
- **registry**: Registry management (get, set, delete, search, recovery URL)
- **logs**: Log management (retrieve logs, supervisor logging)
- **video**: Video output management (modes, EDID, power save, CEC)

### Authentication

The CLI supports several authentication methods:

```bash
# Prompt for password (recommended for security)
bscli 192.168.1.100 info device

# Provide password via flag (not recommended for scripts)
bscli 192.168.1.100 -p mypassword info device

# Custom username (default is 'admin')
bscli 192.168.1.100 -u myuser info device
```

### Debug Mode

Enable debug output to see HTTP requests:

```bash
bscli 192.168.1.100 -d info device
```

### JSON Output

For scripting and automation, use the `--json` flag to get raw JSON output:

```bash
# Human-readable output (default)
bscli 192.168.1.100 info device

# JSON output for scripts
bscli 192.168.1.100 --json info device

# Parse with jq
bscli 192.168.1.100 --json info device | jq '.serial'
```

## Go Library Usage

The `pkg/brightsign` package provides a clean Go API:

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

- Go 1.21 or higher (for building from source)
- BrightSign player with DWS enabled
- Network connectivity to the player

## Testing

Run tests with:

```bash
make test
```

Or run specific tests:

```bash
go test ./pkg/brightsign -v
go test ./internal/cli -v
```

## Development

The project structure:

```
├── cmd/bscli/          # CLI entry point
├── pkg/brightsign/     # Go library package
├── internal/cli/       # CLI implementation
├── bs-api-docs-20250614/ # API documentation
├── Makefile           # Build configuration
└── README.md          # This file
```

## License

See LICENSE.txt for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request
