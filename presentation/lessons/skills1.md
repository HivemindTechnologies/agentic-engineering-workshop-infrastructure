# 💡 Skills

**Skills and Commands.** Two kinds of file hold how your team works: a skill is loaded by the model when its description matches, a command is started by a human.
Both are plain Markdown in the repository, reviewed like code.

----

## Anatomy of a `SKILL.md`

<!-- `text`, not `markdown`: the highlighter reads the line above `---` as a heading. -->
```text
---
name: diagnosing-from-telemetry
description: Use when a service is slow, failing, restarting or
  out of memory, or when someone reports a problem with a running service
---

# Diagnosing From Telemetry

First turn the complaint into one measurable quantity ...
```

- `description`: the only part read before loading, so it does the matching
- Body: the procedure, most important first

----

## Command or skill

Both are a Markdown file with instructions. The difference is who starts it:

- **Command**: you start it, by typing `/` and its name
- **Skill**: the agent, when the task matches

A command is a skill with `disable-model-invocation: true`.

----

## Prose has no compiler

<img src="assets/images/monkeyuser-273-natural-language-instructions.png" alt="an AI machine on a rooftop fires a laser at the city; a developer at the laptop: here's the problem, you wrote evil instead of eval" class="comic">

<p class="credit">MonkeyUser · “Natural Language Instructions” · monkeyuser.com/2024/natural-language-queries</p>

---

## 🛠️ Skills

**Every change in its own git worktree: write it down once, for every session after this one.**

1. With `just claude`, write `skills/working-in-a-git-worktree/SKILL.md`: start a change in a worktree under `.worktrees/`, commit there, merge into `main`, clean up. Its description says when to use it.
2. Review: would you work this way yourself? Commit it.
3. A fresh `just claude` (it copies `skills/` at start): "Run scoring with three replicas." Does it load the skill and work in a worktree? Push the result.
