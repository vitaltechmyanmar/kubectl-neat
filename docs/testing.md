# Testing

This project has three test layers. Unit tests run anywhere; e2e and integration tests use
[bats](https://github.com/bats-core/bats-core) and need a live `kubectl`/cluster.

## Test targets

| Target | Command | Requires |
| --- | --- | --- |
| Unit tests | `make test-unit` | Go toolchain |
| E2E (CLI) | `make test-e2e` | bats |
| Integration (kubectl + krew) | `make test-integration` | bats, kubectl, cluster, krew |
| Everything | `make test` | all of the above |

## Unit tests

```bash
make test-unit
# or
go test -v ./...
```

Covers:

- `cmd`: root command (stdin/file), `get` wrap (using a stubbed `kubectl`), and every neat
  step (`metadata`, `scheduler`, `serviceAccount`, `workloadTemplate`, `empty`, `Neat`)
- `pkg/defaults`: default-value detection via Kubernetes object-model defaulting

Fixtures live under `test/fixtures/` (`*-raw.*` inputs, `*-neat.json` expected outputs).

## E2E CLI tests

Requires `bats`:

```bash
make test-e2e
```

`test/e2e-cli.bats` exercises the built binary:

- invalid arguments are rejected (exit code 1 + error message)
- a local file is neat-ed and YAML is emitted

## Integration tests

These run the real `kubectl get` pipeline against a cluster and exercise a real krew install,
so they are only meaningful in a cluster-capable environment:

```bash
make test-integration
```

- `test/e2e-kubectl.bats` — `kubectl neat get` as a plugin (JSON and YAML output)
- `test/e2e-krew.bats` — installs the plugin through krew and runs `get`

Note: `make test-integration` depends on the goreleaser build artifacts (`dist/*.tar.gz`,
`dist/checksums.txt`), so it expects a prior `make goreleaser`.

## CI

`.github/workflows/ci.yml` runs automatically on every push to `main`/`master` and on PRs:

1. `gofmt -s -d` must produce no output
2. `go vet ./...` must pass
3. `make build` must succeed
4. `make test-unit` must pass