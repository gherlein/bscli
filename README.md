# BrightSign SCP (bscp)

A command-line tool for copying files to BrightSign players using the Diagnostic Web Server (DWS) API. This tool provides an `scp`-like interface for uploading files directly to the SD card storage of BrightSign players over HTTP.

## Purpose

`bscp` enables developers and system administrators to easily transfer files to BrightSign players without requiring SSH access or physical SD card removal. It uses the player's built-in DWS API with HTTP Digest Authentication to securely upload files to the player's storage.

## Features

- **scp-like syntax**: Familiar command-line interface similar to the standard Unix `scp` command
- **HTTP Digest Authentication**: Secure authentication using the DWS password
- **SD Card targeting**: All files are automatically uploaded to the SD card (`/storage/sd`)
- **Upload verification**: Automatically verifies that uploaded files exist on the destination
- **Cross-platform**: Works on Linux, macOS, and Windows
- **IPv4 and hostname support**: Connect using IP addresses or hostnames

## Usage

### Basic Syntax

```bash
bscp <source_file> <host:destination_path>
```

### Examples

Upload a file to the root of the SD card:
```bash
bscp presentation.ppt player.local:/storage/sd/
```

Upload a file to a specific location:
```bash
bscp video.mp4 192.168.1.100:/storage/sd/content/video.mp4
```

Upload an autorun :
```bash
bscp autorun.brs 10.0.1.50:/storage/sd/autorun.brs
```

Upload using relative paths (automatically prefixed with `/storage/sd/`):
```bash
bscp config.json player.local:config.json
# Equivalent to: bscp config.json player.local:/storage/sd/config.json

bscp video.mp4 192.168.1.100:content/videos/
# Equivalent to: bscp video.mp4 192.168.1.100:/storage/sd/content/videos/
```

### Authentication

The tool will prompt for the DWS password when connecting to the player:

```
Password for player.local:
```

Enter the password configured in the player's DWS settings.

### Path Rules

- **Absolute paths** (starting with `/`) are used as-is
- **Relative paths** (not starting with `/`) are automatically prefixed with `/storage/sd/`
- **All files are stored under `/storage/sd`** on the player
- **Directory destinations** (ending with `/`) will use the source filename
- **File destinations** will use the specified filename

## Installation

### Option 1: Using Make

```bash
make install
```

This will build the binary and copy it to `/usr/local/bin/bscp` (requires `sudo`).

### Option 2: Manual Installation

1. Build the binary:
   ```bash
   make build
   ```

2. Copy to your preferred location:
   ```bash
   cp bscp /usr/local/bin/
   # or
   cp bscp ~/bin/
   ```

### Option 3: Run from Source

```bash
go run main.go <source> <host:destination>
```

## Building from Source

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenience)

### Build Commands

#### Using Make (Recommended)

```bash
# Build binary
make build

# Run tests
make test

# Build and test
make all

# Clean build artifacts
make clean

# Install to /usr/local/bin (requires sudo)
make install

# Uninstall from /usr/local/bin (requires sudo)
make uninstall
```

### Cross-compilation

Build for different platforms:

```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bscp-linux .

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bscp.exe .

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bscp-macos .
```

## Technical Details

### API Compatibility

This tool is compatible with BrightSign players running firmware that supports the DWS API (most modern firmware versions). It uses the following API endpoints:

- `PUT /api/v1/files/{storage_device}/` - File upload
- `GET /api/v1/files/{storage_device}/` - Directory listing (for verification)

### Authentication Method

The tool uses HTTP Digest Authentication with the following parameters:
- **Username**: `admin` (fixed)
- **Password**: DWS password configured on the player
- **Algorithm**: MD5 (as per HTTP Digest spec)

### Error Handling

The tool provides detailed error messages for common issues:
- **Connection failures**: Network connectivity problems
- **Authentication errors**: Incorrect password or DWS disabled
- **Path errors**: Invalid destination paths or permissions
- **Upload failures**: Server errors or insufficient storage space

## Troubleshooting

### Common Issues

**"Connection refused"**
- Ensure the player is powered on and connected to the network
- Verify the IP address or hostname is correct
- Check that DWS is enabled on the player

**"401 Unauthorized"**
- Verify the DWS password is correct
- Ensure DWS authentication is enabled on the player

**"Invalid destination format"**
- Destination must include both host and path: `host:/path`
- Path must be absolute (start with `/`)

**"Upload succeeded but file not found"**
- This may occur due to timing issues or storage problems
- Check available space on the SD card
- Verify SD card is properly inserted and functional

### Getting Help

```bash
bscp --help
```

## Development

### Project Structure

```
├── main.go                 # Entry point
├── internal/
│   ├── cli/               # Command-line interface
│   │   ├── cli.go         # CLI logic and argument parsing
│   │   └── cli_test.go    # CLI tests
│   └── dws/               # DWS API client
│       ├── client.go      # HTTP client with digest auth
│       └── client_test.go # Client tests
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── README.md             # This file
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/dws
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Submit a pull request

## License

Apache 2.0.  See LICENSE.txt
