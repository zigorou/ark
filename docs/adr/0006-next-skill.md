# ADR-0006: `/next` Skill for SDD Workflow Navigation

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Type** | Process |
| **Issue** | closes #9 |
| **Refs** | [ADR-0004](0004-auto-commit-push.md) — prior process ADR |

## Context

In the SDD workflow, deciding "what to work on next" requires manually:
1. Listing all issues and filtering by label and state
2. Cross-referencing closed spec issues against implementation issues to detect triage gaps
3. Applying the harness-health-first priority rule (process > implementation > new spec)

This manual process is error-prone and slow. A `/next` Claude Code skill can automate this triage
and surface the highest-priority item in seconds.

Options considered:
- **A**: Legacy single `.md` skill file under `.claude/commands/next.md`
- **B**: Multi-file skill package (not supported by current CC skill format for project-local skills)

## Discussion

> **For** option A (single `.md` file)
> - Project-local skills use `.claude/commands/<name>.md` format — matches CC conventions
> - No external dependencies; runs entirely via `gh` CLI and built-in tools
> - Simple to maintain; the prompt is the implementation

> **Against** option A
> - All logic is expressed as natural language instructions rather than code — harder to unit-test
> - If the prompt grows complex, it may need refactoring into a structured skill package later

The skill needs only three tools: `Bash` (for `gh issue list`), `Read` (for `docs/concept.md`),
and `Grep` (for cross-reference verification). This is a good fit for the single-file format.

The triage gap detection logic (checking whether any implementation issue body references `#N`)
is expressed as step-by-step instructions in the prompt, which is sufficient for the current scale
of the project (≤100 issues per fetch).

## Decision

- [x] Implement `/next` as `.claude/commands/next.md` (legacy single-file skill)
- [x] Allowed tools: `Bash`, `Read`, `Grep`
- [x] Priority order: triage gap → open process → open implementation → concept.md proposal
- [x] Use `AskUserQuestion` when multiple same-priority candidates exist and choice is non-obvious
- [x] Skill is read-only: recommends actions, never executes them

## Consequences

- Developers can run `/next` at the start of any session to get an immediate triage recommendation
- Triage gaps (closed spec without implementation issue) are surfaced automatically, reducing drift
- The priority rule (harness health before features) is enforced by the skill, not by memory
- If the project grows beyond 100 issues, the `--limit 100` cap in `gh issue list` must be raised
