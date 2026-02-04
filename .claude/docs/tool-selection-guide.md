# Tool Selection Guide

This guide explains when and where to use the various Claude Code capabilities: MCP servers, plugins, agents, skills, and commands.

## Tool Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER INVOCATION                                 │
│  /sdlc, /code-review, /test-file                                           │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            SLASH COMMANDS                                    │
│  Trigger workflows, run predefined operations                               │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                ▼                ▼
┌───────────────────────┐  ┌─────────────────┐  ┌─────────────────────────────┐
│       PLUGINS         │  │     AGENTS      │  │          SKILLS             │
│  SDLC phase workflows │  │  Specialized    │  │  Reusable capabilities      │
│  spec, test, code...  │  │  subprocesses   │  │  (auto-invoked)             │
└───────────────────────┘  └─────────────────┘  └─────────────────────────────┘
                    │                │                │
                    └────────────────┼────────────────┘
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            MCP SERVERS                                       │
│  External tools: AWS, GitHub, spec-workflow, OpenTofu, etc.                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Quick Reference: What to Use When

| Need | Use | Why |
|------|-----|-----|
| Start a new feature | `/sdlc` command | Orchestrates full workflow |
| Write specs/requirements | `spec-workflow` MCP | Creates requirements.md, design.md, tasks.md |
| Write tests first (TDD) | `test-engineer` agent | Specializes in test strategy |
| Implement code | `implementation-agent` | Full write access, focused context |
| Review code quality | `code-reviewer` skill | Automatic, focused analysis |
| Security audit | `security-auditor` agent | OWASP-aligned checks |
| Generate docs | `documentation-generator` skill/agent | API docs, CLAUDE.md |
| Look up AWS docs | `aws-documentation-mcp-server` | Official AWS documentation |
| Look up library docs | `context7` MCP | Up-to-date library docs |
| **Initialize OpenTofu deployment** | `infrastructure-deployer` skill | First-time state backend setup |
| Run AWS CLI commands | `aws-api-mcp-server` | Execute and validate AWS commands |
| Design DynamoDB tables | `dynamodb-mcp-server` | Data modeling expertise |
| Build containers | `finch-mcp-server` | Finch/Docker container builds |
| Search OpenTofu registry | `opentofu` MCP | Provider and module docs |
| GitHub operations | `github` MCP | Issues, PRs, code search |
| Browser automation | `playwright` MCP | E2E testing, screenshots |

---

## MCP Servers

MCP (Model Context Protocol) servers provide external tool access.

### AWS Servers

| Server | Purpose | When to Use |
|--------|---------|-------------|
| `aws-documentation-mcp-server` | Search and read AWS docs | Looking up AWS service configuration, best practices |
| `aws-knowledge-mcp-server` | AWS knowledge base queries | Regional availability, service limits, recommendations |
| `aws-api-mcp-server` | Execute AWS CLI commands | Deploy, query, or manage AWS resources |
| `dynamodb-mcp-server` | DynamoDB data modeling | Design single-table schemas, access patterns, validation |
| `syntheticdata-mcp-server` | Generate test data | Create realistic test datasets for development |
| `finch-mcp-server` | Container builds | Build and push Docker images with Finch |

**Example Usage:**
```
Task: "Design a DynamoDB table for user sessions"
→ Use dynamodb-mcp-server for schema design
→ Use aws-api-mcp-server to create the table
→ Use syntheticdata-mcp-server to generate test data
```

### Infrastructure Servers

| Server | Purpose | When to Use |
|--------|---------|-------------|
| `opentofu` | OpenTofu/Terraform registry | Look up providers, resources, modules |
| `spec-workflow` | Specification management | Create specs, track approval, manage tasks |

**Example Usage:**
```
Task: "Add an S3 bucket with versioning"
→ Use opentofu MCP to look up aws_s3_bucket resource
→ Use aws-documentation-mcp-server for S3 best practices
```

### Development Servers

| Server | Purpose | When to Use |
|--------|---------|-------------|
| `github` | GitHub operations | Create issues, PRs, search code, manage repos |
| `context7` | Library documentation | Look up npm/Go/Python package docs |
| `docs-mcp-server` | Custom docs indexing | Index and search project-specific docs |
| `playwright` | Browser automation | E2E tests, UI verification, screenshots |
| `daisyui-blueprint` | UI components | DaisyUI/Tailwind component snippets |

**Example Usage:**
```
Task: "Implement React component with DaisyUI"
→ Use daisyui-blueprint for component snippets
→ Use context7 for React hooks documentation
→ Use playwright to test the component
```

---

## Agents

Agents are specialized subprocesses spawned via the Task tool.

### When to Spawn Agents

| Agent | Trigger | Context to Provide |
|-------|---------|-------------------|
| `implementation-agent` | "Implement this feature" | Spec files, related code, acceptance criteria |
| `test-engineer` | "Write tests for X" | Design docs, interfaces, test requirements |
| `code-review` | "Review these changes" | File paths, PR context, focus areas |
| `security-auditor` | "Security audit X" | Code paths, threat model, compliance requirements |
| `documentation-generator` | "Document this module" | Source files, API contracts, existing docs |

### Agent Spawning Syntax

```
Use Task tool with subagent_type='{agent-name}'
```

**Example:**
```
Task: "Write unit tests for the account service"

Use Task tool with subagent_type='test-engineer'
Prompt: "Write comprehensive unit tests for services/account-vending/internal/handlers/handlers.go
covering all CRUD operations. Follow TDD principles - tests should initially fail."
```

### Agent vs Direct Implementation

| Scenario | Use Agent | Do Directly |
|----------|-----------|-------------|
| Complex multi-file implementation | ✅ `implementation-agent` | |
| Simple single-file edit | | ✅ Direct edit |
| Comprehensive test suite | ✅ `test-engineer` | |
| Add one test case | | ✅ Direct write |
| Full code review | ✅ `code-review` | |
| Quick sanity check | | ✅ Read + comment |
| Security audit | ✅ `security-auditor` | |
| Check one vulnerability | | ✅ Direct grep/read |

---

## Skills

Skills are reusable capabilities automatically available to Claude Code.

### Available Skills

| Skill | Purpose | Auto-Invoked When |
|-------|---------|-------------------|
| `code-reviewer` | Code quality analysis | Reviewing PRs, checking code patterns |
| `documentation-generator` | Generate documentation | Creating API docs, updating CLAUDE.md |
| `infrastructure-deployer` | Initialize OpenTofu state backend | First-time deployments, state migration |

### Skill vs Agent

| Aspect | Skill | Agent |
|--------|-------|-------|
| Invocation | Automatic | Explicit (Task tool) |
| Scope | Single focused task | Multi-step workflow |
| Context | Limited | Full subprocess context |
| Best for | Quick analysis | Complex implementation |

**Example:**
```
# Skill (automatic) - Quick code review insight
"What are the issues with this function?"
→ code-reviewer skill activated automatically

# Agent (explicit) - Full code review with report
"Do a comprehensive code review of the account service"
→ Spawn code-review agent with Task tool
```

---

## Plugins (SDLC Phases)

Plugins represent phases in the SDLC workflow, invoked via `/sdlc` or individually.

### Plugin-Phase Mapping

| Phase | Plugin | Primary Tool | Output |
|-------|--------|--------------|--------|
| 1. Spec | `spec-writer` | spec-workflow MCP | requirements.md, design.md, tasks.md |
| 2. Test | `test-writer` | test-engineer agent | *_test.go, *.test.ts |
| 3. Code | `code-implementer` | implementation-agent | Source files |
| 4. Build | `builder` | bash (go build, npm build) | Compiled artifacts |
| 5. Security | `security-checker` | security-auditor agent | Security report |
| 6. Docs | `docs-generator` | documentation-generator | API docs, CLAUDE.md |

### When to Use Each Plugin

```
/sdlc                     → Full workflow (all 6 phases)
"Run spec phase"          → Phase 1 only
"Write tests for X"       → Phase 2 only
"Implement feature Y"     → Phase 3 only
"Build and verify"        → Phase 4 only
"Security scan"           → Phase 5 only
"Generate docs"           → Phase 6 only
```

---

## Slash Commands

Commands are user-invoked operations starting with `/`.

### Available Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `/sdlc` | Full SDLC workflow | Starting new feature implementation |
| `/update-claudemd` | Update CLAUDE.md | After significant code changes |
| `/code-review` | Run code review | Before PR submission |
| `/test-file <path>` | Generate tests | Adding tests for existing code |

### Command Usage Patterns

```bash
# Start a new feature
/sdlc

# Generate tests for a specific file
/test-file services/account-vending/internal/handlers/handlers.go

# Run code review on a directory
/code-review services/account-vending/

# Update documentation after changes
/update-claudemd
```

---

## Decision Trees

### "I need to implement something"

```
Is it a new feature with requirements?
├── YES → /sdlc (full workflow)
└── NO → Is it complex (multi-file, multi-step)?
    ├── YES → Spawn implementation-agent
    └── NO → Direct edit
```

### "I need documentation"

```
What kind of documentation?
├── AWS service docs → aws-documentation-mcp-server
├── Library/package docs → context7 MCP
├── OpenTofu/Terraform → opentofu MCP
├── Project API docs → documentation-generator agent
└── CLAUDE.md update → /update-claudemd command
```

### "I need to work with AWS"

```
What operation?
├── Look up how to do X → aws-documentation-mcp-server
├── Check regional availability → aws-knowledge-mcp-server
├── Execute CLI command → aws-api-mcp-server
├── Design DynamoDB schema → dynamodb-mcp-server
├── Generate test data → syntheticdata-mcp-server
└── Build/push container → finch-mcp-server
```

### "I need to deploy infrastructure"

```
First-time deployment to this directory?
├── YES → infrastructure-deployer skill
│   └── Configures S3 backend, role assumption
│   └── Outputs state location summary
└── NO → Routine update?
    ├── YES → Direct tofu plan/apply
    └── NO → State migration needed?
        ├── YES → infrastructure-deployer skill with --migrate
        └── NO → Direct tofu commands
```

### "I need to write tests"

```
Is this part of TDD for a new feature?
├── YES → test-engineer agent (via /sdlc Phase 2)
└── NO → Quick test addition?
    ├── YES → Direct write
    └── NO → test-engineer agent
```

### "I need to review code"

```
Comprehensive review needed?
├── YES → code-review agent (or /code-review command)
└── NO → code-reviewer skill (automatic)
```

---

## Best Practices

### 1. Start with the Right Abstraction

```
High-level task → /sdlc command → Orchestrates everything
Medium task → Specific agent → Focused expertise
Low-level task → Direct tools → Quick and efficient
```

### 2. Let MCP Servers Handle External Operations

Don't manually construct AWS CLI commands - use the appropriate MCP server.

```
❌ "Run this command: aws dynamodb create-table..."
✅ "Create a DynamoDB table with this schema" → dynamodb-mcp-server handles it
```

### 3. Use Agents for Complex Multi-Step Work

```
❌ Multiple back-and-forth edits for one feature
✅ Spawn implementation-agent with full context
```

### 4. Prefer Skills for Quick Analysis

```
❌ Spawn code-review agent for "is this function okay?"
✅ Let code-reviewer skill activate automatically
```

### 5. Chain Tools Appropriately

```
Spec → Test → Code → Build → Security → Docs
  ↓      ↓      ↓       ↓        ↓        ↓
MCP   Agent  Agent   Bash    Agent    Agent
```

---

## Troubleshooting

### MCP Server Not Responding
1. Check if the server is configured in `.mcp.json`
2. Verify environment variables are set
3. Try restarting Claude Code session

### Agent Not Producing Expected Output
1. Provide more context in the Task prompt
2. Specify exact files and requirements
3. Include acceptance criteria

### Skill Not Activating
1. Skills activate based on task context
2. Be explicit: "review this code for quality issues"
3. If needed, use the agent version instead

### Command Not Found
1. Check `.claude/commands/` for available commands
2. Verify command file has correct frontmatter
3. Check COMMANDS.md for usage examples
