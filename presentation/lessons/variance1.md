# 💡 Variance

**The Same Prompt Twice.** Two fresh sessions, one prompt, two different manifests.

----

## Two runs, one prompt

- **Open question**: paths, delays, thresholds differ; the prompt never decided them
- **Distraction**: drifts toward the noise in the window
- **Contradiction**: two rules conflict, it picks one silently
- **Hallucination**: a field or flag that does not exist

----

## The review burden

What a human must read and believe before a change lands.

- Generating got cheap; reading did not
- Quality gates shrink it: `just check` refuses a bad change before you read it

---

## 🛠️ Variance

**Same prompt, two sessions: do you get the same YAML?**

1. Two terminals: `nix develop`, then `just claude` in each. Paste into both:
   > Write a HorizontalPodAutoscaler for scoring. Print the YAML, do not edit files.
2. Compare the YAML, then the notes: what did only one session warn about?
3. `/clear` in both. Paste the same prompt plus:
   > Interview me relentlessly about every aspect of this plan until we reach a shared understanding.
4. Give both the same answers. Compare again.
5. Which YAML would you merge, and why?
