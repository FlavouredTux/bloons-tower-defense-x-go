# Bloons TDX Go Port - Build System
# Original fan game by Ramaf Party

GMX_DIR = ../gmx_project/Bloons TDX.gmx
ASSETS_DIR = ./assets

.PHONY: all extract build run clean help

help: ## Show this help
	@echo "BTDX Go Port - Build Targets"
	@echo ""
	@echo "  make extract   - Extract all assets from GMX project into ./assets"
	@echo "  make build     - Build the game binary"
	@echo "  make run       - Build and run the game"
	@echo "  make list      - List all rooms in the game"
	@echo "  make clean     - Clean build artifacts"
	@echo ""
	@echo "First-time setup:"
	@echo "  1. make extract"
	@echo "  2. make run"

extract: ## Extract assets from GameMaker project
	@echo "=== Extracting BTDX assets ==="
	go run ./cmd/extract -gmx "$(GMX_DIR)" -out "$(ASSETS_DIR)"

build: ## Build the game binary
	@echo "=== Building BTDX ==="
	go build -o btdx ./cmd/btdx

run: build ## Build and run the game
	./btdx -assets "$(ASSETS_DIR)"

run-room: build ## Run a specific room (usage: make run-room ROOM=Monkey_Meadows_Norm)
	./btdx -assets "$(ASSETS_DIR)" -room "$(ROOM)"

list: build ## List all available rooms
	./btdx -assets "$(ASSETS_DIR)" -list-rooms

clean: ## Clean build artifacts
	rm -f btdx
	go clean ./...

# Development targets
vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test ./...

fmt: ## Format code
	gofmt -w .
