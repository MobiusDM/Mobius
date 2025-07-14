# Fixes Required for osquery-perf Tool

This document outlines the remaining issues in the `osquery-perf` tool after the major restructuring of the Mobius backend. These issues should be addressed to make the tool compatible with the new architecture.

## Issues in `agent.go`

1. **Import Path Changes**
   - Replace old import paths with the new internal paths
   - Ensure no duplicate imports exist

2. **Type Compatibility Issues**
   - `mdm.ErrorChain` should be replaced with `mdmtest.ErrorChain`
   - Fix `doDeclarativeManagement` function to accept `*mdmtest.CommandPayload` instead of `*mdm.Command`
   - Update HTTP response handling in declarative management code since the return type is now `map[string]interface{}` not an HTTP response

3. **Missing or Renamed Methods**
   - Add missing stat methods to `Stats` struct in `osquery_perf/stats.go`
   - Update installer functions to work with the new client API
   - Replace `orbitClient` references with equivalent `client` functionality or stub implementations

4. **SyncML Command Type Mismatch**
   - Fix type mismatch between `mobius.SyncMLCmd` and `mdmtest.SyncMLCmd`
   - Update field references to match the new structure

5. **Parameter Type Issues**
   - Fix `installerID` parameter in `SoftwareInstallDetails` - convert from string to uint
   - Fix error handling where variable `err` is undefined

## General Approach

1. For `mdm.ErrorChain`, update the import and rename all references to `mdmtest.ErrorChain`.
2. For declarative management handling, modify functions to handle the new return types.
3. For missing `orbitClient` methods, create stub implementations that log actions instead of executing them.
4. For SyncML commands, update the `AppendResponse` call to use the proper type.
5. For installer functions, parse the string ID to uint and handle errors properly.

The osquery-perf tool requires substantial updates to work with the new architecture. It may be more efficient to rewrite critical sections rather than trying to fix each issue individually.
