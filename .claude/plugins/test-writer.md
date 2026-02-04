---
name: test-writer
description: TDD test creation before implementation
phase: 2-testing
skills: []
agents: [test-engineer]
mcp_servers: []
---

# Test Writer Plugin

Creates comprehensive test suites BEFORE implementation code, following Test-Driven Development principles.

## Phase Position

```
1. SPEC → [2. TEST] → 3. CODE → 4. BUILD → 5. SECURITY → 6. DOCS
              ▲
              YOU ARE HERE
```

## Prerequisites

From `spec-writer` plugin:
- `requirements.md` - User stories and acceptance criteria
- `design.md` - Data models and API contracts
- `tasks.md` - Implementation breakdown

## TDD Workflow

### The Three Laws

1. **Write a failing test** before writing any production code
2. **Write only enough test** to demonstrate a failure
3. **Write only enough code** to make the test pass

### Test Categories

| Type | Purpose | Coverage Target |
|------|---------|-----------------|
| Unit | Individual functions | 80%+ |
| Integration | API endpoints, DB operations | Key paths |
| Contract | API response validation | All endpoints |
| E2E | User workflows | Critical flows |

## Workflow

### Step 1: Analyze Design Artifacts

```
1. Read data models from design.md
2. Extract API contracts (request/response schemas)
3. Identify validation rules and constraints
4. List error scenarios from acceptance criteria
```

### Step 2: Create Test Structure

```
tests/
├── unit/
│   └── {feature}/
│       ├── {entity}.test.ts      # Entity logic tests
│       └── {service}.test.ts     # Service layer tests
├── integration/
│   └── {feature}/
│       └── api.test.ts           # API endpoint tests
└── e2e/
    └── {feature}/
        └── workflow.test.ts      # User journey tests
```

### Step 3: Write Tests from Contracts

For each API endpoint in design.md:

```typescript
describe('POST /api/{resource}', () => {
  describe('valid input', () => {
    it('creates resource and returns 201', async () => {
      const input = createValidInput();
      const response = await api.post('/api/resource', input);

      expect(response.status).toBe(201);
      expect(response.body).toMatchSchema(ResourceSchema);
      expect(response.body.field).toBe(input.field);
    });
  });

  describe('validation', () => {
    it('returns 400 when required field missing', async () => {
      const input = { /* missing required field */ };
      const response = await api.post('/api/resource', input);

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('field');
    });

    it('returns 400 when field exceeds max length', async () => {
      const input = { field: 'x'.repeat(51) }; // max is 50
      const response = await api.post('/api/resource', input);

      expect(response.status).toBe(400);
    });
  });

  describe('business rules', () => {
    it('returns 409 when duplicate exists', async () => {
      await createResource({ name: 'existing' });
      const response = await api.post('/api/resource', { name: 'existing' });

      expect(response.status).toBe(409);
    });
  });
});
```

### Step 4: Write Unit Tests from Data Models

For each entity in design.md:

```typescript
describe('EntityName', () => {
  describe('validation', () => {
    it('accepts valid entity', () => {
      const entity = createValidEntity();
      expect(validateEntity(entity)).toBe(true);
    });

    it('rejects invalid status', () => {
      const entity = createValidEntity({ status: 'invalid' });
      expect(() => validateEntity(entity)).toThrow('Invalid status');
    });
  });

  describe('transformations', () => {
    it('converts to API response format', () => {
      const entity = createValidEntity();
      const response = toApiResponse(entity);

      expect(response).not.toHaveProperty('internalField');
      expect(response.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}/);
    });
  });
});
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| Unit tests | `tests/unit/{feature}/` | Function-level validation |
| Integration tests | `tests/integration/{feature}/` | API contract validation |
| Test fixtures | `tests/fixtures/{feature}/` | Reusable test data |
| Schema validators | `tests/schemas/{feature}/` | JSON Schema for responses |

## Test Frameworks

| Language | Framework | Assertion | Mocking |
|----------|-----------|-----------|---------|
| TypeScript | Vitest | expect | vi.mock |
| Go | Testify | assert/require | testify/mock |
| Python | pytest | assert | unittest.mock |

## Subagent Delegation

Spawn `test-engineer` agent for complex test scenarios:

```
Use Task tool with subagent_type='test-engineer'
Provide:
- Data models from design.md
- API contracts to test against
- Specific test category needed
```

## Handoff to Next Phase

**Checklist before proceeding to code-implementer:**
- [ ] All tests written and FAILING (Red phase - this is expected!)
- [ ] Test coverage targets defined (80%+ for critical code)
- [ ] Test fixtures created for all scenarios
- [ ] Tests validate ALL acceptance criteria from requirements.md
- [ ] Tests cover ALL API contracts from design.md

**VERIFICATION**: Run tests to confirm they fail:
```bash
npm test        # Should show FAILING tests
go test ./...   # Should show FAILING tests
```

If tests PASS, you wrote production code already - this violates TDD. Start over.

**NEXT**: Pass to `code-implementer` plugin to make tests pass
