# 💡 Subagents

**Subagents.** A subagent is a second context window for work that would clutter yours.
It can check a diff against the same rules that guided the writing.

----

## Subagents

- **Starts cold:** it sees `AGENTS.md` and the brief your session wrote, nothing of your chat
- **Returns only its answer:** its file reads never enter your window
- **Skip it for a one-liner:** writing the brief costs more than doing the task

----

## Anatomy of a subagent file

```text
---
name: reviewing-a-diff
description: Use after a change is committed, to check the diff
  against the rule file that guided the writing
tools: Read, Grep, Glob
---

You review one diff. You never saw the conversation that wrote it.

Your brief holds the diff and names the rule file. Read the rule file.
Report each violation as file:line, rule, one sentence.
No violation: say so. Do not fix anything.
```

- `.claude/agents/<name>.md`: one file per subagent
- `name` and `description`: the only required fields
- `tools`: the allowlist; with only read tools it cannot change a file
- Body: its system prompt

----

## Notify the fleet

<img src="assets/images/yegge-wheelhouse-1-the-user-has-ruled.jpg" alt="an agent records an offhand 'I guess so?' as a permanent ruling and sends it to a fleet of agents" class="comic">

<p class="credit">Steve Yegge · The Wheelhouse #1 “The User Has Ruled” · yegge.ai/comics/the-user-has-ruled</p>
