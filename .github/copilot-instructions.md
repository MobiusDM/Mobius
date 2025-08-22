# Mobius Mobile Device Management Platform

**CRITICAL: Always follow these instructions first. Only fall back to additional search and context gathering if the information here is incomplete or found to be in error.**

Mobius is a comprehensive Mobile Device Management (MDM) platform written in Go with a Svelte frontend. It provides device management, policy enforcement, and application distribution across Windows, macOS, Linux, iOS, and Android devices.

## Working Effectively

### Bootstrap and Build Process

**CRITICAL BUILD TIMING - NEVER CANCEL COMMANDS:**
- Go dependency download: **NEVER CANCEL** - takes 60-70 seconds. Set timeout to 120+ seconds.
- CLI build: **NEVER CANCEL** - takes 60-75 seconds. Set timeout to 120+ seconds.
- Server build: **NEVER CANCEL** - takes 10-20 seconds. Set timeout to 60+ seconds.
- Frontend build: **NEVER CANCEL** - takes 15-20 seconds. Set timeout to 60+ seconds.
- Complete Makefile build: **NEVER CANCEL** - takes 15-20 seconds total. Set timeout to 60+ seconds.
- Go tests: **NEVER CANCEL** - takes 30-35 seconds. Set timeout to 90+ seconds.

### Required Environment
```bash
# Check versions
go version  # Must be Go 1.24.4+
node --version  # Node.js 20+ required
npm --version   # npm 10+ required
```

### Build Commands (Execute in Order)
```bash
# 1. Sync Go workspace and download dependencies
go work sync  # ~5 seconds

# 2. Download Go dependencies for all modules - NEVER CANCEL: 60-70 seconds
cd mobius-server && go mod download
cd ../mobius-cli && go mod download  
cd ../mobius-client && go mod download
cd ../mobius-cocoon && go mod download
cd ../shared && go mod download
cd ..

# 3. Install frontend dependencies - NEVER CANCEL: 8-10 seconds
cd mobius-web && npm ci && cd ..

# 4. Build frontend - NEVER CANCEL: 15-20 seconds
cd mobius-web && npm run build && cd ..

# 5. Build all Go components - NEVER CANCEL: Total 60-90 seconds
mkdir -p build
cd mobius-server && go build -o ../build/mobius-api ./cmd/api-server  # 10-20 seconds
cd ../mobius-cli && go build -o ../build/mobiuscli ./cmd/mobiuscli    # 60-75 seconds
cd ../mobius-client && go build -o ../build/mobius-client ./cmd/client  # <1 second
cd ../mobius-cocoon && go build -o ../build/mobius-cocoon ./cmd/cocoon  # <1 second
cd ..

# Alternative: Use Makefile for complete build - NEVER CANCEL: 15-20 seconds
make clean && make build
```

### Testing
```bash
# Run Go tests for all modules - NEVER CANCEL: 30-35 seconds
go test -count=1 ./mobius-server/... ./mobius-cli/... ./mobius-client/... ./mobius-cocoon/... ./shared/...

# Run frontend tests - NEVER CANCEL: 2-3 seconds
cd mobius-web && npm test && cd ..

# Run frontend type checking
cd mobius-web && npm run check && cd ..
```

### Running the Application
```bash
# Start the API server (includes frontend)
./build/mobius-api serve --port 8081
# OR from mobius-server directory:
cd mobius-server && ./mobius-api serve --port 8081

# Default credentials:
# Email: admin@mobius.local
# Password: admin123

# Server starts at: http://localhost:8081
```

### CLI Usage
```bash
# Test CLI functionality
./build/mobiuscli --help

# Key CLI commands:
./build/mobiuscli login          # Authenticate with server
./build/mobiuscli get devices    # List devices
./build/mobiuscli get policies   # List policies
./build/mobiuscli apply          # Apply configurations
./build/mobiuscli query          # Run live queries
```

## Validation

### Manual Testing Scenarios
Always test these core workflows after making changes:

1. **API Server Health Check**
```bash
# Start server: ./build/mobius-api serve --port 8081
curl http://localhost:8081/api/v1/health
# Should return: {"status":"healthy",...}
```

2. **Authentication Flow**
```bash
# Login and get token
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@mobius.local","password":"admin123"}'
# Should return: {"token":"token_admin-1_...","user":{...}}
```

3. **Core API Endpoints**
```bash
TOKEN="your_token_here"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/license/status
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/devices
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/policies
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/applications
```

4. **Frontend Integration**
```bash
# Verify frontend is served
curl http://localhost:8081/
# Should return HTML starting with: <!doctype html>
```

5. **CLI Functionality**
```bash
./build/mobiuscli --help  # Should show full command list
./build/mobiuscli login   # Test login flow
```

### CI/CD Validation
Always run these before committing:
```bash
# Frontend validation
cd mobius-web && npm run check  # Type checking
cd mobius-web && npm test       # Unit tests

# Go validation  
go test -count=1 ./...          # All Go tests

# Build validation
make clean && make build        # Complete build process
```

## Repository Structure

```
/
├── mobius-server/          # Core API server
│   ├── cmd/api-server/     # Main server entry point
│   ├── api/                # HTTP handlers and routing
│   ├── pkg/service/        # Business logic
│   └── static/             # Built frontend files (generated)
├── mobius-cli/             # Command-line interface
│   └── cmd/mobiuscli/      # CLI entry point
├── mobius-client/          # Device client agents
├── mobius-cocoon/          # Enterprise web portal
├── mobius-web/             # Svelte frontend application
│   ├── src/                # Frontend source code
│   └── build/              # Built frontend (generated)
├── shared/                 # Common Go libraries
├── tests/                  # Comprehensive test suite
│   ├── test_mdm_functionality.sh    # 29 test scenarios
│   ├── test_websocket_functionality.sh  # 6 scenarios
│   └── run_all_tests.sh    # Test runner
└── .github/workflows/      # CI/CD pipelines
```

## Key Architecture Components

### Backend (Go)
- **mobius-server**: Core MDM server with REST API (builds to ~10MB binary)
- **mobius-cli**: Administrative CLI tool (builds to ~49MB binary) 
- **mobius-client**: Device client agent (builds to ~8.5MB binary)
- **mobius-cocoon**: Enterprise portal service (builds to ~7.9MB binary)
- **shared**: Common libraries used across components

### Frontend (Svelte)
- **mobius-web**: Admin web interface built with SvelteKit
- Built using Vite with TypeScript support
- Static files served by the API server at runtime

### Database & Storage
- Currently uses mock/in-memory services for development
- Designed for PostgreSQL/MySQL in production
- File storage abstraction for various backends

## Troubleshooting

### Common Build Issues
- **Go workspace sync failures**: Run `go work sync` first
- **Frontend build failures**: Ensure Node.js 20+ and run `npm ci` first
- **Binary not found**: Check that build completed in `build/` directory
- **Server won't start**: Verify port 8081 is available

### Performance Issues
- **Slow builds**: Normal - CLI build takes 60-75 seconds, be patient
- **Test timeouts**: Tests can take 30+ seconds, never cancel early
- **Large binaries**: Expected - CLI is 49MB due to embedded dependencies

### Development Tips
- Use `make dev` for development with auto-reload
- Frontend served at `/` when server running
- API endpoints at `/api/v1/*`
- Default admin credentials: `admin@mobius.local` / `admin123`
- WebSocket endpoint available at `/ws` for real-time features

## Docker & Deployment

### Container Builds
```bash
# Build individual component containers
docker build -f mobius-server/Dockerfile .
docker build -f mobius-cli/Dockerfile .
docker build -f mobius-client/Dockerfile . 
docker build -f mobius-cocoon/Dockerfile .
```

### Production Deployment
- Use GitHub Actions workflows in `.github/workflows/`
- Multi-arch builds for linux/amd64 and linux/arm64
- Signed container images with Cosign
- SBOM generation with Syft

## Security Considerations

- JWT-based authentication with role-based access control
- HTTPS/TLS required for production
- Rate limiting and CORS protection built-in
- Security scanning via Trivy in CI/CD
- Secrets management via environment variables

## Important Files and Locations

### Configuration
- `go.work` - Go workspace configuration
- `package.json` - Frontend dependencies in `mobius-web/`
- `Makefile` - Build automation
- `.github/workflows/build-and-deploy.yml` - Main CI/CD pipeline

### Documentation
- `README.md` - Project overview
- `docs/MASTER_PLAN.md` - Comprehensive development plan
- `SECURITY.md` - Security policies and procedures

### Testing
- `tests/run_all_tests.sh` - Comprehensive test runner
- `tests/test_mdm_functionality.sh` - 29 MDM test scenarios
- `tests/test_websocket_functionality.sh` - 6 WebSocket test scenarios

Remember: **ALWAYS** validate changes by running the complete build and test process. The platform is mission-critical infrastructure - never skip validation steps.