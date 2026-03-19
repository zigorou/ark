# ADR-0002: `ark init` Flow, `--repo` Format, and Clone Implementation

| | |
|---|---|
| **Date** | 2026-03-20 |
| **Status** | ✅ Accepted |
| **Issues** | closes #1, closes #2 |

## Context

`ark init` is the entry point for setting up a vault. Two questions needed resolving:

1. Should `ark init` support creating a new GitHub repo, or only cloning an existing one? (#2)
2. What URL formats should `--repo` accept, and how should cloning be implemented? (#1)

## Discussion

### Init flow: clone-only vs create+clone (#2)

The concept defines `git remote origin` as **strongly recommended but not required**.
This means a local-only vault (no GitHub remote) is a valid use case.

> **For** clone-only (v1 creates nothing on GitHub)
> - No GitHub API dependency
> - Simpler implementation

> **Against** clone-only
> - User must manually create the GitHub repo before running `ark init`
> - Since GitHub is ark's storage backend by design, repo creation is a natural responsibility

Given that GitHub is a first-class concept in ark (not an optional integration),
delegating repo creation entirely to the user felt inconsistent. `go-github` as a
Go library dependency is appropriate — it's the same class of decision as embedding
`filippo.io/age` or `getsops/sops`.

Auth for GitHub API (repo creation) is handled via `ARK_GITHUB_TOKEN`.
When the token is absent, ark exits with an actionable error rather than silently failing.

**Conclusion**: `ark init` supports both scenarios, gated on `--repo` and token presence.

### `--repo` format (#1)

Three formats were considered:

| Input | Handling |
|-------|---------|
| `github.com/user/repo` | normalized to `https://github.com/user/repo` |
| `https://github.com/user/repo` | used as-is |
| `git@github.com:user/repo` | used as-is |

All three are accepted. Short form normalizes to HTTPS.
SSH users can specify `git@github.com:user/repo` directly — no config toggle needed.

### Clone implementation: `git` CLI vs `go-git` (#1)

| | `git` CLI shelling out | `go-git` embedded |
|---|---|---|
| SSH keys & credential helpers | ✅ user's git config applies transparently | ❌ requires custom implementation; some SSH key formats unsupported |
| External dependency | ⚠️ requires `git` in PATH | ✅ pure Go, single binary |
| Implementation cost | ✅ low | ❌ high — many edge cases |
| Reliability | ✅ battle-tested | ⚠️ SSH support is immature |

For a developer-targeted tool, `git` in PATH is a safe assumption.
`go-git`'s SSH limitations would be a critical gap for v1 users.

**Conclusion**: shell out to `git clone`.

## Decision

- [x] `ark init` (no flags) runs `git init` in `ARK_VAULT_DIR` and warns about missing remote on every subsequent invocation
- [x] `ark init --repo` clones the repo if it exists on GitHub
- [x] `ark init --repo` creates a **private** GitHub repo via `go-github` if it does not exist and `ARK_GITHUB_TOKEN` is set
- [x] `ark init --repo` exits with error + instructions if repo does not exist and `ARK_GITHUB_TOKEN` is unset
- [x] `--repo` accepts `github.com/user/repo`, `https://...`, and `git@github.com:...`
- [x] `github.com/user/repo` is normalized to `https://github.com/user/repo`
- [x] Clone is implemented by shelling out to `git clone`
- [x] Missing `git` in PATH → exit 1 with actionable error message

## Consequences

- `go-github` is added as a dependency (scoped to `ark init` repo creation path)
- `ARK_GITHUB_TOKEN` is a new recognized environment variable; document in README
- `ark init` on a non-empty `ARK_VAULT_DIR` exits with an error (see Vault Directory rules in `concept.md`)
- Multi-machine rekeying flow (adding a second machine's key) is out of v1 scope — see [ADR-0001](0001-vault-structure.md)
