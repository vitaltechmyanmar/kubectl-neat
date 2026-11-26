# RFC 0001: Interactive Diff Mode (`--diff`)

- **Status**: In Review / Approved
- **Author(s)**: @sc3pt3r
- **Created**: 2026-11-26
- **Updated**: 2026-11-26
- **Target Release**: v3.1.0
- **Tracking Issue**: #0001

---

## 1. Summary

Add an interactive unified diff mode (`--diff`) to `kubectl-neat`. When enabled, instead of printing only the cleaned manifest, `kubectl-neat` produces a colored unified diff showing exactly which fields, metadata, tolerations, and status entries are stripped.

## 2. Motivation & Problem Statement

When operators and developers clean up Kubernetes manifests before committing them to GitOps repositories (ArgoCD, Flux, Helm), they often need to inspect and verify what `kubectl-neat` actually removed. Today, users must redirect the output to a temporary file and manually run `diff -u` or `git diff`.

Providing a first-class `--diff` flag provides immediate visual feedback and improves confidence that critical annotations, labels, or configurations were not unintentionally pruned.

### User Stories
- **As a** Kubernetes platform engineer,
- **I want to** run `kubectl neat -f pod.yaml --diff` or `kubectl neat get -- pod mypod --diff`,
- **So that** I can instantly inspect the removed boilerplate before committing the neat output.

## 3. Goals & Non-Goals

### Goals
- Provide a `--diff` boolean flag for both standalone CLI (`rootCmd`) and `kubectl get` wrapper mode (`getCmd`).
- Display standard unified diff format (`--- raw`, `+++ neat`, `@@ ... @@`).
- Support ANSI terminal colors (Red for pruned lines, Green for added/retained lines, Cyan for headers) with automatic TTY detection.
- Respect `--output` format (`yaml` or `json`).
- Provide clean exit code behavior (exit code 0 on successful diff inspection).

### Non-Goals
- Interactive patching or line-by-line acceptance (like `git add -p`).
- Replacing external full-cluster diff tools like `kubectl diff`.

## 4. Proposed Technical Design

### CLI Interface
```bash
# Compare a local YAML file against its neat version
kubectl neat -f ./test/fixtures/pod1-raw.yaml --diff

# Stdin pipeline diff
kubectl get pod mypod -o yaml | kubectl neat --diff

# Wrapper mode diff
kubectl neat get -- pod mypod -oyaml --diff
```

### Core Architecture
1. **Normalization**:
   - For a fair diff, both the input and the neat output must be serialized into the same representation (YAML or JSON).
   - If input is YAML, normalize both input and neat output to formatted YAML so differences reflect actual stripped keys rather than whitespace or key-ordering variations.
2. **Unified Diff Engine (`pkg/diff/diff.go`)**:
   - Implement or utilize a lightweight, dependency-free Myers/LCS diff algorithm to generate unified diff hunks.
   - Detect if output is directed to a terminal (TTY) or piped (`TERM != dumb` / file descriptor check).
   - Color lines with ANSI escape sequences when writing to a TTY.
3. **Integration Point (`cmd/cmd.go`)**:
   - Add `showDiff *bool` flag.
   - In `rootCmd.RunE` and `getCmd.RunE`, if `*showDiff` is true, invoke `diff.UnifiedDiff(normalizedIn, out, "raw", "neat")`.

## 5. Backward Compatibility & Edge Cases

- **Existing Workflows**: Default behavior remains unchanged. When `--diff` is omitted, standard neat manifest is emitted to stdout.
- **Identical Manifests**: If an input manifest has zero redundant fields, `--diff` prints an empty diff or informs that the manifest is already neat.
- **Piping / CI**: When stdout is not a TTY (or `NO_COLOR` is set), ANSI escape codes are omitted to ensure clean text output for loggers.

## 6. Testing & Quality Plan

- [ ] **Unit Tests (`pkg/diff/diff_test.go`)**: Test diff generation on identical, modified, and empty inputs.
- [ ] **Command Tests (`cmd/cmd_test.go`)**: Verify `--diff` flag invocation.
- [ ] **Bats E2E CLI Tests (`test/e2e-cli.bats`)**: Verify CLI execution with `--diff` on fixtures.
- [ ] **Documentation**: Document `--diff` in `docs/usage.md` and `Readme.md`.
