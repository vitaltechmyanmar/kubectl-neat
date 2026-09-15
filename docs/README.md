# Documentation

This index covers everything you need to work with **kubectl-neat**.

## Contents

| Document | Description |
| --- | --- |
| [Setup](setup.md) | Requirements, installation (krew, binary, source), and verification |
| [Usage](usage.md) | Command reference, flags, modes of operation, and examples |
| [Release](release.md) | Tagging a release and publishing it (GoReleaser + krew) |
| [Testing](testing.md) | Running unit tests, e2e tests, and the CI setup |
| [Architecture](architecture.md) | How kubectl-neat works internally |

## Quick links

- Source: <https://github.com/vitaltechmyanmar/kubectl-neat>
- Releases: <https://github.com/vitaltechmyanmar/kubectl-neat/releases>
- Upstream (maintained fork of): [itaysk/kubectl-neat](https://github.com/itaysk/kubectl-neat)

## Cheatsheet

```bash
# Install as a kubectl plugin (from a local build)
make release
kubectl krew install --manifest=dist/kubectl-neat.yaml \
                     --archive=dist/kubectl-neat_linux_amd64.tar.gz

# Neat a resource (stdin, file, or kubectl pipeline)
kubectl get pod mypod -o yaml | kubectl neat
kubectl neat -f ./my-pod.yaml
kubectl neat get -- pod mypod -oyaml

# Output format: yaml (default), json, or auto-detect the input
kubectl neat -o json -f ./my-pod.yaml

# Version
kubectl neat version
```