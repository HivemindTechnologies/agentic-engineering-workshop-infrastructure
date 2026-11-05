# 💡 Context

**Context Engineering.** You cannot change the weights, only what the model reads.
Rules in `AGENTS.md`, procedures in skills, requests as goal, constraints, acceptance.

----

## What fills a context window

<img src="assets/images/context-injection.svg" alt="AGENTS.md, always loaded, and skills, loaded on demand, enter the short-term context together with your task, domain, standards, business needs and tools; beside it, the static training holds what the model knows about infrastructure and tools in general" class="drawing">

----

## Loading is not obeying

Where a rule lives, from most to least skippable:

1. **Prompt**: what you type, this conversation only
2. **Skill**: when the description matches
3. **`AGENTS.md`**: every turn, still advice
4. **Script**: enforced when someone remembers to run it
5. **Hook**: fires by itself, skippable with a flag
6. **CI**: the repository refuses the merge

1–3 advice. 4–5 enforced, but skippable. 6 needs a review to skip.

----

## Code is context

<img src="assets/images/managerdev-broken-windows-theory.png" alt="a maintained office building, one agent smashing a single window to keep it simple for now, and a month later every window broken while robots copy the pattern" class="comic">

<p class="credit">Anton Zaides · “The broken windows theory of coding agents” · manager.dev/newsletter/the-broken-windows-theory-of-coding-agents</p>

----

## Context rot

- Stale text in the window lowers answer quality
- Symptoms: forgotten decisions, repeated file reads, loops
- `/context` shows what fills the window
- Next task unrelated: `/clear`
- Same task, window full: `/compact focus on the NetworkPolicy`

----

## Context is precious

<img src="assets/images/gollum-context-is-precious.jpg" alt="Gollum, smiling: we just greps the one error line, precious; Gollum, screaming: no, we loads the whole Loki stream" class="comic">

<p class="credit">Template “Gollum good/bad” · imgflip.com/meme/173713637 · film still: The Lord of the Rings, New Line Cinema</p>
