# RFC: [Proposal Title]

- **Status**: [ Draft | In Review | Approved | Implemented | Deferred | Rejected ]
- **Author(s)**: [@username]
- **Created**: [YYYY-MM-DD]
- **Updated**: [YYYY-MM-DD]
- **Target Release**: [e.g. v3.1.0]
- **Tracking Issue**: [#issue-number]

---

## 1. Summary

<!-- A 2-3 sentence executive summary explaining what this feature is and why it is needed. -->

## 2. Motivation & Problem Statement

<!-- Detailed description of user pain points, limitations in kubectl-neat today, or new Kubernetes API developments prompting this change. -->

### User Stories
- **As a** Kubernetes cluster operator / developer,
- **I want to** [action],
- **So that** [value/outcome].

## 3. Goals & Non-Goals

### Goals
- [Clear, measurable objective 1]
- [Objective 2]

### Non-Goals
- [Out of scope item 1]
- [Out of scope item 2]

## 4. Proposed Technical Design

<!-- Detailed architecture and technical proposal. -->

### CLI / API Changes
- New flags or arguments:
  ```bash
  kubectl neat [flags] ...
  ```
- Output changes (JSON/YAML):

### Pipeline & Defaulting Architecture
- Interaction with `cmd/neat.go` steps (1 through 9):
- Interaction with `pkg/defaults/defaults.go` scheme registry:

### Dependency Changes
- Are there new dependencies or Kubernetes API machinery bumps?

## 5. Backward Compatibility & Edge Cases

- **Existing CLI workflows**: Will existing scripts or pipe commands break?
- **Unknown/CRD resources**: How does this handle unregistered GVKs?
- **Performance**: What is the impact on large cluster dumps or multi-megabyte `List` manifests?

## 6. Testing & Quality Plan

- [ ] **Unit Tests**: Specific functions in `cmd/neat_test.go` or `pkg/defaults/defaults_test.go`.
- [ ] **Test Fixtures**: Raw and neat manifest pairs to be added under `test/fixtures/`.
- [ ] **E2E Bats Tests**: New test scenarios in `test/e2e-cli.bats` or `test/e2e-kubectl.bats`.

## 7. Implementation Milestones

1. **Phase 1**: Core logic and unit tests
2. **Phase 2**: CLI flags and e2e integration
3. **Phase 3**: Documentation and user guides
4. **Phase 4**: Release inclusion and changelog entry

## 8. Alternatives Considered

<!-- What other designs were evaluated, and why was this approach chosen? -->
