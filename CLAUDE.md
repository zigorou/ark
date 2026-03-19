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

### Formatting
- All code must be `gofmt`-formatted. Use `goimports` (import grouping: stdlib → external → internal).
- Local import prefix for grouping: `github.com/zigorou/ark`

### Naming
- Packages: lowercase, single word, no underscores (`internal/vault`, not `internal/vault_store`)
- Exported types and functions: PascalCase with godoc comment
- Error variables: `var ErrXxx = errors.New("...")`
- Error strings: lowercase, no trailing punctuation (`"failed to open file"`, not `"Failed to open file."`)
- Interfaces: single-method interfaces end in `-er` (`Reader`, `Encrypter`)

### Error Handling
- Always wrap errors with context: `fmt.Errorf("vault: decrypt %s: %w", path, err)`
- Never use `_` to discard errors silently
- Return errors rather than panicking (panics only for truly unrecoverable states)
- No naked returns

### Package Design
- Define interfaces in the **consumer** package, not the producer
- Prefer small, focused packages over large ones
- `internal/` packages for implementation details not intended for external import

### Comments
- Every exported symbol needs a godoc comment starting with the symbol name
- Do not add comments for unexported, self-evident code

---

## Testing Requirements

All new functionality must be accompanied by tests. CI blocks merge on test failure.

### Standards
- Use **table-driven tests** with `t.Run(name, ...)` for multiple cases
- Test file: `foo_test.go` alongside `foo.go` in the same package
- External test package (`package foo_test`) for black-box testing exported APIs
- Use `t.Helper()` in all test helper functions
- Use `t.Parallel()` where tests are independent
- Run with `-race` flag to detect data races
- No global mutable state in tests; use `t.TempDir()` for temp files

### Coverage Goals
- Core logic packages (`internal/*`): aim for ≥ 80% coverage
- CLI command packages (`cmd/*`): integration-style tests preferred over unit tests

### Test Naming
```go
func TestFunctionName(t *testing.T) { ... }        // basic
func TestFunctionName_scenario(t *testing.T) { ... } // scenario variant
```

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

This project uses Claude Code with harness engineering for development:

- **CLAUDE.md** (this file): the harness — defines constraints, conventions, and guardrails
- **Hooks** (`.claude/settings.json`): automated guardrails triggered on tool calls (e.g., SYNOPSIS auto-update)
- **Plan mode**: use for any non-trivial implementation before writing code
- **Agents**: use for parallelizable research tasks or isolated sub-problems
- **`AskUserQuestion`**: required when multiple viable approaches exist — always include pros/cons and a recommendation

See `docs/development.md` for the full development process documentation.
