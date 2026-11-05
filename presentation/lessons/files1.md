# 💡 Files

**Files as the Interface.** Chat is a transcript; a file is a contract.
Spec, environment, sandbox and skills become reviewable artifacts: every change is a diff, and a lost session costs nothing.

----

## Specs as files

The spec is declared in `docs/SPEC.md`:

```markdown
# Only customers reach scoring

Decided:
- NetworkPolicy, ingress only
- namespace traffic may connect, nothing else

Open:
- may the kubelet's probes still reach scoring?
```

- A changed decision is a diff you review
- You, your colleague and the next session read the same text
- `/compact` and `/clear` leave it as it is

----

## Environments as files

The dev shell is declared in `flake.nix`:

```nix
{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  outputs = { nixpkgs, ... }:
    let pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      devShells.x86_64-linux.default = pkgs.mkShellNoCC {
        packages = with pkgs; [ just k3d kubectl kustomize kube-linter nono ];
      };
    };
}
```

- The agent adds a missing tool as a diff you review
- You, your colleague and the agent enter the same shell: `nix develop`
- CI enters it too: `nix develop --command just check`

----

## Sandboxes as files

The boundary is declared in `nono.json`:

```json
{
  "extends": "default",
  "groups": { "include": ["deny_credentials", "deny_shell_history"] },
  "filesystem": { "allow": ["$WORKDIR"] },
  "environment": { "deny_vars": ["SSH_AUTH_SOCK", "SECRET_TOKEN"] }
}
```

- Widening the boundary is a diff you review
- You, your colleague and the agent run inside the same walls: `just claude`

----

## Paths as an index

- The agent reads directory names before code
- `exercises/`, `deploy/`, `scoring/` beats `misc/`, `stuff/`, `src2/`
- One responsibility per folder

---

## 💬 Discuss: Manual systems?

| Option | You gain | You pay |
| --- | --- | --- |
| By hand | nothing to set up | no record |
| Script, run once | reviewed, in Git | drifts after the run |
| Declared, applied by a pipeline | diff review, `git revert` | upfront work |

Chat: your pick, and what would change your mind.

----

## What we chose, and why

**Declared.** Argo CD already is the pipeline.

→ `setup/bootstrap/application.yaml`

Switch, for a device without a pipeline: a reviewed script
