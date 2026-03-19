#!/usr/bin/env bash
# scripts/hooks/post-tool-use-synopsis.sh
#
# Claude Code PostToolUse hook.
# Triggers after Write or Edit tool calls and regenerates README.md SYNOPSIS
# when a file under cmd/ is modified.
#
# Hook input (JSON via stdin):
#   { "tool_name": "Write"|"Edit", "tool_input": { "file_path": "...", ... }, ... }

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Parse the modified file path from hook JSON
file_path=$(python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    print(d.get('tool_input', {}).get('file_path', ''))
except Exception:
    print('')
" 2>/dev/null || echo "")

# Only act on cmd/*.go changes
if [[ "$file_path" != */cmd/*.go ]]; then
  exit 0
fi

echo "[hook] cmd file changed: $file_path — updating SYNOPSIS..." >&2

# Run in background so the hook doesn't block Claude Code
bash "$REPO_ROOT/scripts/update-synopsis.sh" &
