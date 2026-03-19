# ark — Concept

## Overview

- **Repository name**: `ark` (`github.com/zigorou/ark`)
- **Inspired by**: 1Password CLI (`op`)
- **Implementation language**: Go (easy single-binary distribution)

---

## Motivation

1Password CLI (`op`) offers excellent functionality and UX, but is hard to adopt for the following reasons:

- **Paid**: Too costly for personal use
- **Centralized**: The vault depends on 1Password's servers
- **Service-dependent**: Unavailable offline or during outages

SOPS + age is secure and git-friendly, but the UX is too low-level.

**The goal is to build an `op`-like CLI tool that uses a GitHub private repo as the vault**, combining the best of both worlds.

---

## Design Goals

- Vault storage is **encrypted files in a GitHub private repo** (SOPS + age)
- CLI UX close to `op` (`run`, `read`, `inject`, URI scheme)
- Free, open-source, self-hosted
- Targeted at individuals and small teams

---

## CLI Interface

### URI Scheme

```
ark://<category>/<item>/<field>
```

Examples:
```
ark://obsidian/work/api_key
ark://aws/prod/access_key_id
ark://github/personal/token
```

### Commands

```bash
# Read a secret value
ark read ark://obsidian/work/api_key

# Inject secrets as environment variables and run a command (equivalent to op run)
ark run \
  --env OBSIDIAN_API_KEY=ark://obsidian/work/api_key \
  --env AWS_ACCESS_KEY_ID=ark://aws/prod/access_key_id \
  -- some-command

# Template expansion (equivalent to op inject)
ark inject -i config.tpl -o config.json

# Edit a vault file (decrypt → open in editor → re-encrypt on save)
ark edit obsidian

# Add or update a secret
ark set ark://obsidian/work/new_key "value"

# Initialize the vault
ark init --repo github.com/zigorou/my-vault
```

---

## Vault Structure

A GitHub private repo (e.g. `zigorou/my-vault`) contains encrypted YAML files:

```
my-vault/
  .sops.yaml
  obsidian.enc.yaml
  aws.enc.yaml
  github.enc.yaml
```

Decrypted YAML structure:

```yaml
# obsidian.enc.yaml
work:
  api_key: 9aff4dc9...
  base_url: https://127.0.0.1:27124
personal_core:
  api_key: 8292aaf9...
  base_url: https://127.0.0.1:27125
```

---

## Implementation

| Layer      | Technology                                                     |
|------------|----------------------------------------------------------------|
| Encryption | `filippo.io/age` + `github.com/getsops/sops/v3` (Go libraries) |
| Storage    | GitHub private repo (git)                                      |
| CLI        | Go                                                             |

Both age and SOPS are embedded as Go libraries, allowing distribution as a single binary.
No separate installation of the `sops` or `age` CLI tools is required.
ark acts as a SOPS wrapper + URI resolver + UX layer for `run`/`inject`.

---

## Comparison with Existing Tools

|                    | SOPS          | 1Password CLI (`op`) | ark                  |
|--------------------|---------------|----------------------|----------------------|
| Storage            | Local git     | Central server       | GitHub private repo  |
| UX                 | Low-level     | High-level           | High-level           |
| URI scheme         | None          | `op://vault/item/field` | `ark://category/item/field` |
| Multi-machine sync | Manual        | Auto sync            | `git pull`           |
| Free               | ✓             | ✗                    | ✓                    |
| Offline            | ✓             | ✗                    | △ (if already cloned) |

---

## Out of Scope (v1)

- Team-level access control (vault access is delegated to git repository permissions)
- Web UI
- Browser integration
- SSH key management (equivalent to `op ssh-agent`)
- TOTP / OTP

---

## Security Decisions

- **Pre-commit hook**: A pre-commit hook will be added to the vault repo to prevent accidentally committing unencrypted secrets.
- **Rekey deferral**: Automatic rekeying (re-encrypting vault files with new age keys) is deferred to a future version. For v1, rekeying requires manually re-encrypting files.

---

## Open Questions

- Timing of automatic vault repo clone/pull (pull on every startup, or manual pull?)
- Key management strategy (age keys are local to each machine; backup via Google Passwords or similar as a string)
