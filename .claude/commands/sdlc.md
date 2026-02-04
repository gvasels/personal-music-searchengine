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

## MCP Servers by Phase

**REQUIRED**: Use the correct MCP servers at each phase. The orchestrator (main Claude session) MUST use these to gather context before spawning agents and to pass relevant information in agent prompts.

| Phase | MCP Servers | Purpose |
|-------|-------------|---------|
| **1. SPEC** | `spec-workflow`, `github`, `awslabs.dynamodb-mcp-server` | Spec creation/approval, issue tracking, data modeling for DynamoDB entities |
| **2. TEST (Frontend)** | `context7`, `docs-mcp-server`, `playwright` | Vitest/React Testing Library APIs, framework docs, E2E browser automation |
| **2. TEST (Backend)** | `context7`, `docs-mcp-server` | Testify/httptest APIs, framework docs |
| **3. CODE (Frontend)** | `context7`, `daisyui-blueprint`, `docs-mcp-server` | React/TanStack/Zustand APIs, DaisyUI component snippets, library docs |
| **3. CODE (Backend)** | `context7`, `awslabs.dynamodb-mcp-server`, `awslabs.aws-documentation-mcp-server` | Echo/Go library APIs, DynamoDB access patterns, AWS service docs |
| **3. CODE (Infra)** | `opentofu`, `awslabs.aws-documentation-mcp-server`, `aws-knowledge-mcp-server` | OpenTofu module registry, AWS service docs, AWS recommendations |
| **4. VERIFY** | `playwright` | E2E test execution |
| **5. DOCS** | `github`, `spec-workflow` | PR management, spec status verification |
| **Code Review** | `github` | PR reading, diff review, comment management |
| **Security Audit** | `awslabs.aws-documentation-mcp-server`, `aws-knowledge-mcp-server` | AWS security best practices, compliance guidance |

### MCP Usage Rules

1. **Before spawning agents**: Use relevant MCPs to gather context (API signatures, patterns, examples)
2. **In agent prompts**: Include MCP-gathered context so agents have the information they need
3. **For @claude automation**: Specify MCPs per task in GitHub issue body (see Interactive Sessions section)
4. **Prefer `context7`** for any library API lookup (Go, TypeScript, React, testing frameworks)
5. **Use `daisyui-blueprint`** for ALL frontend UI component work (not just styling)

---

## Phase 1: Specification

**Goal**: Define requirements and design before code.

**MCP Servers**: `spec-workflow`, `github`, `awslabs.dynamodb-mcp-server`

**Tools**:
- `spec-workflow` MCP: `spec-workflow-guide`, `spec-status`, `approvals`
- `github` MCP: Search existing issues, reference prior work, check for duplicates
- `awslabs.dynamodb-mcp-server` MCP: `dynamodb_data_modeling` for DynamoDB entity design (if feature involves data changes)

**Workflow**:
1. Load spec-workflow guide
2. Search GitHub issues for related work or prior specs
3. Create `requirements.md` with user stories
4. Create `design.md` with data models and API contracts
   - Use `dynamodb_data_modeling` for DynamoDB access pattern design
5. Create `tasks.md` with implementation breakdown
   - **MUST** separate unit test tasks, E2E test tasks, and implementation tasks
6. Request approval for each document

**Artifacts**: `requirements.md`, `design.md`, `tasks.md`

---

## Phase 2: Test

**Goal**: Write failing tests BEFORE implementation (TDD Red phase).

**ENFORCEMENT**: This phase is REQUIRED for TDD. If you skip spawning test-engineer agent, you are violating the SDLC workflow.

### CRITICAL TDD RULES

**You MUST write ALL tests BEFORE ANY implementation code.**

#### Step 2.1: Write Unit Tests FIRST

```
# For frontend components - MUST create:
frontend/src/routes/__tests__/{component}.test.tsx

# For backend - MUST create:
backend/internal/{package}/{file}_test.go
```

#### Step 2.2: Write E2E Tests (if UI component)

```
# MUST create:
frontend/e2e/{feature}.spec.ts
```

#### Step 2.3: VERIFY TESTS FAIL (RED)

**You MUST run tests and VERIFY they FAIL before writing any implementation.**

**⚠️ STOP IF TESTS PASS. Tests are not testing the right thing.**

### Pre-Agent MCP Context Gathering

**BEFORE spawning test-engineer**, the orchestrator MUST:

1. **Frontend tests**: Use `context7` to look up Vitest and React Testing Library APIs
2. **Backend tests**: Use `context7` to look up Go Testify and httptest APIs
3. **E2E tests**: Use `context7` to look up Playwright test patterns
4. **Framework docs**: Use `docs-mcp-server` for any framework-specific testing patterns
5. Include gathered API examples in the agent prompt below

### Agent Spawning

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
    - Tests MUST FAIL initially (Red phase)
    - You MUST run tests and verify they FAIL before completing

    For Frontend (TypeScript):
    - Test framework: Vitest + React Testing Library
    - Unit tests: frontend/src/routes/__tests__/{component}.test.tsx
    - E2E tests: frontend/e2e/{feature}.spec.ts (Playwright)
    - [Include context7 API examples here]

    For Backend (Go):
    - Test framework: Testify + httptest
    - Unit tests: backend/internal/{package}/{file}_test.go
    - Integration tests: backend/test/{feature}_test.go
    - [Include context7 API examples here]

    For Infrastructure (OpenTofu):
    - Create test fixtures (tfvars files)
    - Validation tests will run in Phase 4
```

**Verification**: After test-engineer completes:
- [ ] Unit test file exists
- [ ] Unit test was run and FAILED before implementation
- [ ] E2E test file exists (if UI component)
- [ ] E2E test was run and FAILED before implementation
- [ ] Test coverage metrics available (even though tests fail)

**Artifacts**: Test files, fixtures

**IMPORTANT**: If you write production code before writing tests, you have violated TDD principles. Delete the production code and restart from this phase.

---

## Phase 3: Code

**Goal**: Write minimal code to make tests pass (TDD Green phase).

**You MUST write ONLY the minimal code to make tests pass.**

### CRITICAL CODE RULES

#### Step 3.1: Implement Feature
- Read test file to understand exact requirements
- Write ONLY code needed to pass tests
- NO extra features, NO premature optimization
- NEVER add features not covered by tests

#### Step 3.2: VERIFY TESTS PASS (GREEN)
- Run ALL tests (unit + E2E)
- ALL tests MUST pass before marking phase complete

### Pre-Agent MCP Context Gathering

**BEFORE spawning implementation-agent**, the orchestrator MUST:

1. **Frontend code**: Use `context7` for React/TanStack Router/TanStack Query/Zustand APIs
2. **Frontend UI**: Use `daisyui-blueprint` for DaisyUI component snippets and patterns
3. **Backend code**: Use `context7` for Echo v4/Go library APIs
4. **Backend DynamoDB**: Use `awslabs.dynamodb-mcp-server` for DynamoDB access patterns
5. **Backend AWS**: Use `awslabs.aws-documentation-mcp-server` for AWS service integration docs
6. **Infrastructure**: Use `opentofu` for module registry lookups
7. **Infrastructure AWS**: Use `aws-knowledge-mcp-server` for AWS service recommendations
8. Include gathered API examples and patterns in the agent prompt below

### Agent Spawning

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
    - Read test files FIRST to understand exact requirements
    - Implement data models from design
    - Create API endpoints per contract (if applicable)
    - ONLY write code needed to pass tests
    - NEVER add features not covered by tests
    - Follow existing patterns in codebase
    - Run tests and verify ALL PASS before completing

    For Frontend (TypeScript):
    - React 18 + TanStack Router/Query + Zustand + DaisyUI 5
    - [Include context7 API examples here]
    - [Include daisyui-blueprint component snippets here]

    For Backend (Go):
    - Go 1.22+ + Echo v4 + AWS SDK v2 + DynamoDB single-table
    - [Include context7 API examples here]
    - [Include DynamoDB access patterns here]

    For Infrastructure (OpenTofu):
    - Module location: infrastructure/modules/[name]/
    - Required files: main.tf, variables.tf, outputs.tf, README.md, CLAUDE.md
    - Follow existing module patterns
    - [Include opentofu module examples here]
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

1. **Update CLAUDE.md for changed directories** (REQUIRED):
   ```bash
   # Find all directories with changes in this task
   git diff --name-only HEAD~1 | xargs -I {} dirname {} | sort -u

   # For EACH directory with code changes:
   # - If CLAUDE.md exists: UPDATE it with new/changed files, functions, exports
   # - If CLAUDE.md does NOT exist: CREATE it following template in documentation-generator agent
   # - Exception: .claude/ directory uses .claude/docs/ instead
   ```

   **⚠️ DO NOT PROCEED to commit if any changed directory is missing CLAUDE.md.**

   Verification:
   ```bash
   # List directories changed in this task that lack CLAUDE.md
   for dir in $(git diff --name-only HEAD~1 | xargs -I {} dirname {} | sort -u); do
     if [ -d "$dir" ] && [ ! -f "$dir/CLAUDE.md" ] && [[ "$dir" != .claude* ]]; then
       echo "MISSING: $dir/CLAUDE.md"
     fi
   done
   # Output MUST be empty before proceeding
   ```

2. **Update tasks.md** (REQUIRED):
   ```bash
   # Mark task as complete in .spec-workflow/specs/{spec-name}/tasks.md
   # Change: - [ ] Task N.X: Description
   # To:     - [x] Task N.X: Description
   ```

3. **Commit changes** to task branch:
   ```bash
   git add .
   git commit -m "feat(task-N.X): [description]

   - Implemented [feature]
   - Tests passing with X% coverage
   - CLAUDE.md updated for changed directories
   - Updated tasks.md

   Closes #[issue-number]"
   ```

4. **Merge task to group** branch:
   ```bash
   git checkout group-N/{name}
   git merge task-N.X/{description}
   git branch -d task-N.X/{description}
   ```

5. **Close task issue** (if using GitHub Issues):
   ```bash
   gh issue close {issue_number} --comment "✅ Tests passing (X% coverage), CLAUDE.md updated, merged to group branch, tasks.md updated"
   ```

**DO NOT SKIP** updating tasks.md or CLAUDE.md - these track implementation progress and documentation for the entire epic.

---

## Phase 5: Documentation (Epic Completion Only)

**Goal**: Validate and update documentation before marking epic complete.

**When to Run**: Only when ALL tasks in an epic are complete (not per-task).

**CRITICAL**: Epic documentation updates are MANDATORY before creating PR to dev. These are NOT optional checklists - they are REQUIRED steps.

**MCP Servers**: `github`, `spec-workflow`

**REQUIRED**: Follow epic completion checklist:

```bash
# Reference the checklist
cat .claude/docs/epic-completion-checklist.md
```

---

### Step 1: CLAUDE.md Discovery and Verification (MANDATORY)

**⚠️ DO NOT PROCEED to any other documentation step until this is complete.**

**You MUST discover ALL directories that were created or modified during this epic, and verify each one has an up-to-date CLAUDE.md.**

#### Step 1.1: Discover Changed Directories

```bash
# Find ALL directories changed since branch diverged from base
git diff --name-only $(git merge-base HEAD main) HEAD | xargs -I {} dirname {} | sort -u | grep -v '^\.$'
```

#### Step 1.2: Check for Missing CLAUDE.md

```bash
# For EACH changed directory, check if CLAUDE.md exists
# Exception: .claude/ uses .claude/docs/ instead
MISSING=""
for dir in $(git diff --name-only $(git merge-base HEAD main) HEAD | xargs -I {} dirname {} | sort -u | grep -v '^\.$'); do
  if [ -d "$dir" ] && [ ! -f "$dir/CLAUDE.md" ] && [[ "$dir" != .claude* ]] && [[ "$dir" != .spec-workflow* ]] && [[ "$dir" != .github* ]]; then
    MISSING="$MISSING\n  MISSING: $dir/CLAUDE.md"
  fi
done
if [ -n "$MISSING" ]; then
  echo "⚠️ STOP - Missing CLAUDE.md files:$MISSING"
  echo "You MUST create these before proceeding."
else
  echo "✅ All changed directories have CLAUDE.md"
fi
```

**⚠️ STOP IF ANY CLAUDE.md IS MISSING. Create it before proceeding.**

#### Step 1.3: Spawn Documentation Generator for Missing/Outdated CLAUDE.md

**REQUIRED**: For each directory that needs CLAUDE.md created or updated, spawn documentation-generator:

```
Task tool parameters:
- description: "Generate/update CLAUDE.md for [directories]"
- subagent_type: "documentation-generator"
- prompt: |
    Generate or update CLAUDE.md files for the following directories:

    Directories that need CLAUDE.md CREATED:
    - [list directories missing CLAUDE.md]

    Directories that need CLAUDE.md UPDATED:
    - [list directories with code changes since last CLAUDE.md update]

    For each directory, the CLAUDE.md MUST include:
    1. Overview - Brief description (1-3 sentences)
    2. Files - Table listing each file and purpose
    3. Key Functions/Exports - Public API with signatures
    4. Dependencies - Internal and external
    5. Usage Examples - Common patterns

    Template: Read .claude/agents/documentation-generator.md for full template.

    IMPORTANT:
    - Read ALL source files in each directory before writing CLAUDE.md
    - Do NOT copy-paste from other CLAUDE.md files without verifying accuracy
    - Include ONLY functions/exports that actually exist in the code
```

#### Step 1.4: Verify All CLAUDE.md Files Exist

```bash
# Re-run the check - output MUST be empty
for dir in $(git diff --name-only $(git merge-base HEAD main) HEAD | xargs -I {} dirname {} | sort -u | grep -v '^\.$'); do
  if [ -d "$dir" ] && [ ! -f "$dir/CLAUDE.md" ] && [[ "$dir" != .claude* ]] && [[ "$dir" != .spec-workflow* ]] && [[ "$dir" != .github* ]]; then
    echo "STILL MISSING: $dir/CLAUDE.md"
  fi
done
# If ANY output appears, GO BACK TO STEP 1.3. DO NOT PROCEED.
```

---

### Step 2: Epic Status Documentation (MANDATORY)

**Step 2.1: Update epics-user-stories.md (REQUIRED)**
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

**Step 2.2: Update Root CLAUDE.md (REQUIRED)**
```bash
# Edit root CLAUDE.md
# Update project status line with latest epic
# Example: "Epic 4 (CI/CD Pipelines) complete. Epic 5 in progress."
```

**Step 2.3: Update CHANGELOG.md (REQUIRED)**
```bash
# Add comprehensive epic entry with:
# - All deliverables (modules, services, pipelines)
# - Security enhancements (IAM, secrets, etc.)
# - Architecture decisions (versioned paths, canary, etc.)
# - Breaking changes (if any)
```

**Step 2.4: Verify Documentation Consistency (REQUIRED)**
- [ ] All acceptance criteria checked: `- [x]`
- [ ] Epic has completion date: `✅ **COMPLETE** (YYYY-MM-DD)`
- [ ] Implementation notes added to criteria (account IDs, module names, etc.)
- [ ] Root CLAUDE.md project status reflects epic completion
- [ ] CHANGELOG.md has detailed epic entry

**IMPORTANT**: The epics-user-stories.md update is NOT optional. This is how we track epic completion across the entire project. Skipping this step means the epic is NOT considered complete.

---

### Step 3: Planning Document Sync (REQUIRED)

- [ ] Review `planning-docs/*.md` for:
  - Fictional account names → Update to actual account IDs
  - Architecture diagrams → Reflect current state
  - Technology stack references → Current implementation
- [ ] Cross-reference account names with `infrastructure/docs/aws-organizations.md`
- [ ] Mark future accounts explicitly with "future:" prefix
- [ ] Update API documentation (OpenAPI specs if applicable)

---

### Step 4: Automated Validation (MANDATORY)

**REQUIRED**: Run doc-consistency-checker agent before marking epic complete:

```
Task tool parameters:
- description: "Validate documentation consistency"
- subagent_type: "general-purpose"
- prompt: |
    Use the doc-consistency-checker agent instructions to validate
    documentation consistency across the codebase.

    Agent file: .claude/agents/doc-consistency-checker.md

    Focus on:
    1. CLAUDE.md completeness - every directory with code changes has CLAUDE.md
    2. Account name validation (vs infrastructure/docs/aws-organizations.md)
    3. Epic status validation (vs implementation-plan/epics-user-stories.md)
    4. Technology stack validation
    5. Future vs. Current distinction
    6. Module documentation completeness

    Generate report with:
    - Summary statistics
    - Issues by category (CLAUDE.md missing is HIGH severity)
    - Detailed issue list with severity
    - Recommended fixes
```

### Documentation Validation Checklist

- [ ] **ALL changed directories have CLAUDE.md** (verified in Step 1.4)
- [ ] No HIGH severity issues from doc-consistency-checker
- [ ] All account names match `aws-organizations.md`
- [ ] Epic statuses consistent across all docs
- [ ] No outdated technology references
- [ ] Deployed accounts have IDs, future accounts marked
- [ ] Root CLAUDE.md updated
- [ ] CHANGELOG.md updated

### Exit Criteria

**⚠️ DO NOT create PR until ALL of these pass:**

1. **Step 1.4 verification passes** - zero missing CLAUDE.md files
2. Doc-consistency-checker reports no HIGH severity issues
3. All documentation artifacts committed
4. Root CLAUDE.md and CHANGELOG.md updated

**Artifacts**: Updated CHANGELOG.md, root CLAUDE.md, per-directory CLAUDE.md files, planning docs

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
MCPs: spec-workflow, github
```

**For Frontend Test tasks (Phase 2):**
```
@claude write failing tests for this feature following TDD
MCPs: context7 (Vitest, React Testing Library), docs-mcp-server
```

**For Backend Test tasks (Phase 2):**
```
@claude write failing tests for this feature following TDD
MCPs: context7 (Testify, httptest), docs-mcp-server
```

**For Frontend Code tasks (Phase 3):**
```
@claude implement this feature to make tests pass
MCPs: context7 (React, TanStack, Zustand), daisyui-blueprint
```

**For Backend Code tasks (Phase 3):**
```
@claude implement this Go Lambda function to make tests pass
MCPs: context7 (Echo, Go), awslabs.dynamodb-mcp-server, awslabs.aws-documentation-mcp-server
```

**For Infrastructure tasks:**
```
@claude implement this OpenTofu module
MCPs: opentofu, awslabs.aws-documentation-mcp-server, aws-knowledge-mcp-server
```

**For Security Audit tasks:**
```
@claude audit this code for security vulnerabilities
MCPs: awslabs.aws-documentation-mcp-server, aws-knowledge-mcp-server
```

---

## Agent Spawning Reference

| Phase | Agent | subagent_type | MCP Servers (Pre-Agent Context) |
|-------|-------|---------------|-------------------------------|
| 2. Test (FE) | test-engineer | `test-engineer` | `context7`, `docs-mcp-server`, `playwright` |
| 2. Test (BE) | test-engineer | `test-engineer` | `context7`, `docs-mcp-server` |
| 3. Code (FE) | implementation-agent | `implementation-agent` | `context7`, `daisyui-blueprint`, `docs-mcp-server` |
| 3. Code (BE) | implementation-agent | `implementation-agent` | `context7`, `awslabs.dynamodb-mcp-server`, `awslabs.aws-documentation-mcp-server` |
| 3. Code (Infra) | implementation-agent | `implementation-agent` | `opentofu`, `awslabs.aws-documentation-mcp-server`, `aws-knowledge-mcp-server` |
| 5. Docs | doc-consistency-checker | `general-purpose` | `github`, `spec-workflow` |
| Code Review | code-reviewer | `code-reviewer` | `github` |
| Security | security-auditor | `security-auditor` | `awslabs.aws-documentation-mcp-server`, `aws-knowledge-mcp-server` |

**Note**:
- MCP servers are used by the **orchestrator** (main session) to gather context BEFORE spawning agents
- Agents receive MCP-gathered context via their prompt, not direct MCP access
- Phase 5 (Documentation) runs locally ONLY for epic completion, not per-task
- Additional security and documentation checks run in CI/CD after PR creation

---

## Status Tracking

Use TodoWrite throughout:

```
[ ] Phase 1: Specification
    [ ] requirements.md created
    [ ] design.md created
    [ ] tasks.md created (with separate unit test / E2E test / implementation tasks)
[ ] Phase 2: Test (TDD Red) - MANDATORY
    [ ] MCP context gathered (context7, docs-mcp-server, playwright)
    [ ] Spawn test-engineer agent
    [ ] Unit test file exists
    [ ] Unit test was run and FAILED before implementation
    [ ] E2E test file exists (if UI component)
    [ ] E2E test was run and FAILED before implementation
    [ ] ⚠️ STOP: If any test PASSES, it is not testing the right thing
[ ] Phase 3: Code (TDD Green) - MANDATORY
    [ ] MCP context gathered (context7, daisyui-blueprint/dynamodb/opentofu)
    [ ] Verified ALL test files exist before coding
    [ ] Spawn implementation-agent
    [ ] Implementation reads test files FIRST
    [ ] ONLY code needed to pass tests written
    [ ] ALL tests now PASS (Green)
[ ] Phase 4: Verify
    [ ] Tests passing with coverage threshold met
    [ ] **CLAUDE.md created/updated** for ALL changed directories
    [ ] ⚠️ STOP: If any changed directory missing CLAUDE.md
    [ ] **tasks.md updated** (mark task [x] complete)
    [ ] Changes committed to task branch
    [ ] Merged to group branch
    [ ] Task issue closed
[ ] Phase 5: Documentation (Epic Complete Only)
    [ ] Step 1: CLAUDE.md discovery - find ALL changed dirs
    [ ] Step 1: CLAUDE.md verification - zero missing files
    [ ] ⚠️ STOP: DO NOT PROCEED if any CLAUDE.md missing
    [ ] Step 1: Spawn documentation-generator for gaps
    [ ] Step 1: Re-verify - zero missing CLAUDE.md
    [ ] Step 2: epics-user-stories.md updated
    [ ] Step 2: Root CLAUDE.md updated
    [ ] Step 2: CHANGELOG.md updated
    [ ] Step 3: Planning docs synced
    [ ] Step 4: Doc-consistency-checker validation passing
    [ ] Step 4: Zero HIGH severity issues
```

**CRITICAL REMINDERS**:
- The CLAUDE.md update in Phase 4 is NOT optional. Every changed directory MUST have CLAUDE.md before commit.
- The tasks.md update in Phase 4 is NOT optional. It tracks epic progress.
- Phase 5 Step 1 MUST complete with zero missing CLAUDE.md before ANY other Phase 5 step.

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
