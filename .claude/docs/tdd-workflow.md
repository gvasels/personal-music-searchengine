# Test-Driven Development (TDD) Workflow

This document provides detailed guidance on TDD practices for the Personal Music Search Engine.

## Tasks Template

**All `tasks.md` files MUST follow the TDD-structured template** at `.spec-workflow/templates/tasks-template.md`. Key rules:
- Test tasks (N.1, N.2) MUST come BEFORE implementation tasks (N.3)
- Each test task MUST include **MUST RUN AND FAIL** requirement
- Each implementation task MUST include **ALL TESTS MUST PASS** requirement
- Unit tests and E2E tests are SEPARATE tasks from implementation
- Last task group MUST include CLAUDE.md verification and `make local` validation

## CRITICAL TDD RULES

**You MUST write ALL tests BEFORE ANY implementation code.**

### Phase 2: TEST (TDD Red) - MANDATORY STEPS

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

### Phase 3: CODE (TDD Green) - MANDATORY STEPS

**You MUST write ONLY the minimal code to make tests pass.**

#### Step 3.1: Implement Feature
- Read test file to understand exact requirements
- Write ONLY code needed to pass tests
- NO extra features, NO premature optimization

#### Step 3.2: VERIFY TESTS PASS (GREEN)
- Run ALL tests (unit + E2E)
- ALL tests MUST pass before marking phase complete

---

## ENFORCEMENT CHECKLIST

Before marking ANY task complete, MUST verify:
- [ ] Unit test file exists
- [ ] Unit test was run and FAILED before implementation
- [ ] Unit test now PASSES
- [ ] E2E test file exists (if UI)
- [ ] E2E test was run and FAILED before implementation
- [ ] E2E test now PASSES
- [ ] TypeScript check passes (frontend) / go vet passes (backend)
- [ ] CLAUDE.md updated
- [ ] Git commit created

---

## TDD Anti-Patterns (NEVER DO THESE)

These anti-patterns were identified from actual enforcement failures. Each one bypasses the Red-Green cycle.

| # | Anti-Pattern | What Happens | Detection |
|---|---|---|---|
| 1 | **Combined agent prompt** — telling one agent to "write tests and implementation" | Red phase is skipped entirely; agent writes tests that pass immediately | Check: were `test-engineer` and `implementation-agent` spawned as separate Task calls? |
| 2 | **Single commit with tests + implementation** | No proof that tests ever failed | Check: does git log show a Red commit (tests only) before a Green commit (implementation only)? |
| 3 | **Skipping Red phase on session resume** | Context loss leads to "just get it done" mentality | Check: did the session follow the Session Recovery Protocol before writing code? |
| 4 | **Group-level task tracking** | Tracking Tasks 1-7 instead of sub-tasks 1.1, 1.2, 1.3 hides whether Red preceded Green | Check: are TodoWrite tasks at the sub-task level (N.1 test, N.2 e2e, N.3 implementation)? |
| 5 | **Background agents with combined prompts** | Parallel agents each doing test+code bypass sequential Red→Green requirement | Check: were background agents given test-only or implementation-only prompts? |

## Separate Commits Required

TDD enforcement requires **separate git commits** for the Red and Green phases:

### Red Phase Commit (Phase 2)
```
test(task-N.X): Red phase - failing tests for [feature]

- N tests written, all FAILING as expected
- Test failures: [list failing test names]
- No implementation code included
```

### Green Phase Commit (Phase 3)
```
feat(task-N.X): Green phase - implementation for [feature]

- All N tests now PASSING
- Minimal implementation to satisfy tests
```

**Why separate commits matter**: A single commit containing both tests and implementation provides zero evidence that the Red phase occurred. The git history IS the audit trail.

## Session Recovery Protocol

When resuming after context loss (new session, context window exceeded, etc.):

1. **Read `.spec-workflow/specs/{spec-name}/tasks.md`** — determine current sub-task status
2. **Run `git log --oneline -10`** — check if Red phase commits exist
3. **Determine recovery point**:

| If you find... | Then... |
|---|---|
| No test files on disk | Start Phase 2 from scratch |
| Test files exist but no Red commit | Run tests, verify they fail, create Red commit, THEN proceed to Phase 3 |
| Red commit exists but no Green commit | Start Phase 3 (spawn `implementation-agent`) |
| Both Red and Green commits exist | Proceed to Phase 4 (Verify) |

4. **NEVER assume previous session completed TDD correctly** — verify via git log

## Verification Evidence

Every Red phase MUST produce verifiable evidence that tests failed. This evidence MUST be included in the Red phase commit message.

**Acceptable evidence formats:**
- Test runner output showing failure count: `FAIL: 5 tests failed`
- List of failing test names: `TestSearchTracks, TestGetFeaturedTracks, ...`
- Exit code from test runner: `exit status 1`

**NOT acceptable:**
- "Tests were written" (no proof they failed)
- "Tests should fail" (assumption, not verification)
- No mention of test results at all

---

## TDD Cycle

```
Phase 1 (Spec)                    Phase 2 (Test)      Phase 3 (Code)
─────────────────────────────────────────────────────────────────────
Requirements → Data Models → API Contracts → Tests → Implementation → Refactor
                                              │            │
                                              ▼            ▼
                                           RED ────────► GREEN ────► REFACTOR
                                        (failing)     (passing)    (clean up)
```

## The Three Laws of TDD

1. **You MUST write a failing test** before writing any production code
2. **You MUST write only enough test** to demonstrate a failure
3. **You MUST write only enough production code** to make the test pass

## Phase 1 (Spec) Deliverables

Before writing any tests or code, define these during the Spec phase:

### 1. User Stories with Acceptance Criteria

```gherkin
Feature: Service Registration
  As a platform admin
  I want to register new services
  So that they can be deployed to AWS accounts

  Scenario: Register a valid service
    Given I have valid service metadata
    When I submit the registration request
    Then the service is created
    And I receive a service ID
```

### 2. Data Models

Define schemas before implementation:

```typescript
// TypeScript example
interface Service {
  id: string;                    // UUID v4
  name: string;                  // 3-50 chars, alphanumeric + hyphens
  productId: string;             // References Product.id
  owner: string;                 // Team email
  status: 'active' | 'deprecated' | 'archived';
  createdAt: ISO8601Timestamp;
  updatedAt: ISO8601Timestamp;
}
```

```go
// Go example
type Service struct {
    ID        string    `json:"id" dynamodbav:"PK"`
    Name      string    `json:"name" dynamodbav:"name"`
    ProductID string    `json:"productId" dynamodbav:"productId"`
    Owner     string    `json:"owner" dynamodbav:"owner"`
    Status    string    `json:"status" dynamodbav:"status"`
    CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
```

### 3. API Contracts

Define inputs/outputs before implementation:

```yaml
POST /api/services:
  request:
    body:
      name: string (required, 3-50 chars)
      productId: string (required, valid UUID)
      owner: string (required, valid email)
  response:
    201:
      body: Service
    400:
      body: { error: string, details: ValidationError[] }
    409:
      body: { error: "Service with this name already exists" }
```

### 4. Database Schema

- Table/collection definitions
- Indexes and access patterns
- Relationships and constraints

---

## Phase 2 (Test) Requirements

Write tests **before** implementation based on specifications.

### MCP Servers for Test Phase

**The orchestrator MUST use these before spawning test-engineer:**

| MCP Server | Frontend Tests | Backend Tests | E2E Tests |
|------------|---------------|---------------|-----------|
| `context7` | Vitest, React Testing Library APIs | Testify, httptest APIs | Playwright APIs |
| `docs-mcp-server` | Advanced framework patterns | Advanced framework patterns | - |
| `playwright` | - | - | Browser automation, selectors |

### Test Types

#### Unit Tests
- Test individual functions with mocked dependencies
- Cover happy paths, edge cases, and error scenarios
- Validate data transformations and business logic

#### Integration Tests
- Test API endpoints against real (local) databases
- Validate request/response contracts
- Test error handling and status codes

#### Contract Tests
- Validate API responses match defined schemas
- Test backward compatibility when modifying APIs

### Example Test Structure

**TypeScript/Jest:**

```typescript
describe('POST /api/services', () => {
  it('creates a service with valid input', async () => {
    const input = {
      name: 'my-service',
      productId: 'uuid',
      owner: 'team@example.com'
    };
    const response = await api.post('/api/services', input);

    expect(response.status).toBe(201);
    expect(response.body).toMatchSchema(ServiceSchema);
    expect(response.body.name).toBe('my-service');
  });

  it('returns 400 for invalid name', async () => {
    const input = {
      name: 'x', // Too short
      productId: 'uuid',
      owner: 'team@example.com'
    };
    const response = await api.post('/api/services', input);

    expect(response.status).toBe(400);
    expect(response.body.error).toContain('name');
  });

  it('returns 409 for duplicate name', async () => {
    // Create first service
    await api.post('/api/services', {
      name: 'duplicate-service',
      productId: 'uuid',
      owner: 'team@example.com'
    });

    // Try to create duplicate
    const response = await api.post('/api/services', {
      name: 'duplicate-service',
      productId: 'uuid',
      owner: 'team@example.com'
    });

    expect(response.status).toBe(409);
    expect(response.body.error).toContain('already exists');
  });
});
```

**Go:**

```go
func TestCreateService_Success(t *testing.T) {
    // Arrange
    repo := mocks.NewMockServiceRepository()
    handler := NewServiceHandler(repo)

    input := CreateServiceRequest{
        Name:      "my-service",
        ProductID: "valid-uuid",
        Owner:     "team@example.com",
    }

    // Act
    result, err := handler.Create(context.Background(), input)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    assert.Equal(t, "my-service", result.Name)
    assert.Equal(t, "active", result.Status)
}

func TestCreateService_InvalidName(t *testing.T) {
    handler := NewServiceHandler(nil)

    input := CreateServiceRequest{
        Name:      "x", // Too short
        ProductID: "valid-uuid",
        Owner:     "team@example.com",
    }

    _, err := handler.Create(context.Background(), input)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name")
}
```

---

## Phase 3 (Code) Implementation

Only after tests are written and verified FAILING:

### MCP Servers for Code Phase

**The orchestrator MUST use these before spawning implementation-agent:**

| MCP Server | Frontend | Backend | Infrastructure |
|------------|----------|---------|----------------|
| `context7` | React, TanStack, Zustand | Echo, Go AWS SDK | - |
| `daisyui-blueprint` | DaisyUI component snippets | - | - |
| `awslabs.dynamodb-mcp-server` | - | DynamoDB access patterns | - |
| `awslabs.aws-documentation-mcp-server` | - | AWS service docs | AWS service docs |
| `opentofu` | - | - | Module registry |
| `aws-knowledge-mcp-server` | - | - | AWS recommendations |

### Steps

1. **Run tests** - they MUST fail ("Red")
2. **Write minimal code** to pass tests ("Green")
3. **Refactor** while keeping tests passing
4. **Repeat** for next feature/behavior

### Implementation Guidelines

- Write the **simplest code** that passes the test
- NEVER add features not covered by tests
- Refactor after each passing test
- Keep functions small and focused
- MUST read test files first to understand exact requirements

---

## Test Coverage Requirements

| Type | Minimum Coverage | Focus Areas |
|------|------------------|-------------|
| Unit Tests | 80% | Business logic, data transformations |
| Integration Tests | Key paths | API contracts, database operations |
| Contract Tests | All FE↔BE boundaries | Param names, response field names |
| Wiring/Smoke Tests | All routes | Endpoint accessibility verification |
| E2E Tests | Critical flows | User journeys, cross-service interactions |

## Contract Test Limitations of Mocked Unit Tests

**⚠️ Mocked unit tests CANNOT validate frontend↔backend API contracts.**

When frontend and backend are tested independently with mocks, each side defines its own mock that may not match the other side's actual behavior. This was discovered in iteration 2 of hello-world validation:

| Side | What it mocked | What it expected | Actual |
|------|---------------|------------------|--------|
| Frontend | `axios.get('/v1/hello/search', { params: { query: 'jazz' } })` | `response.data.tracks` | Backend reads `q`, returns `items` |
| Backend | `MockHelloService.Search(ctx, "jazz")` | Returns `HelloTrack[]` | Frontend sends `query` not `q` |

Both sides' tests passed independently. The mismatch was only caught by `make local` curl testing.

### How to Catch Contract Mismatches

1. **Contract alignment check** (REQUIRED in tasks.md template — Task N.1): Grep both codebases for param/field names
2. **`make local` validation** (REQUIRED in Phase 4): End-to-end curl tests against real local stack
3. **Shared contract types** (recommended): Define API types in a shared location or OpenAPI spec that both FE and BE reference
4. **Integration tests** (when available): Tests that hit actual endpoints catch mismatches that mocks hide

---

## Wiring Verification Tests (CRITICAL)

**Why This Matters**: Unit tests with mocks can pass even when:
- Service isn't added to Services struct
- Handler isn't created in main.go
- Routes aren't registered
- Environment variables are missing from Lambda
- Frontend routes aren't registered in router

Only integration/smoke tests that hit actual endpoints catch these wiring issues.

### Backend Wiring Test Pattern

For every new route, add a smoke test that verifies the route is accessible:

```go
// integration_test.go
func TestAdminRoutes_Accessible(t *testing.T) {
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/admin/users?search=test", nil)
    req.Header.Set("Authorization", "Bearer "+validToken)

    router.ServeHTTP(resp, req)

    // Should NOT be 404 - route must be registered
    assert.NotEqual(t, 404, resp.Code, "Route not registered!")
}
```

### Frontend Route Test Pattern

For every new page, add a test that verifies the route renders:

```typescript
// routes.test.tsx
describe('Admin Routes', () => {
  it('renders admin users page at /admin/users', async () => {
    render(<Router initialEntries={['/admin/users']} />);

    // Should NOT show 404/NotFound
    expect(screen.queryByText('Page not found')).not.toBeInTheDocument();
    // Should show expected page content
    expect(await screen.findByText('User Management')).toBeInTheDocument();
  });
});
```

### Wiring Checklist

**ALWAYS follow the wiring checklist** when adding new features:
- See: `.claude/docs/wiring-checklist.md`

The checklist covers:
- New Service Checklist
- New Handler Checklist
- New Route Checklist
- New Frontend Route Checklist
- Verification Steps

---

## Testing Tools

| Layer | Tool | Purpose |
|-------|------|---------|
| Unit (JS/TS) | Jest / Vitest | Fast isolated tests |
| Unit (Go) | go test | Standard Go testing |
| API | Supertest / httptest | HTTP endpoint testing |
| Contract | JSON Schema | Response validation |
| E2E | Playwright | Browser automation |
| Infrastructure | Terratest | OpenTofu module testing |

---

## Test Organization

### Directory Structure (Go)

```
services/
└── account-vending/
    └── internal/
        ├── handlers/
        │   ├── handlers.go
        │   └── handlers_test.go    # Tests alongside code
        └── testutil/
            └── fixtures.go         # Shared test fixtures
```

### Directory Structure (TypeScript)

```
src/
├── services/
│   └── account.service.ts
└── __tests__/
    ├── services/
    │   └── account.service.test.ts
    └── fixtures/
        └── accounts.ts
```

---

## Common Patterns

### Table-Driven Tests (Go)

```go
func TestValidateName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid name", "my-service", false},
        {"too short", "ab", true},
        {"too long", strings.Repeat("a", 51), true},
        {"invalid chars", "my_service!", true},
        {"empty", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateName(%q) error = %v, wantErr %v",
                    tt.input, err, tt.wantErr)
            }
        })
    }
}
```

### Mocking Dependencies

```go
// Interface for mocking
type ServiceRepository interface {
    Create(ctx context.Context, svc *Service) error
    GetByName(ctx context.Context, name string) (*Service, error)
}

// Mock implementation
type MockServiceRepository struct {
    mock.Mock
}

func (m *MockServiceRepository) Create(ctx context.Context, svc *Service) error {
    args := m.Called(ctx, svc)
    return args.Error(0)
}
```

---

## Quick Reference

| Step | Action | Tool | MCP Servers |
|------|--------|------|-------------|
| 1 | Define requirements | spec-workflow MCP | `spec-workflow`, `github` |
| 2 | Define data models | design.md | `awslabs.dynamodb-mcp-server` |
| 3 | Define API contracts | design.md / OpenAPI | `awslabs.aws-documentation-mcp-server` |
| 4 | Write failing tests | test-engineer agent | `context7`, `docs-mcp-server`, `playwright` |
| 5 | Implement code | implementation-agent | `context7`, `daisyui-blueprint`, `opentofu` |
| 6 | Verify tests pass | `go test ./...` / `npm test` | `playwright` (E2E) |
| 7 | Refactor | Keep tests green | - |
