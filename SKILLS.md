# .claude/skills/ - Reusable Capabilities

## Overview

Skills are reusable capabilities that Claude Code can use autonomously. Unlike agents (which are explicitly spawned), skills are automatically available and can be invoked by the model when needed.

## Available Skills

| Skill | Purpose | Tools |
|-------|---------|-------|
| `code-reviewer` | Comprehensive code review with security, performance, and quality analysis | read, grep, diff, lint_runner |
| `documentation-generator` | Generate API documentation, OpenAPI specs, and code documentation | read, write, bash, grep, glob |

## Skill vs Agent

| Aspect | Skill | Agent |
|--------|-------|-------|
| Invocation | Automatic (model decides) | Explicit (Task tool) |
| Scope | Focused capability | Complete workflow |
| State | Stateless | Can maintain context |
| Use case | Single task | Multi-step process |

## How Skills Work

1. Skills are defined as markdown files in this directory
2. The main agent can load skills as needed
3. Skills provide focused expertise for specific tasks
4. Agents can declare skills to auto-load in their frontmatter

## Skill Structure

```markdown
---
description: What this skill does
tools: tool1, tool2
---

# Skill Name

## Purpose
What this skill is for.

## Capabilities
What it can do.

## Usage
How to use it.
```

## Adding New Skills

1. Create `{skill-name}.md` in this directory
2. Include frontmatter with description and tools
3. Document the skill's capabilities
4. Update this CLAUDE.md
