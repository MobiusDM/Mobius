# Mobius Device Management

![Mobius logo](Mobius-Logo-Text_1.png)

[![Build & Deploy](https://github.com/NotAwar/Mobius/actions/workflows/build-and-deploy.yml/badge.svg)](https://github.com/NotAwar/Mobius/actions/workflows/build-and-deploy.yml)
[![Unit Tests](https://github.com/NotAwar/Mobius/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/NotAwar/Mobius/actions/workflows/unit-tests.yml)

Mobius is a proprietary platform for managing computers and mobile devices.
It combines osquery-based visibility with Ansible automation to help
organizations monitor and secure their devices.

## Repository Structure

This repository is a Go workspace with separate modules per product:

```text
mobius-server/          # Server-side product
├── cmd/mobius/         # Main server application
├── server/             # Core server implementation
├── api/                # API schemas
├── tools/              # Server-specific tools

mobius-cli/             # Administrative CLI  
├── cmd/mobiuscli/      # CLI application

mobius-client/          # Client / testing utilities
├── cmd/client/

mobius-cocoon/          # Future storefront service
├── cmd/cocoon/

shared/                 # Shared libraries across modules
└── pkg/
```

## Products

### Mobius Server (`mobius-server/`)

The core backend server that provides:

- **Device Management**: osquery orchestration and MDM protocols
- **REST API**: Complete API for device management operations
- **Web Interface**: Admin GUI is a separate React app that talks to the API
- **Security**: Vulnerability scanning and compliance monitoring
- **Multi-tenancy**: Team-based device organization

**Target Environment**: Deployed on servers/cloud infrastructure

### Mobius CLI (`mobius-cli/`)

Command-line interface for:

- **Configuration Management**: GitOps-style device policy management
- **Server Administration**: Remote server management
- **Data Analysis**: Query execution and data export
- **Automation**: Scripting and integration support

**Target Environment**: Administrator workstations and CI/CD pipelines

### Shared Libraries (`shared/`)

Common utilities used by both products:

- Certificate management
- HTTP client libraries  
- File operations
- Cryptographic utilities

## Installation & Usage

Each product can be built and deployed independently:

```bash
# Build server
cd mobius-server && go build -o ../build/mobius ./cmd/mobius

# Build CLI  
cd mobius-cli && go build -o ../build/mobiuscli ./cmd/mobiuscli
```

## Development

The products are designed with clear separation:

- **Server**: Handles device connections, data storage, and management logic
- **CLI**: Provides administrative interface and automation capabilities
- **Shared**: Common code that both products depend on

This structure enables:

- Independent releases and versioning
- Clear product boundaries
- Focused development teams
- Simplified deployment scenarios


## License

Mobius is not open source.
