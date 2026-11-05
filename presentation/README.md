# Presentation

The slides of Part 1, a `reveal-md` presentation. The tooling, theme and
most lessons come from the software track
(`agentic-engineering-workshop-software/presentation/`), so a diff against it
shows every infrastructure change.

The spine, `deck-infrastructure-part1.md`, holds the title slides, the ordered
lesson list and the wrap-up; each lesson's slides are one flat
`lessons/<id>.md`. `build-deck` splices a spine into `.build/<stem>.md`.

## Running

```sh
cd presentation
nix develop    # reveal-md, python3, just, entr
just serve 1   # Part 1 on localhost:8000; edits to a lesson reload the browser
just html 1    # Part 1 as a static site in dist/
just check     # lesson format, the deck builds, comic credits, no local paths
```

For a PDF, open `localhost:8000/deck-infrastructure-part1.md?print-pdf` in
Chrome and print to a file.

## Running the Workshop

- The menu is the `<!-- lessons ... -->` list in the spine. Each line is one
  lesson id; a leading `#` skips it, `== <title>` opens a part, `-- <title>`
  opens a block with its own divider slide, `~~ Break` ends a 50-minute
  session with a break slide. **Position is the order.**
- Each block opens with a divider showing the part's agenda, current block
  highlighted. The bar at the top spans every slide of the deck; a ☕
  marks each break.
- Every 🛠️ slide runs in `workbench/`; `workbench/exercises/part1.md` holds
  the same steps, one section per 🛠️ slide.
- During a 🛠️ slide: show at most the first step, then stay quiet with the
  exercise slide on screen. Answer questions, add small comments ("Argo CD
  needs about a minute"). No demo: in the software track people watched
  the demo instead of running the exercise.

## Lessons

One flat `lessons/<id>.md` per lesson, no frontmatter:

- First line is `# 💡 <Name>`; `<Name>` is the slide label and the agenda
  entry.
- A long-form intro paragraph starting `**<Full Title>.** …`, then `----`
  sub-slides for the theory.
- Optionally a `---` then `## 🛠️ <Exercise>`, a name of its own or the lesson's: 3 to 6 steps, a "done when"
  checklist only where the steps do not already state the outcome. The
  exercise fits one slide.
- Optionally a `## 💬 Discuss: <question>?` slide, after a `---` or below the
  exercise with `----` (theory left, practice in one column right): a table Option | You
  gain | You pay with 2 to 4 options and no recommendation, one line naming
  the file that decides, and the chat prompt. The next slide, `What we
  chose, and why`, gives our pick, the excerpt that sets it, and how to
  switch.
- A drawing belongs to the lesson it motivates, as a `----` sub-slide with an
  `<img class="drawing">`.

The number in an id is the part the lesson runs in, as in the software
track: `harness1` in Part 1.

| Lessons | Origin |
| --- | --- |
| the other lessons | software track, examples swapped for this cluster |
| `rollback1` | the earlier infrastructure deck (`agentic-engineering-workshop`) |
| `setup1`, `gitops1`, `qualitygate1`, `secrets1`, `debug1`, `debugagent1` | written for this track |

## Drawings and Comics

The Excalidraw sources of the drawings live in the software track
(`docs/drawings/`); only the exported SVGs are here, in `assets/images/`.
Each comic names its author, title and source URL on its own slide;
`just check` fails otherwise.
