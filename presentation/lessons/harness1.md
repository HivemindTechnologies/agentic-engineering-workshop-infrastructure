# 💡 Harness

**The Harness.** The program around the model. It sets four things: system prompt, context, tools, model.
Everything today changes one of these four.
Today: Claude Code.

----

## What is actually in there

<img src="assets/images/xkcd-1838-machine-learning.png" alt="a pile of linear algebra you stir until the answers look right" class="comic">

<p class="credit">xkcd 1838 “Machine Learning” · CC BY-NC 2.5 · xkcd.com/1838</p>

----

## Anatomy

- **Model**: Haiku, Sonnet or Opus
- **System prompt**: Claude Code's own instructions, fixed per session
- **Context**: all other text the model sees: `AGENTS.md`, your messages, files it read, tool output
- **Tools**: the actions the model can call: `Read`, `Edit`, `Bash`, …

----

## Pick the model per task

| model | per million tokens, in / out | use it for |
| --- | --- | --- |
| Haiku 5.5 | $0.10 / $0.50, prompts up to 100k tokens | renames, boilerplate, search, summaries |
| Sonnet 5.5 | $2 / $10 | everyday coding |
| Opus 5.5 | $4 / $20 | ambiguous specs, hard bugs, design |

- Cheap for mechanical work, a stronger model for judgement
- Review at the writer's tier or above
- `/model` switches, `/usage` shows what the session spent

<p class="credit">Prices: platform.claude.com/docs/en/about-claude/pricing, October 2026</p>

----

## Innocent tasks

<div class="chat">
<p class="chat-user">Clear the project cache.</p>
<div class="chat-agent">
<img src="assets/images/meta-muse-jolly-avatar.png" alt="Jolly, the plush mascot of Meta's agent Muse" class="chat-avatar">
<span class="chat-typing">• • •</span>
</div>
</div>

----

## What happened

| Asked | Happened | Missing |
| --- | --- | --- |
| Antigravity: clear a cache | deleted the whole D: drive | Sandbox |
| PocketOS: fix staging credentials | used a stray token: database and backups gone in 9 s | Secrets |
| DataTalks.Club: clean up Terraform | ran `terraform destroy` on an old state: production database gone | GitOps |
| Kiro: a change in AWS (reported) | deleted and recreated the environment: 13 h down | Setup |

Amazon: “user error, specifically misconfigured access controls, not AI.” **That is our job.**

<p class="credit">Tom's Hardware · aishippingblog.com · Financial Times via The Guardian</p>
