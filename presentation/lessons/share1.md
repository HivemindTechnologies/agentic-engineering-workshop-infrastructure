# 💡 Share

**Copy or Share?** You know the question from Terraform modules and Helm charts: one shared module, or a copy in each repository? A skill asks it again. **Start with a copy.** Extract only when two copies need the same fix.

----

## Extract or duplicate

A library breaks **loudly**: the compiler goes red.<br>
A plugin breaks **silently**: reworded prose, nothing goes red.

- Extract when call sites share knowledge, not when they look alike
- Three similar skills beat one wrong abstraction
- Still learning what works? Keep it local: edit, rerun, no release
- Knobs multiply: each setting added for one team is one more for every team

----

## Why sharing a skill is hard

- Its interface is prose read by a model: no compiler, no type, no red build
- Unpinned, a reworded sentence changes the behaviour of every user at once
- So one team runs a new version first and reports what it caught

----

## Review a change like code

Per repository, a skill drifts. As a versioned plugin, it changes by a version bump.

- **Owner:** one named team, in the skill's `README.md`
- **Review:** a reworded `SKILL.md` is a pull request
- **Test, in this repo:** a fresh `just claude` loads the changed skill; run the same prompt with and without it

----

## A skill served over MCP?

- **Over MCP, today:** a tool, not a skill: no local files, no frontmatter; a change is live for everyone at once
- **Over MCP, future:** [Skills over MCP](https://modelcontextprotocol.io/extensions/skills/overview): real skills, changes still live for everyone; not yet in Claude Code or Cursor
- **As a plugin:** a change is a pull request; users get it with `claude plugin update`, or from their next launch once they turn on auto-update

----

## What to commit

- **Commit:** `AGENTS.md`, skills, specs, research and survey docs
- **Keep private:** personal global rules, credentials
