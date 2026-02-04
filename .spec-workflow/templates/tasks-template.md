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

## Task N: Integration and Documentation

- [ ] N.1 Validate full stack locally
  - **Validation**: `make local` + manual verification
  - **Acceptance Criteria**:
    - Feature works end-to-end on local stack
    - API endpoints respond correctly
    - UI renders and functions as designed

- [ ] N.2 Update CLAUDE.md for all changed directories
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
