---
name: agents-overview
description: overview of agents
tools: read, grep, diff
---

# .claude/agents/ - Specialized Subagents

## Overview

This directory contains specialized subagents that can be spawned by the main Claude Code agent to handle specific tasks. Each agent has focused expertise and a defined set of tools.

## Available Agents

| Agent | Purpose | Tools |
|-------|---------|-------|
| `implementation-agent` | Full feature implementation | read, write, bash, grep, edit, glob |
| `code-review` | Code quality analysis | read, grep, diff |
| `test-engineer` | Test strategy and TDD | read, write, bash, grep |
| `security-auditor` | Security vulnerability analysis | read, grep, bash |
| `documentation-generator` | Documentation generation | read, write, bash, grep, glob |

## Agent Invocation

Agents are spawned using the Task tool:

```
Use Task tool with subagent_type='{agent-name}'
Provide context about:
- What needs to be done
- Relevant files/code
- Expected outputs
```

## Agent Structure

Each agent file follows this structure:

```markdown
---
name: agent-name
description: What the agent does
tools: tool1, tool2, tool3
skills: optional-skills-to-load
---

# Agent Title

## Purpose
What this agent specializes in.

## Workflow
Steps the agent follows.

## Inputs
What context to provide.

## Outputs
What the agent produces.
```

## When to Spawn Agents

| Scenario | Agent |
|----------|-------|
| Implementing a feature | `implementation-agent` |
| Reviewing code changes | `code-review` |
| Writing or improving tests | `test-engineer` |
| Security audit or vulnerability check | `security-auditor` |
| Generating documentation | `documentation-generator` |

## Adding New Agents

1. Create `{agent-name}.md` in this directory
2. Include frontmatter with name, description, tools
3. Document the agent's workflow
4. Update this CLAUDE.md
