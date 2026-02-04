# Claude Code Architecture - How Components Work Together

This document describes how Claude Code's agents, skills, plugins, commands, and MCP servers work together to implement a spec-driven, TDD development workflow.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           USER REQUEST                                          │
│                      "Implement feature X"                                      │
└─────────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        CLAUDE CODE AGENT                                        │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                     ENTRY POINTS                                         │   │
│  │  /sdlc          → Orchestrates full SDLC workflow                        │   │
│  │  /code-review   → Code quality analysis                                  │   │
│  │  /test-file     → Generate tests for specific file                       │   │
│  │  /update-claudemd → Update docs from git changes                         │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                    ┌───────────────┼───────────────┐                           │
│                    ▼               ▼               ▼                           │
│  ┌──────────────────────┐ ┌──────────────┐ ┌─────────────────────┐            │
│  │     PLUGINS          │ │   AGENTS     │ │   MCP SERVERS       │            │
│  │ (Phase Orchestration)│ │ (Task Exec)  │ │ (External Tools)    │            │
│  │                      │ │              │ │                     │            │
│  │ • spec-writer        │ │ • test-      │ │ • spec-workflow     │            │
│  │ • test-writer        │ │   engineer   │ │ • github            │            │
│  │ • code-implementer   │ │ • impl-      │ │ • context7          │            │
│  │ • builder            │ │   agent      │ │ • dynamodb-mcp      │            │
│  │ • security-checker   │ │ • code-      │ │ • aws-docs-mcp      │            │
│  │ • docs-generator     │ │   review     │ │ • opentofu          │            │
│  │ • deploy-verifier    │ │ • security-  │ │ • playwright        │            │
│  │                      │ │   auditor    │ │                     │            │
│  └──────────────────────┘ │ • docs-gen   │ └─────────────────────┘            │
│                           │ • doc-       │                                     │
│                    ┌──────│   checker    │──────┐                              │
│                    │      └──────────────┘      │                              │
│                    ▼                            ▼                              │
│  ┌──────────────────────────────┐  ┌──────────────────────────────────┐       │
│  │         SKILLS               │  │     DOCUMENTATION                │       │
│  │   (Reusable Capabilities)    │  │  (Process Guidance)              │       │
│  │                              │  │                                  │       │
│  │ • code-reviewer/             │  │ • .claude/docs/*.md              │       │
│  │ • documentation-generator/   │  │ • tdd-workflow.md                │       │
│  │ • infrastructure-deployer    │  │ • wave-assignment.md             │       │
│  └──────────────────────────────┘  │ • epic-completion-checklist.md   │       │
│                                    └──────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         ARTIFACTS                                               │
│                                                                                 │
│  Specs:           .spec-workflow/specs/{epic}/requirements.md, design.md,       │
│                   tasks.md                                                      │
│  Tests:           backend/**/*_test.go, frontend/**/*.test.ts                   │
│  Code:            backend/**, frontend/**                                       │
│  Infra:           infrastructure/**/*.tf                                        │
│  Docs:            CLAUDE.md (root + subdirs), CHANGELOG.md                      │
│  Tracking:        implementation-plan/epics-user-stories.md                     │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### 1. Commands (Entry Points)

**Location**: `.claude/commands/`

Commands are user-invocable workflows that orchestrate complex multi-step processes.

| Command | Purpose | Invokes |
|---------|---------|---------|
| `/sdlc` | Full SDLC workflow (Spec→Test→Code→Verify→Docs) | Plugins + Agents |
| `/code-review` | Comprehensive code quality analysis | code-reviewer skill |
| `/test-file <path>` | Generate tests for specific file | test-engineer agent |
| `/update-claudemd` | Update CLAUDE.md from git changes | documentation-generator |
| `/analyze-codebase` | Generate codebase analysis | Exploration tools |

**Usage**: Type the command in Claude Code (e.g., `/sdlc`) to start the workflow.

### 2. Plugins (Phase Orchestration)

**Location**: `.claude/plugins/`

Plugins define how each SDLC phase should be executed. They specify which agents to spawn, what artifacts to create, and validation criteria.

| Plugin | Phase | Spawns | Artifacts |
|--------|-------|--------|-----------|
| `spec-writer` | 1. Specification | - (uses MCP) | requirements.md, design.md, tasks.md |
| `test-writer` | 2. Test | test-engineer | *_test.go, *.test.ts |
| `code-implementer` | 3. Code | implementation-agent | Source code |
| `builder` | 4. Build | - (bash) | Build output, coverage reports |
| `security-checker` | - | security-auditor | Security report |
| `docs-generator` | 5. Docs | documentation-generator | CLAUDE.md, CHANGELOG.md |
| `deploy-verifier` | - | - (CI/CD) | Deployment verification |

### 3. Agents (Task Execution)

**Location**: `.claude/agents/`

Agents are specialized subagents spawned via the `Task` tool to execute specific types of work.

| Agent | Purpose | Tools Available |
|-------|---------|-----------------|
| `test-engineer` | Write failing tests (TDD Red) | read, write, bash, grep, glob |
| `implementation-agent` | Write code to pass tests (TDD Green) | read, write, bash, grep, edit, glob |
| `code-review` | Code quality and security analysis | read, grep, diff, lint |
| `security-auditor` | OWASP vulnerability scanning | read, grep, bash, glob |
| `documentation-generator` | API docs, CLAUDE.md generation | read, write, bash, grep, glob |
| `doc-consistency-checker` | Validate documentation consistency | read, grep, glob |

**Spawning an Agent** (via Task tool):
```
Task tool:
- description: "Write tests for upload handler"
- subagent_type: "test-engineer"
- prompt: "Write tests for backend/internal/handlers/upload.go..."
```

### 4. Skills (Reusable Capabilities)

**Location**: `.claude/skills/`

Skills are reusable capabilities that can be invoked by commands or directly.

| Skill | Purpose | Contains |
|-------|---------|----------|
| `code-reviewer/` | Code review automation | Templates, analysis scripts |
| `documentation-generator/` | Doc generation | Generation scripts |
| `infrastructure-deployer` | OpenTofu deployment | Backend config guide |

### 5. MCP Servers (External Tools)

**Configuration**: `.mcp.json`

MCP servers provide specialized capabilities through external tools.

| Server | Purpose | Key Tools |
|--------|---------|-----------|
| `spec-workflow` | Specification management | spec-status, approvals, log-implementation |
| `github` | GitHub operations | Issues, PRs, code search |
| `context7` | Library documentation | Up-to-date API docs |
| `dynamodb-mcp-server` | DynamoDB design | Data modeling guidance |
| `aws-documentation-mcp-server` | AWS docs | Service documentation |
| `opentofu` | OpenTofu/Terraform docs | Provider, module docs |
| `playwright` | Browser automation | E2E testing support |

### 6. Documentation (Process Guidance)

**Location**: `.claude/docs/`

Documentation files provide process guidance and best practices.

| Document | Purpose |
|----------|---------|
| `tdd-workflow.md` | Test-driven development patterns |
| `wave-assignment.md` | Parallel execution planning |
| `wave-planning.md` | Wave sequencing strategy |
| `epic-completion-checklist.md` | Epic completion requirements |
| `integration-cadence.md` | Integration checkpoint guidance |
| `tool-selection-guide.md` | When to use which tool |
| `*-lessons.md` | Technology-specific patterns |

## Workflow Execution Flow

### Example: Implementing a New Feature

```
1. USER: "/sdlc"
   └── Command loads .claude/commands/sdlc.md

2. PHASE 1 - SPECIFICATION (Plugin: spec-writer)
   ├── Use spec-workflow MCP to create requirements.md
   ├── Define data models and API contracts in design.md
   └── Break into tasks in tasks.md

3. PHASE 2 - TEST (Plugin: test-writer)
   ├── Spawn test-engineer agent via Task tool
   ├── Agent writes failing tests from design.md
   └── Verify tests are RED (failing)

4. PHASE 3 - CODE (Plugin: code-implementer)
   ├── Spawn implementation-agent via Task tool
   ├── Agent implements minimal code
   └── Verify tests are GREEN (passing)

5. PHASE 4 - BUILD (Plugin: builder)
   ├── Run linting: go vet, eslint
   ├── Run type checking: go build, tsc
   ├── Run tests with coverage: go test, npm test
   └── Build artifacts: binaries, bundles

6. PHASE 5 - DOCS (Plugin: docs-generator)
   ├── Update epics-user-stories.md
   ├── Update CLAUDE.md files
   ├── Update CHANGELOG.md
   └── Run doc-consistency-checker
```

## Key Integration Points

### Spec-Workflow MCP ↔ Tasks.md ↔ Epics-User-Stories.md

```
.spec-workflow/specs/{epic}/
├── requirements.md  ────► User stories, acceptance criteria
├── design.md        ────► Data models, API contracts
└── tasks.md         ────► Task breakdown with status
                              │
                              ▼
implementation-plan/epics-user-stories.md
└── Centralized epic tracking with completion dates
```

### Agents ↔ Plugins ↔ Commands

```
/sdlc (Command)
  │
  ├─► spec-writer (Plugin) ──► spec-workflow MCP
  │
  ├─► test-writer (Plugin) ──► test-engineer (Agent)
  │                             └── Writes *_test.go, *.test.ts
  │
  ├─► code-implementer (Plugin) ──► implementation-agent (Agent)
  │                                   └── Writes source files
  │
  └─► docs-generator (Plugin) ──► documentation-generator (Agent)
                                    └── Updates CLAUDE.md, CHANGELOG.md
```

## Critical Enforcement Points

### TDD Enforcement

```
Phase 2 (Test):
├── test-engineer MUST be spawned
├── Tests MUST exist before Phase 3
└── Tests MUST be FAILING (Red)

Phase 3 (Code):
├── implementation-agent spawned
├── Code implements ONLY what tests require
└── Tests MUST be PASSING (Green)

VIOLATION: Writing code before tests
└── Delete code, restart from Phase 2
```

### Task Tracking Enforcement

```
After Phase 4 (tests pass):
├── Mark task [x] in tasks.md (MANDATORY)
├── Update acceptance criteria in epics-user-stories.md
└── Merge task branch to group branch

After all tasks complete:
├── Run doc-consistency-checker
├── Update CHANGELOG.md
├── Mark epic complete with date
└── Create PR to dev
```

## File Locations Summary

```
.claude/
├── commands/           # Entry points (/sdlc, /code-review, etc.)
├── agents/             # Specialized subagents (test-engineer, etc.)
├── plugins/            # SDLC phase orchestration
├── skills/             # Reusable capabilities
├── docs/               # Process documentation
└── settings.local.json # Local settings

.spec-workflow/
├── specs/              # Feature specifications
├── approvals/          # Approval tracking
├── templates/          # Spec templates
└── steering/           # Project steering docs

implementation-plan/
├── epics-user-stories.md    # Centralized epic tracking
└── claude-code-architecture.md  # This document
```

## Quick Reference Card

| Need | Use |
|------|-----|
| Start new feature | `/sdlc` command |
| Write tests | Spawn `test-engineer` agent |
| Write code | Spawn `implementation-agent` agent |
| Review code | `/code-review` or `code-reviewer` skill |
| Update docs | `/update-claudemd` or `documentation-generator` agent |
| Check spec status | `spec-workflow` MCP → `spec-status` tool |
| Deploy infra | `infrastructure-deployer` skill |
| Search AWS docs | `aws-documentation-mcp-server` |
| Design DynamoDB | `dynamodb-mcp-server` |
| Library docs | `context7` MCP |
