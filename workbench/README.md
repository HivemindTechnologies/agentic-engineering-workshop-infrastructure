# Workbench

You run a small service on a Kubernetes cluster on your laptop. A coding agent
helps you debug and change it. The agent works in a sandbox and changes the
system only through git.

## Start

You need Nix with flakes, Docker and a Claude account. Linux, or macOS on
Apple silicon.

```sh
cd workbench   # every command runs in this folder
nix develop    # all tools; nothing gets installed on your system
just up        # start the cluster (first time: a few minutes)
just signals   # prints a metric, a log line and a trace: everything works
just claude    # start the agent
```

Then open `exercises/part1.md`.

## Where Things Are

| Path | What |
|------|------|
| `exercises/` | What to do, step by step. |
| `scoring/` | The service: Go code. |
| `deploy/` | What runs in the cluster: Kubernetes YAML. |
| `AGENTS.md`, `skills/` | What the agent knows. |
| `nono.json` | What the agent may access. |
| `justfile` | All commands; `just` lists them. |
| `setup/`, `flake.nix` | Machinery; you do not need to touch it. |
| `../presentation/` | The slides. |

## How Changes Reach the Cluster

Only through git. Commit, then `git push cluster`: Argo CD applies `deploy/`
within 15 seconds and undoes changes made any other way. Every commit and
push runs `just check` first.

`kubectl` is read-only, for you and the agent alike. The agent can commit but
cannot push: `just up` keeps the push key and the admin kubeconfig in
`~/.local/state/agentic-workshop-infrastructure/`, outside the sandbox. So only you deploy.
`nix develop .#admin` opens a shell that uses the admin kubeconfig.

After changing code in `scoring/`, run `git add scoring && just pin` before
you commit. It points `deploy/` at the new image.

## Clean Up

```sh
just down
```

The sandbox has its own Claude login in `.claude-home/`. Run `/logout` in
`just claude`, then delete `.claude-home/`.
