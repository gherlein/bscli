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

.PHONY: all build clean test deps install uninstall help

all: test build

build:
	$(CGO_FLAGS) $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/bscli

build-linux:
	$(CGO_FLAGS) GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_UNIX) ./cmd/bscli

test:
	$(CGO_FLAGS) $(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

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
	@echo "  all       - Run tests and build"
	@echo "  build     - Build the binary"
	@echo "  build-linux - Build Linux binary"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  install   - Install binary to /usr/local/bin (requires sudo)"
	@echo "  uninstall - Remove binary from /usr/local/bin (requires sudo)"
	@echo "  help      - Show this help message"