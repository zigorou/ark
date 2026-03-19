---
paths:
  - "**/*.go"
---

# Go Style Conventions

## Formatting

- All code must be `gofmt`-formatted.
- Use `goimports` for import grouping: stdlib → external → internal.
- Local import prefix for grouping: `github.com/zigorou/ark`

## Naming

- Packages: lowercase, single word, no underscores (`internal/vault`, not `internal/vault_store`)
- Exported types and functions: PascalCase with godoc comment
- Error variables: `var ErrXxx = errors.New("...")`
- Error strings: lowercase, no trailing punctuation (`"failed to open file"`, not `"Failed to open file."`)
- Interfaces: single-method interfaces end in `-er` (`Reader`, `Encrypter`)

## Comments

- Every exported symbol needs a godoc comment starting with the symbol name.
- Do not add comments for unexported, self-evident code.
