BINARY_NAME=cbak

# Build and install directory
INSTALL_DIR=/usr/local/bin

# Go-related variables
GO=go

# Default target
all: build install

# Build the binary
build:
	@echo "Building the binary..."
	@mkdir -p ./build
	$(GO) build -o ./build/$(BINARY_NAME)

# Install the binary
install: build
	@echo "Installing the binary to $(INSTALL_DIR)..."
	sudo cp ./build/$(BINARY_NAME) $(INSTALL_DIR)
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f ./build/$(BINARY_NAME)

# Run the binary (for testing)
run: build
	@echo "Running the binary..."
	./$(BINARY_NAME)

# Add this target to print debug info
debug:
	@echo "BINARY_NAME: $(BINARY_NAME)"
	@echo "SRC: $(SRC)"
	@echo "INSTALL_DIR: $(INSTALL_DIR)"

# Help message
help:
	@echo "Makefile commands:"
	@echo "  make          - Build and install the binary"
	@echo "  make build    - Build the binary"
	@echo "  make install  - Install the binary"
	@echo "  make clean    - Remove the built binary"
	@echo "  make run      - Build and run the binary"
	@echo "  make help     - Show this help message"

.PHONY: all build install clean run help
