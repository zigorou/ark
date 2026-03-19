# ADR-0005: Uninitialized Vault Detection — Design

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Type** | Design |
| **Issue** | closes #6 |
| **Refs** | [ADR-0003](0003-uninitialized-vault-error.md), spec #5 |

## Context

ADR-0003 decided that all commands except `ark init` must check vault initialization at
startup and exit 1 with a prescribed message. This ADR records the design choices made
during implementation: package structure, error type, Cobra hook wiring, and error output
control.

## Discussion

### Package: `internal/vault`

The vault check logic is placed in `internal/vault` — the first `internal/` package in
the project. Three functions are exported:

- `CheckInitialized(dir string) error` — checks dir, `.git/`, and `.sops.yaml`
- `ResolveDir(flagValue string) (string, error)` — resolves vault dir by priority
- `DefaultDir() (string, error)` — returns `~/.ark`

### Return type: `error` vs `bool`

> **For** returning `error`
> - Each missing condition (dir / `.git/` / `.sops.yaml`) can be wrapped with context
> - `errors.Is(err, vault.ErrVaultNotInitialized)` enables precise assertions in tests
> - Future per-condition error messages become straightforward to add

> **Against** returning `error`
> - Slightly more verbose call sites compared to `if !vault.IsInitialized(dir)`

**Conclusion**: return `error` wrapping `ErrVaultNotInitialized`.

### Error type: sentinel error

```go
var ErrVaultNotInitialized = errors.New("vault not initialized")
```

A sentinel error is sufficient — no structured fields are needed at this stage.
`errors.Is` provides clean matching in both production code and tests.

### Cobra hook wiring: `PersistentPreRunE`

The check is installed on `rootCmd.PersistentPreRunE` so it runs before every subcommand.
`ark init` is excluded by overriding `PersistentPreRunE` on `initCmd` with a no-op:

```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
```

> **For** PersistentPreRunE override on initCmd
> - In Cobra v1, a child's `PersistentPreRunE` fully replaces the parent's — no special
>   exclusion list needed in the parent
> - `initCmd` is the only command that must skip the check; the pattern is easy to extend

> **Against**
> - Requires every future command that skips the check to declare its own no-op, which
>   could be forgotten

**Conclusion**: override on `initCmd` only; document the pattern in code comments.

### Error output control

`rootCmd.SilenceErrors = true` prevents Cobra from printing the error string automatically.
The `PersistentPreRunE` handler calls `cmd.PrintErrln()` directly to emit the ADR-0003
message to stderr, then sets `cmd.SilenceUsage = true` to suppress usage output on error.
This gives full control over what appears on stderr.

## Decision

- [x] Vault check logic lives in `internal/vault` with `CheckInitialized`, `ResolveDir`, and `DefaultDir`
- [x] Return `error` wrapping `ErrVaultNotInitialized` (not `bool`)
- [x] Sentinel error `var ErrVaultNotInitialized = errors.New("vault not initialized")`
- [x] Check installed via `rootCmd.PersistentPreRunE`; `initCmd` overrides with a no-op
- [x] `SilenceErrors: true` on `rootCmd`; error message emitted via `cmd.PrintErrln()`

## Consequences

- `internal/vault` is the first `internal/` package; future vault-related helpers should
  be added here before creating new packages
- Any command added in the future that must bypass the vault check must declare its own
  `PersistentPreRunE` no-op — this is not enforced automatically
- `cmd/init.go` still has `MarkFlagRequired("repo")` which contradicts ADR-0002
  (repo-less init is valid); this will be resolved when `ark init` is fully implemented
