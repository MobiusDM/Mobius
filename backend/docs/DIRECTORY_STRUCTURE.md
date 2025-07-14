# Mobius Backend Directory Structure and Components

This document explains the directory structure of the restructured Mobius backend and details what each component does.

## 📁 Root Directory Structure

```
backend/
├── cmd/                    # Main applications
├── internal/              # Private application code
├── pkg/                   # Public library code
├── api/                   # API definitions and schemas
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
├── tools/                 # Development and operational tools
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md              # Main documentation
```

## 🚀 Applications (`cmd/`)

### `cmd/mobius/` - Main Server Application

**Purpose**: The core Mobius server that provides all device management functionality.

**Key Files**:

- `main.go` - Application entry point and initialization
- `serve.go` - HTTP server setup and routing
- `config_dump.go` - Configuration validation and debugging
- `prepare.go` - Database migration and setup commands
- `cron.go` - Background job scheduling and execution
- `version.go` - Version information and build details
- `vuln_process.go` - Vulnerability scanning coordination

**Responsibilities**:

- HTTP/HTTPS server management
- REST API endpoint handling
- Database connection management
- Background job processing
- MDM protocol handling
- Osquery communication
- Web UI serving

### `cmd/mobiuscli/` - Command Line Interface

**Purpose**: Administrative CLI tool for managing Mobius remotely.

**Key Files**:

- `main.go` - CLI application entry point
- `mobiuscli/mobiuscli.go` - Core CLI functionality and commands
- `mobiuscli/preview.go` - Development/testing server functionality
- `mobiuscli/vulnerability_data_stream.go` - Vulnerability data management

**Available Commands**:

- Device management and querying
- Policy creation and deployment
- Live query execution
- User and team management
- Configuration management
- Bulk operations and automation

## 🏠 Internal Code (`internal/`)

### `internal/server/` - Core Server Logic

**Purpose**: Contains all server-side business logic and implementations.

#### `internal/server/mobius/` - Core Entities and Interfaces

- Device models and management
- User and team structures
- Policy definitions
- Query and result handling
- Core business logic interfaces

#### `internal/server/service/` - Business Logic Services

- Device enrollment and management
- Query execution and scheduling
- Policy evaluation and enforcement
- User authentication and authorization
- Team and permission management

#### `internal/server/datastore/` - Data Persistence

- **`mysql/`** - MySQL database implementation
  - Connection management
  - Query builders and transactions
  - Database migrations
  - Performance optimization
- **`redis/`** - Redis cache implementation
  - Session management
  - Real-time data caching
  - Distributed locking

#### `internal/server/contexts/` - Request Context Management

- Database context propagation
- Error handling and reporting
- Request tracing and logging
- Authentication context

#### `internal/server/config/` - Configuration Management

- Configuration file parsing
- Environment variable handling
- Validation and defaults
- Feature flags

#### `internal/server/mdm/` - Mobile Device Management

- **`apple/`** - Apple MDM implementation (macOS, iOS)
- **`microsoft/`** - Microsoft MDM implementation (Windows)
- **`scep/`** - Certificate enrollment protocols
- **`nanomdm/`** - MDM protocol handling

#### `internal/server/vulnerabilities/` - Security Scanning

- **`nvd/`** - National Vulnerability Database integration
- CVE data processing and analysis
- Vulnerability assessment and reporting
- Security compliance checking

#### `internal/server/authz/` - Authorization

- Role-based access control (RBAC)
- Permission checking and enforcement
- Team-based resource isolation
- API authorization middleware

#### `internal/server/cron/` - Background Jobs

- Scheduled task management
- Vulnerability scanning jobs
- Data cleanup and maintenance
- Report generation

### `internal/server/goose/` - Database Migrations

- Database schema versioning
- Migration execution and rollback
- Schema change management
- Database initialization

## 📚 Public Libraries (`pkg/`)

### Core Utilities

#### `pkg/open/` - Cross-platform File Operations

- File and URL opening across operating systems
- Platform-specific implementations
- Browser and application launching

#### `pkg/file/` - File System Utilities

- File path manipulation
- Directory operations
- Cross-platform file handling
- Temporary file management

#### `pkg/mobiushttp/` - HTTP Client Library

- Standardized HTTP client for Mobius APIs
- Authentication handling
- Request/response processing
- Error handling and retries

#### `pkg/scripts/` - Script Execution

- Cross-platform script running
- Command execution utilities
- Process management
- Output capturing

#### `pkg/secure/` - Security Utilities

- Cryptographic operations
- Certificate management
- Secure random generation
- Key derivation functions

#### `pkg/download/` - Download Management

- File downloading with progress tracking
- Integrity verification
- Resume capability
- Concurrent downloads

### Specialized Libraries

#### `pkg/certificate/` - Certificate Operations

- X.509 certificate parsing and validation
- Certificate chain verification
- Key pair generation
- Certificate signing

#### `pkg/buildpkg/` - Package Building

- Software package creation
- Installer generation
- Cross-platform packaging
- Dependency management

#### `pkg/retry/` - Retry Logic

- Configurable retry strategies
- Exponential backoff
- Circuit breaker patterns
- Error classification

#### `pkg/mdm/` - MDM Utilities

- **`mdmtest/`** - MDM testing utilities and mocks
- Device simulation for testing
- MDM protocol testing helpers

## 🛠️ Development Tools (`tools/`)

### Build and Development Tools

#### `tools/api/` - API Development

- API documentation generation
- OpenAPI specification tools
- API testing utilities

#### `tools/ci/` - Continuous Integration

- Build automation scripts
- Testing infrastructure
- Deployment pipelines

#### `tools/release/` - Release Management

- Version tagging and bumping
- Release note generation
- Binary distribution
- Package publishing

### Operational Tools

#### `tools/backup_db/` - Database Backup

- Automated database backups
- Backup verification
- Restore procedures
- Backup retention policies

#### `tools/dbutils/` - Database Utilities

- Database maintenance scripts
- Performance analysis tools
- Schema comparison utilities
- Data migration helpers

#### `tools/nvd/` - Vulnerability Data Management

- NVD data feed synchronization
- CVE database updates
- Vulnerability data processing
- Feed format conversion

### Testing and Performance Tools

#### `tools/osquery-perf/` - Performance Testing

- Osquery load testing
- Performance benchmarking
- Scalability testing
- Resource utilization monitoring

#### `tools/file-server/` - Development File Server

- Local file serving for development
- Asset hosting for testing
- Mock service implementations

### Integration Tools

#### `tools/terraform/` - Infrastructure as Code

- Infrastructure deployment templates
- Cloud resource provisioning
- Configuration management
- Environment setup automation

#### `tools/webhook/` - Webhook Testing

- Webhook endpoint simulation
- Event payload testing
- Integration testing helpers

## 📋 API Definitions (`api/`)

### `api/schema/` - API Schemas

- **`osquery_mobius_schema.json`** - Osquery table definitions
- API request/response schemas
- Data validation rules
- Documentation specifications

### `api/tables/` - Database Table Definitions

- Osquery table schemas
- Custom table implementations
- Table relationship definitions

## 📖 Documentation (`docs/`)

- **`SYSTEM_ARCHITECTURE.md`** - Overall system architecture
- **`DIRECTORY_STRUCTURE.md`** - This document
- **`README.md`** - Developer documentation
- **`architecture.md`** - Technical architecture details
- **`pkg-readme.md`** - Package documentation
- **`restructuring-summary.md`** - Restructuring notes

## 🔧 Build Scripts (`scripts/`)

- **`update-imports.sh`** - Import path updating during restructuring
- **`fix-remaining-imports.sh`** - Import path fixes
- **`update-version.sh`** - Version management
- Build automation and CI/CD helpers

## 🎯 Key Design Patterns

### Dependency Injection

- Services injected through interfaces
- Testable and mockable components
- Clear separation of concerns

### Repository Pattern

- Data access abstraction
- Database-agnostic interfaces
- Consistent data operations

### Middleware Pattern

- HTTP request/response processing
- Authentication and authorization
- Logging and monitoring

### Event-Driven Architecture

- Asynchronous job processing
- Event publishing and subscription
- Loose coupling between components

## 🚀 Getting Started

### Building the Project

```bash
# Build the server
go build ./cmd/mobius

# Build the CLI
go build ./cmd/mobiuscli

# Run tests
go test ./...

# Install dependencies
go mod tidy
```

### Running the Server

```bash
# Start with default configuration
./mobius serve

# Start with custom config
./mobius serve --config /path/to/config.yml

# Run database migrations
./mobius prepare db
```

### Using the CLI

```bash
# Configure CLI
./mobiuscli config set --address https://your-server.com

# Login
./mobiuscli login

# List devices
./mobiuscli get hosts

# Run live query
./mobiuscli query "SELECT * FROM processes"
```

This directory structure follows Go best practices and provides a clean separation of concerns, making the codebase maintainable, testable, and scalable.
