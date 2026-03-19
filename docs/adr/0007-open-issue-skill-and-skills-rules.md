# ADR-0007: `/open-issue` Skill and Skill Authoring Rules

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Type** | Process |
| **Issue** | — (not tied to a tracked issue) |
| **Refs** | [ADR-0006](0006-next-skill.md) — `/next` skill (prior skill addition) |

## Context

Two harness gaps were identified after the `/next` skill was added:

1. **No skill for creating GitHub issues** — creating spec/impl/process issues required
   manually composing `gh issue create` with the right label, title prefix, and template body.
   This is repetitive and error-prone.

2. **No authoring rules for skill files** — the `.claude/commands/` directory had no
   guardrails on frontmatter structure, description style, tool permissions, or the
   requirement to use `skill-creator` for non-trivial changes. A bad description reduces
   trigger accuracy; over-permissive `allowed-tools` violates least-privilege.

## Discussion

### Issue creation: single skill vs. type-specific skills

> **For** single `/open-issue` skill
> - One entry point is easier to discover; the user doesn't need to remember three skill names
> - Type is passed as a positional arg — `AskUserQuestion` fills the gap if omitted
> - Shared logic (title prefixing, template loading, AC reminder) lives in one place

> **Against** single `/open-issue` skill
> - Slightly more prompt complexity to branch on three types
> - Longer `allowed-tools` list than any individual type-specific skill would need

> **For** three separate skills (`open-spec`, `open-impl`, `open-process`)
> - Each skill is simpler and more focused
> - `allowed-tools` can be tighter per skill

> **Against** three separate skills
> - Three files to maintain; shared changes (e.g. AC reminder wording) must be applied in all three
> - User must remember which skill name maps to which type

Single skill chosen: discoverability and maintainability outweigh the slight prompt complexity.

### Skills rules: which checks to enforce

Four checks were selected:

| Check | Rationale |
|-------|-----------|
| Frontmatter required fields | Skills without `description` or `allowed-tools` fail silently or load with wrong permissions |
| `description` imperative style | Consistent style enables accurate trigger matching; article/pronoun starts reduce clarity |
| `skill-creator` for non-trivial changes | Ensures evals are run and trigger accuracy is validated before a skill ships |
| `allowed-tools` least privilege | Read-only skills must not carry `Edit`/`Write`; prevents accidental mutations |

## Decision

- [x] Add `.claude/skills/open-issue/SKILL.md` — single skill handling spec/impl/process types
- [x] Skill reads the corresponding `.github/ISSUE_TEMPLATE/*.md`, gathers missing info via
      `AskUserQuestion`, then runs `gh issue create`
- [x] Skill is non-destructive beyond issue creation: does not start Plan mode or write code
- [x] Add `.claude/rules/skills.md` scoped to `.claude/skills/**/SKILL.md`
- [x] Rules enforce: frontmatter required fields, description imperative style,
      `skill-creator` for non-trivial changes, `allowed-tools` least privilege
- [x] Migrated `/next` and `/open-issue` from `.claude/commands/` (legacy) to `.claude/skills/`
      (recommended format) per issue #10

## Consequences

- Developers can run `/open-issue spec|impl|process` to create a correctly-labelled,
  templated issue without touching `gh` directly
- The AC reminder after issue creation reinforces the SDD guardrail (no Plan mode without AC)
- New skill authors get immediate in-context guidance on the four required properties
- `skill-creator` is now the canonical path for skill authoring, reducing ad-hoc edits
