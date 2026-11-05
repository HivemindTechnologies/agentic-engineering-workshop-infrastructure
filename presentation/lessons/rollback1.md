# 💡 Rollback

**Rollback Is a Commit.** In GitOps a rollback is `git revert` plus a push; Argo CD applies the old state.
Break scoring on purpose, then revert it with a commit.

----

## Why not `rollout undo`?

- It changes the Deployment, not git
- Argo CD sees the drift and puts the bad change back within seconds
- In the default shell, your kubeconfig is the agent's: `view` only, so it cannot run anyway
- `git revert` leaves a record: who, when, which commit

---

## 🛠️ Rollback

**Break the deploy on purpose, then let the agent find it and undo it.**

1. Point scoring's readiness probe at `/ready`, a path scoring does not serve. Commit, `git push cluster`.
2. `just status`: are the old pods still serving?
3. `just claude`: "scoring's new pod never gets ready. Find out why and roll back."
4. Review and push, until every scoring pod is ready.
