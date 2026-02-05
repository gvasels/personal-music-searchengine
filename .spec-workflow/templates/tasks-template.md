# Tasks Document

**Spec**: [spec-name]
**Design**: [design.md](./design.md)
**Requirements**: [requirements.md](./requirements.md)

---

## CRITICAL: TDD Task Structure

**Every feature MUST follow this ordering:**
1. **Test tasks FIRST** (TDD Red) - Write failing tests
2. **Implementation tasks SECOND** (TDD Green) - Write code to make tests pass

**DO NOT** put implementation tasks before their corresponding test tasks.
**DO NOT** combine test and implementation in a single task.

---

## Task 1: [Feature Group Name]

- [ ] 1.1 Write unit tests (TDD Red)
  - **Unit Test**: `[path/to/test_file]` (create)
  - **Acceptance Criteria**:
    - [Test assertion 1]
    - [Test assertion 2]
    - **MUST RUN AND FAIL** before implementation
  - _Leverage: [existing test utils, fixtures]_
  - _Requirements: [requirement IDs]_

- [ ] 1.2 Write E2E tests (TDD Red) [if UI component]
  - **E2E Test**: `frontend/e2e/[feature].spec.ts` (create)
  - **Acceptance Criteria**:
    - [E2E scenario 1]
    - [E2E scenario 2]
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: [requirement IDs]_

- [ ] 1.3 Implement feature (TDD Green)
  - **Implementation**: `[path/to/source_file]` (create/modify)
  - **Acceptance Criteria**:
    - [Implementation requirement 1]
    - [Implementation requirement 2]
    - **ALL TESTS MUST PASS**
  - _Leverage: [existing code to reuse]_
  - _Requirements: [requirement IDs]_

---

## Task 2: [Next Feature Group]

- [ ] 2.1 Write unit tests (TDD Red)
  - **Unit Test**: `[path/to/test_file]` (create)
  - **Acceptance Criteria**:
    - [Test assertions]
    - **MUST RUN AND FAIL** before implementation
  - _Requirements: [requirement IDs]_

- [ ] 2.2 Implement feature (TDD Green)
  - **Implementation**: `[path/to/source_file]` (create/modify)
  - **Acceptance Criteria**:
    - [Implementation requirements]
    - **ALL TESTS MUST PASS**
  - _Requirements: [requirement IDs]_

---

## Task N: Integration, Validation, and Documentation

- [ ] N.1 Contract alignment check (REQUIRED for specs with both FE and BE)
  - **Purpose**: Verify frontend and backend use the same API contract
  - **Acceptance Criteria**:
    - Frontend API client sends the correct query param names (verify with grep)
    - Frontend expects the correct response field names (verify with grep)
    - Backend handler reads the correct query param names (verify with grep)
    - Backend returns the correct JSON field names (verify with grep)
  - **Why this exists**: Mocked unit tests on FE and BE pass independently even when they disagree on param/field names. This was discovered in iteration 2 of hello-world validation (FE used `query`/`tracks`, BE used `q`/`items`).
  - **Skip condition**: Only skip if spec has no frontend OR no backend changes.

- [ ] N.2 Local stack validation (`make local`)
  - **Validation**: `make local` + curl + browser verification
  - **Acceptance Criteria**:
    - Feature works end-to-end on local stack
    - API endpoints respond correctly (curl verification)
    - UI renders and functions as designed
  - **Skip condition**: Only skip for infrastructure-only, test-only, or docs-only changes.

- [ ] N.3 BUILD validation
  - **Validation**: All build commands pass with zero errors
  - **Acceptance Criteria**:
    - Backend: `go vet ./...`, `go build ./cmd/api`, `go test ./...`
    - Frontend: `tsc --noEmit`, `eslint .`, `npm run build`, `npm test -- --run`

- [ ] N.4 Update CLAUDE.md for all changed directories
  - **Documentation**: Create/update CLAUDE.md files
  - **Directories requiring CLAUDE.md**: [list after implementation]
  - **Acceptance Criteria**:
    - Each CLAUDE.md has: Overview, Files, Key Functions, Dependencies
    - Zero missing CLAUDE.md for changed directories

---

## Task Structure Rules

| Rule | Requirement |
|------|-------------|
| **Test before code** | Unit test tasks (N.1) MUST come before implementation tasks (N.3) |
| **Separate test tasks** | Unit tests and E2E tests are separate tasks from implementation |
| **MUST RUN AND FAIL** | Every test task must include this requirement |
| **ALL TESTS MUST PASS** | Every implementation task must include this requirement |
| **Leverage existing code** | Each task should note what existing code it reuses |
| **Requirement traceability** | Each task links to requirement IDs |
| **CLAUDE.md final task** | Last task group must include CLAUDE.md verification |

---

## SDLC Gap Tracking

[Document any gaps found during this spec's lifecycle. Include: what was expected, what actually happened, and proposed fix.]

| # | Gap | Where | Impact | Proposed Fix |
|---|-----|-------|--------|--------------|
| - | - | - | - | - |
