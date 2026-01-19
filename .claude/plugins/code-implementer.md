---
name: code-implementer
description: Implementation code to make tests pass
phase: 3-implementation
skills: []
agents: [implementation-agent]
mcp_servers: []
---

# Code Implementer Plugin

Writes production code to make failing tests pass, following the design specifications.

## Phase Position

```
1. SPEC → 2. TEST → [3. CODE] → 4. BUILD → 5. SECURITY → 6. DOCS
                        ▲
                        YOU ARE HERE
```

## Prerequisites

**REQUIRED from previous phases:**
- `design.md` - Data models and API contracts (from spec-writer)
- **Failing tests** - Test suites expecting implementation (from test-writer)

**CRITICAL VERIFICATION**: Before starting implementation, verify tests exist and are failing:

```bash
# Verify test files exist
ls -la tests/unit/*/
ls -la tests/integration/*/

# Verify tests are RED (failing)
npm test        # Should show FAILURES
go test ./...   # Should show FAILURES
```

**If no tests exist or tests are passing, STOP**:
- You are violating TDD workflow
- Return to Phase 2 (test-writer plugin)
- Write failing tests FIRST, then return here

## TDD Implementation Flow

```
┌─────────────────────────────────────────┐
│           TESTS ARE RED                  │
│        (All tests failing)               │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│     Write MINIMAL code to pass ONE      │
│              test at a time             │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│           RUN TESTS                      │
│    Does the targeted test pass?         │
└─────────────────────────────────────────┘
          │                    │
         YES                  NO
          │                    │
          ▼                    ▼
┌─────────────────┐   ┌─────────────────┐
│    REFACTOR     │   │   DEBUG/FIX     │
│  (Keep green)   │   │   (Stay red)    │
└─────────────────┘   └─────────────────┘
          │                    │
          └────────────────────┘
                    │
                    ▼
            Next failing test
```

## Workflow

### Step 1: Identify Failing Tests

```bash
# Run tests to see what needs implementation
npm test -- --reporter=verbose
go test ./... -v
pytest -v
```

### Step 2: Implement Data Models

From `design.md`, create type definitions:

```typescript
// src/types/{feature}.ts
export interface EntityName {
  id: string;
  name: string;
  status: 'active' | 'inactive';
  createdAt: Date;
  updatedAt: Date;
}

export const EntityNameSchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(3).max(50),
  status: z.enum(['active', 'inactive']),
  createdAt: z.date(),
  updatedAt: z.date(),
});
```

### Step 3: Implement API Endpoints

From API contracts in `design.md`:

```typescript
// src/routes/{feature}.ts
import { Hono } from 'hono';
import { zValidator } from '@hono/zod-validator';

const app = new Hono();

app.post(
  '/api/resource',
  zValidator('json', CreateResourceSchema),
  async (c) => {
    const input = c.req.valid('json');

    // Check for duplicates (409 case)
    const existing = await db.findByName(input.name);
    if (existing) {
      return c.json({ error: 'Resource already exists' }, 409);
    }

    // Create resource
    const resource = await db.create(input);
    return c.json(resource, 201);
  }
);
```

### Step 4: Implement Business Logic

```typescript
// src/services/{feature}.ts
export class ResourceService {
  constructor(private db: Database) {}

  async create(input: CreateResourceInput): Promise<Resource> {
    // Validation handled by schema
    const resource: Resource = {
      id: generateUUID(),
      ...input,
      status: 'active',
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    await this.db.insert(resource);
    return resource;
  }
}
```

### Step 5: Run Tests After Each Change

```bash
# After implementing each piece, verify tests pass
npm test -- --watch

# Check coverage
npm test -- --coverage
```

## Code Style Guidelines

| Language | Style Guide | Formatter |
|----------|-------------|-----------|
| TypeScript/React | Airbnb React | Prettier |
| Go | Google Go Style | gofmt |
| Python | PEP 8 | Black |

## Implementation Principles

1. **Minimal Code** - Only write what's needed to pass the test
2. **No Premature Optimization** - Make it work, then make it right
3. **Follow the Contract** - Match API specs exactly
4. **Refactor After Green** - Clean up only when tests pass

## Subagent Delegation

Spawn `implementation-agent` for complex features:

```
Use Task tool with subagent_type='implementation-agent'
Provide:
- Failing test file location
- Design specs to implement
- Code style requirements
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| Type definitions | `src/types/` | Data models |
| API routes | `src/routes/` | Endpoint handlers |
| Services | `src/services/` | Business logic |
| Repositories | `src/repositories/` | Data access |

## Handoff to Next Phase

**Checklist before proceeding to builder/verify:**
- [ ] All tests pass (Green phase - changed from Red!)
- [ ] Test coverage meets threshold (80%+ for critical code)
- [ ] Code follows design specifications exactly
- [ ] No hardcoded values or shortcuts
- [ ] No over-engineering (only code needed to pass tests)

**CRITICAL**: Update tasks.md IMMEDIATELY:
```bash
# Edit .spec-workflow/specs/{spec-name}/tasks.md
# Change: - [ ] Task N.X: Description
# To:     - [x] Task N.X: Description
```

**VERIFICATION**: Run tests to confirm GREEN:
```bash
npm test -- --coverage   # Should show PASSING + coverage
go test ./... -v         # Should show PASSING
```

**NEXT**: Pass to `builder` plugin for build verification, then close task
