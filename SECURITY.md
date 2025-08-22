# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| Previous minor | :white_check_mark: |
| Older versions | :x: |

## Reporting a Vulnerability

The Mobius team takes security bugs seriously. We appreciate your efforts to responsibly disclose your findings.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by emailing: **security@mobius.local**

Include the following information in your report:
- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

### What to Expect

- We will acknowledge your email within 48 hours
- We will send a more detailed response within 7 days indicating the next steps
- We will keep you informed of our progress toward a fix and full announcement
- We may ask for additional information or guidance

### Disclosure Policy

- We ask that you give us reasonable time to investigate and mitigate an issue before making any information public
- We will make a good faith effort to resolve security issues in a timely manner
- We will publicly acknowledge your responsible disclosure, unless you prefer to remain anonymous

## Security Best Practices

### For Contributors

- Never commit secrets, API keys, or credentials to the repository
- Use environment variables or secure secret management for sensitive data
- Follow the principle of least privilege in code design
- Validate all inputs and sanitize outputs
- Use parameterized queries to prevent SQL injection
- Keep dependencies updated and monitor for vulnerabilities

### For Deployments

- Always use HTTPS/TLS in production
- Implement proper authentication and authorization
- Use strong, unique passwords and enable MFA where possible
- Regularly update and patch all systems
- Monitor logs for suspicious activity
- Implement network segmentation and firewalls
- Regular security assessments and penetration testing

## Security Features

Mobius implements several security features:

- **Input Validation**: Comprehensive input sanitization and validation
- **Authentication**: Secure token-based authentication
- **Authorization**: Role-based access control (RBAC) with Open Policy Agent
- **Encryption**: TLS encryption for all communications
- **Audit Logging**: Comprehensive audit trails for security events
- **Vulnerability Scanning**: Automated dependency and container scanning
- **Secret Management**: Secure handling of sensitive configuration data

## Security Scanning

This repository includes automated security scanning:

- **Dependency Scanning**: Dependabot and dependency-review-action
- **Vulnerability Scanning**: Trivy scanner for containers and code
- **Code Quality**: golangci-lint with security-focused rules
- **OSV Scanner**: Google's Open Source Vulnerability scanner

Results are automatically uploaded to GitHub Security tab for tracking and remediation.