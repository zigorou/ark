# ark

> A 1Password-like secrets manager CLI powered by GitHub private repos, SOPS, and age.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Synopsis

<!-- SYNOPSIS:START -->
```
ark [--vault <path>] <command> [arguments]

ark edit <category>
ark init --repo <repository>
ark inject -i <template> -o <output>
ark read <uri>
ark run [--env <NAME=ark://...>]... -- <command> [args...]
ark set <uri> <value>
```
<!-- SYNOPSIS:END -->

## Description

**ark** stores secrets in a **GitHub private repository** as SOPS-encrypted YAML files,
and provides a [1Password CLI (`op`)](https://developer.1password.com/docs/cli/)-like
UX for reading, injecting, and editing those secrets.

Key differences from existing tools:

| | SOPS | 1Password CLI | ark |
|---|---|---|---|
| Storage | Local git | Central server | GitHub private repo |
| UX | Low-level | High-level | High-level |
| URI scheme | None | `op://vault/item/field` | `ark://category/item/field` |
| Multi-machine | Manual | Auto sync | `git pull` |
| Free | ✓ | ✗ | ✓ |

## Installation

```bash
go install github.com/zigorou/ark@latest
```

Or build from source:

```bash
git clone https://github.com/zigorou/ark.git
cd ark
go build -o ark .
```

## URI Scheme

```
ark://<category>/<item>/<field>
```

Examples:

```
ark://aws/prod/access_key_id
ark://obsidian/work/api_key
ark://github/personal/token
```

## Usage

### Initialize vault

Link a GitHub private repository as your vault:

```bash
ark init --repo github.com/yourname/my-vault
```

### Read a secret

```bash
ark read ark://aws/prod/access_key_id
```

### Run a command with injected secrets

```bash
ark run \
  --env AWS_ACCESS_KEY_ID=ark://aws/prod/access_key_id \
  --env AWS_SECRET_ACCESS_KEY=ark://aws/prod/secret_access_key \
  -- aws s3 ls
```

### Inject secrets into a template

Given a template file `config.tpl`:

```
api_key: ark://obsidian/work/api_key
base_url: ark://obsidian/work/base_url
```

Run:

```bash
ark inject -i config.tpl -o config.json
```

### Edit a vault file

Opens the decrypted file in `$EDITOR`, re-encrypts on save:

```bash
ark edit obsidian
```

### Set a secret

```bash
ark set ark://obsidian/work/new_api_key "my-secret-value"
```

## Vault Structure

The GitHub private repository contains SOPS-encrypted YAML files:

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
```

## Configuration

ark looks for a configuration file at `~/.config/ark/config.yaml`:

```yaml
vault: ~/repos/my-vault   # path to the cloned vault repository
```

## Development

This project is developed with [Claude Code](https://claude.ai/claude-code) using
a harness engineering approach — see [docs/development.md](docs/development.md).

### Local setup

After cloning, install the git hooks and recommended tools:

```bash
bash scripts/install-hooks.sh
brew install gitleaks golangci-lint
```

The pre-commit hook runs `gitleaks protect --staged` to prevent accidentally
committing secrets.

## License

[MIT](LICENSE)
