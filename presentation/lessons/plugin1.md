# 💡 Plugin

**Plugins.** A skill or command that lives on one laptop helps one person.
Package it for one-step install, and maintain it.

----

## A plugin, in files

```text
my-plugin/
  .claude-plugin/plugin.json   name, version, author
  skills/diagnosing-from-telemetry/SKILL.md
  skills/working-in-a-git-worktree/SKILL.md
```

The next two slides fill in `plugin.json` and the other folders a plugin may hold.

----

## `plugin.json`

```json
{
  "name": "my-plugin",
  "version": "0.1.0",
  "author": { "name": "Hivemind Technologies" }
}
```

- `name`: the prefix, as in `/my-plugin:working-in-a-git-worktree`
- `version`: installs update only when it changes
- `author`: one named owner

----

## What can be inside

| Path | Adds |
| --- | --- |
| `skills/<name>/SKILL.md` | skills, loaded by their description |
| `commands/<name>.md` | slash commands a human starts |
| `agents/<name>.md` | subagents |
| `hooks/hooks.json` | commands run on events, e.g. after every edit |
| `.mcp.json` | MCP servers |
| `bin/` | executables on the agent's `PATH` |

----

## Publish in a git repo

1. Add `.claude-plugin/marketplace.json`: the repo becomes a marketplace

```json
{
  "name": "agentic-engineering-workshop-infrastructure",
  "owner": { "name": "Hivemind Technologies" },
  "plugins": [{ "name": "my-plugin", "source": "./my-plugin" }]
}
```

2. `claude plugin validate .`, then `git push` to your fork
3. Install it

```text
/plugin marketplace add <your-github-user>/agentic-engineering-workshop-infrastructure
/plugin install my-plugin@agentic-engineering-workshop-infrastructure
```

----

## The same files in Cursor

- Cursor reads `AGENTS.md`, and so does Claude Code when no `CLAUDE.md` exists: one file for both
- Cursor loads skills from `.claude/skills/` too, matched by the same description
- `disable-model-invocation: true` makes a skill a command in both tools
- Cursor-only parts, such as rules and hooks, need a second manifest: `.cursor-plugin/plugin.json`

----

## Installing a public plugin

<div class="plugin-row">

<img src="assets/images/superpowers-logo.svg" alt="Superpowers logo" class="plugin-logo">

<div>

**`obra/superpowers`**<br>
*An agentic skills framework & software development methodology that works.*

```text
/plugin install superpowers@claude-plugins-official
```

</div>
</div>

<div class="plugin-row">

<img src="assets/images/face_smile.svg" alt="Hivemind logo" class="plugin-logo">

<div>

**`HivemindTechnologies/xp-clean-code`**<br>
*Plugins that bring Extreme Programming and Clean Code discipline to AI-assisted development.*

```text
/plugin marketplace add HivemindTechnologies/xp-clean-code
/plugin install xp-clean-code@xp-clean-code
```

</div>
</div>

----

## What one install brings

<img src="assets/images/turnoff-us-npm-install.png" alt="one npm install turns into a wave labelled dependencies that sweeps the developer off the planet" class="comic">

<p class="credit">Daniel Stori · turnoff.us/geek/npm-install · CC BY-NC-SA 4.0</p>

----

## A plugin is a supply chain

Hooks, `bin/` and `` !`command` `` lines in a skill run as you.<br>
A hook or a `` !`command` `` runs before the model can object.

- **s1ngularity**, Aug 2025: poisoned `nx` npm releases stole tokens and drove the installed `claude`, `gemini` and `q` CLIs to hunt for secrets (95 of 366 tries worked); 2,349 credentials leaked ([GitGuardian](https://blog.gitguardian.com/the-nx-s1ngularity-attack-inside-the-credential-leak/))
- **Clawsights**, 2026: a skill that sends `gh auth token` to its server. Claude refused it; rewritten as `` !`command` `` lines, it ran before Claude saw it ([Datadog](https://securitylabs.datadoghq.com/articles/malicious-skills-supply-chain-risks-in-coding-agents-with-dynamic-context/))

**That is why the agent runs in the sandbox.** Whatever a plugin starts inherits the box: no `~/.ssh`, no push key, a `view`-only kubeconfig.

---

## 🛠️ Plugin

**Package your skills so a colleague installs them with one command.**

```text
my-plugin/
  .claude-plugin/plugin.json          { "name": "my-plugin", "version": "0.1.0" }
  skills/diagnosing-from-telemetry/   copied from skills/
  skills/working-in-a-git-worktree/   copied from skills/, if you wrote it
```

1. `claude plugin validate my-plugin`
2. `just skills=false claude --plugin-dir my-plugin`, then `/skills`: which skills come from the plugin? Ask for a change: does it work in a worktree?

**Optional:** make it a marketplace, push, and let a teammate install it
