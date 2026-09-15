# Architecture

`kubectl-neat` is a small Go program built with [Cobra](https://github.com/spf13/cobra) on top
of Kubernetes' own API machinery. It treats a resource as plain JSON and surgically removes
redundant fields.

```
main.go
  └─ cmd.Execute()                      entry point

cmd/cmd.go
  ├─ rootCmd     (kubectl-neat)          stdin/file mode
  ├─ getCmd      (kubectl-neat get ...)  kubectl get wrapper
  ├─ versionCmd  (kubectl-neat version)  prints Version
  └─ NeatYAMLOrJSON()                    converts yaml<->json around Neat()

cmd/neat.go
  └─ Neat()                              the neat pipeline

pkg/defaults/defaults.go
  └─ NeatDefaults()                      Kubernetes object-model defaulting
```

## Command layer (`cmd/`)

### rootCmd — local mode

Reads input from `-f/--file` (default `-` = stdin), calls `NeatYAMLOrJSON`, prints the result.
The `-o/--output` flag is `yaml`/`json`; if the user does not set it, the output format
is `same` (matching the input).

### getCmd — `kubectl get` wrapper

Builds `kubectl get -o json <args...>` and runs it. It always requests JSON from kubectl and
then neats it. Output defaults to `json` unless the neat `-o` flag is set (kubectl-side
`-o yaml` only changes the read format).

### NeatYAMLOrJSON

Detects whether the input is YAML vs JSON (is it `{`-prefixed), converts YAML→JSON, runs
`Neat`, and converts back if the requested output format is YAML.

## The neat pipeline (`cmd/neat.go`)

`Neat(json)` runs these steps in order:

1. **List handling** — if `kind == List`, recursively neat each item, then strip list metadata.
2. **Default values** — `defaults.NeatDefaults` removes fields equal to the Kubernetes object
   model's defaults.
3. **Scheduler** — remove `spec.nodeName`.
4. **Tolerations** — `neatTolerations` removes system-injected tolerations
   (`DefaultTolerationSeconds`: `not-ready`/`unreachable` NoExecute + 300s; `TaintNodesByCondition`:
   memory/disk/pid pressure, unschedulable, network-unavailable) from `spec.tolerations` and
   `spec.template.spec.tolerations`.
5. **RuntimeClass** — `neatRuntimeClass` removes the admission-computed `spec.overhead` and
   `spec.template.spec.overhead`.
6. **Pod** — `neatServiceAccount`: remove `default-token-*` volumes and volumeMounts and the
   deprecated `spec.serviceAccount`.
   **Other workloads** — `neatWorkloadTemplate`: remove
   `spec.template.metadata.creationTimestamp`.
7. **Metadata** — drop `kubectl.kubernetes.io/last-applied-configuration`, keep only
   `name`, `namespace`, `labels`, `annotations`.
8. **Status** — remove the whole `status` block.
9. **Empties** — recursively delete empty arrays/objects, re-checking parents so no empty
   shells remain.

## Default-value removal (`pkg/defaults/`)

`NeatDefaults` decides default-ness the same way Kubernetes would:

1. Recognize the `GroupVersionKind` against a scheme registered with `core/v1`, `apps/v1`, and
   `batch/v1`.
2. Flatten `spec` into leaf paths.
3. For each leaf, delete it from the object, run it through the Kubernetes decoder +
   `myscheme.Default`, marshal back, and read the value that the object model would set.
4. If the incoming value equals the computed default, it is removable.

Because defaulting only applies to registered kinds, resources from other groups (e.g. CRDs)
pass through unchanged.

## Common mutating controllers

Here are the [recommended](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#what-does-each-admission-controller-do)
admission controllers, and their relation to kubectl-neat:

| controller | description | neat |
| --- | --- | --- |
| NamespaceLifecycle | rejects operations on resources in namespaces being deleted | ignore |
| LimitRanger | set default values for resource requests and limits | ignore |
| ServiceAccount | set default service account and assign token | Remove `default-token-*` volumes. Remove deprecated `spec.serviceAccount` |
| TaintNodesByCondition | automatically taint a node based on node conditions | Remove node-condition tolerations |
| Priority | validate priority class and add it's value | ignore |
| DefaultTolerationSeconds | configure pods to temporarily tolerate not-ready and unreachable taints | Remove default `not-ready`/`unreachable` tolerations |
| DefaultStorageClass | validate and set default storage class for new pvc | ignore |
| StorageObjectInUseProtection | prevent deletion of pvc/pv in use by adding a finalizer | ignore |
| PersistentVolumeClaimResize | enforce pvc resizing only for enabled storage classes | ignore |
| MutatingAdmissionWebhook | implement the mutating webhook feature | ignore |
| ValidatingAdmissionWebhook | implement the validating webhook feature | ignore |
| RuntimeClass | add pod overhead according to runtime class | Remove `spec.overhead` |
| ResourceQuota | implement the resource quota feature | ignore |
| Kubernetes Scheduler | assign pods to nodes | Remove `spec.nodeName` |

## Test fixtures

`test/fixtures/` contains `*-raw.*` inputs and `*-neat.json` expected outputs used by unit and
e2e tests, e.g. `pod1-raw.yaml`, `service1-neat.json`.

## Related

- [Kubernetes admission controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- [Kubernetes API conventions / defaulting](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)