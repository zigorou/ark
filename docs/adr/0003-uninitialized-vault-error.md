# ADR-0003: Uninitialized Vault Error Policy

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Type** | Spec |
| **Issue** | closes #5 |

## Context

When a user runs any ark command before `ark init`, the vault directory either does not
exist or is in an incomplete state. Without explicit detection, ark would produce
confusing low-level errors. A clear, actionable error is needed.

## Discussion

### What counts as "uninitialized"?

Three distinct states can indicate an uninitialized vault:

| Condition | Meaning |
|-----------|---------|
| `ARK_VAULT_DIR` does not exist | `ark init` was never run |
| exists but no `.git/` | directory exists but is not a git repo |
| git repo but no `.sops.yaml` | cloned manually or setup was interrupted |

All three are treated as uninitialized. The check runs at startup for every command
except `ark init` itself.

### Exit code: dedicated vs exit 1

> **For** a dedicated exit code (e.g. exit 2)
> - Scripts can distinguish uninitialized state from other errors

> **Against** a dedicated exit code
> - ark targets individuals and small teams — scripting around init state is a rare need
> - `git`, `gh`, and most developer CLI tools use exit 1 as the generic error code
> - Managing and documenting a dedicated exit code adds ongoing maintenance cost

**Conclusion**: exit 1, consistent with the broader CLI ecosystem.

## Decision

- [x] Uninitialized is detected if any of: vault dir missing / no `.git/` / no `.sops.yaml`
- [x] All commands except `ark init` check vault initialization at startup
- [x] On uninitialized vault, print to **stderr** and exit 1:

```
ark: vault is not initialized.
Run 'ark init --repo <repository>' to set up a new vault,
or 'ark init' to initialize a local vault without a remote.
```

## Consequences

- A `vault.IsInitialized()` (or equivalent) function is needed in `internal/` and called from the root command's `PersistentPreRunE`
- `ark init` must be explicitly excluded from the initialization check
