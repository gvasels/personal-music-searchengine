# .claude/commands/ - Slash Commands

## Overview

Slash commands are user-invoked operations that start with `/`. When executed, the command file is expanded into a full prompt.

## Available Commands

| Command | Description |
|---------|-------------|
| `/sdlc` | Full SDLC workflow orchestrator |
| `/update-claudemd` | Update CLAUDE.md from git changes |
| `/code-review` | Run comprehensive code review |
| `/test-file` | Generate tests for a specific file |

## Command Structure

```markdown
---
description: What the command does
allowed-tools: Optional tool restrictions
---

# Command Name

Instructions for the command.
Can reference files with @filename
Can run shell commands with !`command`
```

## Usage

Invoke commands with a slash:

```
/sdlc
/code-review src/services/
/test-file src/services/user.ts
```

## Adding New Commands

1. Create `{command-name}.md` in this directory
2. Include frontmatter with description
3. Write the command prompt/instructions
4. Update this CLAUDE.md
