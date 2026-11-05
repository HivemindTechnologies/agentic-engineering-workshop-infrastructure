#!/usr/bin/env python3
"""Expand assemble-time HTML directives for the reveal-md deck (outer.v24.M1a–M1c).

Pinned call forms (HTML comments, processed by build-deck before splice):

    <!-- @steckbrief(lesson=refine-1) -->
    <!-- @steckbrief(lesson=refine-1, name=Refine, exercise=lesson-refine1) -->
    <!-- @include(path=lessons/_shared/foo.md) -->
    <!-- @include(lesson=harness1) -->

The token after `@` is the template stem under `presentation/bridge/<stem>.md`
(the registry), or `include` for file inlining. Extra `key=value` args are
comma-separated. `steckbrief` reuses `scripts/presentation-bridge` render helpers.

Unknown templates / missing includes fail closed. Nested @include max depth = 4.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PRES = ROOT / "presentation"
BRIDGE_DIR = PRES / "bridge"
DEFAULT_JOURNEY = ROOT / "workshop" / "JOURNEY.md"

DIRECTIVE_RE = re.compile(
    r"<!--\s*@([A-Za-z][A-Za-z0-9_-]*)\s*\((.*?)\)\s*-->",
    re.DOTALL,
)
INCLUDE_MAX_DEPTH = 4
ARG_RE = re.compile(
    r"""^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^,\s]+))\s*$"""
)

SLIDE_META_RE = re.compile(r"<!--\s*slide:\s*(.*?)\s*-->", re.IGNORECASE | re.DOTALL)
KIND_ALIASES = {
    "listen": "listen",
    "work": "work",
    "discuss": "discuss",
    "💡": "listen",
    "theory": "listen",
    "🏝️": "work",
    "journey": "work",
    "steckbrief": "work",
    "🛠️": "work",
    "exercise": "work",
}


class DirectiveError(Exception):
    pass


def registry_dir() -> Path:
    return BRIDGE_DIR


def list_templates() -> list[str]:
    if not BRIDGE_DIR.is_dir():
        return []
    skip = {"bindings", "readme"}
    return sorted(
        p.stem
        for p in BRIDGE_DIR.glob("*.md")
        if p.stem.lower() not in skip and not p.name.startswith("_")
    )


def parse_args(raw: str) -> dict[str, str]:
    raw = raw.strip()
    if not raw:
        return {}
    args: dict[str, str] = {}
    for part in raw.split(","):
        part = part.strip()
        if not part:
            continue
        m = ARG_RE.match(part)
        if not m:
            raise DirectiveError(f"bad directive argument {part!r}")
        key = m.group(1)
        value = m.group(2) or m.group(3) or m.group(4) or ""
        args[key] = value
    return args


def _exercise_num(lesson_id: str, exercise: str | None) -> str:
    if exercise:
        m = re.search(r"\d+[a-z]?$", exercise)
        if m:
            return m.group(0)
    m = re.search(r"(\d+[a-z]?)$", lesson_id.replace("-", ""))
    if m:
        return m.group(1)
    m = re.search(r"(\d+[a-z]?)$", lesson_id)
    if m:
        return m.group(1)
    return "1"


def _load_bridge():
    import runpy
    from types import SimpleNamespace

    path = ROOT / "scripts" / "presentation-bridge"
    if not path.is_file():
        raise DirectiveError(f"cannot load {path}")
    return SimpleNamespace(**runpy.run_path(str(path)))


def expand_steckbrief(args: dict[str, str], *, journey_path: Path = DEFAULT_JOURNEY) -> str:
    lesson_id = args.get("lesson")
    if not lesson_id:
        raise DirectiveError("@steckbrief requires lesson=<journey-id>")
    bridge = _load_bridge()
    template_path = BRIDGE_DIR / "steckbrief.md"
    if not template_path.is_file():
        raise DirectiveError(f"template missing: {template_path}")
    template = template_path.read_text(encoding="utf-8")
    journey_text = journey_path.read_text(encoding="utf-8")
    heading, body = bridge.lesson_section(journey_text, lesson_id)
    name = args.get("name")
    exercise = args.get("exercise")
    if not name or not exercise:
        for binding in bridge.load_bindings(BRIDGE_DIR / "bindings.yaml"):
            if binding.get("lesson") == lesson_id:
                name = name or binding.get("name") or "Lesson"
                exercise = exercise or binding.get("exercise")
                break
    if not exercise:
        compact = lesson_id.replace("-", "")
        exercise = f"lesson-{compact}" if not compact[0].isdigit() else compact
        if not re.search(r"\d+[a-z]?$", exercise):
            exercise = f"{exercise}1"
    if not name:
        name = lesson_id.split("-")[0].capitalize()
    if not re.search(r"\d+[a-z]?$", exercise):
        exercise = f"{exercise}{_exercise_num(lesson_id, None)}"
    return bridge.render_steckbrief(
        template,
        exercise_id=exercise,
        name=name,
        lesson_id=lesson_id,
        heading=heading,
        body=body,
        hand_region=None,
    )


def expand_generic(name: str, args: dict[str, str]) -> str:
    path = BRIDGE_DIR / f"{name}.md"
    if not path.is_file():
        raise DirectiveError(
            f"unknown template {name!r}; registry has: {', '.join(list_templates()) or '(empty)'}"
        )
    text = path.read_text(encoding="utf-8")
    for key, value in args.items():
        text = text.replace("{{" + key.upper() + "}}", value)
        text = text.replace("{{" + key + "}}", value)
    if re.search(r"\{\{[A-Za-z0-9_]+\}\}", text):
        leftover = sorted(set(re.findall(r"\{\{[A-Za-z0-9_]+\}\}", text)))
        raise DirectiveError(f"unfilled placeholders in {name}: {', '.join(leftover)}")
    return text if text.endswith("\n") else text + "\n"


EXPANDERS = {
    "steckbrief": expand_steckbrief,
}


def expand_include(
    args: dict[str, str],
    *,
    journey_path: Path = DEFAULT_JOURNEY,
    stack: list[Path] | None = None,
) -> str:
    stack = list(stack or [])
    if len(stack) >= INCLUDE_MAX_DEPTH:
        raise DirectiveError(
            f"@include nested deeper than {INCLUDE_MAX_DEPTH}: "
            + " -> ".join(str(p) for p in stack)
        )
    rel = args.get("path") or args.get("file")
    if args.get("lesson") and not rel:
        rel = f"lessons/{args['lesson']}.md"
    if not rel:
        raise DirectiveError(
            "@include requires path=<rel-to-presentation/> or lesson=<id>"
        )
    target = (PRES / rel).resolve()
    try:
        target.relative_to(PRES.resolve())
    except ValueError as tip:
        raise DirectiveError(f"@include path escapes presentation/: {rel}") from tip
    if not target.is_file():
        raise DirectiveError(f"@include missing file: {rel}")
    if target in stack:
        parts = []
        for p in stack + [target]:
            try:
                parts.append(str(p.relative_to(PRES)))
            except ValueError:
                parts.append(str(p))
        raise DirectiveError(f"@include cycle: {' -> '.join(parts)}")
    text = target.read_text(encoding="utf-8")
    return expand_text(text, journey_path=journey_path, stack=stack + [target])


def expand_text(
    text: str,
    *,
    journey_path: Path = DEFAULT_JOURNEY,
    stack: list[Path] | None = None,
) -> str:
    def repl(match: re.Match[str]) -> str:
        name = match.group(1)
        try:
            args = parse_args(match.group(2))
            if name == "include":
                return expand_include(
                    args, journey_path=journey_path, stack=stack
                ).rstrip("\n")
            if name in EXPANDERS:
                if name == "steckbrief":
                    return EXPANDERS[name](args, journey_path=journey_path).rstrip("\n")
                return EXPANDERS[name](args).rstrip("\n")
            return expand_generic(name, args).rstrip("\n")
        except DirectiveError:
            raise
        except SystemExit as tip:
            raise DirectiveError(str(tip)) from tip

    return DIRECTIVE_RE.sub(repl, text)


def expand_file(path: Path, *, journey_path: Path = DEFAULT_JOURNEY) -> str:
    return expand_text(path.read_text(encoding="utf-8"), journey_path=journey_path)


def parse_slide_meta(text: str) -> dict[str, str]:
    m = SLIDE_META_RE.search(text)
    if not m:
        return {}
    meta: dict[str, str] = {}
    for part in re.split(r"[,\n]+", m.group(1).strip()):
        part = part.strip()
        if not part:
            continue
        if "=" in part:
            key, _, value = part.partition("=")
        elif ":" in part:
            key, _, value = part.partition(":")
        else:
            continue
        meta[key.strip().lower()] = value.strip().strip("'\"")
    if "kind" in meta:
        kind = meta["kind"]
        meta["kind"] = KIND_ALIASES.get(kind, KIND_ALIASES.get(kind.lower(), kind.lower()))
    return meta


def icon_from_kind(kind: str | None) -> str:
    if kind == "listen":
        return "💡"
    if kind == "work":
        return "🏝️"
    if kind == "discuss":
        return "💬"
    return ""


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} <markdown-file>", file=sys.stderr)
        sys.exit(2)
    try:
        sys.stdout.write(expand_file(Path(sys.argv[1])))
    except DirectiveError as err:
        print(f"directives: {err}", file=sys.stderr)
        sys.exit(1)
