# Contributing to kubectl-neat

Thank you for taking interest in contributing to kubectl-neat!

## Getting Started

1. Fork and clone the repository
2. Ensure your Go toolchain matches `go` in `go.mod` (1.26+)
3. Run `make build` to verify the build
4. Run `make test` to run all tests

## Issues

- Feel free to open issues for any reason as long as you make it clear if this issue is about a bug/feature/question/comment.
- Please use our structured [GitHub Issue Templates](.github/ISSUE_TEMPLATE/) for Feature Requests, Enhancement Proposals, and Bug Reports.
- Please look for existing issues before you open.
- The issue should clearly explain the reason for opening, the proposal if you have any, and any technical information that's relevant.
- For major architectural changes or new CLI commands, follow our [SDLC Lifecycle](docs/sdlc-lifecycle.md) and submit an [RFC Proposal](docs/proposals/template.md).

## Pull Requests

- Every Pull Request should have an associated Issue unless you are fixing a trivial documentation issue.
- Please follow our [Pull Request Template](.github/pull_request_template.md) and link the associated Issue.
- Describe what the PR does using conventional commit prefixes (e.g. `feat:`, `enhance:`, `fix:`, `docs:`).
- For changes to the neat pipeline or defaulting engine, you **must** include test fixtures in `test/fixtures/` (`*-raw.*` and `*-neat.json`).
- There's no need to add or tag reviewers.
- If a reviewer commented on your code, or asked for changes, please remember to mark the discussion as resolved after you address it. PRs with unresolved issues should not be merged (even if the comment is unclear or requires no action from your side).
- Include tests in Go and/or bats if necessary.
- Refer to [docs/sdlc-lifecycle.md](docs/sdlc-lifecycle.md) for complete details on quality gates and release steps.

## Development

```bash
# Build locally
make build

# Run unit tests
make test-unit

# Run all tests (requires bats for e2e)
make test

# Format code
gofmt -s -w .

# Vet code
go vet ./...
```
