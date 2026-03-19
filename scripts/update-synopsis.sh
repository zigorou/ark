#!/usr/bin/env bash
# scripts/update-synopsis.sh
#
# Regenerates the SYNOPSIS section in README.md by parsing the Use: field
# from each cmd/*.go file (excluding root.go). Does not require building the
# binary, so it works even when the code does not compile yet.
#
# The README section is delimited by:
#   <!-- SYNOPSIS:START -->  ...  <!-- SYNOPSIS:END -->

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export REPO_ROOT

python3 - <<'PYEOF'
import os
import re
import glob

repo_root = os.environ['REPO_ROOT']
cmd_dir   = os.path.join(repo_root, 'cmd')
readme    = os.path.join(repo_root, 'README.md')

use_re = re.compile(r'\bUse:\s+"([^"]+)"')

# Collect Use: values from all cmd/*.go except root.go, sorted by filename.
lines = ['ark [--vault <path>] <command> [arguments]', '']
for go_file in sorted(glob.glob(os.path.join(cmd_dir, '*.go'))):
    if go_file.endswith('root.go'):
        continue
    with open(go_file) as f:
        content = f.read()
    m = use_re.search(content)
    if m:
        lines.append(f'ark {m.group(1)}')

new_synopsis = '\n'.join(lines)

with open(readme) as f:
    content = f.read()

pattern     = r'(<!-- SYNOPSIS:START -->).*?(<!-- SYNOPSIS:END -->)'
replacement = f'\\1\n```\n{new_synopsis}\n```\n\\2'
new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

if new_content != content:
    with open(readme, 'w') as f:
        f.write(new_content)
    print('update-synopsis: README.md SYNOPSIS updated')
else:
    print('update-synopsis: README.md SYNOPSIS unchanged')
PYEOF
