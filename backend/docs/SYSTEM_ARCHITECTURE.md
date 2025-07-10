# Mobius Backend Architecture Overview

This document explains the overall architecture of the Mobius backend after the restructuring to remove orbit dependencies and implement an API-first, self-hosted device management system.

## 🏗️ Architecture Overview

Mobius is now a **purely API-first** device management platform that follows Go best practices and eliminates the need for custom client agents.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Client Devices                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   macOS     │  │   Windows   │  │    Linux    │     │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │     │
│  │ │Osquery  │ │  │ │Osquery  │ │  │ │Osquery  │ │     │
│  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │     │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │             │     │
│  │ │MDM      │ │  │ │MDM      │ │  │             │     │
│  │ │Profile  │ │  │ │Client   │ │  │             │     │
│  │ └─────────┘ │  │ └─────────┘ │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
                           │
                    HTTPS/REST APIs
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  Mobius Server                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   REST API  │  │   Web UI    │  │    CLI      │     │
│  │  Endpoints  │  │  Interface  │  │  Interface  │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Device    │  │    MDM      │  │Vulnerability│     │
│  │ Management  │  │   Engine    │  │  Scanner    │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   MySQL     │  │    Redis    │  │   Policy    │     │
│  │  Database   │  │    Cache    │  │   Engine    │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
                           │
                    Management APIs
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                 Management Clients                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ mobiuscli   │  │   Web UI    │  │  3rd Party  │     │
│  │ (Command    │  │ (Browser)   │  │    APIs     │     │
│  │   Line)     │  │             │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

## 🔧 Core Components

### 1. **Mobius Server (`cmd/mobius`)**
- **Role**: Central management server
- **Responsibilities**:
  - Provides REST API endpoints for all operations
  - Manages device enrollment and authentication
  - Handles policy distribution and enforcement
  - Processes osquery results and logs
  - Manages MDM profiles and commands
  - Serves web UI interface
  - Coordinates vulnerability scanning

### 2. **Mobius CLI (`cmd/mobiuscli`)**
- **Role**: Command-line management client
- **Responsibilities**:
  - Remote device management via REST APIs
  - Bulk operations and automation
  - Configuration management
  - Live query execution
  - Reporting and data export

### 3. **Device Agents (Client-side)**
- **Osquery Agent**: System monitoring and data collection
- **MDM Clients**: Native OS device management (macOS, Windows, iOS)
- **No Custom Agent**: Uses standard OS tools and protocols

## 🚀 Key Architectural Principles

### **API-First Design**
- All functionality exposed through REST APIs
- Consistent HTTP/JSON interfaces
- Stateless server design
- Authentication via JWT tokens

### **Self-Hosted**
- No external dependencies or cloud services required
- Complete control over data and infrastructure
- Can be deployed on-premises or in private cloud

### **Platform Native**
- Uses OS-built-in MDM capabilities
- Leverages standard osquery installation
- No proprietary client software required

### **Microservices Architecture**
- Modular internal design
- Separate concerns (MDM, osquery, vulnerabilities)
- Scalable and maintainable codebase

## 🔄 Data Flow

### Device Enrollment
1. Device installs osquery with Mobius server configuration
2. Osquery enrolls with server using enrollment secret
3. Server assigns unique node key to device
4. Device begins reporting to server via HTTP APIs

### Policy Distribution
1. Administrator creates policies via web UI or CLI
2. Policies stored in MySQL database
3. Server distributes policies to relevant devices
4. Devices apply policies and report compliance

### Live Queries
1. Administrator initiates live query via CLI/web UI
2. Server queues query for target devices
3. Devices poll server for distributed queries
4. Results returned via HTTP API and stored in database

### MDM Operations
1. Device enrolls in MDM via configuration profile
2. Server manages device certificates and profiles
3. Commands sent via Apple/Microsoft MDM protocols
4. Device status and results reported back to server

## 🛡️ Security Model

### Authentication
- JWT-based authentication for API access
- Device-specific node keys for osquery clients
- MDM certificate-based device authentication

### Authorization
- Role-based access control (RBAC)
- Team-based device isolation
- Granular permission system

### Communication
- TLS encryption for all communications
- Certificate pinning for osquery clients
- Mutual TLS for MDM communications

## 📊 Data Storage

### MySQL Database
- Device inventory and metadata
- Query results and logs
- User accounts and permissions
- Policy definitions and compliance data

### Redis Cache
- Session management
- Real-time query distribution
- Temporary data storage

### File System
- Log files and audit trails
- MDM certificates and profiles
- Vulnerability databases

## 🔌 Integration Points

### External Systems
- **LDAP/Active Directory**: User authentication
- **SIEM Systems**: Log forwarding
- **Vulnerability Databases**: CVE data feeds
- **Cloud Storage**: Backup and archival

### APIs
- **REST APIs**: Primary interface for all operations
- **WebSocket**: Real-time updates and live queries
- **MDM Protocols**: Apple and Microsoft device management

## 📈 Scalability

### Horizontal Scaling
- Multiple server instances behind load balancer
- Shared MySQL and Redis clusters
- Distributed query processing

### Performance Optimization
- Database query optimization
- Redis caching for frequently accessed data
- Asynchronous processing for bulk operations

## 🔍 Monitoring and Observability

### Metrics
- Prometheus metrics export
- Device health and connectivity
- Query performance and success rates

### Logging
- Structured logging with configurable levels
- Audit trails for all administrative actions
- Integration with external log aggregation systems

### Tracing
- OpenTelemetry support
- Distributed tracing for complex operations
- Performance profiling capabilities

## 🏆 Benefits of This Architecture

1. **Simplified Deployment**: No custom agents to maintain
2. **Enhanced Security**: Uses OS-native security features
3. **Better Performance**: Reduced overhead and complexity
4. **Easier Maintenance**: Standard protocols and tools
5. **Improved Reliability**: Fewer moving parts and dependencies
6. **Cost Effective**: Lower resource requirements
7. **Vendor Independence**: No lock-in to proprietary solutions

## 🚧 Migration from Orbit

The restructuring removed all orbit dependencies while maintaining full functionality:

- **Removed**: Orbit client agent, custom protocols, orbit-specific packaging
- **Replaced**: Standard osquery + native MDM clients
- **Maintained**: All device management capabilities, security features, scalability
- **Improved**: Simplified architecture, better performance, easier maintenance

This architecture provides a robust, scalable, and maintainable foundation for enterprise device management while following modern software engineering best practices.
