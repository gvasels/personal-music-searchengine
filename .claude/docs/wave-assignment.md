# Wave Assignment Decision Tree

Use this decision tree to assign tasks to the correct wave.

## Quick Reference

| Wave | Purpose | Typical Content |
|------|---------|-----------------|
| 0 | Foundation (local) | Contracts, types, interfaces, test fixtures |
| 1 | Independent infra | Databases, storage, base IAM, standalone resources |
| 2 | Dependent infra | API Gateway, compute that needs Wave 1 outputs |
| 3 | Core implementation | Business logic, handlers, services |
| 4 | Integration | Cross-service features, integration tests |
| 5 | Polish | Documentation, final tests, cleanup |

## Decision Tree

```
START: Does this task have ANY dependencies?
│
├─► NO: Does it create shared resources (DB, queues, etc.)?
│   │
│   ├─► YES ──────────────────────────────────► WAVE 1
│   │
│   └─► NO: Is it defining types/interfaces only?
│       │
│       ├─► YES ──────────────────────────────► WAVE 0 (local)
│       │
│       └─► NO: Is it standalone documentation?
│           │
│           ├─► YES ──────────────────────────► WAVE 5
│           │
│           └─► NO ───────────────────────────► WAVE 1
│
└─► YES: What does it depend on?
    │
    ├─► Wave 1 infrastructure only?
    │   │
    │   └─► Is it infrastructure (API GW, compute)?
    │       │
    │       ├─► YES ──────────────────────────► WAVE 2
    │       │
    │       └─► NO: Is it business logic/handlers?
    │           │
    │           └─► ──────────────────────────► WAVE 3
    │
    ├─► Wave 2 outputs?
    │   │
    │   └─► Is it integration/orchestration?
    │       │
    │       ├─► YES ──────────────────────────► WAVE 3 or 4
    │       │
    │       └─► NO ───────────────────────────► WAVE 3
    │
    ├─► Wave 3 outputs?
    │   │
    │   └─► Is it integration testing?
    │       │
    │       ├─► YES ──────────────────────────► WAVE 4
    │       │
    │       └─► NO: Is it documentation?
    │           │
    │           ├─► YES ──────────────────────► WAVE 5
    │           │
    │           └─► NO ───────────────────────► WAVE 4
    │
    └─► Multiple waves?
        │
        └─► Assign to wave AFTER the highest dependency
```

## Wave Assignment Examples

### Infrastructure-Heavy Epic

```
Wave 0 (Local):
  - Define DynamoDB access patterns
  - Define API request/response types

Wave 1:
  - DynamoDB table                    [no deps]
  - S3 bucket for artifacts           [no deps]
  - Base IAM execution role           [no deps]

Wave 2:
  - Lambda function                   [needs: IAM role]
  - API Gateway                       [needs: Lambda ARN]
  - SQS queue                         [no deps, but pairs with Lambda]

Wave 3:
  - Request handlers                  [needs: Lambda, DynamoDB]
  - Event processors                  [needs: Lambda, SQS]
  - Step Functions workflow           [needs: Lambda ARNs]

Wave 4:
  - Integration tests                 [needs: all services deployed]
  - E2E workflow tests                [needs: Step Functions]

Wave 5:
  - API documentation                 [needs: handlers finalized]
  - CLAUDE.md files                   [needs: structure finalized]
```

### Service-Heavy Epic

```
Wave 0 (Local):
  - Define domain models
  - Define service interfaces
  - Create test fixtures

Wave 1:
  - DynamoDB tables                   [no deps]

Wave 2:
  - Repository layer                  [needs: DynamoDB]
  - External API clients              [no deps, but uses types]

Wave 3:
  - Domain services                   [needs: repositories]
  - Business rule validators          [needs: domain models]

Wave 4:
  - HTTP handlers                     [needs: services]
  - GraphQL resolvers                 [needs: services]

Wave 5:
  - Integration tests                 [needs: handlers]
  - Documentation                     [needs: API finalized]
```

## Common Mistakes

### Mistake 1: Everything in Wave 1
```
❌ Wave 1:
   - DynamoDB table
   - Lambda function
   - API Gateway
   - Handlers
   - Tests
```
**Problem:** Lambda needs IAM, API GW needs Lambda, handlers need DynamoDB.

**Fix:** Spread across Waves 1-3 based on dependencies.

### Mistake 2: Over-Sequentialization
```
❌ Wave 1: DynamoDB table
❌ Wave 2: IAM role
❌ Wave 3: Lambda function
❌ Wave 4: API Gateway
❌ Wave 5: Handlers
```
**Problem:** DynamoDB and IAM have no dependency - they can parallel.

**Fix:** DynamoDB + IAM in Wave 1, Lambda + API GW in Wave 2.

### Mistake 3: Ignoring Integration Points
```
❌ Wave 1: Service A (all components)
❌ Wave 2: Service B (all components)
❌ Wave 3: Integration between A and B
```
**Problem:** Services should be built in parallel at each layer.

**Fix:**
```
✅ Wave 1: Infrastructure for both A and B
✅ Wave 2: Core logic for both A and B
✅ Wave 3: Integration between A and B
```

## Wave Capacity Planning

| Wave | Recommended Parallel Tasks | Reason |
|------|---------------------------|--------|
| 0 | N/A (local) | Done by human |
| 1 | 3-5 | Simple, independent infrastructure |
| 2 | 2-4 | Some dependency complexity |
| 3 | 3-5 | Most parallelizable business logic |
| 4 | 2-3 | Integration reduces parallelism |
| 5 | 2-4 | Documentation can parallel |

## Validation Checklist

Before finalizing wave assignments:

- [ ] No task depends on another task in the SAME wave
- [ ] Every task in Wave N depends on at least one task in Wave N-1 (except Wave 1)
- [ ] Wave 0 has NO deployable infrastructure
- [ ] Wave 1 has NO cross-task dependencies
- [ ] Documentation tasks are in Wave 4 or 5
- [ ] Integration tests are after the code they test
- [ ] Total waves ≤ 5 (merge if more)

## Dependency Notation

When documenting dependencies in GitHub issues:

```markdown
## Dependencies
- #42 (DynamoDB table from Wave 1)
- #43 (IAM role from Wave 1)
```

The automation will:
1. Add `blocked` label if dependencies are open
2. Add `ready` label when all dependencies close
3. Trigger @claude when `ready` label is added
