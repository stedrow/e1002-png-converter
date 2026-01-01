.PHONY: build run clean install test help

BINARY_NAME=e1002-convert
BUILD_DIR=bin
INPUT_DIR=./input
OUTPUT_DIR=./output

help:
	@echo "E1002 PNG Converter - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the binary"
	@echo "  make install        - Install binary to /usr/local/bin"
	@echo "  make run INPUT=file - Convert a PNG file"
	@echo "  make test           - Run a test conversion"
	@echo "  make clean          - Remove built binaries"
	@echo ""
	@echo "Example:"
	@echo "  make run INPUT=input/myimage.png"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

run: build
	@if [ -z "$(INPUT)" ]; then \
		echo "Error: INPUT not specified."; \
		echo "Usage: make run INPUT=input/myimage.png"; \
		exit 1; \
	fi
	@mkdir -p $(OUTPUT_DIR)
	@echo "Converting $(INPUT)..."
	./$(BUILD_DIR)/$(BINARY_NAME) "$(INPUT)"

test: build
	@if [ ! -f "$(INPUT_DIR)/wallpaper-sample.png" ]; then \
		echo "Error: test image not found in $(INPUT_DIR)/"; \
		echo "Please add a wallpaper-sample.png file to the input directory"; \
		exit 1; \
	fi
	@echo "Running test conversion..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(INPUT_DIR)/wallpaper-sample.png
	@echo "✓ Test complete! Check output/wallpaper-sample_dithered.png"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

# Development commands
dev-deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

# Cross-compilation targets
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Built macOS binaries"

build-all: build-linux build-windows build-darwin
	@echo "✓ Built all platform binaries"
