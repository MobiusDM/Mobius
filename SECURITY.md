# Security Policy

## Supported Versions

We actively support and provide security updates for the following versions of Mobius:

| Version | Supported          | End of Support |
| ------- | ------------------ | -------------- |
| 1.x.x   | :white_check_mark: | TBD            |
| < 1.0   | :x:                | Immediate      |

**Note**: As Mobius is under active development, we recommend using the latest release for optimal security and features.

## Reporting a Vulnerability

We take security vulnerabilities seriously and appreciate your help in responsibly disclosing security issues.

### How to Report

**For security vulnerabilities, please do NOT create a public issue.** Instead:

1. **Email**: Send details to `security@domain.local` (replace with actual security contact)
2. **GitHub**: Use [GitHub Security Advisory](https://github.com/NotAwar/Mobius/security/advisories/new) (preferred)
3. **Subject Line**: Include "SECURITY" and a brief description

### What to Include

Please provide as much information as possible:

- **Description**: Clear description of the vulnerability
- **Impact**: Potential impact and attack scenarios
- **Reproduction**: Step-by-step instructions to reproduce
- **Environment**: Version, OS, configuration details
- **Proof of Concept**: Code, screenshots, or logs (if safe to share)
- **Suggested Fix**: Any ideas for remediation (optional)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Triage**: Within 5 business days
- **Status Updates**: Weekly until resolution
- **Fix Timeline**: Critical issues resolved within 30 days, others within 90 days

### Disclosure Policy

- We follow responsible disclosure practices
- We'll work with you to understand and resolve the issue
- Public disclosure only after fix is available and deployed
- Credit will be given in security advisories (unless you prefer anonymity)

## Security Features

### Authentication & Authorization

- **JWT-based Authentication**: Secure token-based authentication
- **Role-Based Access Control (RBAC)**: Admin, maintainer, and observer roles
- **API Security**: Bearer token authentication for all API endpoints
- **Session Management**: Configurable token expiration

### Network Security

- **HTTPS/TLS**: All communications encrypted in transit
- **CORS Protection**: Configurable cross-origin resource sharing
- **Security Headers**: HSTS, CSP, X-Frame-Options, and other security headers
- **Rate Limiting**: Protection against brute force and DoS attacks

### Device Management Security

- **Certificate-Based Enrollment**: Secure device enrollment using certificates
- **Policy Enforcement**: Cryptographically signed policy enforcement
- **Command Verification**: Authenticated command execution on managed devices
- **Audit Logging**: Comprehensive audit trails for all security-relevant events

### Data Protection

- **Input Validation**: Comprehensive validation and sanitization
- **SQL Injection Protection**: Parameterized queries and ORM protection
- **Secret Management**: Secure handling of certificates and keys
- **Data Encryption**: Sensitive data encrypted at rest

## Security Best Practices

### Deployment Security

#### Production Environment
```yaml
# Recommended security configuration
environment:
  - MOBIUS_HTTPS_ENABLED=true
  - MOBIUS_TLS_CERT_PATH=/etc/ssl/certs/mobius.crt
  - MOBIUS_TLS_KEY_PATH=/etc/ssl/private/mobius.key
  - MOBIUS_JWT_SECRET_KEY=<strong-random-key>
  - MOBIUS_RATE_LIMIT_ENABLED=true
  - MOBIUS_AUDIT_LOG_ENABLED=true
```

#### Network Security
- Deploy behind reverse proxy (nginx, traefik)
- Use WAF (Web Application Firewall)
- Implement network segmentation
- Enable firewall rules restricting access to necessary ports only

#### Database Security
- Use dedicated database credentials
- Enable database encryption at rest
- Regular database backups with encryption
- Network isolation for database connections

### Configuration Security

#### JWT Configuration
```yaml
jwt:
  secret_key: <256-bit-random-key>  # Use strong random key
  expiration: 24h                   # Adjust based on security requirements
  issuer: "mobius-mdm"             # Set appropriate issuer
```

#### Rate Limiting
```yaml
rate_limiting:
  enabled: true
  requests_per_minute: 100
  burst_size: 20
  blocked_duration: 300s
```

### Operational Security

#### Regular Security Tasks
- [ ] Update dependencies monthly
- [ ] Review access logs weekly
- [ ] Rotate JWT secrets quarterly
- [ ] Update TLS certificates before expiration
- [ ] Review user access permissions monthly
- [ ] Backup configuration and certificates regularly

#### Monitoring & Alerting
- Monitor failed authentication attempts
- Alert on policy enforcement failures
- Track unusual device enrollment patterns
- Monitor certificate expiration dates
- Review audit logs for security events

## Vulnerability Management

### Automated Security Scanning

Our repository includes automated security measures:

- **Trivy Scans**: Daily vulnerability scans of Docker images and dependencies
- **Dependency Review**: Automated review of dependency updates for known vulnerabilities
- **Code Analysis**: Static analysis for security issues
- **Container Scanning**: Multi-layer container security scanning

### Manual Security Review

- Security-focused code reviews for all changes
- Regular penetration testing (recommended annually)
- Security architecture reviews for major releases
- Third-party security audits for enterprise deployments

## Compliance & Standards

### Security Frameworks

Mobius is designed to support compliance with:

- **NIST Cybersecurity Framework**
- **ISO 27001** security controls
- **SOC 2** requirements
- **CIS Controls** for endpoint security

### Privacy Considerations

- Minimal data collection by design
- Configurable data retention policies
- Support for data export and deletion
- GDPR compliance support for EU deployments

## Security Contacts

- **Primary Contact**: security@domain.local (update with actual security email)
- **Maintainer Team**: [@NotAwar](https://github.com/NotAwar)
- **Security Advisory**: [GitHub Security](https://github.com/NotAwar/Mobius/security)
- **Repository Issues**: [GitHub Issues](https://github.com/NotAwar/Mobius/issues) (for non-security issues only)

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [CIS Controls](https://www.cisecurity.org/controls/)
- [MDM Security Best Practices](https://csrc.nist.gov/publications/detail/sp/800-124/rev-1/final)

---

**Last Updated**: 2025-08-22
**Security Policy Version**: 1.0.0

For general questions about Mobius security, please see our [documentation](docs/) or create a [public issue](https://github.com/NotAwar/Mobius/issues/new).