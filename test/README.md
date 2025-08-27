# BrightSign CLI Integration Tests

This directory contains comprehensive integration tests for the BrightSign CLI that test against a real BrightSign player.

## Overview

The integration tests verify that every CLI command works correctly with an actual BrightSign player, testing both human-readable output and JSON output modes.

## Test Categories

### 1. **Info Commands**
- `info device` - Device information retrieval
- `info health` - Player health status
- `info time` - Time configuration
- `info video-mode` - Video mode settings
- `info apis` - Available API endpoints

### 2. **File Commands**
- `file list` - Directory listing
- `file upload` - File upload operations
- `file download` - File download operations
- `file delete` - File deletion (with cleanup)

### 3. **Diagnostics Commands**
- `diagnostics run` - Network diagnostics
- `diagnostics ping` - Ping tests
- `diagnostics dns-lookup` - DNS resolution
- `diagnostics interfaces` - Network interface listing

### 4. **Control Commands** (Non-destructive)
- `control dws-password status` - DWS password status
- `control local-dws status` - Local DWS status

### 5. **Registry Commands**
- `registry get-all` - Full registry dump
- `registry set/get/delete` - Registry key operations (with cleanup)

### 6. **Logs Commands**
- `logs get` - Log retrieval
- `logs supervisor get-level` - Supervisor logging level

### 7. **Video Commands** (Safe operations)
- `video output-info` - Video output information (when available)

### 8. **Error Handling Tests**
- Invalid commands
- Invalid file paths
- JSON error formatting

### 9. **JSON Consistency Tests**
- Verifies all commands produce valid JSON with `--json` flag
- Ensures no human-readable text leaks into JSON output

## Running the Tests

### Prerequisites

1. **BrightSign Player**: Access to a BrightSign player on your network
2. **Network Access**: Player must be reachable from your machine
3. **Authentication**: Valid username/password for the player
4. **DWS Enabled**: Player must have Diagnostic Web Server enabled

### Setup

1. Set environment variables:
   ```bash
   export BSCLI_TEST_HOST=192.168.1.100        # Your player's IP or hostname
   export BSCLI_TEST_PASSWORD=yourpassword      # Player password
   export BSCLI_TEST_USERNAME=admin             # Optional, defaults to 'admin'
   ```

2. Run the tests:
   ```bash
   # Using the test script (recommended)
   ./run_integration_tests.sh
   
   # Or run Go tests directly
   go test -v -timeout 10m .
   ```

### Test Script Features

The `run_integration_tests.sh` script provides:

- **Environment validation** - Checks required variables are set
- **User confirmation** - Warns about test operations before running
- **Colored output** - Easy-to-read test results
- **Automatic build** - Builds bscli binary for testing
- **Timeout protection** - 10-minute timeout to prevent hanging

## Test Safety

### Safe Operations
- All read operations (info, diagnostics, registry get)
- File uploads to test locations
- Temporary registry key creation

### Cleanup Operations
- Test files are automatically deleted after upload/download tests
- Temporary registry keys are removed after testing
- No permanent changes are made to the player

### Avoided Operations
- **Reboot commands** - Not tested (destructive)
- **Factory reset** - Not tested (destructive)
- **Firmware updates** - Not tested (risky)
- **System configuration changes** - Not tested (potentially disruptive)

## Example Output

```bash
$ ./run_integration_tests.sh

BrightSign CLI Integration Test Runner
======================================
Test Configuration:
  Host: 192.168.1.100
  Username: admin
  Password: [REDACTED]

WARNING: These tests will:
  - Connect to the BrightSign player
  - Create and delete test files on the SD card
  - Modify registry keys (temporary test keys)
  - Run network diagnostics

Continue with integration tests? (y/N): y

Running Integration Tests...

Starting Go integration tests...
=== RUN   TestInfoCommands
=== RUN   TestInfoCommands/DeviceInfo
=== RUN   TestInfoCommands/Health
=== RUN   TestInfoCommands/Time
=== RUN   TestInfoCommands/VideoMode
=== RUN   TestInfoCommands/APIs
=== RUN   TestFileCommands
=== RUN   TestFileCommands/ListFiles
=== RUN   TestFileCommands/FileUploadDownload
...

✅ All integration tests passed!

Integration test summary:
  - Tested all major CLI commands
  - Verified JSON output consistency
  - Tested file upload/download operations
  - Verified error handling
  - Tested registry operations

Integration tests completed successfully!
```

## Troubleshooting

### Common Issues

1. **Connection failures**
   - Verify player IP/hostname is correct
   - Check network connectivity: `ping <player-ip>`
   - Ensure DWS is enabled on the player

2. **Authentication failures**
   - Verify password is correct
   - Check if default password is the player serial number
   - Ensure username is correct (usually 'admin')

3. **Permission errors**
   - Some operations require administrator access
   - Verify player configuration allows DWS access

4. **Timeout errors**
   - Player may be slow to respond
   - Network latency issues
   - Player may be busy with other operations

### Debug Mode

Run tests with verbose output and debug information:

```bash
BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=password go test -v -timeout 10m .
```

## Adding New Tests

To add tests for new CLI commands:

1. Add the test function to the appropriate test category
2. Test both human-readable and JSON output modes
3. Include cleanup operations for any state changes
4. Add the command to `TestJSONConsistency` if it supports JSON output
5. Update this README with the new test coverage

## Integration with CI/CD

These integration tests can be run in CI/CD pipelines with:

```yaml
# Example GitHub Actions
- name: Run Integration Tests
  env:
    BSCLI_TEST_HOST: ${{ secrets.TEST_PLAYER_HOST }}
    BSCLI_TEST_PASSWORD: ${{ secrets.TEST_PLAYER_PASSWORD }}
  run: |
    cd test
    go test -v -timeout 10m .
```

**Note**: Integration tests require actual hardware and are typically run separately from unit tests.