# CLAUDE.md - Personal Music Search Engine

## HARD REQUIREMENTS

**Non-negotiable rules. Violations require rollback and restart.**

| # | Rule | Requirement | If Violated |
|---|------|-------------|-------------|
| 1 | **SDLC Workflow** | Use `/sdlc` for ALL features: SPEC → TEST → CODE → BUILD → DOCS | Start over with `/sdlc` |
| 2 | **TDD Enforcement** | Spawn `test-engineer` BEFORE `implementation-agent`. Tests must fail first (Red). | Delete code, restart Phase 2 |
| 3 | **Spec-First** | All features need specs in `.spec-workflow/specs/` before coding | Create specs first |
| 4 | **Documentation** | Update CLAUDE.md + CHANGELOG.md in every directory with changes | Update before PR |
| 5 | **Quality Gates** | 80% test coverage, all tests pass, no critical security vulnerabilities | Fix before merge |

---

## Project Overview

| Field | Value |
|-------|-------|
| **Name** | Personal Music Search Engine |
| **Description** | Music library app for uploading, organizing, searching, streaming audio files |
| **Status** | Development |
| **AWS Account** | 887395463840 |
| **AWS Profile** | `gvasels-muza` |
| **Region** | us-east-1 |

## Technology Stack

| Layer | Technologies |
|-------|--------------|
| **Backend** | Go 1.22+, Echo v4, AWS Lambda, dhowden/tag |
| **Frontend** | React 18, TanStack Router/Query, Tailwind + DaisyUI 5, Zustand, Vite |
| **Infrastructure** | AWS, OpenTofu 1.8+ (NOT Terraform), GitHub Actions |
| **Data** | DynamoDB (single-table), S3, Nixiesearch, Cognito, CloudFront |

---

## SDLC Workflow

```
/sdlc → SPEC → TEST → CODE → BUILD → DOCS → PR
         ↓      ↓      ↓       ↓       ↓
       spec-  test-  impl-   lint/  docs-
       workflow engineer agent  test  generator
```

| Phase | Purpose | Artifacts | Quality Gate |
|-------|---------|-----------|--------------|
| 1. SPEC | Requirements & design | `requirements.md`, `design.md`, `tasks.md` | Design approved |
| 2. TEST | TDD Red - failing tests | `*_test.go`, `*.test.ts` | Tests written & failing |
| 3. CODE | TDD Green - minimal code | Source files | All tests passing |
| 4. BUILD | Lint, type check, coverage | Build output | Coverage 80%+ |
| 5. DOCS | API docs, CLAUDE.md | Documentation | Files updated |

**Prerequisites**: Each phase requires previous phase completion.

---

## Repository Structure

```
├── backend/                 # Go Lambda services
│   ├── cmd/{api,indexer,processor}/  # Lambda entrypoints
│   └── internal/{handlers,models,repository,service,metadata,search}/
├── frontend/src/            # React app
│   └── {components,pages,hooks,lib,routes}/
├── infrastructure/          # OpenTofu IaC
│   └── {global,shared,backend,frontend}/
├── .claude/                 # Claude Code automation
│   └── {plugins,agents,skills,commands,docs}/
├── .spec-workflow/specs/    # Feature specifications
└── .github/workflows/       # CI/CD
```

---

## Quick Reference

### Commands

| Command | Purpose |
|---------|---------|
| `/sdlc` | Start SDLC workflow (REQUIRED for features) |
| `/update-claudemd` | Update CLAUDE.md from changes |
| `/code-review` | Run code review |
| `/test-file <path>` | Generate tests for file |

### Build & Test

```bash
# Backend
cd backend && go test ./...
cd backend && go build ./cmd/api

# Frontend
cd frontend && npm test && npm run build

# Infrastructure
cd infrastructure/{global,shared,backend} && tofu plan
```

### Key Agents

| Agent | When to Spawn |
|-------|---------------|
| `test-engineer` | Phase 2 - ALWAYS FIRST |
| `implementation-agent` | Phase 3 - AFTER tests exist |
| `documentation-generator` | Phase 5 |
| `code-review` | Before PR |

### MCP Servers

| Server | Purpose |
|--------|---------|
| `spec-workflow` | Spec management & approvals |
| `context7` | Library documentation |
| `github` | GitHub API |
| `dynamodb-mcp-server` | DynamoDB design |
| `opentofu` | OpenTofu modules |

---

## Branching Strategy

```
main (protected)
  └── group-N/{name}           # Feature group
        ├── task-N.1/{desc}    # Task branches
        └── task-N.2/{desc}
              └── [PR to group → PR to main]
```

| Type | Pattern | Example |
|------|---------|---------|
| Group | `group-N/{name}` | `group-1/foundation-backend` |
| Task | `task-N.X/{desc}` | `task-1.1/domain-models` |
| Hotfix | `hotfix/{desc}` | `hotfix/lambda-timeout` |

---

## DynamoDB Schema (Single-Table)

| Entity | PK | SK |
|--------|----|----|
| User | `USER#{userId}` | `PROFILE` |
| Track | `USER#{userId}` | `TRACK#{trackId}` |
| Album | `USER#{userId}` | `ALBUM#{albumId}` |
| Playlist | `USER#{userId}` | `PLAYLIST#{playlistId}` |
| Upload | `USER#{userId}` | `UPLOAD#{uploadId}` |
| Tag | `USER#{userId}` | `TAG#{tagName}` |

---

## Documentation Standards

### Subdirectory CLAUDE.md Requirements

Every major subdirectory needs a CLAUDE.md with:
1. Overview - Brief description
2. File Descriptions - Purpose of each file
3. Key Functions - Signatures and behavior
4. Dependencies - Internal and external
5. Usage Examples

**Exception**: `.claude/` uses root files: `AGENTS.md`, `COMMANDS.md`, `PLUGINS.md`, `SKILLS.md`

### CHANGELOG.md Format

```markdown
## [Unreleased]
### Added
### Changed
### Fixed
```

---

## Reference Docs

Located in `.claude/docs/`:
- `tdd-workflow.md` - TDD practices
- `wave-planning.md` - Parallel @claude execution
- `epic-completion-checklist.md` - Epic completion
- `task-granularity.md` - Task breakdown rules

Current specs: `.spec-workflow/specs/music-search-engine/`
