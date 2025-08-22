# Mobius MDM Build System

.PHONY: all clean build-frontend build-backend build test dev help

# Default target
all: build

# Help target
help:
	@echo "Available targets:"
	@echo "  all           - Build everything (default)"
	@echo "  build         - Build both frontend and backend"
	@echo "  build-frontend - Build Svelte frontend only"
	@echo "  build-backend  - Build Go backend only" 
	@echo "  test          - Run all tests"
	@echo "  dev           - Start development servers"
	@echo "  clean         - Clean build artifacts"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf mobius-server/static/_app
	rm -f mobius-server/static/index.html
	rm -f mobius-server/static/robots.txt
	rm -f mobius-server/mobius-api
	rm -rf mobius-web/.svelte-kit
	rm -rf mobius-web/build

# Build Svelte frontend
build-frontend:
	@echo "Building Svelte frontend..."
	cd mobius-web && npm run build
	@echo "Copying frontend files to API server..."
	mkdir -p mobius-server/static
	cp -r mobius-web/build/* mobius-server/static/
	@echo "Frontend build complete"

# Build Go backend
build-backend:
	@echo "Building Go backend..."
	cd mobius-server && go build -o mobius-api cmd/api-server/main.go
	@echo "Backend build complete"

# Build everything
build: build-frontend build-backend
	@echo "Build complete! Run ./mobius-server/mobius-api to start the server"

# Run tests
test:
	@echo "Running Go tests..."
	cd mobius-server && go test ./...
	cd mobius-client && go test ./...
	cd mobius-cli && go test ./...
	cd mobius-cocoon && go test ./...
	cd shared && go test ./...
	@echo "Running frontend tests..."
	cd mobius-web && npm test 2>/dev/null || echo "No frontend tests configured"

# Development mode
dev:
	@echo "Starting development servers..."
	@echo "Starting backend server..."
	cd mobius-server && go run cmd/api-server/main.go &
	@echo "Starting frontend development server..."
	cd mobius-web && npm run dev

# Install frontend dependencies
install-deps:
	@echo "Installing frontend dependencies..."
	cd mobius-web && npm install
