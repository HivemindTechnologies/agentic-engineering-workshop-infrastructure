# 💡 MCP

**Model Context Protocol.** The agent sees only what is in its context. Your tickets, docs and dashboards are not, so you copy them in.
MCP is the standard way for an agent to fetch them itself.

----

## Why it exists

Before MCP, every agent built its own integration for every tool.

```
without MCP: every agent × every tool

Claude  ──┬── Jira
Cursor  ──┼── GitHub
Copilot ──┴── Postgres

with MCP: one server per tool, every agent uses it

Claude  ──┐       ┌── Jira server
Cursor  ──┼─ MCP ─┼── GitHub server
Copilot ──┘       └── Postgres server
```

- Anthropic published it in November 2024 as an open standard
- Most agents speak it today: Claude Code, Cursor, Copilot, ChatGPT, Gemini CLI
- The same idea as the Language Server Protocol, between editors and languages

----

## How it works

```
Claude Code  ──MCP──▶  MCP server  ──API──▶  the system
(the agent)            (an adapter program)   (Jira, GitHub, a database)
```

1. At start the agent asks the server `tools/list`: each tool's name, description and input schema
2. The model picks a tool and calls it: `tools/call search_issues {"query": "login fails"}`
3. The server calls the real API; the result lands in the context

- The messages are JSON-RPC, over one of two transports:
  - **stdio**: a program on your laptop; the agent starts it and talks over stdin/stdout
  - **HTTP**: a URL, usually at the vendor
- Besides tools, a server can offer resources (data to read) and prompts (templates). Agents mostly use tools

----

## Adding a server

One command per transport:

```sh
# HTTP: Sentry runs the server, you sign in in the browser
claude mcp add --transport http sentry https://mcp.sentry.dev/mcp

# stdio: claude starts the program, it runs next to the agent
claude mcp add db -- npx -y @bytebase/dbhub@1.3.1 \
  --dsn "postgresql://readonly:pass@localhost:5432/app"
```

In the session, `/mcp` shows whether each server connected and which tools it offers.

----

## What it costs

Everything a server tells the agent lands in the context window, in three steps:

1. **At start**: only tool names, plus the server's instructions
2. **On search**: the agent calls `ToolSearch`, and a matching tool's description and input schema load
3. **On call**: the result lands in the context, capped at 25k tokens by default

----

## Connect everything?

<img src="assets/images/xkcd-2044-sandboxing-cycle.png" alt="a cycle: connect everything, so many connections it creates security holes, enclose everything in sandboxes, wish the parts could communicate" class="comic">

<p class="credit">xkcd 2044 "Sandboxing Cycle" · CC BY-NC 2.5 · xkcd.com/2044</p>

**Which of your systems would you connect to an agent?**

----

## What it risks

An MCP server is a dependency whose version lives outside your repository.

- Its tool descriptions are text you did not write, and the model can follow them as instructions
- **HTTP**: the vendor deploys when it wants; you cannot pin it
- **stdio**: `npx -y pkg` runs whatever npm serves that day; `pkg@1.3.1` pins it in `.mcp.json`

----

## MCP or CLI

**A CLI** where the model knows one: `gh`, `kubectl`, `psql`, or `curl` against an HTTP API

- the model has seen it in training
- it costs nothing until it runs
- its output can be filtered first

**MCP** where the model knows no CLI

- your in-house system: never in training, thinly documented
- `tools/list` tells the agent what exists, each tool with a typed input schema

The code says what should run, the process what does: "How many rules does scoring use?" finds a generator in the code, the number in the pod.

---

## 💬 Discuss: MCP or CLI?

| Option | You gain | You pay |
| --- | --- | --- |
| CLI recipes in a skill | no extra server; same in your shell | agent writes the queries |
| Grafana MCP | typed tools, fewer calls | another credential; large results |
| Debug MCP in your app | the live process's own view | you build and maintain it |

Chat: your pick, and what would change your mind.

----

## It depends

**Asked once: CLI. Repeated in every incident: MCP.**

→ `skills/diagnosing-from-telemetry`

Try MCP: `mcp-grafana` in `flake.nix` and `.mcp.json`, Grafana token with role Viewer
