# 💡 Spec

**Spec-Driven Development.** A vague request becomes a spec with checkable acceptance criteria; a human reviews it, then the agent implements it.
Reviewing a spec costs minutes. Reviewing manifests from a wrong spec costs a day.

----

## Nom nom

<img src="assets/images/sprouts.svg" alt="a tiny watered spec sprout that gen AI grows into a carnivorous code plant eating a person" class="drawing">

----

## Why a spec

**A spec is a written agreement** between the people who own a system, like the criteria on a ticket.

- **Intent:** what was wanted, and why
- **Unspecified:** whatever it does not say was never agreed, so it is still open
- Without a spec, the agent guesses each gap, differently each run
- Agreeing on paper is cheaper than disagreeing in code

----

## Specs are not new

**We always had specs**, spread over many places:

- A Jira ticket
- A Slack thread
- A talk in a meeting

**A standard structure is not new either:** remember Jira ticket templates?

**What is new:**

- The spec lives in the repo, next to the code
- Processes and tools work on it

----

## The changeflow flips

<img src="assets/images/new-world.svg" alt="old world: requirements, then code, then documentation, one human at each step; new world: requirements, then documentation, then code, with a human and agents at each step and three agents at the code" class="drawing">

----

## When to write a spec

<div class="two-column">
<div class="column">

**Use it**

- A change to what the system guarantees, e.g. who reaches scoring
- The intent matters: what, when and why
- A later step checks against the intent, e.g. an agent review
- Two or more people must agree

</div>
<div class="column">

**Skip it**

- Typo, rename, version bump
- A bug: the failing check is the spec
- A spike you throw away
- A refactor: `just check` is the spec

</div>
</div>

----

## Spec and plan

<div class="two-column">
<div class="column">

**Spec:** what and why

- Agreed with the team
- Each criterion names its check
- A human reads every line
- Kept and updated with `deploy/`

</div>
<div class="column">

**Plan:** how (proposal)

- Written by the agent, for the agent
- Exact files, manifests, commands
- A human skims it
- Done once the change is deployed

</div>
</div>

----

<!-- .slide: class="text-slide" -->

## What users report

Four single cases, one point each:

1. Settling intent before code pays off: the spec fixed per-host cooldowns before any code; changing that later would have meant a refactor ([Builder.io](https://www.builder.io/blog/claude-code-superpowers-plugin), vendor blog)
2. Heavy workflows cost more than they return: 2,577 lines of markdown for 689 lines of code, 3.5 h of review ([Scott Logic](https://blog.scottlogic.com/2025/11/26/putting-spec-kit-through-its-paces-radical-idea-or-reinvented-waterfall.html))
3. One workflow does not fit every size: Kiro turned a small bug into 16 acceptance criteria ([Böckeler](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html))
4. Specs written per change go unread: 94 spec folders, never reopened ([Perez](https://www.alfredo-perez.dev/blog/2026-08-31-i-wrote-238-specs-and-never-read-one-again), tool author)

----

## Nom nom

<img src="assets/images/intent.svg" alt="No spec, no idea: people ask why is it like this and is this normal, while the plant from Nom nom eats a person" class="drawing">

---

## 🛠️ Spec

**Every pod in the cluster can reach scoring. Only the customers should.**

1. Probe scoring from namespace `default` (command in `exercises/part1.md`): it prints `reachable`
2. `just claude`, then:
   > Only the customers (namespace traffic) may reach scoring. Interview me about this, one question at a time, then write deploy/demo/SPEC.md. Each acceptance criterion names the command that checks it. Do not implement yet.
3. Review the spec before any code: what else talks to scoring, or does scoring talk to?
4. "Implement deploy/demo/SPEC.md and commit." Then `git push cluster`
5. Prove each criterion with its command; the probe from step 1 now fails
