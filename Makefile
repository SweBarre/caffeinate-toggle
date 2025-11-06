PYTHON := python3
REQUIRED_PYTHON_VERSION := 3.11
VENV_NAME := ct_venv
VENV_PATH := ./$(VENV_NAME)
APP_NAME := Caffeinate\ Toggle
APP_BUNDLE := dist/$(APP_NAME).app
INSTALL_PATH := /Applications/$(APP_NAME).app

# Default target
.DEFAULT_GOAL := build

# ----------------------------------------------
# Helper: check Python version and venv
# ----------------------------------------------
define check_env
	@if [ "$${VIRTUAL_ENV##*/}" != "$(VENV_NAME)" ]; then \
		echo "⚠️  You are not in the $(VENV_NAME) virtual environment."; \
		echo "👉  Run 'source $(VENV_PATH)/bin/activate' and try again."; \
		exit 1; \
	fi; \
	CURRENT_PYTHON_VERSION=$$($(PYTHON) -c 'import sys; print(f"{sys.version_info.major}.{sys.version_info.minor}")'); \
	if [ "$$CURRENT_PYTHON_VERSION" != "$(REQUIRED_PYTHON_VERSION)" ]; then \
		echo "❌ Python version $$CURRENT_PYTHON_VERSION found, but $(REQUIRED_PYTHON_VERSION) is required."; \
		echo "👉  Please recreate your virtual environment with Python $(REQUIRED_PYTHON_VERSION)."; \
		exit 1; \
	fi
endef

.PHONY: dev clean build install

# ----------------------------------------------
# Install dev dependencies in the virtual env
# ----------------------------------------------
dev:
	@$(call check_env)
	@echo "📦 Installing requirements..."
	@pip install -r requirements.txt
	@echo "✅ Development environment ready."

# ----------------------------------------------
# Clean up build artifacts
# ----------------------------------------------
clean:
	@echo "🧹 Cleaning build and dist directories..."
	@rm -rf build dist
	@echo "✅ Clean complete."

# ----------------------------------------------
# Build the macOS app bundle using PyInstaller
# ----------------------------------------------
build: clean dev
	@$(call check_env)
	@echo "🏗️  Building $(APP_NAME)..."
	@pyinstaller \
		--noconfirm \
		--windowed \
		--name "Caffeinate Toggle" \
		--osx-bundle-identifier "nu.rre.caffeinate-toggle" \
		caffeinate_toggle.py
	@/usr/libexec/PlistBuddy -c "Add :LSUIElement bool true" "dist/Caffeinate Toggle.app/Contents/Info.plist" 2>/dev/null || \
		/usr/libexec/PlistBuddy -c "Set :LSUIElement true" "dist/Caffeinate Toggle.app/Contents/Info.plist"
	@echo "✅ Build complete: $(APP_BUNDLE)"

# ----------------------------------------------
# Install the app to /Applications
# ----------------------------------------------
install:
	@echo "📦 Installing $(APP_NAME) to /Applications..."
	@if [ ! -d $(APP_BUNDLE) ]; then \
		echo "❌ Build not found! Run 'make build' first."; \
		exit 1; \
	fi
	@rm -rf "$(INSTALL_PATH)"
	@echo "cp -R $(APP_BUNDLE) /Applications/"
	@echo "✅ Installed to /Applications."
