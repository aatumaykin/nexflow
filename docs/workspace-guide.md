# Workspace Guide for Nexflow

## Overview

Nexflow uses a **workspace-based** system for managing agent personality, user context, and memory. The workspace is automatically loaded and injected into the agent's context at the start of each session.

## Workspace Location

By default: `~/nexflow`

Can be customized through an environment variable:
```bash
export NEXFLOW_WORKSPACE="/custom/path/to/workspace"
```

Or in configuration:
```yaml
agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
```

## Bootstrap Files

These files define who the agent is, who it helps, and how it behaves.

### SOUL.md - Agent Identity
Defines the agent's personality, boundaries, and core truths.

**Contents:**
- Core Truths (how to behave)
- Identity (name, emoji)
- Boundaries (what to do/not do)
- Vibe (communication style)

**When updated:** Tell the user — it's the agent's "soul"

### USER.md - User Profile
Defines who the agent is helping.

**Contents:**
- User's name
- What to call them
- Pronouns
- Timezone
- Context (projects, preferences)

**Purpose:** The more the agent knows, the better it can help.

### AGENTS.md - Agent Instructions
Defines how the agent works and handles memory.

**Contents:**
- Session startup routine (what files to read)
- Memory management (daily vs long-term)
- Safety rules
- External vs internal actions

**Key rule:** "Text > Brain" — write things down, don't rely on mental notes.

### NOTES.md - Local Context
Contains environment-specific information not in code.

**Examples:**
- Camera names and locations
- SSH hosts and aliases
- TTS voice preferences
- Device nicknames

**Purpose:** Keeps local context separate from skill documentation.

## Memory System

### Daily Logs: `memory/YYYY-MM-DD.md`
- Raw logs of events and conversations
- Created automatically
- Used for recent context (today + yesterday)
- Location: `memory/` directory in workspace

### Long-term Memory: `memory/memory.md`
- Curated memories and decisions
- Only loaded in **main session** (direct chat)
- NOT loaded in groups or shared contexts
- Updated periodically from daily logs

## Session Types

### Main Session
- Direct chat with the human
- **Loads:** SOUL.md, USER.md, NOTES.md, MEMORY.md, daily logs
- Can read/write to all memory files

### Other Sessions
- Group chats, shared contexts
- **Loads:** SOUL.md, USER.md, NOTES.md, daily logs
- **Does NOT load MEMORY.md** (security)

## Configuration

### Disable Bootstrap
```yaml
agents:
  defaults:
    skip_bootstrap: true
```
Useful for pre-seeded workspaces.

### File Size Limits
```yaml
agents:
  defaults:
    bootstrap_max_chars: 20000
```
Files larger than this are truncated.

### Custom Workspace Path
```bash
# Default
export NEXFLOW_WORKSPACE=""

# Custom
export NEXFLOW_WORKSPACE="/path/to/workspace"

# Or in config
agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
```

## Initial Setup

Run the setup wizard:
```bash
nexflow setup
```

This creates a workspace with templates and guides you through initial configuration.

**Setup questions:**
1. Workspace path (or use default `~/nexflow`)
2. User's name
3. What to call the user
4. Timezone
5. Agent's name
6. Agent's emoji

The wizard will:
- Create workspace directory structure
- Generate bootstrap files from templates
- Create `memory/` directory
- Initialize empty `memory/memory.md`
- Create `NOTES.md` ready for your local notes
