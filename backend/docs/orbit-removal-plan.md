# Orbit Removal Plan - API-First Architecture

## Overview

This document tracks the systematic removal of Orbit dependencies from the
Mobius backend to achieve a pure Go backend with API-first architecture.

## Goals

- **Remove Orbit Agent**: Eliminate all Orbit agent dependencies
- **API-First Design**: Direct communication via REST APIs only
- **Frontend Rebuild**: Separate web application consuming APIs
- **Simplify Architecture**: Go backend + API + Modern frontend

## Progress

### ✅ Phase 1: Orbit Endpoint Removal (Completed)

- [x] Commented out all `/api/mobius/orbit/*` endpoints in handler.go
- [x] Removed orbit enrollment endpoint
- [x] Build verification passed

### ✅ Phase 2: Orbit Service Cleanup (Completed)

- [x] Remove `internal/server/service/orbit.go` (large file with many orbit handlers)
- [x] Create stub implementations for orbit methods in Service interface (`orbit_stubs.go`)
- [x] Fix import paths in stub implementations
- [x] osquery-perf agent.go orbit simulation maintained for load testing
- [x] Build verification passed - all orbit methods stubbed successfully

### ✅ Phase 3: Configuration & Schema Cleanup (Completed)

- [x] Remove `api/schema/tables/orbit_info.yml`
- [x] Remove orbit entries from `api/schema/osquery_mobius_schema.json`
- [x] Remove orbit channel configs from deployment files
- [x] Remove orbit table reference from collect-mobiusdaemon-information.yml query
- [x] Clean up orbit references in documentation

### ✅ Phase 4: Database & Interface Cleanup (Completed)

- [x] Fix SQL syntax errors in `LoadHostByOrbitNodeKey` function
- [x] Remove `directIngestOrbitInfo` function from osquery_utils/queries.go
- [x] Remove orbit_info query definition from osquery table mappings  
- [x] Remove `SetOrUpdateHostOrbitInfo` method implementation
- [x] Create stub implementation for `GetHostOrbitInfo` method
- [x] Add orbit method signatures to Datastore interface for compatibility
- [x] Update mock datastore implementations with orbit stubs
- [x] Remove unused database/sql import from osquery_utils
- [x] Fix database schema references and method calls
- [x] Build verification passed - all orbit database dependencies resolved

### ✅ Phase 5: Final Cleanup & Testing (Completed)

- [x] Remove `orbitEnroll` enum value from enrollment types
- [x] Simplify enrollment logic to remove orbit-specific code paths
- [x] Update `matchHostDuringEnrollment` to remove orbit_node_key column usage
- [x] Clean orbit references in enrollment comments and documentation
- [x] Update MDM integration to use LoadHostByNodeKey instead of LoadHostByOrbitNodeKey
- [x] Ensure GetHostOrbitInfo calls are properly stubbed in microsoft_mdm.go
- [x] Maintain OrbitConfig types for interface compatibility with stubs
- [x] Build verification passed - orbit removal functionally complete

### ✅ Phase 6: Final Verification & Documentation (Completed)

- [x] Remove obsolete orbit references from comments and stubs
- [x] Update agent simulation comments for clarity
- [x] Clean up function parameter names (hostOrbitNodeKey → hostNodeKey)
- [x] Update endpoint middleware to return clear deprecation errors
- [x] Update deployment documentation for orbit-free architecture
- [x] Final build verification passed - all orbit references cleaned

## Files to Remove/Modify

### Core Service Files

- `internal/server/service/orbit.go` - **REMOVE ENTIRE FILE**
- `internal/server/mobius/orbit.go` - **REMOVE ORBIT STRUCTS**
- `internal/server/service/handler.go` - ✅ **ENDPOINTS REMOVED**

### ✅ API Schema Files

- `api/schema/tables/orbit_info.yml` - **REMOVED**
- `api/schema/osquery_mobius_schema.json` - **ORBIT ENTRIES REMOVED**

### ✅ Deployment Configurations

- `deployments/it-and-security/teams/workstations.yml` - **ORBIT CONFIGS REMOVED**
- `deployments/it-and-security/teams/compliance-exclusions.yml` - **ORBIT CONFIGS REMOVED**
- `deployments/it-and-security/teams/workstations-canary.yml` - **ORBIT CONFIGS REMOVED**

### ✅ Testing & Tools

- `cmd/tools/osquery-perf/agent.go` - **ORBIT SIMULATION MAINTAINED FOR TESTING**
- Various tools under `tools/` directory - **CLEANED**

### ✅ Final Problem Resolution

- **Code Issues**: All Go compilation errors resolved
- **Comment Cleanup**: Obsolete orbit references removed from comments
- **Build System**: VS Code tasks.json properly configured for Go builds
- **Module Management**: go.mod and dependencies properly maintained
- **Testing**: All main packages compile and test successfully

## Technical Decisions

### Orbit Endpoint Replacement Strategy

Instead of removing orbit endpoints completely, they have been commented out with clear documentation explaining the architectural change. This approach:

1. **Preserves History**: Shows what endpoints existed before the change
2. **Documentation**: Clear explanation of the API-first transition
3. **Reversibility**: Could theoretically be restored if needed
4. **Safety**: Gradual removal approach reduces risk

### Future Agent Architecture

The new architecture will use:

- **Direct osquery**: osquery communicates directly with Mobius APIs
- **REST APIs**: All functionality exposed via standard HTTP APIs
- **No Agent Layer**: Remove the orbit intermediary layer
- **Modern Frontend**: Separate web application built with modern frameworks

## Next Steps

1. **Remove orbit.go service file**: This is the largest file and contains all orbit handlers
2. **Clean up orbit data structures**: Remove from mobius.go
3. **Update configurations**: Remove orbit references from deployment configs
4. **Test thoroughly**: Ensure no broken references remain

## Risk Assessment

- **Low Risk**: Orbit endpoints were already disabled, no functional impact
- **Medium Risk**: Removing service files may reveal hidden dependencies
- **Mitigation**: Gradual removal with build verification at each step
