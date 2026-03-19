#!/usr/bin/env bash
# scripts/hooks/pre-commit.sh
#
# Pre-commit hook: runs gitleaks on staged changes to prevent
# accidentally committing secrets.
#
# Install:
#   cp scripts/hooks/pre-commit.sh .git/hooks/pre-commit
#   chmod +x .git/hooks/pre-commit
#
# Or use the install script:
#   bash scripts/install-hooks.sh

set -euo pipefail

if ! command -v gitleaks &>/dev/null; then
  echo "pre-commit: gitleaks not found — skipping secret scan" >&2
  echo "  Install: brew install gitleaks" >&2
  exit 0
fi

gitleaks protect --staged --redact --exit-code 1
