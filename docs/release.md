# Release

Releases are created by pushing a git tag. GitHub Actions then runs [GoReleaser]
to build platform archives and publish them to GitHub Releases.

## Versioning

This project follows [Semantic Versioning](https://semver.org/) (`MAJOR.MINOR.PATCH`).
Tags are `v`-prefixed, for example `v3.0.0`. The version number is injected into the binary
at build time via the `cmd.Version` ldflag (see `.goreleaser.yml`).

## Prerequisites (local)

For a fully local release you need:

| Tool | Purpose |
| --- | --- |
| `git` | tagging and pushing |
| `goreleaser` | building archives + checksums (snapshot without publish) |
| `yq` and `jq` | `make release` krew-manifest generation |
| GitHub token with `repo` scope | optional, for `make release` before pushing the tag |

CI itself only needs a git tag — the workflow installs GoReleaser.

## Release workflow

### 1. Update version references

Bump the version in `krew-template.yaml` before tagging:

```yaml
spec:
  version: "v3.0.0"          # <-- new version
  homepage: https://github.com/vitaltechmyanmar/kubectl-neat
```

The platform `uri` entries also point at the release download URL:

```
https://github.com/vitaltechmyanmar/kubectl-neat/releases/download/v3.0.0/kubectl-neat_linux_amd64.tar.gz
```

### 2. Update the changelog

Add a section to `CHANGELOG.md` describing the new version's changes (features, fixes,
docs). Keep sections in date order, newest first.

### 3. (Optional) Update supported versions

Keep `SECURITY.md`'s supported-versions table in sync with what you'll support.

### 4. Commit

```bash
git add CHANGELOG.md krew-template.yaml SECURITY.md
git commit -m "release: prepare v3.0.1"
```

### 5. Create and push the tag

```bash
git tag -a v3.0.0 -m "kubectl-neat v3.0.0"
git push origin v3.0.0
```

Pushing the tag triggers `.github/workflows/release.yml`, which runs:

```bash
goreleaser release --clean
```

This builds **linux** and **darwin** for **amd64** and **arm64**, produces the
`*.tar.gz` archives and `checksums.txt`, and publishes them as a GitHub Release.

### 5. Generate the krew manifest (recommended)

After the release is published, regenerate the krew manifest using the published SHA256
checksums:

```bash
make release            # runs goreleaser locally and merges platform manifests into dist/kubectl-neat.yaml
```

The resulting `dist/kubectl-neat.yaml` can then be submitted to the
[krew-index](https://github.com/kubernetes-sigs/krew-index) or hosted in your own krew index.

> `make goreleaser` (without `publish=1`) runs in snapshot mode and does **not** publish.
> `make release` sets `publish=1`, so make sure the tag exists on the remote before running it.

## Local dry-run (snapshot)

Preview what a release would look like without pushing anything:

```bash
make goreleaser            # snapshot build in ./dist, nothing published
```

## CI pipeline summary

| Workflow | Trigger | What it runs |
| --- | --- | --- |
| `.github/workflows/ci.yml` | push to `main`/`master` or PR | gofmt check, `go vet`, `make build`, `make test-unit` |
| `.github/workflows/release.yml` | push of a `v*` tag | `goreleaser release --clean` |

## Checklist

- [ ] `CHANGELOG.md` section added for the new version
- [ ] `krew-template.yaml` version and `uri`s updated
- [ ] `SECURITY.md` supported versions up to date
- [ ] `cmd.Version` will be injected (check the goreleaser ldflag path)
- [ ] tests pass (`make test-unit`)
- [ ] tag created and pushed: `git push origin vMYVERSION`
- [ ] GitHub Release published with archives + `checksums.txt`
- [ ] krew manifest regenerated (`make release`) and submitted/updated in the krew index
- [ ] GitHub Release description summarizes the changes (GoReleaser changelog excludes
      `docs:`/`test:` commits)