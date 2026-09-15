# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Refactored `Readme.md` into a concise landing page; moved the deep design reference to
  `docs/architecture.md`.
- Corrected the krew install instructions in `docs/setup.md` (the official index `neat`
  installs the upstream project; document installing this fork from a local manifest).

## [3.0.1] - 2026-09-15

### Added
- Remove system-injected tolerations:
  - `DefaultTolerationSeconds`: `node.kubernetes.io/not-ready` and
    `node.kubernetes.io/unreachable` tolerations (`NoExecute`, 300s).
  - `TaintNodesByCondition`: node-condition tolerations (memory/disk/pid pressure,
    unschedulable, network-unavailable).
- Remove `spec.overhead` computed by the `RuntimeClass` admission controller, from Pods and
  workload pod templates.
- `-v`/`--version` flag to print the version, alongside the existing `version` subcommand.
- Full `docs/` folder: setup, usage, release, testing, and architecture guides.

### Fixed
- Removed a no-op `creationTimestamp` in workload pod templates.

## [3.0.0] - 2026-09-15

### Added
- Maintained fork of `itaysk/kubectl-neat` v2.0.4 with modernized dependencies.
- Default-value removal extended to `apps/v1` and `batch/v1` (in addition to `core/v1`).
- `version` subcommand.
- GoReleaser-based release workflow and krew packaging.
- Expanded unit and e2e test coverage.

[Unreleased]: https://github.com/vitaltechmyanmar/kubectl-neat/compare/v3.0.1...HEAD
[3.0.1]: https://github.com/vitaltechmyanmar/kubectl-neat/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/vitaltechmyanmar/kubectl-neat/compare/v2.0.4...v3.0.0