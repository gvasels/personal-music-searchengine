# Test-Driven Development (TDD) Workflow

This document provides detailed guidance on TDD practices for the oopo platform.

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

1. **Write a failing test** before writing any production code
2. **Write only enough test** to demonstrate a failure
3. **Write only enough production code** to make the test pass

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

Only after tests are written:

1. **Run tests** - they should fail ("Red")
2. **Write minimal code** to pass tests ("Green")
3. **Refactor** while keeping tests passing
4. **Repeat** for next feature/behavior

### Implementation Guidelines

- Write the **simplest code** that passes the test
- Don't add features not covered by tests
- Refactor after each passing test
- Keep functions small and focused

---

## Test Coverage Requirements

| Type | Minimum Coverage | Focus Areas |
|------|------------------|-------------|
| Unit Tests | 80% | Business logic, data transformations |
| Integration Tests | Key paths | API contracts, database operations |
| E2E Tests | Critical flows | User journeys, cross-service interactions |

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

| Step | Action | Tool |
|------|--------|------|
| 1 | Define requirements | spec-workflow MCP |
| 2 | Define data models | design.md |
| 3 | Define API contracts | design.md / OpenAPI |
| 4 | Write failing tests | test-engineer agent |
| 5 | Implement code | implementation-agent |
| 6 | Verify tests pass | `go test ./...` |
| 7 | Refactor | Keep tests green |
