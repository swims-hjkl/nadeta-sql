APP_NAME := sqlmigrator
DIST_DIR := dist
CONFIG_FILE := .sqlmigratorconfig.json

# .PHONY means even if there is a file with trigger name, don't consider it and run the command
.PHONY: build clean run test fmt help start rmsql

build:
	@echo "building $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	@go build -o "$(DIST_DIR)/$(APP_NAME)"
	@echo "built $(APP_NAME) at $(DIST_DIR)/$(APP_NAME)"

clean:
	@echo "cleaning up..."
	@rm -rf $(DIST_DIR)
	@rm -rf $(CONFIG_FILE)
	@rm -rf "SQLMigrator"
	@echo "done..."

test:
	@echo "Running tests..."
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
