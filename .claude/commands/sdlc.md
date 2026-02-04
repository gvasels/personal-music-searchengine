---
name: sdlc
description: SDLC workflow - spec, test, code, verify tests pass
---

# SDLC Workflow

This workflow guides development from specification to verified implementation.

## Workflow Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LOCAL DEVELOPMENT (Per Task)                              │
├─────────┬─────────┬─────────┬─────────┬─────────────────────────────────────┤
│ 1.SPEC  │ 2.TEST  │ 3.CODE  │ 4.VERIFY│ 5.DOCS (Epic Complete Only)         │
│─────────│─────────│─────────│─────────│─────────────────────────────────────│
│ MCP:    │ Agent:  │ Agent:  │ Run     │ Agent: doc-consistency-checker      │
│ spec-   │ test-   │ implmnt-│ tests   │ Checklist: epic-completion-         │
│ workflow│ engineer│ agent   │ pass →  │ checklist.md                        │
│         │         │         │ Close   │ Update: CHANGELOG, CLAUDE.md        │
│         │         │         │ task    │                                     │
└─────────┴─────────┴─────────┴─────────┴─────────────────────────────────────┘
                                        │
                    All tasks complete? │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       CREATE PR TO DEV                                       │
│   PR triggers GitHub Actions: Build, Code Review, Security, Docs             │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     CI/CD PIPELINE (Automated)                               │
├─────────────────────┬─────────────────────┬─────────────────────────────────┤
│   GITHUB ACTIONS    │   DEV DEPLOYMENT    │   BUILDKITE VERIFY              │
│─────────────────────│─────────────────────│─────────────────────────────────│
│  - Build            │   Merge to dev      │  - Health checks                │
│  - Code review      │   triggers deploy   │  - Smoke tests                  │
│  - Security scans   │                     │  - Rollback on failure          │
│  - Documentation    │                     │                                 │
└─────────────────────┴─────────────────────┴─────────────────────────────────┘
```

## Branching Strategy

```
main (prod) ←── staging ←── dev ←── group-N/{name}
     │                              │
     │                              └── task-N.X/{description}
     │
     └── release/vX.Y.Z (for rollbacks)
```

| Branch | Purpose | Deploys To |
|--------|---------|------------|
| `main` | Production releases | Prod accounts |
| `staging` | Pre-production validation | Staging accounts |
| `dev` | Integration branch | Dev accounts |
| `group-N/{name}` | Feature group work | - |
| `task-N.X/{description}` | Individual task | - |
| `release/vX.Y.Z` | Release tags | Rollback reference |

---

## Phase 1: Specification

**Goal**: Define requirements and design before code.

**Tools**: Use `spec-workflow` MCP server:
- `spec-workflow-guide` - Load workflow instructions
- `spec-status` - Check specification progress
- `approvals` - Request document approvals

**Workflow**:
1. Load spec-workflow guide
2. Create `requirements.md` with user stories
3. Create `design.md` with data models and API contracts
4. Create `tasks.md` with implementation breakdown
5. Request approval for each document

**Artifacts**: `requirements.md`, `design.md`, `tasks.md`

---

## Phase 2: Test

**Goal**: Write failing tests BEFORE implementation (TDD Red phase).

**ENFORCEMENT**: This phase is REQUIRED for TDD. If you skip spawning test-engineer agent, you are violating the SDLC workflow.

**REQUIRED**: Spawn test-engineer agent via Task tool:

```
Task tool parameters:
- description: "Write tests for [feature]"
- subagent_type: "test-engineer"
- prompt: |
    Write tests for [feature name].

    Design spec: [path to design.md]

    Requirements:
    - Unit tests for data model validation
    - Integration tests for API endpoints (if applicable)
    - Test all success and error scenarios from design.md
    - Tests should FAIL initially (Red phase)

    For Application code:
    - Test framework: Vitest/Jest (TS) or Testify (Go)
    - Output: tests/[unit|integration]/[feature]/

    For Infrastructure (OpenTofu):
    - Create test fixtures (tfvars files)
    - Validation tests will run in Phase 4
```

**Verification**: After test-engineer completes:
- [ ] Test files exist in appropriate test directories
- [ ] Tests run and FAIL (Red phase - this is correct!)
- [ ] Test coverage metrics available (even though tests fail)

**Artifacts**: Test files, fixtures

**IMPORTANT**: If you write production code before writing tests, you have violated TDD principles. Delete the production code and restart from this phase.

---

## Phase 3: Code

**Goal**: Write minimal code to make tests pass (TDD Green phase).

**REQUIRED**: Spawn implementation-agent via Task tool:

```
Task tool parameters:
- description: "Implement [feature]"
- subagent_type: "implementation-agent"
- prompt: |
    Implement [feature name] to make tests pass.

    Design spec: [path to design.md]
    Tests location: [path to test files]

    Requirements:
    - Implement data models from design
    - Create API endpoints per contract (if applicable)
    - Only write code needed to pass tests
    - Follow existing patterns in codebase

    For Infrastructure (OpenTofu):
    - Module location: infrastructure/modules/[name]/
    - Required files: main.tf, variables.tf, outputs.tf, README.md, CLAUDE.md
    - Follow existing module patterns
```

**Artifacts**: Source files, module files

---

## Phase 4: Verify

**Goal**: Confirm all tests pass, then update tasks.md and close task issue.

### For Application Code

```bash
# Run tests
npm test
# or: go test ./...

# Verify coverage meets threshold (80%+)
npm test -- --coverage
```

### For Infrastructure (OpenTofu)

```bash
# Format check
tofu fmt -check -recursive infrastructure/modules/[name]/

# Validate
cd infrastructure/modules/[name]
tofu init -backend=false
tofu validate

# Clean up
rm -rf .terraform .terraform.lock.hcl
```

### Task Completion Checklist

**CRITICAL**: When tests pass, you MUST complete ALL of these steps:

1. **Update tasks.md** (REQUIRED):
   ```bash
   # Mark task as complete in .spec-workflow/specs/{spec-name}/tasks.md
   # Change: - [ ] Task N.X: Description
   # To:     - [x] Task N.X: Description
   ```

2. **Commit changes** to task branch:
   ```bash
   git add .
   git commit -m "feat(task-N.X): [description]

   - Implemented [feature]
   - Tests passing with X% coverage
   - Updated tasks.md

   Closes #[issue-number]"
   ```

3. **Merge task to group** branch:
   ```bash
   git checkout group-N/{name}
   git merge task-N.X/{description}
   git branch -d task-N.X/{description}
   ```

4. **Close task issue** (if using GitHub Issues):
   ```bash
   gh issue close {issue_number} --comment "✅ Tests passing (X% coverage), merged to group branch, tasks.md updated"
   ```

**DO NOT SKIP** updating tasks.md - this tracks implementation progress for the entire epic.

---

## Phase 5: Documentation (Epic Completion Only)

**Goal**: Validate and update documentation before marking epic complete.

**When to Run**: Only when ALL tasks in an epic are complete (not per-task).

**CRITICAL**: Epic documentation updates are MANDATORY before creating PR to dev. These are NOT optional checklists - they are REQUIRED steps.

**REQUIRED**: Follow epic completion checklist:

```bash
# Reference the checklist
cat .claude/docs/epic-completion-checklist.md
```

### Epic Completion Checklist (MANDATORY)

**CRITICAL**: When ALL tasks in an epic are complete, you MUST complete these steps BEFORE creating PR:

**Step 1: Update epics-user-stories.md (REQUIRED)**
```bash
# Edit implementation-plan/epics-user-stories.md
# For the completed epic:
# 1. Add completion status: ✅ **COMPLETE** (YYYY-MM-DD)
# 2. Check ALL acceptance criteria: - [x]
# 3. Add implementation notes to criteria (what was built, account IDs, etc.)
# 4. Move epic to "Completed Epics" section if it exists
```

**Example:**
```markdown
## Epic 4: CI/CD Pipelines ✅ **COMPLETE** (2024-12-22)

#### US-4.1: Buildkite + CodeBuild Infrastructure
**Acceptance Criteria**:
- [x] Buildkite organization configured _(3 webhook handlers deployed to accounts 634945387634, 543613944458, 471544433440)_
- [x] CodeBuild projects created _(infrastructure, lambda, microfrontend runners)_
- [x] Webhook integration triggering CodeBuild _(HMAC-verified Lambda at POST /webhook)_
```

**Step 2: Update CLAUDE.md (REQUIRED)**
```bash
# Edit root CLAUDE.md
# Update project status line with latest epic
# Example: "Epic 4 (CI/CD Pipelines) complete. Epic 5 in progress."
```

**Step 3: Update CHANGELOG.md (REQUIRED)**
```bash
# Add comprehensive epic entry with:
# - All deliverables (modules, services, pipelines)
# - Security enhancements (IAM, secrets, etc.)
# - Architecture decisions (versioned paths, canary, etc.)
# - Breaking changes (if any)
```

**Step 4: Verify Documentation Consistency (REQUIRED)**
- [ ] All acceptance criteria checked: `- [x]`
- [ ] Epic has completion date: `✅ **COMPLETE** (YYYY-MM-DD)`
- [ ] Implementation notes added to criteria (account IDs, module names, etc.)
- [ ] CLAUDE.md project status reflects epic completion
- [ ] CHANGELOG.md has detailed epic entry

**IMPORTANT**: The epics-user-stories.md update is NOT optional. This is how we track epic completion across the entire project. Skipping this step means the epic is NOT considered complete.

### Additional Documentation Updates

**2. Planning Document Sync**
- [ ] Review `planning-docs/*.md` for:
  - Fictional account names → Update to actual account IDs
  - Architecture diagrams → Reflect current state
  - Technology stack references → Current implementation
- [ ] Cross-reference account names with `infrastructure/docs/aws-organizations.md`
- [ ] Mark future accounts explicitly with "future:" prefix or 🔮 emoji

**3. Code Documentation**
- [ ] Generate/update `CLAUDE.md` files for modified directories:
  ```bash
  # Use documentation-generator agent
  /documentation-generator infrastructure/modules/{new-module}/
  ```
- [ ] Update API documentation (OpenAPI specs if applicable)

### Automated Validation

**REQUIRED**: Run doc-consistency-checker agent before marking epic complete:

```bash
# Spawn agent via Task tool
Task tool parameters:
- description: "Validate documentation consistency"
- subagent_type: "general-purpose"
- prompt: |
    Use the doc-consistency-checker agent instructions to validate
    documentation consistency across the codebase.

    Agent file: .claude/agents/doc-consistency-checker.md

    Focus on:
    1. Account name validation (vs infrastructure/docs/aws-organizations.md)
    2. Epic status validation (vs implementation-plan/epics-user-stories.md)
    3. Technology stack validation
    4. Future vs. Current distinction
    5. Module documentation completeness

    Generate report with:
    - Summary statistics
    - Issues by category
    - Detailed issue list with severity
    - Recommended fixes
```

### Documentation Validation Checklist

- [ ] No HIGH severity issues from doc-consistency-checker
- [ ] All account names match `aws-organizations.md`
- [ ] Epic statuses consistent across all docs
- [ ] No outdated technology references (e.g., EC2 agents instead of CodeBuild)
- [ ] Deployed accounts have IDs, future accounts marked
- [ ] All new modules have `CLAUDE.md` files

### Exit Criteria

Documentation phase passes when:
1. Epic completion checklist fully executed
2. Doc-consistency-checker reports no HIGH severity issues
3. All documentation artifacts committed

**Artifacts**: Updated CHANGELOG.md, CLAUDE.md, planning docs, module docs

---

## Group Completion: PR to Dev

When all tasks in a group are complete:

```bash
# Create PR to dev
gh pr create \
  --base dev \
  --head group-N/{name} \
  --title "Feature: [Group description]" \
  --body "$(cat <<'EOF'
## Summary
[Brief description of the feature group]

## Tasks Completed
- [x] Task N.1: [description]
- [x] Task N.2: [description]
- [x] Task N.3: [description]

## Testing
- All unit tests passing
- Integration tests passing (if applicable)
- OpenTofu validation passing (if applicable)

## Specs
- Design: `.spec-workflow/specs/{spec}/design.md`
- Tasks: `.spec-workflow/specs/{spec}/tasks.md`

---
🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

---

## CI/CD Pipeline (Automated)

These run automatically after PR creation - **do not run locally**.

### GitHub Actions (on PR)

| Check | Workflow | Purpose |
|-------|----------|---------|
| Build | `build.yml` | Lint, type check, compile |
| Code Review | `claude-code-review.yml` | AI-assisted review |
| Security | `security.yml` | gitleaks, checkov, dependency scan |
| Documentation | `docs.yml` | Verify CLAUDE.md, README updates |

### Buildkite (on Deploy to Dev)

| Check | Purpose |
|-------|---------|
| Health checks | Verify service is responding |
| Smoke tests | Basic functionality verification |
| Contract tests | API contract validation |
| Rollback | Automatic on failure |

---

## Workflow Modes

| Mode | Detect By | Key Differences |
|------|-----------|-----------------|
| **Application** | `.ts`, `.tsx`, `.go`, `.py` files | Unit/integration tests |
| **Infrastructure** | `.tf` files | tofu validate + security scan |

---

## Interactive Sessions vs @claude Automation

### Interactive `/sdlc` (Local Claude Code)

All MCP servers available:
- `spec-workflow` - Specification management
- `github` - Issue and PR management
- `context7` - Documentation lookup
- `aws-*` - AWS documentation and tools
- `terraform` - OpenTofu/Terraform docs

Use full MCP capabilities for exploration and complex tasks.

### @claude Automation (GitHub Issues)

Specify MCPs per task to minimize context. Include in issue body:

**For Spec tasks:**
```
@claude implement this task
MCPs: spec-workflow
```

**For Application Test/Code tasks:**
```
@claude implement this task following TDD
MCPs: context7 (for library docs)
```

**For Infrastructure tasks:**
```
@claude implement this OpenTofu module
MCPs: terraform, aws-documentation
```

**For Go services:**
```
@claude implement this Go Lambda function
MCPs: context7, godoc
```

---

## Agent Spawning Reference

| Phase | Agent | subagent_type |
|-------|-------|---------------|
| 2. Test | test-engineer | `test-engineer` |
| 3. Code | implementation-agent | `implementation-agent` |
| 5. Docs (Epic Complete) | doc-consistency-checker | `general-purpose` |

**Note**:
- Phase 5 (Documentation) runs locally ONLY for epic completion, not per-task
- Additional security and documentation checks run in CI/CD after PR creation

---

## Status Tracking

Use TodoWrite throughout:

```
[ ] Phase 1: Specification
    [ ] requirements.md created
    [ ] design.md created
    [ ] tasks.md created
[ ] Phase 2: Test
    [ ] Spawn test-engineer agent
    [ ] Tests written and failing (Red)
    [ ] Verified tests run and show FAILURES
[ ] Phase 3: Code
    [ ] Verified tests exist before coding
    [ ] Spawn implementation-agent
    [ ] Code implemented
    [ ] Verified tests now PASS (Green)
[ ] Phase 4: Verify
    [ ] Tests passing with coverage threshold met
    [ ] **tasks.md updated** (mark task [x] complete)
    [ ] Changes committed to task branch
    [ ] Merged to group branch
    [ ] Task issue closed
[ ] Phase 5: Documentation (Epic Complete Only)
    [ ] Epic completion checklist executed
    [ ] Doc-consistency-checker validation passing
    [ ] CHANGELOG.md updated
    [ ] CLAUDE.md project status updated
    [ ] Planning docs synced with implementation
```

**CRITICAL REMINDER**: The tasks.md update in Phase 4 is NOT optional. It is required to track epic progress. Update it IMMEDIATELY after tests pass, before merging.

---

## Quick Reference

### Start a Task
```bash
git checkout group-N/{name}
git checkout -b task-N.X/{description}
```

### Complete a Task
```bash
# Tests pass
git add . && git commit -m "feat: [description]"
git checkout group-N/{name}
git merge task-N.X/{description}
git branch -d task-N.X/{description}
```

### Complete a Group
```bash
gh pr create --base dev --head group-N/{name}
```

### Promote Through Environments
```
dev → staging → main
```

Each promotion is a PR with required reviews and CI checks.

---

## Begin Workflow

When user invokes `/sdlc`:

1. **Detect mode**: Application or Infrastructure?
2. **Check for existing spec** or start Phase 1
3. **Execute Phases 1-4** for each task
4. **Close task issues** when tests pass
5. **When epic complete**: Execute Phase 5 (Documentation validation)
6. **Create PR to dev** when group complete

**Remember**:
- Phases 1-4 run locally per task
- Phase 5 runs locally ONLY for epic completion (not per-task)
- Build, security, and additional documentation checks run in CI/CD after PR creation
