# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=bscli
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
BUILD_FLAGS=-ldflags "-s -w"
CGO_FLAGS=CGO_ENABLED=0

.PHONY: all build clean test deps install uninstall example run-example help

all: test build

build:
	$(CGO_FLAGS) $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/bscli

build-linux:
	$(CGO_FLAGS) GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_UNIX) ./cmd/bscli

example:
	@echo "Building example program..."
	$(CGO_FLAGS) $(GOBUILD) -o examples/basic_usage ./examples/basic_usage.go
	@echo "Example program built:"
	@echo "  - examples/basic_usage (uses environment variables)"
	@echo ""
	@echo "To run the example:"
	@echo "  export BSCLI_TEST_HOST=192.168.1.100"
	@echo "  export BSCLI_TEST_PASSWORD=yourpassword"
	@echo "  ./examples/basic_usage [command]"

run-example: example
	@echo "Running example program..."
	@if [ -z "$(BSCLI_TEST_HOST)" ] || [ -z "$(BSCLI_TEST_PASSWORD)" ]; then \
		echo "Error: Please set BSCLI_TEST_HOST and BSCLI_TEST_PASSWORD environment variables"; \
		echo "Example: make run-example BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=mypassword"; \
		exit 1; \
	fi
	@./examples/basic_usage $(ARGS)

test:
	$(CGO_FLAGS) $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f examples/basic_usage

deps:
	$(GOMOD) download
	$(GOMOD) tidy

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully!"

uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully!"

help:
	@echo "Available targets:"
	@echo "  all         - Run tests and build"
	@echo "  build       - Build the binary"
	@echo "  build-linux - Build Linux binary"
	@echo "  example     - Build the example program"
	@echo "  run-example - Build and run the example program (requires env vars)"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  install     - Install binary to /usr/local/bin (requires sudo)"
	@echo "  uninstall   - Remove binary from /usr/local/bin (requires sudo)"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Example usage:"
	@echo "  make build                     # Build the CLI"
	@echo "  make example                   # Build the example program"
	@echo "  make run-example BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=pass"
	@echo "  make run-example BSCLI_TEST_HOST=192.168.1.100 BSCLI_TEST_PASSWORD=pass ARGS=diagnostics"