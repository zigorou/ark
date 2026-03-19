# ADR-0001: Vault File Structure and SOPS Configuration

- **Date**: 2026-03-20
- **Status**: Accepted
- **Issue**: #4

## Context

ark uses SOPS + age to encrypt secrets stored in a GitHub private repository.
Two structural decisions were needed before implementation could begin:

1. How vault files should be named
2. Where and how `.sops.yaml` should be managed

## Discussion

### File naming: `{category}.enc.yaml` vs `{category}.yaml`

The URI scheme `ark://category/item/field` maps `category` directly to a filename.
The question was whether to include `.enc.` in the filename as a signal that the file is encrypted.

Arguments for `.enc.`:
- Visual indicator in `ls` output that the file is encrypted

Arguments against `.enc.`:
- Users interact via URI only (`ark read`, `ark edit category`) — the filename is never visible
- SOPS-encrypted files are unreadable as plaintext regardless of the extension
- Simpler `path_regex` in `.sops.yaml` without needing to exclude the config file itself
- `{category}.yaml` is cleaner and sufficient

**Conclusion**: `.enc.` is redundant given the CLI abstraction. Dropped.

### `.sops.yaml` location: vault repo vs ark config

Two options were considered:

| | vault repo (`~/.ark/.sops.yaml`) | ark config (`~/.config/ark/`) |
|---|---|---|
| SOPS CLI compatibility | ✅ auto-detected by SOPS | ❌ requires `--config` on every invocation |
| Vault self-containment | ✅ data and encryption rules co-located | ❌ separated |
| Versioned in git | ✅ | ❌ |
| Public key in repo | △ (public key only, no security risk) | ✅ not in repo |

The decisive factor was SOPS CLI compatibility. ark is a SOPS wrapper; power users
should be able to fall back to raw `sops` commands without friction.

**Conclusion**: `.sops.yaml` lives in the vault repo.

### `path_regex` pattern

`\.yaml$` would technically match `.sops.yaml` itself (though SOPS does not encrypt
its own config file). To be explicit and match only category files:

```
^[a-z][a-z0-9_-]*\.yaml$
```

This matches `obsidian.yaml`, `aws.yaml` etc. and excludes dotfiles.

### Multi-machine support (multiple age keys)

Adding a second machine requires:
1. `age-keygen` on the new machine
2. Adding the new public key to `.sops.yaml`
3. Running `sops updatekeys` on every vault file (rekeying)

This is the standard SOPS workflow. For v1, single-machine use is assumed.
A future `ark rekey` command will automate step 3.

## Decision

- Vault files are named **`{category}.yaml`** (no `.enc.` infix)
- **`.sops.yaml` lives in the vault repo** (`~/.ark/`)
- `ark init` auto-generates `.sops.yaml` with `path_regex: ^[a-z][a-z0-9_-]*\.yaml$`
  and the age public key derived from `identity_file`
- **Multi-machine (multi-key) support is out of v1 scope**

## Consequences

- `ark init` must implement age public key derivation from the identity file
- `.sops.yaml` is committed to the vault repo on `ark init`
- Category names must match `^[a-z][a-z0-9_-]*$` (validated at `ark set` / `ark edit`)
- Future ADR needed for `ark rekey` command design
