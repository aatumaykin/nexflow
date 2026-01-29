# Memory System Guide for Nexflow

## Overview

Nexflow uses a **hybrid memory system** combining:
- Daily logs in Markdown
- Long-term curated memory
- Automatic session context
- (v1.0+) Heartbeats for proactive memory maintenance

## Daily Logs: `memory/YYYY-MM-DD.md`

### Purpose
Raw, unfiltered logs of what happens each day.

### When to Use
- Recording events as they happen
- Logging decisions made
- Capturing conversations or tasks

### Format
Free-form Markdown:
```markdown
# 2026-01-30

## Morning
User asked about deploying X project. Discussed deployment strategy with Vercel.

## Afternoon
Decided to use Railway instead. Key factors: cost, ease of setup, automatic SSL.

## Evening
Updated NOTES.md with new project structure.
```

### Automatic Loading
- **Today** and **yesterday** files are loaded into context
- Files are created automatically if they don't exist
- Location: `memory/` directory in workspace

## Long-term Memory: `memory/memory.md`

### Purpose
Curated, distilled memories worth keeping long-term.

### What Goes Here
- **Decisions** (e.g., "Use Railway for X project")
- **Preferences** (e.g., "Prefer Nova voice for TTS")
- **Lessons learned** (e.g., "Always test deployments before 5 PM")
- **Important context** (e.g., "User works at Company X on project Y")
- **Open loops** (e.g., "Need to decide between A and B")

### What Does NOT Go Here
- Secrets (unless explicitly requested)
- Trivial daily events
- Temporary information
- Duplicate information

### Security Rule
**MEMORY.md is ONLY loaded in main session** (direct chat).

**Does NOT load in:**
- Group chats
- Shared contexts
- Sessions with other people

This prevents personal information from leaking to strangers.

## Memory Synthesis (v1.0+)

When heartbeats are enabled, the agent periodically:
1. Reads recent `memory/YYYY-MM-DD.md` files
2. Identifies significant events, lessons, or insights
3. Updates `memory/memory.md` with distilled learnings
4. Removes outdated information from memory.md

## "Text > Brain" Rule

**Critical:** Agent's memory is limited between sessions. If you want to remember something:
- **WRITE IT TO A FILE** (`memory/YYYY-MM-DD.md` or `memory/memory.md`)
- Do NOT rely on "mental notes" — they don't survive session restarts

## Example Workflow

### During a session:
```bash
User: "Remember to use Railway for X project"
Agent: [Updates memory/2026-01-30.md with decision]
```

### At the end of the day:
```bash
User: "Summarize today"
Agent: [Reads memory/2026-01-30.md]
Agent: [Updates memory/memory.md with key decisions]
```

### In a future session:
```bash
Agent: [Loads memory/memory.md]
Agent: "Based on past decisions, you prefer Railway for X project deployments. Should I use it now?"
```

## Best Practices

1. **Write immediately** — Don't wait until "later"
2. **Be specific** — Capture context, not just conclusions
3. **Update regularly** — Review and prune outdated info
4. **Respect security** — Don't load personal memory in shared contexts
5. **Curate wisely** — MEMORY.md should be distilled wisdom, not raw logs

## Configuration

Memory location:
```yaml
agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
```

Disable memory loading:
```yaml
agents:
  defaults:
    skip_bootstrap: true
```
