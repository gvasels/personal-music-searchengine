# SDLC Workflow Plugins

## Overview

This directory contains SDLC (Software Development Lifecycle) plugins that orchestrate the development workflow from specification to deployment. Each plugin represents a phase in the development process and can be invoked independently or as part of the full `/sdlc` workflow.

**Location:** `.claude/plugins/`

## Plugin Architecture

```
Main Agent (Claude Code)
       │
       ├── /sdlc command ─────────────────┐
       │                                   │
       ▼                                   ▼
┌─────────────┐                    ┌─────────────┐
│   Skills    │◄──────────────────►│   Plugins   │
│ (Reusable)  │                    │ (Workflows) │
└─────────────┘                    └─────────────┘
       │                                   │
       ▼                                   ▼
┌─────────────┐                    ┌─────────────┐
│   Agents    │◄───── spawns ──────│  MCP Tools  │
│ (Subagents) │                    │  (External) │
└─────────────┘                    └─────────────┘
```

## Plugins

| File | Phase | Description |
|------|-------|-------------|
| `spec-writer.md` | 1 | Requirements gathering and technical design |
| `test-writer.md` | 2 | TDD test creation before implementation |
| `code-implementer.md` | 3 | Implementation code to make tests pass |
| `builder.md` | 4 | Build verification, linting, type checking |
| `security-checker.md` | 5 | Security vulnerability scanning and audit |
| `docs-generator.md` | 6 | API and code documentation generation |

## SDLC Flow

```
┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
│  SPEC   │──►│  TEST   │──►│  CODE   │──►│  BUILD  │──►│SECURITY │──►│  DOCS   │
│ Phase 1 │   │ Phase 2 │   │ Phase 3 │   │ Phase 4 │   │ Phase 5 │   │ Phase 6 │
└─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘
     │             │             │             │             │             │
     ▼             ▼             ▼             ▼             ▼             ▼
 spec-workflow  test-engineer  implementation  bash/lint   security-    documentation-
    MCP           agent          agent        tools        auditor        generator
                                                           agent          agent
```

## Plugin-to-Agent Mapping

| Plugin | Primary Agent | Skills Used | MCP Servers |
|--------|---------------|-------------|-------------|
| spec-writer | - | - | spec-workflow |
| test-writer | test-engineer | - | - |
| code-implementer | implementation-agent | - | - |
| builder | - | - | - |
| security-checker | security-auditor | code-reviewer | - |
| docs-generator | documentation-generator | documentation-generator | - |

## Usage

### Full Workflow

```
/sdlc
```

This invokes the complete SDLC workflow, guiding through all 6 phases.

### Individual Phases

Each plugin can be invoked independently:

- **Spec phase**: "Run the spec phase for feature X"
- **Test phase**: "Write tests for this design"
- **Code phase**: "Implement this feature"
- **Build phase**: "Build and verify the code"
- **Security phase**: "Run security scan"
- **Docs phase**: "Generate documentation"

## Phase Prerequisites

| Phase | Requires |
|-------|----------|
| 1. Spec | Feature request or GitHub issue |
| 2. Test | requirements.md, design.md |
| 3. Code | Failing tests |
| 4. Build | Implementation code |
| 5. Security | Passing build |
| 6. Docs | Security-approved code |

## Quality Gates

| Phase | Gate Criteria |
|-------|---------------|
| Spec | Design approved via spec-workflow |
| Test | Tests written and failing (Red) |
| Code | All tests passing (Green) |
| Build | Lint, Types, Coverage 80%+ |
| Security | 0 critical/high vulnerabilities |
| Docs | OpenAPI valid, CLAUDE.md updated |

## Failure Handling

Each plugin documents how to handle failures:

- **Spec fails** → Revise requirements with stakeholder
- **Tests fail to cover** → Add more test cases
- **Code fails tests** → Fix implementation
- **Build fails** → Fix lint/type/coverage issues
- **Security fails** → Remediate vulnerabilities
- **Docs incomplete** → Generate missing documentation

## Adding New Plugins

To add a new SDLC plugin:

1. Create `{plugin-name}.md` in `.claude/plugins/`
2. Include frontmatter with phase number and dependencies
3. Document workflow steps
4. Specify agents/skills to leverage
5. Define inputs/outputs
6. Update this PLUGINS.md
