---
paths:
  - ".claude/skills/**/SKILL.md"
---

# Skill Authoring Rules

These rules apply to all Claude Code skill files under `.claude/commands/`.

---

## 1. Frontmatter — required fields

Every skill file **must** have a YAML frontmatter block with at least these two fields:

```yaml
---
description: <one-sentence description>
allowed-tools: <comma-separated tool list>
---
```

A skill missing either field is invalid and will not load correctly.

---

## 2. `description` style

- Must be a single English sentence
- Must start with an **imperative verb** (e.g. `Recommend`, `Create`, `Fetch`, `Run`, `List`)
- Must not start with "This skill…", "A skill that…", or any article/pronoun
- Should describe **what the skill does**, not how it works internally
- Maximum ~120 characters

**Good:**
```
description: Create a typed GitHub issue (spec / impl / process) from the project's SDD templates
description: Recommend the next task to work on based on GitHub issue state and SDD workflow
```

**Bad:**
```
description: This skill creates issues          # starts with article
description: issue creation helper              # not a sentence, not imperative
description: Creates, manages, and tracks...   # too vague
```

---

## 3. Creating or modifying skills — use `skill-creator`

When asked to **create a new skill** or **substantially modify an existing one**:

- Use the `skill-creator` skill instead of editing the file directly
- `skill-creator` runs evals, checks trigger accuracy, and validates description quality
- Only use direct `Edit`/`Write` for trivial fixes (typos, formatting) that do not change skill behavior

**Trigger phrases that require `skill-creator`:**
- "create a new skill for…"
- "add a skill that…"
- "update the /foo skill to also…"
- "rewrite the /bar skill"

---

## 4. `allowed-tools` — principle of least privilege

Only list tools the skill actually uses. Do not add tools speculatively.

| If the skill… | Include | Exclude |
|---------------|---------|---------|
| Only reads files / issues | `Bash`, `Read`, `Grep` | `Edit`, `Write` |
| Creates/modifies files | `Edit`, `Write` | only if needed |
| Asks clarifying questions | `AskUserQuestion` | — |
| Enters plan mode | `EnterPlanMode`, `ExitPlanMode` | — |

**Never add** `Edit` or `Write` to a skill that is declared read-only in its prompt.
If a skill's prompt says "do not make any changes", its `allowed-tools` must not include `Edit` or `Write`.

---

## Checklist for new skill files

Before saving a new skill, verify:

- [ ] Frontmatter has `description` and `allowed-tools`
- [ ] `description` starts with an imperative verb and is ≤120 chars
- [ ] `allowed-tools` contains only tools the skill actually invokes
- [ ] Skill was authored or reviewed via `skill-creator` (unless trivial fix)
