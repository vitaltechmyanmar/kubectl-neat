## Description

<!-- Provide a concise summary of the changes introduced by this pull request. -->

Fixes / Closes #(issue)

## Type of Change

- [ ] `feat`: A new feature or capability
- [ ] `enhance`: An enhancement/improvement to existing functionality
- [ ] `fix`: A bug fix
- [ ] `perf`: A performance optimization
- [ ] `refactor`: Code refactoring without behavior changes
- [ ] `docs`: Documentation updates or additions
- [ ] `chore`: Tooling, CI, dependencies, or repository maintenance

## Scope of Impact

- [ ] Core pipeline (`cmd/neat.go`)
- [ ] Defaulting engine (`pkg/defaults/defaults.go`)
- [ ] CLI flags and options (`cmd/cmd.go`)
- [ ] Krew packaging & release (`.goreleaser.yml`, `krew-package.sh`, Makefile)
- [ ] Documentation (`docs/`)

## Quality & SDLC Checklist

- [ ] **Issue Linked**: An approved Issue exists and is linked above.
- [ ] **Tests Added / Updated**:
  - [ ] Unit tests (`make test-unit`)
  - [ ] Bats E2E CLI tests (`make test-e2e`) if CLI flags or UX changed
  - [ ] Fixture files added under `test/fixtures/` (`*-raw.*` and `*-neat.json`) for pipeline changes
- [ ] **Local Verification**:
  - [ ] `gofmt -s -w .` applied
  - [ ] `go vet ./...` passes
  - [ ] `make build` and `make test-unit` pass cleanly
- [ ] **Documentation**:
  - [ ] Updated `docs/usage.md` or relevant documentation
  - [ ] Updated `CHANGELOG.md` under `## [Unreleased]`
- [ ] **Backward Compatibility**:
  - [ ] This change is backwards compatible and does not unexpectedly break existing manifest pruning.
