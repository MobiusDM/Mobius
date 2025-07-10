# GitHub Actions Workflow Cleanup Plan

## Summary
After analyzing all 23 GitHub Actions workflows in the repository, I've identified workflows that need to be removed, updated, or kept as-is based on the new API-first, self-hosted architecture without orbit dependencies.

## Workflows to Remove (Orbit-dependent)

### 1. `generate-desktop-targets.yml`
- **Purpose**: Generates Mobius Desktop targets for Orbit
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: References "orbit-*" tags and builds desktop applications for Orbit

### 2. `mobiusdaemon-tuf.yml`
- **Purpose**: Updates TUF mobiusdaemon components
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: References orbit/TUF.md and orbit/old-TUF.md files

### 3. `verify-mobiusdaemon-base.yml`
- **Purpose**: Verifies mobiusdaemon-base files
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: Checks mobiusdaemon-base.msi and .pkg files used for orbit distribution

### 4. `mobiuscli-preview.yml`
- **Purpose**: Tests mobiuscli preview command
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: References orbit.log and orbit directory, expects orbit to be running

### 5. `integration.yml`
- **Purpose**: Tests enrolling agents using orbit
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: Tests "mobiuscli package" command and orbit agent enrollment

### 6. `check-tuf-timestamps.yml`
- **Purpose**: Checks TUF signatures for orbit distribution
- **Status**: ❌ REMOVE - Orbit-dependent
- **Reason**: Checks tuf.mobiuscli.com which is used for orbit distribution

## Workflows to Review/Update

### 7. `generate-osqueryd-targets.yml`
- **Purpose**: Generates osqueryd targets for MobiusDaemon
- **Status**: ❓ REVIEW - May be obsolete
- **Reason**: No Makefile targets found, may have been for orbit packaging
- **Action**: Remove unless osqueryd binaries are still needed for API-first architecture

### 8. `generate-nudge-targets.yml`
- **Purpose**: Generates Nudge targets for Mobiusd
- **Status**: ❓ REVIEW - May be obsolete
- **Reason**: No Makefile targets found, may have been for orbit packaging
- **Action**: Remove unless nudge binaries are still needed for API-first architecture

### 9. `build-and-deploy.yml`
- **Purpose**: Build and deploy Mobius backend
- **Status**: ⚠️ UPDATE - Needs path fixes
- **Reason**: References `./cmd/mobius` but should reference `./backend/cmd/mobius`
- **Action**: Update build paths to reflect new directory structure

## Workflows to Keep (Core functionality)

### 10. `golangci-lint.yml`
- **Purpose**: Go linting and code quality
- **Status**: ✅ KEEP - Essential for code quality

### 11. `trivy-scan.yml`
- **Purpose**: Security vulnerability scanning
- **Status**: ✅ KEEP - Essential for security

### 12. `build-and-check-mobiuscli-docker-and-deps.yml`
- **Purpose**: Build mobiuscli container and check vulnerabilities
- **Status**: ✅ KEEP - Essential for CLI distribution

### 13. `code-sign-windows.yml`
- **Purpose**: Windows binary attestation
- **Status**: ✅ KEEP - Essential for Windows support

### 14. `update-osquery-versions.yml`
- **Purpose**: Update osquery versions in UI
- **Status**: ✅ KEEP - Essential for osquery integration

### 15. `scorecards-analysis.yml`
- **Purpose**: Security scorecard analysis
- **Status**: ✅ KEEP - Essential for security compliance

### 16. `dependency-review.yml`
- **Purpose**: Dependency security review
- **Status**: ✅ KEEP - Essential for security

### 17. `randokiller-go.yml`
- **Purpose**: Random number security testing
- **Status**: ✅ KEEP - Essential for security

### 18. `tfvalidate.yml`
- **Purpose**: Terraform validation
- **Status**: ✅ KEEP - Essential for infrastructure

### 19. `close-stale-eng-initiated-issues.yml`
- **Purpose**: Issue management
- **Status**: ✅ KEEP - Essential for maintenance

### 20. `ingest-maintained-apps.yml`
- **Purpose**: App ingestion
- **Status**: ✅ KEEP - Essential for app management

### 21. `check-updates-timestamps.yml`
- **Purpose**: General update checking
- **Status**: ✅ KEEP - Essential for maintenance

### 22. `check-vulnerabilities-in-released-docker-images.yml`
- **Purpose**: Docker security scanning
- **Status**: ✅ KEEP - Essential for security

### 23. `update-certs.yml`
- **Purpose**: Certificate updates
- **Status**: ✅ KEEP - Essential for security

## Implementation Plan

1. **Remove orbit-dependent workflows** (6 workflows)
2. **Review and potentially remove** osqueryd/nudge workflows (2 workflows)
3. **Update build-and-deploy workflow** to fix paths
4. **Test remaining workflows** to ensure they work with new structure

## Expected Benefits

- **Reduced complexity**: Remove 6-8 obsolete workflows
- **Better maintainability**: Keep only workflows relevant to API-first architecture
- **Improved CI/CD**: Fix build paths and ensure workflows work with new structure
- **Enhanced security**: Keep all security-related workflows intact
