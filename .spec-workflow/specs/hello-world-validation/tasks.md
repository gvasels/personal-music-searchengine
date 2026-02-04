# Hello World Validation - Tasks

**Purpose**: Reference template demonstrating proper TDD task structure with explicit test-before-code ordering.

---

## Task 1: HelloPage Component with TDD

- [ ] 1.1 Write unit tests (TDD Red)
  - **Unit Test**: `frontend/src/routes/__tests__/hello.test.tsx` (create)
  - **Acceptance Criteria**:
    - Test verifies "Hello World" heading renders
    - Test verifies clicking button shows alert
    - **MUST RUN AND FAIL** before implementation

- [ ] 1.2 Write E2E tests (TDD Red)
  - **E2E Test**: `frontend/e2e/hello.spec.ts` (create)
  - **Acceptance Criteria**:
    - Test navigates to `/hello` route
    - Test clicks button and verifies alert
    - **MUST RUN AND FAIL** before implementation

- [ ] 1.3 Implement HelloPage (TDD Green)
  - **Implementation**: `frontend/src/routes/hello.tsx` (create)
  - **Acceptance Criteria**:
    - Component renders correctly
    - **ALL TESTS MUST PASS**

---

## Task Template (Copy for New Features)

```markdown
- [ ] N.1 Write unit tests (TDD Red)
  - **Unit Test**: `path/to/test_file` (create)
  - **Acceptance Criteria**:
    - [specific test assertions]
    - **MUST RUN AND FAIL** before implementation

- [ ] N.2 Write E2E tests (TDD Red) [if UI component]
  - **E2E Test**: `frontend/e2e/{feature}.spec.ts` (create)
  - **Acceptance Criteria**:
    - [specific E2E scenarios]
    - **MUST RUN AND FAIL** before implementation

- [ ] N.3 Implement feature (TDD Green)
  - **Implementation**: `path/to/source_file` (create/modify)
  - **Acceptance Criteria**:
    - [implementation requirements]
    - **ALL TESTS MUST PASS**
```

---

## SDLC Violation Log

### Previous Attempt Violations:

#### Violation 1: Unit tests created during BUILD phase
- **Phase**: TEST
- **What happened**: Unit test created after implementation
- **What should have happened**: Unit test MUST be created and run (failing) BEFORE implementation
- **Root cause**: Task only specified E2E tests, not unit tests
- **Prevention**: Tasks MUST explicitly list unit test as separate task before implementation

#### Violation 2: TDD Red not verified
- **Phase**: TEST
- **What happened**: Tests written but not run to verify failure
- **What should have happened**: MUST run tests and see "FAIL" output before writing code
- **Root cause**: Workflow doc said "write tests" but didn't require verification
- **Prevention**: Add explicit "MUST RUN AND FAIL" requirement to each test task

#### Violation 3: Code written before tests
- **Phase**: CODE
- **What happened**: First attempt wrote hello.tsx before any tests
- **What should have happened**: Tests MUST exist and fail before any implementation
- **Root cause**: No enforcement mechanism in workflow
- **Prevention**: Separate test tasks from implementation tasks with explicit ordering
