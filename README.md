# BrightSign CLI (bscli)

## UNOFFICIAL TOOL - NOT AN OFFICIAL BRIGHTSIGN TOOL - PERSONAL PROJECT ONLY

An experimental command-line tool and Go library for managing BrightSign players via their Diagnostic Web Server (DWS) API.

## Features

- **Complete DWS API Coverage**: Implements all available Local DWS API endpoints
- **Go Library**: Clean, well-structured Go package for programmatic use
- **CLI Tool**: User-friendly command-line interface
- **Authentication**: Built-in digest authentication support
- **Comprehensive Testing**: Unit tests for all components

## Installation

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

### TLS/HTTPS Support

For BrightSign players using locally signed certificates (common in newer firmware):

```bash
# Use HTTPS with insecure TLS (skip certificate verification)
bscli 192.168.1.100 --local info device
bscli 192.168.1.100 -l info device

# Can be combined with other flags
bscli 192.168.1.100 --local -p mypassword info device
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

For detailed information about using the Go library programmatically, see [docs/library-use.md](docs/library-use.md).

## Example Program

The `examples/` directory contains a comprehensive example program that demonstrates how to use the Go library:

### Building the Example

```bash
# Manual build (recommended)
cd examples
CGO_ENABLED=0 go build -o basic_usage basic_usage.go

# Or using Make (if working)
make example
```

### Running the Example

Set the required environment variables:

```bash
export BSCLI_TEST_HOST=192.168.1.100
export BSCLI_TEST_PASSWORD=yourpassword
export BSCLI_TEST_USERNAME=admin       # Optional, defaults to "admin"
export BSCLI_TEST_DEBUG=true           # Optional, enables debug output
```

Run basic information gathering:

```bash
./examples/basic_usage
```

Run with specific commands:

```bash
./examples/basic_usage list-files    # List files on SD card
./examples/basic_usage diagnostics   # Run network diagnostics
./examples/basic_usage registry      # Test registry operations
./examples/basic_usage video         # Get video output information
```

The example program demonstrates:
- Device information retrieval
- Player health checking
- File listing and management
- Network diagnostics (ping, DNS lookup)
- Registry operations (get, set, delete with cleanup)
- Video output configuration
- Proper error handling and environment variable usage

For more details, see [examples/README.md](examples/README.md).

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
