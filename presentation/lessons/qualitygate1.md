# 💡 Quality Gate

**A Quality Gate on Every Commit.** The spec holds today.
Deterministic checks in the git hooks keep it true for every commit, yours and the agent's.

----

## One quality gate: `just check`

`just check` runs from `.githooks/` on every commit and every push:

| check | refuses |
| --- | --- |
| pins | an image tag that is not the git tree of its source |
| render | `deploy/` that `kustomize build` cannot render |
| schema | a field Kubernetes does not know (`kubeconform`) |
| lint | a risky workload, e.g. no resource limits (`kube-linter`) |
| go | a failing `go vet` or `go test` |

----

## Review only once

<img src="assets/images/hitl-final-iteration.svg" alt="Human-in-the-loop after many iterations: the loop tightens into a working refine and implement cycle" class="drawing">

---

## 🛠️ Quality Gate

**The spec holds today. Make it hold after the next commit too.**

1. `kube-linter checks list` has a check for pods that no NetworkPolicy selects. Have `just claude` enable it in `.kube-linter.yaml`, make `just check` pass, and commit.
2. Every other workload gets a policy, or an ignore annotation with its reason. Is each reason true?
3. Delete scoring's NetworkPolicy: `just check` must fail. Restore it, `git push cluster`.
