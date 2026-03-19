#!/usr/bin/env bash
# scripts/install-hooks.sh
#
# Installs git hooks for local development.
# Run once after cloning the repository.
#
# Usage:
#   bash scripts/install-hooks.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

install_hook() {
  local name="$1"
  local src="$REPO_ROOT/scripts/hooks/${name}.sh"
  local dst="$HOOKS_DIR/$name"

  if [[ ! -f "$src" ]]; then
    echo "  skip: $src not found"
    return
  fi

  cp "$src" "$dst"
  chmod +x "$dst"
  echo "  installed: .git/hooks/$name"
}

echo "Installing git hooks..."
install_hook pre-commit
echo "Done."

echo ""
echo "Recommended local tools:"
echo "  brew install gitleaks   # secret scanning (used by pre-commit hook)"
echo "  brew install golangci-lint  # linting"
