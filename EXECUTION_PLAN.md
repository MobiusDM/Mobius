# Mobius Modernization and Cleanup Plan

This living plan tracks changes to refactor, simplify, and harden the
repository for API-first operations, small images, and clear licensing.

## Objectives

- Clarify architecture and boundaries (DESIGN.md).
- Remove stale/duplicated artifacts; keep only supported components.
- Harden CI/CD: consistent Go versions, deterministic builds, slim images.
- Define licensing skeleton and wire checks across APIs.
- Ship a minimal agent heartbeat path end-to-end.

## Phase 1 — Baseline hygiene

- Align Go versions across go.work and Dockerfiles (done).
- Fix ShellCheck errors and scripted pitfalls (done).
- Ensure each component has a valid multi-stage Dockerfile (done).
- Add root lint config and markdown CI checks (existing; enforce).

## Phase 2 — API-first and licensing

- Codify API-only usage for orbit/client; remove internal-only code paths.
- Introduce licensing context and enforce feature gates consistently (skeleton
	present; verify and extend where missing).
- Add admin-visible license endpoints (status, apply) behind RBAC.

## Phase 3 — Agent minimal loop

- Keep existing osquery enroll/config endpoints.
- Add documented heartbeat endpoint and rate-limiting.
- Provide sample agent config and minimal systemd/service manifests per OS
	(todo).

## Phase 4 — Cocoon separation

- Ensure Cocoon reads/writes only through API (no direct DB access).
- Tag endpoints and configs specific to Cocoon.

## Phase 5 — Packaging and releases

- Ensure reproducible builds via pinned toolchains.
- Provide docker-compose for local run and Helm/Kustomize for K8s
	(deployments/).
- Automate SBOMs and container signing (todo).

## CI/CD Workflows (added and enforced)

- Lint: golangci-lint across all Go modules on push/PR (existing).
- Unit tests: fast per-module test workflow (unit-tests.yml) on push/PR
	touching Go files (added).
- Build & release: build-and-deploy.yml runs tests, builds multi-arch Docker
	images, and attaches binaries on release (existing, verified).
- Security: dependency review and Trivy scans for container images (existing).

## Housekeeping (continuous)

- Delete/archival of legacy artifacts under root that are superseded by
	mobius-*/ projects.
- Keep README badges and quick-start up to date.

This plan will be kept current as changes land. See commit messages referencing
EXECUTION_PLAN.md updates.
