# AGENTS.md

A risk scoring service (`scoring/`, Go) on a local k3d cluster. Argo CD
applies `deploy/` from the in-cluster git server. `setup/` holds the cluster
bootstrap and `setup/loadgen/`, the outside world sending customer traffic.
Leave `setup/` alone.

## Changing the System

- The only way into the cluster is git. Edit files, run `just check`, commit.
  A human deploys with `git push cluster`; you cannot push, do not try.
- Your kubeconfig is read-only. Never try to work around that.
- After changing `scoring/`: `git add scoring`, `just pin`, then commit code
  and pin together, so one commit is one deployable state.
- One commit, one change, one line of message.
- A rollback is `git revert <commit>`: the old state reaches the cluster the
  way every change does.

## Looking at the System

- `kubectl get`, `describe`, `logs`.
- Prometheus `localhost:9090`, Loki `localhost:3100`, Tempo `localhost:3200`,
  Pyroscope `localhost:4040`.

## Notes

Add what you had to explain twice.
