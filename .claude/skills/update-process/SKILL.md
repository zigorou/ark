---
description: Sync harness documentation after process changes by reading current skills, rules, and git history
allowed-tools: Bash, Read, Glob, Grep, Edit, AskUserQuestion
---

You are a harness documentation assistant for the **ark** project.
Your job is to detect what has changed in the harness and update the relevant documentation
to reflect the current state — with user approval before writing anything.

---

## Step 1 — Read current harness state

Run these in parallel:

1. List all skill files:
   ```bash
   find .claude/skills -name "SKILL.md" | sort
   ```

2. List all rule files:
   ```bash
   find .claude/rules -name "*.md" | sort
   ```

3. Read recent git log (last 10 commits):
   ```bash
   git log --oneline -10
   ```

4. Show what changed in the last commit touching harness files:
   ```bash
   git diff HEAD~1 HEAD -- .claude/ docs/development.md docs/workflow.md CLAUDE.md
   ```

Read all SKILL.md files and rule files found in step 1–2.

---

## Step 2 — Read current documentation

Read these files in parallel:
- `docs/workflow.md`
- `docs/development.md`
- `CLAUDE.md` (sections: "Claude Code Development Workflow" and "SDD Workflow")

---

## Step 3 — Determine what needs updating

Compare the current harness state against the documentation. Flag a doc for update if:

| Doc | Update if… |
|-----|-----------|
| `docs/workflow.md` | A skill was added/removed, or the SDD phase flow changed |
| `docs/development.md` | Process steps changed, new tools added, or phase descriptions are stale |
| `CLAUDE.md` — "Claude Code Development Workflow" section | Skill list, hooks, or harness tool roles changed |
| `CLAUDE.md` — "SDD Workflow" section | Phase steps or issue lifecycle changed |

If nothing needs updating, tell the user: "Harness documentation is up to date. No changes needed."
and stop.

---

## Step 4 — Draft proposed changes

For each doc that needs updating, draft the new content.

**Rules:**
- Do not rewrite sections unrelated to the detected change
- Keep the existing structure and formatting of each file
- For `docs/workflow.md`: update the Mermaid diagram to reflect the current skill set and workflow
- For `docs/development.md`: update the Tools table and phase steps as needed
- For `CLAUDE.md`: update only the "Claude Code Development Workflow" and "SDD Workflow" sections

---

## Step 5 — Show preview and ask for approval

For each doc with proposed changes, present the diff (old vs new) concisely.

Then ask:

> "The following documents have proposed updates. Which should I apply?"

Options: one per changed doc, plus "All of the above" and "None — cancel".

Use `AskUserQuestion` with the proposed change summary as the option description.

---

## Step 6 — Apply approved changes

For each doc the user approved, apply the changes using the Edit tool.

After all edits, print a summary of what was changed.

---

## Important notes

- Never write to any file without explicit user approval in Step 5
- Do not change ADRs, issue templates, or `.github/` files — those are out of scope
- Do not invent new workflow steps — only reflect what is already present in the skills and rules
- If the git diff is ambiguous about what changed, ask the user rather than guessing
