# ADR-0008: `/update-process` Skill for Harness Documentation Sync

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Type** | Process |
| **Issue** | closes #11 |
| **Refs** | [ADR-0007](0007-open-issue-skill-and-skills-rules.md) — prior skill additions |

## Context

After process changes (new skills, updated rules, workflow modifications), multiple
documentation files must be kept in sync:

- `docs/workflow.md` — Mermaid workflow diagram
- `docs/development.md` — development process reference
- `CLAUDE.md` — harness config (specifically "Claude Code Development Workflow" and "SDD Workflow" sections)

This sync has been done manually, leading to drift. The `/next` and `/open-issue` skills
reduced friction in triage and issue creation respectively; a `/update-process` skill
closes the remaining manual step in the process improvement loop.

## Discussion

> **For** a single skill that covers all three docs
> - One invocation updates everything; no need to remember which doc tracks which concern
> - The skill can read git diff to infer what changed, reducing the user's cognitive load
> - Consistent with the project's pattern of single-purpose harness skills

> **Against** a single skill
> - Slightly larger `allowed-tools` surface (needs `Edit` in addition to read-only tools)
> - Risk of over-writing if the inference is wrong — mitigated by mandatory preview step

The approval gate (Step 5) makes the write risk acceptable: the skill never touches a file
without explicit user confirmation.

**Scope of CLAUDE.md updates** was narrowed to two sections:
- "Claude Code Development Workflow" — skill list, hooks, harness tool roles
- "SDD Workflow" — phase steps, issue lifecycle

Other CLAUDE.md sections (Go conventions, CI guardrails, anti-drift rules) are stable and
out of scope to avoid unintended mutations.

## Decision

- [x] Add `.claude/skills/update-process/SKILL.md`
- [x] Skill reads all SKILL.md and rules files plus recent `git diff` to infer changes
- [x] Skill targets: `docs/workflow.md`, `docs/development.md`, and two CLAUDE.md sections
- [x] Skill shows per-doc preview via `AskUserQuestion` before any write
- [x] Skill applies edits only after explicit user approval

## Consequences

- The process improvement loop (observe → open issue → fix → sync docs → commit) is now
  fully skill-assisted at every step: `/next` → `/open-issue` → implement → `/update-process`
- Documentation drift after process changes becomes the exception rather than the norm
- If new documentation files are added to the harness, the skill's target list must be updated
