# 💡 GitOps

**The GitOps Loop.** The cluster runs what `deploy/` on `main` declares.
A change reaches it by `git push cluster` only; Argo CD applies it and undoes changes made any other way.

----

## One path in

| step | who | what runs |
| --- | --- | --- |
| edit, commit | agent | pre-commit hook: `just check` |
| `git push cluster` | you | pre-push hook: `just check`, `just images` |
| sync every 15 s | Argo CD | applies `deploy/`, undoes drift, deletes what git dropped |

----

## Beyond the laptop: GitOps

- **Files are the interface:** knowledge and desired state live in Git
- **A reconciler converges:** makes the real world match each merged change (Argo CD, Flux)
- **Read-only access:** the reconciler writes to production; the agent only reads logs and metrics

---

## 🛠️ GitOps

**One change from prompt to running pod, and nobody types `kubectl apply`.**

1. `just claude`: a change in `deploy/`, e.g. "Run scoring with two replicas."
2. Second terminal, `nix develop`: review the commit, then `git push cluster`. What did the hooks run?
3. `just status`: how long until Argo CD runs your commit?
4. Now skip Git: scale to three with `kubectl`, first in this shell, then in `nix develop .#admin`. Watch the deployment for a minute: what happened each time, and who did it?

----

## 💬 Discuss: Agent commits?

| Option | You gain | You pay |
| --- | --- | --- |
| Read only | every commit is yours | you write every commit |
| Commits, no push (today) | checked commits; you deploy | history you did not write |
| Commits and push | agent deploys itself | Argo CD applies whatever it pushes |

Chat: your pick, and what would change your mind.

----

## What we chose, and why

**Commits, no push.** Only you deploy.

→ `nono.json`, `filesystem.allow`: `$WORKDIR/../.git`

Switch: move it to `filesystem.read`, and only you commit
