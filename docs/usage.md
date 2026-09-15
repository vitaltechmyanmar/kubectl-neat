# Usage

`kubectl-neat` reads a Kubernetes resource (YAML or JSON), removes the clutter that Kubernetes
adds to it, and prints a readable version.

There are two input modes plus two subcommands.

## Commands

```
kubectl-neat                Neat a resource from a file or stdin
kubectl-neat get <args...>  Run `kubectl get` and neat its output
kubectl-neat version        Print the installed version
```

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-f`, `--file <path>` | `-` (stdin) | Path to the YAML/JSON file to neat. Use `-` for stdin. |
| `-o`, `--output <fmt>` | auto | Output format: `yaml` or `json`. Auto-detected from the input if omitted. |
| `-v`, `--version` | `false` | Print the `kubectl-neat` version and exit. |

Flags may appear around the `get` subcommand. See [Output format](#output-format) below.

## Mode 1: Local file or stdin

This is the default mode. Input is read from a local file or from stdin. The file must be a
valid Kubernetes resource in YAML or JSON.

```bash
# from stdin (explicit)
kubectl neat -f - < ./my-pod.json

# from stdin (default)
kubectl neat < ./my-pod.json

# from a file
kubectl neat -f ./my-pod.yaml

# pipe a kubectl result
kubectl get pod mypod -o yaml | kubectl neat

# change the output format
kubectl get pod mypod -oyaml | kubectl neat -o json
kubectl neat -f ./my-pod.json --output yaml
```

## Mode 2: kubectl get wrapper

`kubectl neat get ...` is a convenience that runs `kubectl get` and neats the output in one
command. It accepts any argument `kubectl get` accepts and forwards them. It needs `kubectl`
on the `PATH`.

Internally it always asks `kubectl get` for JSON; the final output is converted to JSON unless
you pass `-o yaml`. Any `-o/--output` you pass to the `kubectl` side is respected for the
read step.

```bash
kubectl neat get -- pod mypod -oyaml
kubectl neat get -- svc -n default myservice --output json
kubectl neat get -- all -A
```

> Note: `kubectl neat get` outputs **JSON by default**. Add `-o yaml` (a neat flag) to get
> YAML, e.g. `kubectl neat get -o yaml -- pod mypod`.

## Output format

| Command | Default output |
| --- | --- |
| `kubectl neat` (local mode) | Same as input format (auto-detected) |
| `kubectl neat get` | `json` |
| `kubectl neat -o yaml` | `yaml` |
| `kubectl neat -o json` | `json` |

## Version

```bash
kubectl neat version
kubectl neat --version
kubectl neat -v
```

All three print the installed version. On locally built binaries (without goreleaser
ldflags) this shows `v0.0.0+unknown`.

## What gets removed

The exact list of fields that are neat-ed out is documented in
[architecture.md](architecture.md). In short:

- `status` and other runtime information
- scheduler-assigned `spec.nodeName`
- Pod service-account token volumes (`default-token-*`) and deprecated `spec.serviceAccount`
- system-added tolerations: the `DefaultTolerationSeconds` entries (`node.kubernetes.io/not-ready`,
  `node.kubernetes.io/unreachable`, 300s) and node-condition tolerations
  (memory/disk/pid pressure, unschedulable, network-unavailable)
- RuntimeClass-computed pod overhead (`spec.overhead`)
- `creationTimestamp` in workload pod templates (`spec.template.metadata`)
- the `kubectl.kubernetes.io/last-applied-configuration` annotation
- default values for the `v1` (core), `apps/v1`, and `batch/v1` API groups
- empty arrays and objects

Resources from API groups without defaulting support (e.g. CRDs) pass through unchanged.

## Exit codes and errors

The process exits non-zero and prints an error to stderr when:

- an unknown flag or command is given
- input JSON is empty or invalid
- the input file cannot be read
- YAML/JSON conversion fails
- `kubectl` (in `get` mode) fails

Example diagnostics:

```
$ kubectl neat --foo
Error: unknown flag: --foo

$ kubectl neat -f ./not-a-resource.json
Error: error converting from yaml to json : ...