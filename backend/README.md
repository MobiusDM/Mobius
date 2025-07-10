# Mobius Backend

A modern, API-first Mobile Device Management (MDM) platform built in Go.

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- MySQL 8.0 or later
- Redis 6.0 or later

### Running Locally

1. **Start Dependencies**:

   ```bash
   docker-compose up -d mysql redis
   ```

2. **Configure Environment**:

   ```bash
   cp mobius.yml.example mobius.yml
   # Edit mobius.yml with your configuration
   ```

3. **Initialize Database**:

   ```bash
   go run ./cmd/mobius prepare db
   ```

4. **Start Server**:

   ```bash
   go run ./cmd/mobius serve
   ```

5. **Use CLI**:

   ```bash
   go run ./cmd/mobiuscli --help
   ```

## Architecture

This codebase follows Go best practices with a clear separation between:

- **Public APIs** (`pkg/`): Reusable libraries
- **Internal Logic** (`internal/`): Application-specific code  
- **Applications** (`cmd/`): Server and CLI entry points
- **Deployment** (`deployments/`): Infrastructure as code

See [docs/architecture.md](docs/architecture.md) for detailed information.

## API-First Design

All functionality is exposed through REST APIs:

- **Server** (`cmd/mobius`): Provides REST APIs for all MDM operations
- **CLI** (`cmd/mobiuscli`): Interacts with server exclusively via APIs
- **No Direct DB Access**: All clients use APIs, ensuring consistency

### Example API Usage

```bash
# Authenticate
mobiuscli login --server https://mobius.example.com

# List devices  
mobiuscli get devices

# Apply policy
mobiuscli apply -f policy.yaml
```

## Development

### Project Structure

```text
cmd/           # Applications and tools
├── mobius/    # Main server
├── mobiuscli/ # CLI tool  
└── tools/     # Utilities

pkg/           # Public libraries
internal/      # Private application code
api/           # API definitions
deployments/   # Infrastructure configs
```

### Building

```bash
# Build server
go build -o bin/mobius ./cmd/mobius

# Build CLI
go build -o bin/mobiuscli ./cmd/mobiuscli

# Build all tools
go build -o bin/ ./cmd/tools/...
```

### Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...
```

## Deployment

### Docker

```bash
docker build -t mobius .
docker run -p 8080:8080 mobius
```

### Kubernetes

```bash
helm install mobius ./deployments/charts/mobius
```

### Ansible

```bash
ansible-playbook ./deployments/ansible-mdm/site.yml
```

## Contributing

1. Follow Go best practices
2. All functionality must be API-accessible
3. Update documentation for user-facing changes
4. Add tests for new features

## License

See [OPEN_SOURCE_CREDITS.md](../OPEN_SOURCE_CREDITS.md) for third-party licenses.
