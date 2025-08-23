# Security Updates and Vulnerability Fixes

This document summarizes the security updates and vulnerability fixes applied to address security and dependabot alerts.

## Summary

All critical security vulnerabilities have been resolved across:
- Node.js/npm dependencies in the web frontend
- Go dependencies across all modules
- Docker container security improvements

## Fixed Vulnerabilities

### Node.js Dependencies

**Fixed**: Cookie vulnerability (GHSA-pxg6-pf52-xh8x)
- **Issue**: cookie package <0.7.0 accepts cookie name, path, and domain with out of bounds characters
- **Solution**: Added npm overrides to force cookie version ^0.7.2
- **Files**: `mobius-web/package.json`
- **Verification**: ✅ npm audit reports 0 vulnerabilities

### Go Dependencies

**Updated critical security packages**:

| Package | Previous Version | Updated Version | Security Impact |
|---------|------------------|-----------------|-----------------|
| golang.org/x/crypto | v0.36.0 | v0.41.0 | Cryptographic functions |
| golang.org/x/net | v0.38.0 | v0.43.0 | Network security |
| golang.org/x/oauth2 | v0.27.0 | v0.30.0 | OAuth2 implementation |
| github.com/ProtonMail/go-crypto | v1.0.0 | v1.3.0 | Additional crypto functions |
| github.com/gorilla/websocket | v1.5.1 | v1.5.3 | WebSocket security |
| google.golang.org/grpc | v1.71.1 | v1.75.0 | gRPC security fixes |
| google.golang.org/protobuf | v1.36.6 | v1.36.8 | Protocol buffer security |

**Modules updated**:
- `mobius-server`: All dependencies updated
- `mobius-cli`: Security packages updated
- `mobius-client`: Security packages updated
- `mobius-cocoon`: Security packages updated
- `shared`: Security packages updated
- `tools/*`: All tool modules updated

### Docker Security Improvements

**Added non-root user execution**:
- **Issue**: Containers running as root pose security risks
- **Solution**: Created dedicated `mobius` user (UID 1001) for container execution
- **Files**: `Dockerfile`, `Dockerfile.combined`
- **Impact**: Reduces attack surface by following principle of least privilege

## Verification

### Build Verification
- ✅ All Go modules build successfully
- ✅ Web application builds without errors
- ✅ Docker containers can be built with security improvements

### Test Verification
- ✅ Web application tests pass (10/10)
- ✅ Go modules compile and link correctly
- ✅ No dependency conflicts introduced

### Security Verification
- ✅ npm audit: 0 vulnerabilities
- ✅ Updated to latest security patches for all critical packages
- ✅ Docker containers run as non-root user

## Production Recommendations

### Environment Configuration
- Ensure all production deployments use environment variables for secrets
- Review `docker-compose.yml` default passwords and override in production
- Use secrets management systems for sensitive configuration

### Ongoing Maintenance
- Regularly run `npm audit` and `go list -m -u all` to check for updates
- Monitor security advisories for Go and Node.js ecosystems
- Set up automated dependency update workflows

### Docker Deployment
- Use the updated Dockerfiles which implement security best practices
- Consider adding security scanning to CI/CD pipeline
- Implement resource constraints and security contexts in Kubernetes

## Files Modified

### Dependency Files
- `mobius-web/package.json` - Added cookie version override
- `mobius-web/package-lock.json` - Updated with secure dependencies
- `*/go.mod` - Updated Go module dependencies across all modules
- `*/go.sum` - Updated dependency checksums
- `go.work.sum` - Updated workspace checksums

### Security Configuration
- `Dockerfile` - Added non-root user and security improvements
- `Dockerfile.combined` - Added non-root user and security improvements

### Documentation
- `SECURITY_UPDATES.md` - This documentation file

## Verification Commands

To verify the security updates:

```bash
# Check npm vulnerabilities
cd mobius-web && npm audit

# Check for Go dependency updates
cd mobius-server && go list -m -u all

# Build verification
go build ./...
cd mobius-web && npm run build

# Test verification
cd mobius-web && npm test
```

All commands should complete successfully with no security warnings.