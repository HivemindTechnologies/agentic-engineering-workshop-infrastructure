# 💡 Debug MCP

**A Debug MCP in the Live Pod.** Logs and traces show what the code chose to tell you.
A debug MCP lets the agent ask the running pod, through a tunnel only you can open.

----

## A debug MCP in your app

An MCP server built into scoring itself, no vendor, no API in between.

`just debug` opens it, `setup/debug-mcp.json` connects the agent:

| Tool | What it does |
| --- | --- |
| `recent_requests` | the slowest requests, with the time per step. Read only |
| `replay` | scores one event with the pod's rules, grown to n items. Runs code |
| `runtime` | goroutines, heap, GC. Read only |
| `config` | the rules the pod loaded. Read only |

----

## Into the live pod, through kubectl

The debug MCP has no login. Kubernetes does the login for it.

```
claude  ──▶  127.0.0.1:7070  ──kubectl port-forward──▶  pod: 127.0.0.1:7070
(sandbox)    (your laptop)      (your kubeconfig)         (no Service, no port)
```

- In the pod it listens on loopback only: no Service, no Ingress, no other pod reaches it
- `just debug` port-forwards with your kubeconfig. The agent's `view` kubeconfig gets a 403
- The audit log shows who opened it, Loki shows every tool call
- Off: remove `DEBUG_MCP_ADDR` from `deploy/demo/scoring.yaml`. `just check` fails if it listens beyond loopback or a port opens it
- Read only means no state anyone else sees changes. `replay` still burns CPU next to real traffic: one at a time, 2000 items at most

Whoever may port-forward usually may exec into the pod too. The debug port gives them little they did not have.
