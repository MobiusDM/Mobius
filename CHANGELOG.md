# Changelog

All notable changes to Mobius will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CodeQL security analysis workflow
- Comprehensive issue templates (bug report, feature request, security)
- Automatic labeling system for PRs and issues
- Conventional commit enforcement for PR titles
- SECURITY.md with vulnerability reporting process
- CONTRIBUTING.md with development guidelines
- Enhanced dependabot configuration with security updates enabled

### Changed
- Updated golangci-lint configuration to use current linter versions
- Enabled dependabot dependency updates (was previously disabled)

### Fixed
- Deprecated linter configuration that was blocking CI/CD
- Missing security and workflow files for better repository management

### Security
- Enabled automated dependency security updates via dependabot
- Added CodeQL analysis for vulnerability detection
- Created formal security policy and reporting process

## Previous Releases

*Historical releases and changes will be documented here as the project evolves.*

---

## Categories

- **Added** for new features
- **Changed** for changes in existing functionality  
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** for vulnerability fixes and security improvements