---
description: Create a typed GitHub issue (spec / impl / process) from the project's SDD templates
allowed-tools: Bash, AskUserQuestion, Read
---

You are a project assistant for the **ark** project. Your job is to create a new GitHub issue
using the SDD workflow templates.

## Step 1 — Parse args

The user invokes this skill as:

```
/open-issue <type> ["<title>"] [--refs "<#N, ADR-XXXX>"]
```

- `<type>` must be one of: `spec`, `impl`, `process`
- `"<title>"` is optional — the text after the type and before any `--` flags
- `--refs` is only relevant for `impl` issues

If `<type>` is missing or not one of the three values, ask the user:

> "Which issue type do you want to create?"
> Options: `spec`, `impl`, `process`

## Step 2 — Gather required information

### For all types

If no title was provided in the args, ask:
> "What is the issue title? (will be prefixed with `<type>: ` automatically)"

Keep the question short. Do not ask for the full body yet.

### For `impl` only

Ask (or parse from `--refs`):
> "Which spec issue(s) and ADR(s) does this implementation reference? (e.g. `#3, ADR-0004` — or leave blank if none yet)"

### For `spec` only

Ask:
> "What is the core design question this spec addresses? (1–2 sentences)"

### For `process` only

Ask:
> "What problem or friction does this process issue address? (1–2 sentences)"

Collect all answers before proceeding.

## Step 3 — Read the template

Read the corresponding template file:

| type    | template path |
|---------|--------------|
| spec    | `.github/ISSUE_TEMPLATE/spec.md` |
| impl    | `.github/ISSUE_TEMPLATE/implementation.md` |
| process | `.github/ISSUE_TEMPLATE/process.md` |

## Step 4 — Build the issue body

Construct the body by filling in the template's placeholder sections with the information gathered.

**Rules:**
- Prefix the title with the type: `spec: <title>`, `implementation: <title>`, `process: <title>`
- For `impl`, insert the refs line: `Refs: <refs>`
- Keep placeholder sections that require more detail as-is (e.g. `## Acceptance Criteria` with empty checkboxes) — the user will fill them in after creation
- Strip the YAML frontmatter from the template (lines between `---` markers)

## Step 5 — Create the issue

Run:

```bash
gh issue create \
  --title "<prefixed title>" \
  --body "<constructed body>" \
  --label "<type>"
```

Use a heredoc or `printf` to handle multi-line body safely.

After creation, output the issue URL so the user can navigate directly to it.

## Step 6 — Remind the user

After creating the issue, print this reminder:

> "Issue created. Before starting Plan mode, make sure the `## Acceptance Criteria` section
> has at least one checkbox filled in — issues without AC are not ready for implementation."

## Important notes

- Do not start Plan mode or begin implementation — this skill only creates the issue
- Do not invent Acceptance Criteria — leave them blank for the user to define
- For `impl` issues, do not invent refs — leave them blank if the user didn't provide them
