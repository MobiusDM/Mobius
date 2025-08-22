# GitHub Workflows Analysis Report

## Overview
This document provides a comprehensive review of all GitHub workflows in the Mobius MDM platform repository, analyzing their functionality, dependencies, triggers, and current status.

## Workflow Categories

### 1. Core CI/CD Workflows

#### 1.1 `golangci-lint.yml` - Linting Workflow
- **Purpose**: Static code analysis and linting for Go code
- **Triggers**: 
  - Push to main/patch-*/prepare-* branches (Go files only)
  - Pull requests (Go files or workflow changes)
  - Manual dispatch
- **Matrix Strategy**: Multi-OS (Ubuntu, macOS, Windows)
- **Components Tested**: mobius-server, mobius-cli, mobius-client
- **Key Features**:
  - Uses latest golangci-lint version
  - Platform-specific dependencies (GTK for Ubuntu)
  - Includes cloner-check tool validation
- **Status**: ✅ Active and functional

#### 1.2 `unit-tests.yml` - Unit Testing Workflow
- **Purpose**: Run unit tests across all Go modules with coverage reporting
- **Triggers**: 
  - Push to main/develop (Go files, go.work, go.mod/sum changes)
  - Pull requests to main
  - Manual dispatch
- **Matrix Strategy**: Tests each module independently (mobius-server, mobius-cli, mobius-client, mobius-cocoon, shared)
- **Key Features**:
  - Go workspace synchronization
  - Race condition detection (-race flag)
  - Coverage profile generation and aggregation
  - Comprehensive coverage summary in GitHub step summary
- **Artifacts**: Individual and merged coverage reports
- **Status**: ✅ Active and functional

#### 1.3 `build-and-deploy.yml` - Main Build and Deployment
- **Purpose**: Comprehensive build, test, and deployment pipeline
- **Triggers**: 
  - Push to main/develop
  - Pull requests to main
  - Release events (published)
  - Manual dispatch
- **Architecture**: Multi-stage pipeline (test → build → release)
- **Components**:
  - **Test Stage**: Multi-platform testing (Ubuntu, macOS, Windows)
  - **Build Stage**: Multi-arch Docker image builds (linux/amd64, linux/arm64)
  - **Release Stage**: Cross-platform binary builds and release asset uploads
- **Key Features**:
  - Frontend (Node.js) and backend (Go) integration
  - GHCR (GitHub Container Registry) deployment
  - SBOM (Software Bill of Materials) generation with Syft
  - Container signing with Cosign (keyless)
  - Multi-platform binary releases
- **Security**: ID token and attestation permissions for signing
- **Status**: ✅ Active and comprehensive

#### 1.4 `build-and-deploy-new.yml` - Legacy/Disabled Workflow
- **Purpose**: Older version of build and deploy workflow
- **Status**: 🚫 DISABLED - Explicitly disabled to prevent conflicts
- **Note**: Should be removed if no longer needed

### 2. Security Workflows

#### 2.1 `trivy-scan.yml` - Vulnerability Scanning
- **Purpose**: Security vulnerability scanning using Trivy
- **Triggers**: 
  - Push to main (Terraform files)
  - Pull requests (Terraform files)
  - Manual dispatch
  - Scheduled daily at 4 AM UTC
- **Features**:
  - Custom Trivy DB from AWS ECR
  - SARIF format output
  - GitHub Security tab integration
  - Covers CRITICAL, HIGH, MEDIUM, LOW severities
- **AWS Integration**: Uses AWS credentials for ECR access
- **Status**: ✅ Active with proper security integration

#### 2.2 `dependency-review.yml` - Dependency Security
- **Purpose**: Review dependencies in pull requests for known vulnerabilities
- **Triggers**: Pull requests only
- **Features**:
  - Scans dependency manifest changes
  - Blocks PRs with vulnerable dependencies (if set as required)
  - Uses GitHub's dependency-review-action
- **Status**: ✅ Active and protecting PRs

#### 2.3 `build-and-check-mobiuscli-docker-and-deps.yml` - Docker Security
- **Purpose**: Build and scan Docker dependencies for CLI tools
- **Triggers**: 
  - Manual dispatch
  - Scheduled daily at 6 AM UTC
- **Components Scanned**:
  - ghcr.io/notawar/wix
  - ghcr.io/notawar/bomutils  
  - ghcr.io/notawar/mobiuscli
- **Features**:
  - VEX (Vulnerability Exploitability eXchange) file support
  - Custom Trivy installation and scanning
  - Exit on critical/high vulnerabilities
- **Status**: ✅ Active with advanced vulnerability management

#### 2.4 `code-sign-windows.yml` - Windows Code Signing
- **Purpose**: Attest Windows binaries with GitHub's free signing
- **Type**: Reusable workflow (workflow_call)
- **Features**:
  - GitHub attestation for build provenance
  - DigiCert KeyLocker integration (configured but using GitHub free signing)
  - Artifact download/upload pattern
- **Status**: ✅ Available for use by other workflows

### 3. Release and Distribution Workflows

#### 3.1 `release-helm.yaml` - Helm Chart Release
- **Purpose**: Publish Helm charts to GitHub Pages
- **Triggers**: 
  - Release events (released, not pre-releases)
  - Manual dispatch
- **Features**:
  - Uses stefanprodan/helm-gh-pages action
  - Publishes to charts directory
  - Linting disabled
- **Status**: ✅ Active for Helm chart distribution

### 4. Maintenance and Automation Workflows

#### 4.1 `update-osquery-versions.yml` - Version Automation
- **Purpose**: Automatically update OSQuery version options in UI
- **Triggers**: 
  - Daily scheduled run at midnight UTC
  - Manual dispatch
- **Features**:
  - Python script execution
  - Automatic PR creation with updates
  - Solves issue #21431
- **Dependencies**: Python 3.13.1
- **Status**: ✅ Active automation

#### 4.2 `close-stale-eng-initiated-issues.yml` - Issue Management
- **Purpose**: Automatically manage stale engineering-initiated issues
- **Triggers**: 
  - Daily at 8:10 PM CDT (1:10 AM UTC)
  - Manual dispatch
- **Configuration**:
  - Marks issues stale after 365 days
  - Closes after 14 days of being stale
  - Only affects `~engineering-initiated` labeled issues
  - Processes up to 200 operations per run
- **Status**: ✅ Active for repository maintenance

#### 4.3 `check-updates-timestamps.yml` - TUF Signature Monitoring
- **Purpose**: Monitor TUF (The Update Framework) signature expiration
- **Targets**: https://updates.mobiusmdm.com
- **Triggers**: 
  - Twice daily (10 AM and 10 PM UTC)
  - Manual dispatch
  - PR changes to workflow file
- **Monitored Files**:
  - timestamp.json (4-day warning)
  - snapshot.json (30-day warning)
  - targets.json (30-day warning)
  - root.json (30-day warning)
- **Status**: ✅ Active monitoring for update infrastructure

### 5. Development and Testing Workflows

#### 5.1 `randokiller-go.yml` - Stress Testing
- **Purpose**: Stress test and debug flaky Go tests
- **Triggers**: 
  - Push to branches ending with `-randokiller`
  - Manual dispatch
- **Features**:
  - Configurable test matrix via randokiller.json
  - Multiple test run iterations
  - Comprehensive environment setup (MySQL, Redis, etc.)
  - Failure analysis and summary generation
- **Dependencies**: Docker Compose, ZSH, TLS certificates
- **Status**: ✅ Active debugging tool

#### 5.2 `tfvalidate.yml` - Terraform Validation
- **Purpose**: Validate Terraform configurations
- **Triggers**: 
  - Push to main (Terraform files)
  - Pull requests (Terraform files)
  - Manual dispatch
- **Validated Modules**:
  - ./infrastructure/loadtesting/terraform
  - ./infrastructure/infrastructure/terraform
  - ./infrastructure/dogfood/terraform/aws-tf-module
  - ./terraform (root module)
- **Terraform Version**: 1.9.0
- **Status**: ✅ Active for infrastructure validation

## Workflow Dependencies and Integration

### External Services
- **AWS**: ECR for Trivy databases, IAM roles for credentials
- **GitHub Container Registry (GHCR)**: Docker image storage
- **GitHub Security**: SARIF upload, dependency review
- **DigiCert**: KeyLocker for code signing (configured)
- **Docker Hub**: Various container registries

### Internal Dependencies
- **Go Workspace**: All workflows use go.work for module management
- **Test Infrastructure**: Docker Compose setup for databases and services
- **Frontend Integration**: Node.js build process integrated into deployment
- **Configuration Files**: 
  - .golangci.yml for linting configuration
  - .github/workflows/config/ for workflow-specific configs

## Security Considerations

### Permissions Model
- Most workflows follow least-privilege principle
- Proper separation of read/write permissions
- ID token and attestation permissions for signing workflows

### Secrets Management
- GITHUB_TOKEN for standard operations
- MOBIUS_RELEASE_GITHUB_PAT for network tests
- DigiCert credentials for code signing
- AWS IAM roles for ECR access

### Supply Chain Security
- SBOM generation for all Docker images
- Container signing with Cosign
- Dependency vulnerability scanning
- TUF signature monitoring for updates

## Performance and Optimization

### Concurrency Control
- All workflows implement proper concurrency groups
- Cancel-in-progress enabled to prevent resource waste
- Strategic timing for scheduled workflows to avoid conflicts

### Caching Strategies
- Go module caching in build workflows
- Docker layer caching in build pipelines
- npm cache for frontend builds

### Matrix Strategies
- Multi-OS testing for comprehensive coverage
- Multi-architecture builds for broader deployment
- Module-based testing for parallel execution

## Recommendations

### Immediate Actions
1. **Remove `build-and-deploy-new.yml`** - It's explicitly disabled and should be cleaned up
2. **Add missing config files** - Some workflows reference config files that may be missing
3. **Verify secret availability** - Ensure all referenced secrets are properly configured

### Optimization Opportunities
1. **Workflow consolidation** - Consider combining related workflows to reduce complexity
2. **Enhanced caching** - Implement more aggressive caching strategies
3. **Parallel execution** - Optimize job dependencies for better performance

### Security Enhancements
1. **Regular security audits** - Schedule periodic reviews of workflow permissions
2. **Enhanced monitoring** - Add more comprehensive failure notifications
3. **Supply chain verification** - Consider additional attestation mechanisms

## Conclusion

The Mobius MDM platform has a comprehensive and well-structured CI/CD pipeline with strong security considerations. The workflows cover all essential aspects of modern software development including testing, building, deploying, security scanning, and maintenance automation. The modular approach with clear separation of concerns makes the system maintainable and reliable.

**Overall Status**: ✅ Excellent - The workflow infrastructure is production-ready with comprehensive coverage of all development lifecycle phases.