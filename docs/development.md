# Development Process

This document describes the development workflow for ark using Claude Code (CC).

---

## Overview

ark uses **Spec Driven Development (SDD)** with CC as the primary development assistant.
All feature work is traceable to a GitHub issue with Acceptance Criteria.

---

## Tools & Their Roles

| CC Feature | Role in this project |
|------------|----------------------|
| `CLAUDE.md` | Harness — encodes constraints, conventions, and guardrails that CC always follows |
| `docs/concept.md` | Source of truth for design goals — consulted before any non-trivial implementation |
| `docs/adr/` | Decision log — records why choices were made (Spec, Design, Implementation, Process types) |
| GitHub Issues | Unit of work — `spec`, `implementation`, and `process` issues, each with Acceptance Criteria |
| Plan mode | Bridge between spec and implementation — design before writing code |
| `AskUserQuestion` | Decision gate — surfaces trade-offs and gets explicit sign-off before proceeding |
| Tasks | In-conversation progress tracking for multi-step implementations |
| Memory | Persists decisions and context across conversations |
| Agents | Parallelizable research or isolated sub-problems |

---

## Development Workflow

### Phase 1 — Spec

1. Open a GitHub issue (label: `spec`) for each design question
2. Discuss options in conversation; use `AskUserQuestion` when trade-offs exist
3. Update `docs/concept.md` with the decision
4. Fill in the `## Acceptance Criteria` checklist on the issue
5. Write a `Type: Spec` ADR in `docs/adr/`
6. Close the spec issue
7. Open a corresponding `implementation` issue referencing the spec issue and ADR

> Issues without Acceptance Criteria are not ready for implementation.
> Do not enter Plan mode until AC is defined and the spec ADR is written.
>
> **Triage rule**: if a closed spec issue has no corresponding open `implementation` issue,
> create one before starting any other work. When asked "what's next", check for this gap
> first and recommend filling it.

### Phase 2 — Design

1. Open a GitHub issue (label: `implementation`) referencing the spec issue and ADR
2. Enter Plan mode (`/plan`)
3. Read the spec issue, AC, and relevant ADRs
4. Consult `docs/concept.md` for architectural context
5. Design the implementation approach; use `AskUserQuestion` if multiple viable approaches exist
6. If architectural choices are made (package structure, interfaces, library selection), write a `Type: Design` ADR
7. Exit Plan mode with user approval

**Plan mode exit checklist** (must be complete before starting implementation):
- [ ] Acceptance Criteria reviewed and understood
- [ ] Design ADR written — or explicitly noted as not needed (purely mechanical change)
- [ ] User has approved the plan

### Phase 3 — Implementation

1. Use Tasks to break the work into verifiable steps
2. Implement against the Acceptance Criteria — not more, not less
3. Run CI locally before pushing:
   ```bash
   go build ./...
   go vet ./...
   golangci-lint run
   go test -race ./...
   ```
4. Commit with `Closes #N` to link the implementation to the issue:
   ```
   feat: implement ark init clone flow

   Closes #N
   ```
5. Push and confirm CI passes on GitHub Actions

### Phase 4 — Review

1. Verify all Acceptance Criteria checkboxes are met by the implementation
2. If a criterion turns out to be wrong or missing, update the issue before closing
3. Issue is closed automatically when the commit lands on master

---

### Process Improvement Workflow

Process issues track friction or failures observed in the development workflow itself.
This is Harness Engineering — the process of iteratively improving the harness.

1. Open a GitHub issue (label: `process`) describing the observed problem
2. Fill in the `## Acceptance Criteria` checklist (what "done" looks like)
3. Implement the change (docs, CLAUDE.md, templates, etc.)
4. Write a `Type: Process` ADR capturing the decision and motivation
5. Commit with `Closes #N` and push

> Process issues follow the same AC-first rule as spec and implementation issues.
> A `Type: Process` ADR is required when the change affects how the harness operates
> (workflow steps, guardrails, templates). Skip only for purely cosmetic doc fixes.

---

## Issue Structure

Use the GitHub issue templates in `.github/ISSUE_TEMPLATE/` — one for each label type.

### spec

```markdown
## Question / Context

(What is the design question or feature being specified?)

## Options / Considerations

(What are the alternatives? What are the trade-offs?)

## Acceptance Criteria

- [ ] Criterion one
- [ ] Criterion two
```

### implementation

```markdown
## Overview

(Brief description of what needs to be implemented)

Refs: #N (spec issue), ADR-XXXX

## Acceptance Criteria

- [ ] Criterion one
- [ ] Criterion two
- [ ] CI passes
```

### process

```markdown
## Problem

(What friction, failure, or inefficiency was observed?)

## Proposed Change

(What will be different after this is resolved?)

## Acceptance Criteria

- [ ] Criterion one
- [ ] Criterion two
```

---

## Guardrails

These rules keep CC on track and prevent common LLM drift patterns:

- **No implementation without AC** — if an issue has no Acceptance Criteria, define them first
- **No silent scope expansion** — implement exactly what the AC specifies, nothing more
- **concept.md is authoritative** — if a request conflicts with the design goals, flag it via `AskUserQuestion` before proceeding
- **Security over convenience** — encryption correctness and credential safety take priority over UX shortcuts
- **CI is the definition of done** — a feature is not complete until all CI checks pass

---

## File Locations Reference

| Path | Purpose |
|------|---------|
| `CLAUDE.md` | CC harness (conventions, guardrails, SDD rules) |
| `docs/concept.md` | Design goals and architecture decisions |
| `docs/development.md` | This file — development process |
| `.github/workflows/ci.yml` | CI pipeline |
| `.golangci.yml` | Linter configuration |
