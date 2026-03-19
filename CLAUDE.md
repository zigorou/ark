# ark

## Overview

**ark** is an open-source secrets manager CLI that uses a GitHub private repository as encrypted vault storage (via SOPS + age). It provides a 1Password-like UX — URI scheme, `run`, `inject`, and `edit` — without the cost or centralized server dependency.

Module path: `github.com/zigorou/ark`

---

## Build & Run

```bash
# Build the binary
go build -o ark .

# Run directly
go run . --help

# Run tests (with race detector)
go test -race ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

---

## Project Structure

```
ark/
  main.go                   # Entry point — calls cmd.Execute()
  cmd/
    root.go                 # Root command (ark), persistent flags
    read.go                 # ark read <uri>
    run.go                  # ark run [--env KEY=ark://...] -- <command>
    inject.go               # ark inject -i <template> -o <output>
    edit.go                 # ark edit <category>
    set.go                  # ark set <uri> <value>
    init.go                 # ark init --repo <repository>
  internal/                 # Private packages (not importable externally)
  scripts/
    update-synopsis.sh      # Regenerates README.md SYNOPSIS from CLI help
    hooks/                  # Claude Code hook scripts
  docs/
    concept.md              # Project design document (source of truth for goals)
  .github/
    workflows/
      ci.yml                # GitHub Actions CI pipeline
  .golangci.yml             # golangci-lint configuration
```

---

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `filippo.io/age` | Encryption (to be added) |
| `github.com/getsops/sops/v3` | SOPS integration (to be added) |

---

## Documentation Language

All documentation and code comments are written in **English**.
Commit messages are written in **English**.

---

## Go Conventions

Follow [Effective Go](https://go.dev/doc/effective_go) and the [Google Go Style Guide](https://google.github.io/styleguide/go/).

> Style conventions (formatting, naming, comments) and testing standards are defined in
> `.claude/rules/go-style.md` and `.claude/rules/go-testing.md` (path-scoped rules).

### Error Handling
- Always wrap errors with context: `fmt.Errorf("vault: decrypt %s: %w", path, err)`
- Never use `_` to discard errors silently
- Return errors rather than panicking (panics only for truly unrecoverable states)
- No naked returns

### Package Design
- Define interfaces in the **consumer** package, not the producer
- Prefer small, focused packages over large ones
- `internal/` packages for implementation details not intended for external import

### Testing Requirement

All new functionality must be accompanied by tests. CI blocks merge on test failure.

---

## CI Guardrails

The CI pipeline (`.github/workflows/ci.yml`) runs on every push and PR:

1. `go build ./...` — compilation check
2. `go vet ./...` — static analysis
3. `golangci-lint run` — linting
4. `go test -race ./...` — tests with race detector

**A feature is not complete until CI passes.**

---

## Anti-Drift Rules

These rules prevent the three primary LLM drift patterns in harness engineering:

### Specification Gaming
> *The model satisfies the letter of the instructions but misses the intent.*

- Before implementing, explicitly confirm: "What was asked" vs "What I'm building" — if they differ, use `AskUserQuestion`
- Implement exactly what was requested — not more, not less
- Do not add features, refactors, comments, tests, or error handling for code **not** touched in the task
- If the simplest correct solution exists, prefer it over an elegant but over-engineered one

### Goal Misgeneralization
> *The model pursues an objective that differs from the actual project goal.*

- `docs/concept.md` is the **source of truth** for design goals — consult it before any non-trivial implementation
- If a user request would cause architectural drift from the stated design goals (e.g., swapping core dependencies, changing the vault model), flag it via `AskUserQuestion` **before** implementing
- Do not silently reinterpret the project's direction based on conversation history alone
- When uncertain whether a change aligns with project goals, ask rather than assume

### Sycophancy
> *The model agrees with the user rather than giving an honest assessment.*

- Never implement a proposal that has security, correctness, or architectural issues just because the user requested it
- Raise concerns **before** starting implementation, not after
- In `AskUserQuestion` options, provide honest pros/cons even if the user appears committed to a specific choice
- Security concerns (credential exposure, encryption correctness) **always** take priority over convenience or user preference
- If the user's framing of a problem appears incorrect, say so directly before proceeding

---

## Claude Code Development Workflow

> **Harness health**: If CLAUDE.md exceeds ~300 lines, suggest extracting repeated procedures (e.g. commit flow, issue creation) into CC skills rather than adding more inline content.
>
> **CLAUDE.md vs rules**: Keep critical guardrails and architectural constraints here. Move style conventions and non-critical guidelines to `.claude/rules/` — preferably with `paths:` frontmatter for conditional loading. If broken it's a style issue, it belongs in rules; if broken it causes bugs or security issues, it belongs here.

This project uses Claude Code with harness engineering for development:

- **CLAUDE.md** (this file): the harness — defines constraints, conventions, and guardrails
- **Hooks** (`.claude/settings.json`): automated guardrails triggered on tool calls (e.g., SYNOPSIS auto-update)
- **Plan mode**: use for any non-trivial implementation before writing code
- **Agents**: use for parallelizable research tasks or isolated sub-problems
- **`AskUserQuestion`**: required when multiple viable approaches exist — always include pros/cons and a recommendation

See `docs/development.md` for the full development process documentation.

---

## Spec Driven Development (SDD)

This project follows a Spec Driven Development workflow. All feature work must be traceable to a GitHub issue with Acceptance Criteria.

### Issue Structure

Every spec issue **must** contain an `## Acceptance Criteria` section with a checklist before implementation begins:

```markdown
## Acceptance Criteria
- [ ] Criterion one
- [ ] Criterion two
```

- Issues without Acceptance Criteria are not ready for implementation.
- Acceptance Criteria must be written and agreed upon **before** entering Plan mode.
- When closing an issue via commit, use `Closes #N` in the commit message.

### SDD Workflow

```
[Spec phase]
spec issue (label: spec)
  → Discuss & decide (AskUserQuestion for trade-offs)
  → Update concept.md with the decision
  → Fill in Acceptance Criteria
  → Write ADR (Type: Spec)
  → Close spec issue
  → Open implementation issue

[Design phase]
implementation issue (label: implementation, refs spec issue)
  → Plan mode (design against AC)
  → Write ADR (Type: Design) if architectural choices are made
  → Exit plan mode with user approval

[Implementation phase]
  → Implement (Tasks for step tracking within a conversation)
  → CI passes → Closes #N in commit

[Process phase]
process issue (label: process)
  → Describe problem and proposed change
  → Fill in Acceptance Criteria
  → Implement (update docs, CLAUDE.md, templates, etc.)
  → Write ADR (Type: Process)
  → Closes #N in commit
```

### Issue Labels

| Label | Purpose |
|-------|---------|
| `spec` | Design question — what the tool does. Close after AC is defined and ADR written. |
| `implementation` | Coding task — how to build it. References spec issue and ADR. Close after CI passes. |
| `process` | Harness improvement — workflow friction or failure. Close after docs updated and Process ADR written. |

### Guardrails

- Do not start Plan mode for an issue that lacks Acceptance Criteria — ask the user to define them first.
- Do not mark an issue closed unless all Acceptance Criteria checkboxes are satisfiable by the implementation.
- If implementation reveals a criterion is wrong or missing, update the issue before proceeding.

### ADR Requirements

**ADR-worthy test** — ask this before writing any ADR:
> "Would a future developer or CC reasonably ask *why we made this choice*, and would the commit message alone fail to answer it?"
> If **No** → skip the ADR.

Changes that do **not** warrant an ADR:
- Renames, moves, or path updates that follow a prior decision
- Changes where the only alternative was not actually viable (false trade-off)
- Anything fully explained by the commit message

An ADR (`docs/adr/NNNN-<slug>.md`) **must** be written when the test passes, in these situations:

1. **Closing a spec issue**: record the decision and discussion as a `Type: Spec` ADR before closing.
2. **Design decisions in Plan mode**: if Plan mode produces architectural choices (package structure, interface design, library selection), record them as a `Type: Design` ADR before implementing.
3. **Closing a process issue**: record the problem, decision, and motivation as a `Type: Process` ADR before closing.
4. **Before merging a PR**: if the PR introduces a non-obvious decision not yet captured in an ADR, write it before merge.

Use `docs/adr/0000-template.md` as the template.
