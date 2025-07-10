# Mobius Backend

This directory contains all the Go backend components for the Mobius device management platform.

## Directory Structure

```
backend/
├── cmd/           # CLI applications and entry points
│   ├── mobius/    # Main Mobius server application
│   ├── mobiuscli/ # Command-line interface
│   └── ...        # Other CLI tools
├── server/        # Core server logic and services
│   ├── mobius/    # Core Mobius service
│   ├── mdm/       # Mobile Device Management
│   ├── sso/       # Single Sign-On
│   └── ...        # Other server components
├── pkg/           # Shared packages and utilities
│   ├── certificate/   # Certificate management
│   ├── mobiushttp/    # HTTP utilities
│   └── ...            # Other shared packages
├── orbit/         # Agent code for managed devices
│   ├── cmd/       # Agent entry points
│   └── pkg/       # Agent packages
├── tools/         # Development and utility tools
├── go.mod         # Go module definition
└── go.sum         # Go module checksums
```

## Development

### Prerequisites

- Go 1.24.4 or later
- Make

### Building

From the backend directory:

```bash
cd backend

# Build main server
go build -o ../build/mobius ./cmd/mobius

# Build CLI
go build -o ../build/mobiuscli ./cmd/mobiuscli

# Run tests
go test ./...
```

### Import Paths

All imports now use the `backend/` prefix:

```go
import (
    "github.com/notawar/mobius/internal/server/mobius"
    "github.com/notawar/mobius/pkg/certificate"
)
```

## Components

### Main Applications (`cmd/`)

- **mobius**: Main Mobius server
- **mobiuscli**: Command-line interface for Mobius management

### Core Services (`server/`)

- **mobius**: Core platform services
- **mdm**: Mobile Device Management functionality
- **sso**: Single Sign-On integration
- **datastore**: Database abstraction layer
- **vulnerabilities**: Security vulnerability management

### Shared Libraries (`pkg/`)

- **certificate**: X.509 certificate utilities
- **mobiushttp**: HTTP client/server utilities
- **file**: File processing and validation
- **secure**: Security utilities

### Agent Code (`orbit/`)

- Device agent that runs on managed endpoints
- Handles communication with Mobius server
- Executes queries and collects system information

### Development Tools (`tools/`)

- Various utilities for development and operations
- Database management tools
- Testing utilities
- Build and deployment helpers
