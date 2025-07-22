VERSION ?= $(shell git describe --tags --abbrev=0)
LD_FLAGS = -X 'main.VERSION=$(VERSION)'
APP_NAME := nadeta-sql
DIST_DIR := dist

# .PHONY means even if there is a file with trigger name, don't consider it and run the command
.PHONY: build clean run test fmt help start rmsql

build:
	@echo "building $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -o "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64.exe"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -o "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64"
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -o "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64"
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LD_FLAGS)" -o "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64"
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LD_FLAGS)" -o "$(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64"
	@echo "built $(APP_NAME) at $(DIST_DIR)/$(APP_NAME)"

clean:
	@echo "cleaning up..."
	@rm -rf $(DIST_DIR)
	@echo "done..."

test:
	@echo "Running tests..."
	@go clean -testcache
	@go test ./...

fmt:
	@echo "Formatting..."
	@go fmt ./...

help:
	@echo "Available targets:"
	@echo "  build   - Build the app"
	@echo "  test    - Run unit tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  fmt     - Format code"
	@echo "  help    - Show this help"
