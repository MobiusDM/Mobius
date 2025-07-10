# Mobius Backend Architecture

This document describes the restructured Mobius backend following Go best practices.

## Directory Structure

```text
backend/
├── cmd/                     # Main applications and tools
│   ├── mobius/             # Main server application
│   ├── mobiuscli/          # CLI tool for interacting with Mobius API
│   └── tools/              # Utility tools (cpe, cve, macoffice, msrc, osquery-perf)
├── pkg/                    # Public libraries (can be imported externally)
│   ├── certificate/        # Certificate handling
│   ├── mdm/               # MDM-related functionality
│   ├── mobiushttp/        # HTTP client utilities
│   └── ...                # Other reusable packages
├── internal/              # Private application code (cannot be imported externally)
│   └── server/            # HTTP server implementation and business logic
│       ├── config/        # Configuration management
│       ├── datastore/     # Database abstraction layer
│       ├── service/       # HTTP service handlers
│       ├── mdm/          # MDM service implementations
│       └── ...           # Other server components
├── api/                   # API definitions and schemas
│   └── schema/           # OpenAPI/database schemas
├── deployments/          # Deployment configurations
│   ├── charts/           # Helm charts
│   ├── ansible-mdm/      # Ansible playbooks
│   └── it-and-security/  # IT policies and security configs
├── tools/                # Development and operational tools
├── scripts/              # Build and maintenance scripts
└── docs/                 # Documentation
```

## Architecture Principles

### API-First Design

- **Mobius Server**: Exposes REST APIs for all functionality
- **Mobius CLI**: Interacts with server exclusively through APIs
- **No Direct Database Access**: All clients use APIs, no direct DB connections

### Separation of Concerns

- **`cmd/`**: Entry points and CLI parsing only
- **`pkg/`**: Reusable libraries that could be imported by other projects
- **`internal/`**: Application-specific code that should not be imported externally
- **`api/`**: API contracts and schemas

### Import Path Strategy

- Public packages: `github.com/notawar/mobius/pkg/...`
- Internal packages: `github.com/notawar/mobius/internal/...`
- Tools and utilities: `github.com/notawar/mobius/cmd/tools/...`

## Main Components

### Mobius Server (`cmd/mobius`)

The main MDM server that provides:

- REST API endpoints for device management
- Authentication and authorization
- MDM protocol implementations (Apple, Google, Microsoft)
- Policy management
- Vulnerability scanning

### Mobius CLI (`cmd/mobiuscli`)

Command-line interface that:

- Communicates with Mobius server via REST APIs
- Provides device management operations
- Supports configuration management
- Enables automation and scripting

### Development Tools (`cmd/tools`)

Utility tools for:

- CVE data processing (`cve`)
- CPE validation (`cpe`)
- Microsoft security updates (`msrc`)
- Performance testing (`osquery-perf`)
- macOS Office updates (`macoffice`)

## API Design

All functionality is exposed through RESTful APIs:

- Authentication via JWT tokens
- Consistent JSON request/response format
- OpenAPI specifications in `api/schema/`
- Client libraries in `pkg/` for reuse

## Deployment

The restructured codebase supports multiple deployment scenarios:

- **Local Development**: Using docker-compose
- **Kubernetes**: Using Helm charts in `deployments/charts/`
- **Ansible**: Using playbooks in `deployments/ansible-mdm/`
- **Cloud Platforms**: Configuration in `deployments/`
