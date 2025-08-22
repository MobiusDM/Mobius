# Copilot Instructions for Mobius MDM Platform

## Repository Overview

Mobius is a modern, API-first Mobile Device Management (MDM) platform designed for self-hosted environments. It provides comprehensive device management, policy enforcement, and application distribution across Windows, macOS, Linux, iOS, and Android devices.

**Key Characteristics:**
- Enterprise MDM solution (not open source)
- Go-based backend with Svelte frontend
- Multi-module monorepo architecture
- Containerized deployment with Docker
- RESTful API with JWT authentication
- Production-ready with comprehensive CI/CD

## Architecture & Structure

### Core Modules (Go Workspace)
```
mobius-server/          # Core API server and business logic
├── api/                # HTTP routing, handlers, middleware
├── pkg/service/        # Business logic implementations  
├── cmd/api-server/     # Standalone API server (current)
└── cmd/mobius/         # Legacy server (deprecated)

mobius-cli/             # Command-line management tool
├── cmd/mobiuscli/      # CLI application
└── pkg/               # CLI business logic

mobius-client/          # Device client agents
├── cmd/client/         # Cross-platform device client
└── pkg/               # Client libraries

mobius-cocoon/          # Enterprise web portal (future)
├── cmd/cocoon/         # Web application server
└── pkg/               # Portal business logic

shared/                 # Common libraries and utilities
└── pkg/               # Shared Go packages (crypto, http, file ops)

mobius-web/             # Svelte frontend application
├── src/               # Svelte source code
├── static/            # Static assets
└── build/             # Build output (copied to mobius-server/static/)
```

### Technology Stack
- **Backend**: Go 1.24.4, RESTful API, JWT auth, CORS, rate limiting
- **Frontend**: Svelte 5, TypeScript, Vite, Vitest for testing
- **Database**: Embedded storage solutions
- **Containerization**: Docker with security hardening
- **CI/CD**: GitHub Actions with 20+ workflows
- **Monitoring**: Health checks, Prometheus metrics, structured logging

## Development Workflow

### Build System
Use the Makefile for common operations:
```bash
make help              # Show available targets
make build            # Build both frontend and backend
make build-frontend   # Build Svelte frontend only
make build-backend    # Build Go backend only
make test             # Run all tests across modules
make dev              # Start development servers
make clean            # Clean build artifacts
```

### Frontend Development (mobius-web/)
```bash
cd mobius-web
npm install           # Install dependencies
npm run dev          # Development server
npm run build        # Production build
npm test             # Run Vitest tests
npm run check        # TypeScript and Svelte checks
```

### Backend Development
```bash
# Each Go module can be built independently
cd mobius-server && go build -o mobius-api cmd/api-server/main.go
cd mobius-cli && go build -o mobiuscli cmd/mobiuscli/main.go

# Run tests
go test ./...         # In any module directory
```

### API Server Quick Start
```bash
cd mobius-server
go run cmd/api-server/main.go
# Server starts on http://localhost:8081
# Default credentials: admin@mobius.local / admin123
```

## Code Standards & Practices

### Go Code Guidelines
- Follow standard Go conventions and formatting
- Use clear package structure with separation of concerns
- Implement proper error handling and logging
- Use dependency injection for testability
- Maintain clean API boundaries between modules

### Frontend Guidelines
- Use TypeScript for type safety
- Follow Svelte 5 best practices
- Implement comprehensive testing with Vitest
- Use semantic commit messages
- Maintain responsive design principles

### Security Requirements
- All dependencies must be kept up-to-date
- Security vulnerabilities must be addressed immediately
- Use secure coding practices (input validation, CORS, etc.)
- Implement proper authentication and authorization
- Follow Docker security best practices (non-root users)

## Testing Strategy

### Go Testing
- Unit tests for all business logic
- Integration tests for API endpoints
- Mock external dependencies
- Use table-driven tests where appropriate

### Frontend Testing
- Component testing with Testing Library
- Unit tests for utility functions
- Integration tests for user workflows
- Visual regression testing where needed

## Deployment & Infrastructure

### Docker Configuration
- `Dockerfile`: Main application container
- `Dockerfile.combined`: Combined services
- `Dockerfile-desktop-linux`: Desktop client
- Security hardening with non-root user (UID 1001)

### Environment Configuration
- Development: Local builds with hot reload
- Production: Containerized deployment
- CI/CD: Automated testing and deployment pipelines

## API Design Principles

### RESTful API Standards
- Use standard HTTP methods and status codes
- Implement consistent error response format
- Use semantic versioning for API versions
- Provide comprehensive OpenAPI 3.1 documentation

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (admin/operator/viewer)
- Secure token handling and refresh mechanisms

## Licensing & Business Context

### Product Tiers
- **Community**: Free tier with basic features
- **Professional**: Advanced features for small/medium teams
- **Enterprise**: Full feature set for large organizations

### Key Features
- Multi-platform device support
- Policy engine and enforcement
- Application distribution
- License management
- Self-hosted deployment

## Common Tasks & Patterns

### Adding New API Endpoints
1. Define handler in `mobius-server/api/handlers/`
2. Add route in `mobius-server/api/router/`
3. Implement business logic in `mobius-server/pkg/service/`
4. Add tests for all layers
5. Update API documentation

### Adding Frontend Features
1. Create Svelte components in `mobius-web/src/`
2. Add TypeScript types
3. Implement API client calls
4. Add comprehensive tests
5. Update navigation and routing

### Security Updates
1. Monitor for dependency vulnerabilities
2. Update packages using npm overrides or go.mod
3. Test thoroughly across all modules
4. Update security documentation
5. Validate with security scanning tools

## CI/CD Workflows

The repository includes comprehensive GitHub Actions workflows:
- Build and deployment pipelines
- Security vulnerability scanning
- Dependency checks and updates
- Multi-platform Docker builds
- Automated testing and quality checks

When contributing, ensure all CI checks pass and follow the established patterns for robust, secure, and maintainable code.

## Troubleshooting Common Issues

### Build Issues
- Ensure Go 1.24.4+ is installed
- Run `go mod tidy` in each module
- Check frontend dependencies with `npm install`
- Use `make clean` to reset build state

### Development Setup
- Use `make dev` for full development environment
- Check API server logs for authentication issues
- Verify frontend builds complete successfully
- Ensure Docker daemon is running for container operations

Remember: This is an enterprise MDM platform focused on security, reliability, and scalability. All changes should maintain these core principles while following established architectural patterns.