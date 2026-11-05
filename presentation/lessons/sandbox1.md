# 💡 Sandbox

**Sandboxing.** The agent reads untrusted text and runs commands in one loop.
A sandbox bounds what those commands reach. This can be: files, network, environment.

----

## Bobby Tables, all grown up

<img src="assets/images/xkcd-327-exploits-of-a-mom.png" alt="a mother names her son so that a school database runs his name as SQL" class="comic">

<p class="credit">xkcd 327 “Exploits of a Mom” · CC BY-NC 2.5 · xkcd.com/327</p>

----

## Injection, then exfiltration

- Untrusted text carries instructions
- Read secret → hide it in a URL, a commit, a PR → out
- It needs both: **secret in reach** + **a way out**

----

## Write access alone

- Your `~/.kube/config` in reach: `kubectl delete` runs against the cluster
- Your `~/.aws` in reach: the agent acts as you on the whole account
- Destroy or alter, through the endpoint the work already uses
- **A write credential in reach is already the damage**

----

## Ways to sandbox

- **Claude Code's `/sandbox`**: built in, bounds Bash only; `Read` and `Edit` go through permissions
- **Docker Sandboxes (`sbx`)**: each agent in its own small VM; network rules set per machine, extended per sandbox
- **A process wrapper (nono)**: the kernel bounds one command and everything it starts

We use nono: it wraps the workflow you already have, and its whole config is one file you diff and commit.

----

## Why this comes first

- [93 % of permission prompts get approved](https://www.anthropic.com/engineering/claude-code-auto-mode)
- [Claude Code's built-in OS sandbox removed 84 % of them](https://www.anthropic.com/research/claude-code-sandboxing)
- [Red team: 24 of 25 credential-theft attempts worked](https://www.anthropic.com/engineering/how-we-contain-claude)
- Only egress and filesystem limits stopped them

**Today no agent runs unsandboxed, the facilitator's included.**

----

## Box any command

```text
nono run --allow . -- ls ~
nono run --allow . --block-net -- curl https://example.com
nono run --allow . --allow-domain example.com -- curl https://example.com
# add --dry-run: print the capabilities, run nothing
```

- `--allow DIR` reads and writes, `--read DIR` only reads
- `--block-net` or `--allow-domain`: what may leave the machine
- NixOS users: add `--read /nix/store`, where every program lives

----

## Ask why it was denied

```text
$ nono why --path ~/.aws/credentials --op read --allow .
DENIED
  Reason: filesystem_deny
  Policy: filesystem.deny

$ nono why --path ~/notes.txt --op read --allow .
DENIED
  Reason: path_not_granted
  Suggested fix: --read-file ~/notes.txt
```

Inside the box, `nono why --self` asks the running sandbox. The agent can ask it too.

----

## Inside the wall: approvals in Claude Code

nono decides what it can reach. Claude Code's permission mode decides who approves.

| mode | who approves |
| --- | --- |
| `manual` | you, every action |
| `acceptEdits` | you, commands only |
| `plan` | you, the plan first |
| `auto` | a classifier |
| `dontAsk` | your `allow` rules |
| `bypassPermissions` | nobody |

A new session starts in `auto`. `Shift+Tab` cycles. The status bar shows the mode.

---

## 🛠️ Sandbox

**Box a command you run every day, and see what the box stops.** Run everything inside `nix develop`.

```text
nono run --allow . -- CMD
```

1. **Run:** `ls ~`, `env`, `curl https://example.com`: what got through?
2. **Network:** add `--block-net`, then `--allow-domain example.com`
3. **Why:** ask `nono why` about one denial
4. **Yours:** `terraform plan`, `kubectl get pods` or your own script: what broke, and which flag fixes it?
