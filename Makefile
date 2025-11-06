APP_NAME = CaffeinateToggle
APP_DIR = $(HOME)/Library/Application\ Support/$(APP_NAME)
AGENT_DIR = $(HOME)/Library/LaunchAgents
AGENT_PLIST = $(AGENT_DIR)/nu.rre.caffeinate-toggle.plist


init:
	go mod init github.com/SweBarre/caffeinate-toggle
	go mod tidy

build:
	go build -o dist/$(APP_NAME) ./cmd/caffeinate-toggle

install: build
	@echo "Installing $(APP_NAME)..."
	@mkdir -p $(APP_DIR)
	@cp dist/$(APP_NAME) $(APP_DIR)/
	@chmod +x $(APP_DIR)/$(APP_NAME)
	@mkdir -p $(AGENT_DIR)
	@cp nu.rre.caffeinate-toggle.plist $(AGENT_PLIST)
	@sed -i '' 's|__BINARY_PATH__|$(APP_DIR)/$(APP_NAME)|g' $(AGENT_PLIST)
	@launchctl load -w $(AGENT_PLIST)
	@echo "$(APP_NAME) installed and set to start on login."

uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@launchctl unload -w $(AGENT_PLIST) || true
	@rm -f $(AGENT_PLIST)
	@rm -rf $(APP_DIR)
	@echo "$(APP_NAME) removed."


run:
	go run ./cmd/caffeinate-toggle

clean:
	rm -rf dist
