---
description: Recommend the next task to work on based on GitHub issue state and SDD workflow
allowed-tools: Bash, Read, Grep
---

You are a project assistant for the **ark** project following the Spec Driven Development (SDD) workflow.
Your job is to determine and recommend the next task to work on.

## Step 1 — Fetch all issues

Run this command and parse the JSON output:

```bash
gh issue list --state all --json number,title,labels,state,body --limit 100
```

## Step 2 — Classify issues

From the JSON, classify each issue into:
- **closed spec issues**: `state == "CLOSED"` and has label `spec`
- **open process issues**: `state == "OPEN"` and has label `process`
- **open implementation issues**: `state == "OPEN"` and has label `implementation`
- **open spec issues**: `state == "OPEN"` and has label `spec`

## Step 3 — Detect triage gaps

For each closed spec issue (number `N`):
1. Check all implementation issues' `body` fields for a reference matching `#N` (as a word boundary, e.g. `#9` not part of `#90`)
2. If **no** open or closed implementation issue references `#N`, that spec has a triage gap

## Step 4 — Determine recommendation

Apply this priority order:

### Priority 1: Triage gap detected

If one or more closed spec issues have no corresponding implementation issue:
- List each gap: spec issue number, title
- Recommend creating an implementation issue for the **lowest-numbered** ungapped spec
- Tell the user: "Spec #N (`<title>`) is closed but has no implementation issue. Recommend creating one."

### Priority 2: Open process issues exist

If there are open process issues (and no triage gap):
- List them all
- Recommend working on the **lowest-numbered** open process issue
- Explain: "Harness health before features — process issues should be resolved first."

### Priority 3: Open implementation issues exist

If there are open implementation issues (and no triage gap, no open process issues):
- List them all
- Recommend the **lowest-numbered** open implementation issue as the default
- If there are **multiple** candidates and the choice is non-obvious, use `AskUserQuestion` to let the user choose, presenting the options with brief descriptions

### Priority 4: No open issues

If there are no open issues and no triage gaps:
- Read `docs/concept.md` using the Read tool
- Identify the most prominent **unresolved design question** or **unimplemented feature** mentioned in concept.md
- Propose opening a new spec issue for it
- Tell the user: "No open issues found. Based on concept.md, the next logical spec might be: `<topic>`. Shall I open a spec issue?"

## Output format

Present your recommendation as:

```
## Next Task

**Recommendation**: <one-sentence summary>

**Reason**: <why this is the highest priority item>

**Action**: <concrete next step — e.g., "Run `/next-impl #N` to start implementation" or "Open spec issue for X">
```

If multiple issues exist at the same priority level, also show a numbered list so the user has context.

## Important notes

- Do not open issues or make any changes — only recommend
- Do not enter Plan mode — this skill is read-only
- Be concise; this is a quick triage tool, not a deep analysis
