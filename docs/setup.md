# Setup

This guide covers requirements and the supported ways of installing `kubectl-neat`.

## Requirements

| Requirement | Needed for | Notes |
| --- | --- | --- |
| `kubectl` | the `get` subcommand and krew installs | must be on the `PATH` for `kubectl neat get` |
| `krew` | installing as a kubectl plugin | optional if you use a standalone binary |
| Go 1.26+ | building from source | `go` in `go.mod` is `1.26.0` |
| `make` (optional) | convenience build targets | the Makefile wraps `go` commands |

## Install from the krew index

If (and once) this fork is published to the official krew index:

```bash
kubectl krew index add default
kubectl krew install neat
```

## Install from a release manifest (this repository)

The `make release` flow produces a krew manifest at `dist/kubectl-neat.yaml`. Install it directly:

```bash
kubectl krew index add vitaltech https://github.com/vitaltechmyanmar/kubectl-neat.git
kubectl krew install vitaltech/neat
```

Or install a locally generated manifest:

```bash
kubectl krew install --manifest=dist/kubectl-neat.yaml --archive=dist/kubectl-neat_linux_amd64.tar.gz
```

## Download a standalone binary

Grab the latest archive for your platform from the
[releases page](https://github.com/vitaltechmyanmar/kubectl-neat/releases):

- `kubectl-neat_linux_amd64.tar.gz`
- `kubectl-neat_linux_arm64.tar.gz`
- `kubectl-neat_darwin_amd64.tar.gz`
- `kubectl-neat_darwin_arm64.tar.gz`

Each release also ships a `checksums.txt` to verify the download:

```bash
tar -xzf kubectl-neat_linux_amd64.tar.gz
sha256sum -c <(grep linux_amd64 checksums.txt) kubectl-neat
```

The binary is named `kubectl-neat`; place it somewhere on your `PATH`. Renaming it to
`kubectl-neat` is required for standalone CLI usage, and krew installs it as `kubectl-neat`
automatically.

## Build from source

```bash
# method 1: plain go
go build -o kubectl-neat .

# method 2: make target (builds dist/kubectl-neat_<os>_<arch>)
make build

# verify
./kubectl-neat version
```

## Verify the installation

```bash
kubectl neat version
# or
kubectl-neat version
```

When the binary is installed as a kubectl plugin, the command is `kubectl neat`.
As a standalone executable it is `kubectl-neat`. Both forms accept the same flags.