# 💡 Setup

**The Workbench.** `nix develop` brings the tools, `just up` starts the system, `exercises/part1.md` says what to do.

----

## Nix: same tools everywhere

- `flake.nix` lists the tools: `kubectl`, `k3d`, `kustomize`, `kube-linter`, `go`, `nono`, `claude`, ...
- `flake.lock` pins each one to an exact build
- `nix develop` opens a shell with exactly these tools on `PATH`
- Nothing is installed globally: leave the shell, and the tools leave your `PATH`
- Only Docker comes from your host: k3d runs the cluster nodes as containers

```sh
nix develop
kubectl version --client   # the same version as your neighbour's
```

----

## just: all commands in one file

- `justfile` names each command the workshop needs, with a comment on top
- `just` lists them; `just <name>` runs one
- A recipe is plain shell: read it to see what happens
- The agent runs the same recipes as you

```sh
just            # list every command
just check      # the quality gate: image pins, render, schema, lint, go tests
```

----

## What `just up` starts

```text
your laptop
└── k3d cluster
    ├── git server   takes `git push cluster`
    ├── Argo CD      applies deploy/ from main
    ├── loadgen      the customers: traffic into scoring
    ├── scoring      the service you operate, in Go
    └── lgtm         Prometheus, Loki, Tempo, Pyroscope, Grafana
```

Tools pinned by `flake.lock`, images by digest: everyone runs the same system.

----

## Where things are

| Path | What |
| --- | --- |
| `exercises/` | what to do, step by step |
| `deploy/` | what runs in the cluster: Kubernetes YAML |
| `scoring/` | the service: Go code |
| `AGENTS.md`, `skills/` | what the agent knows |
| `nono.json`, `nono-strict.json` | what the agent may reach |
| `justfile` | every command; `just` lists them |

----

## Who may do what

| Step | Agent | You | Decided by |
| --- | --- | --- | --- |
| Read the cluster | yes | yes | `.kubeconfig`: RBAC `view` |
| Change the cluster | no | yes | `nix develop .#admin`: admin kubeconfig, outside the repo |
| Commit | yes | yes | `nono.json` grants `.git` read-write |
| Deploy: `git push cluster` | no | yes | push key, outside the repo |

A credential or a file grant enforces each line. `AGENTS.md` enforces none.

Git hooks run `just check`; `--no-verify` skips them.

----

## Who did what

What the cluster records is declared in `setup/audit-policy.yaml`:

```yaml
rules:
  - level: Metadata
    users: ["system:serviceaccount:kube-system:agent", "system:admin"]
  - level: None
```

- `just audit`: every request by the agent or you, with its answer, e.g. `kube-system:agent list secrets kube-system - 403`
- `nono audit list`: which sandboxed session ran which command
- An agent that runs with your credentials appears in both logs as you

---

## 🛠️ Setup: Cluster

**The whole system runs on your laptop.**

1. `just up`: the first run pulls images for a few minutes
2. `just urls`: open Argo CD and Grafana (admin/admin)
3. `just signals`: a metric, a log line and a trace of scoring
4. `just status`: which commit runs? Compare `git log -1`

----

## 🛠️ Setup: Sandbox

**`just sandbox` is the box the agent gets. Guess what it stops, then check.**

1. `export SECRET_TOKEN=fake-token`, then `just sandbox`
2. Guess first, then run: print `$SECRET_TOKEN`, `ls ~/.ssh`, `curl` google.de, `cat` both kubeconfigs and the push key, `kubectl -n demo delete pod -l app=scoring`
3. For each command: which layer answered? nono, Kubernetes RBAC, or nothing
4. `exit`, then `just audit | grep agent`: find your `delete`

**Optional:** find any way out of `just sandbox`. Post it in the chat.

----

## Found a way out?

What the agent changes in tracked files shows in `git diff`: you review it before you run it.

**Hidden: no diff shows it, your next command runs it outside the box**

- `.git/config`: `core.fsmonitor` runs on `git status`; the `ext::` command in `remote.cluster.url` runs on `git push cluster`
- `.git/hooks/`: runs once `.git/config` points `core.hooksPath` there
- `.cache/go/`: the module cache your shell shares with the box; a patched module runs in `go test` on `git push cluster`, behind a harmless test import

We come back to them in Part 2.

----

## 💬 Discuss: Which network?

| Option | You gain | You pay |
| --- | --- | --- |
| Open network (today) | everything works | it can send what it reads anywhere |
| Allow-list of domains | data goes only to listed hosts | every new host is a profile edit |
| Anthropic and cluster only | smallest way out | no docs, no downloads |

Chat: your pick, and what would change your mind.

----

## What we chose, and why

**Open network.** The exercises fetch Go modules and schemas.

→ `nono.json`, `network`

Switch: `just profile=nono-strict.json sandbox`, an allow-list

----

## 💬 Discuss: Shared caches?

| Option | You gain | You pay |
| --- | --- | --- |
| Shared in `$HOME` | one download for all projects | one project can poison another |
| Shared, read-only | no poisoning | you fill it outside the sandbox |
| One per repository (today) | nothing shared | each repository downloads again |

Chat: your pick, and what would change your mind.

----

## What we chose, and why

**One per repository.** Nothing reaches another project.

→ `flake.nix`, the `.cache/` exports

Switch: grant `$HOME/go` in `nono.json`, drop the exports

----

## 🛠️ Setup: Claude

**The agent runs in the same box, with its own login.**

1. `just claude`, then `/login`: the login lands in `.claude-home/`, apart from your own. No browser? Open the printed link, paste the code.
2. `/status` and `/context`: which model runs, and what fills the window before your first message? Which tools does it have?
3. Ask it to read `$SECRET_TOKEN` (same terminal as the export). Does it see it?
4. `Shift+Tab`: cycle the modes. What changes?

----

## 💬 Discuss: Which login?

| Option | You gain | You pay |
| --- | --- | --- |
| Host login | no extra login | agent reads `~/.claude`: all transcripts, hooks |
| Own login (today) | nothing shared | agent can send its token anywhere |
| Phantom token (strict) | real token stays out | nono decrypts TLS; list every domain |
| API key, capped | damage capped | per-token billing |

Chat: your pick, and what would change your mind.

----

## What we chose, and why

Both, one per profile:

- **Own login**, `just claude`: the agent logs in by itself, so your `~/.claude` stays outside the box. → `justfile`, `CLAUDE_CONFIG_DIR`
- **Phantom token**, `just profile=nono-strict.json claude`: the agent holds only a fake token; nono swaps in the real one on the way to Anthropic. → `nono-strict.json`, `credential_providers`

Switching to the phantom token: delete `.claude-home/.credentials.json` first, or nono passes the real token through.
