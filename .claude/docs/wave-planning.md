# Wave Planning Guide

This guide explains how to plan waves for parallel @claude execution on feature epics.

## Overview

Wave-based execution enables multiple @claude instances to work on different tasks simultaneously. Tasks within the same wave have no dependencies on each other and can run in parallel. Tasks in later waves depend on earlier waves being complete.

## Wave Planning Process

### Step 1: List All Tasks

Start by breaking down the feature into discrete tasks. Each task should be:
- **Self-contained**: Can be implemented and tested independently
- **Single-purpose**: Does one thing well
- **Testable**: Has clear acceptance criteria

Example from Account Vending Service:
```
- DynamoDB table and indexes
- Go Lambda function scaffold
- API Gateway configuration
- Step Functions state machine
- Request handlers
- Approval workflow logic
- Notification integration
```

### Step 2: Identify Dependencies

For each task, ask:
1. What files/modules does this task need that don't exist yet?
2. What types/interfaces does this task consume?
3. What infrastructure must exist before this can deploy?

Create a dependency map:
```
Task                          Dependencies
─────────────────────────────────────────────────────
DynamoDB table                None
Go Lambda scaffold            None
API Gateway                   Lambda function ARN
Step Functions                Lambda function ARN
Request handlers              Lambda scaffold, DynamoDB table
Approval workflow             Request handlers, Step Functions
Notification integration      Approval workflow
```

### Step 3: Group into Waves

**Wave 0** (Optional - Local Only)
- Foundation work that doesn't deploy
- Shared types, interfaces, contracts
- Test fixtures and mocks
- Usually done locally, not via GitHub issues

**Wave 1** - Independent Infrastructure
- Tasks with no dependencies
- Can all start immediately
- Usually: databases, storage, base compute

**Wave 2+** - Dependent Tasks
- Group tasks that depend only on previous waves
- Tasks within the same wave must be independent of each other

Example Wave Assignment:
```
Wave 1 (Parallel):
  ├── DynamoDB table           [no deps]
  ├── Go Lambda scaffold       [no deps]
  └── IAM roles                [no deps]

Wave 2 (Parallel, after Wave 1):
  ├── API Gateway              [needs: Lambda ARN]
  ├── Step Functions           [needs: Lambda ARN]
  └── Request handlers         [needs: Lambda, DynamoDB]

Wave 3 (Parallel, after Wave 2):
  ├── Approval workflow        [needs: handlers, Step Functions]
  └── Integration tests        [needs: API Gateway, handlers]

Wave 4 (After Wave 3):
  └── Notification integration [needs: approval workflow]
```

### Step 4: Validate Wave Assignment

Check each wave for violations:

**Violation 1: Circular dependencies within wave**
```
❌ Wave 2:
   - Task A needs Task B's output
   - Task B needs Task A's output

Fix: Break circular dependency or merge into one task
```

**Violation 2: Missing dependency link**
```
❌ Wave 2:
   - Task C needs types from Task D
   - Task D is also in Wave 2

Fix: Move Task D to Wave 1, or merge C and D
```

**Violation 3: Over-sequentialization**
```
❌ Every task in its own wave (Wave 1-10)

Fix: Review dependencies - can tasks truly run in parallel?
```

## Task Granularity Rules

### Bundle Sequential Code
When code changes must happen in sequence (file A before file B), make it ONE task:
```
✅ Good: "Implement API handlers with validation"
   - handlers.go
   - validation.go
   - handlers_test.go

❌ Bad: Three separate tasks that must merge in order
```

### Separate by Integration Boundary
Create separate tasks when:
- Different deployment targets (Lambda vs. API Gateway)
- Different runtime environments (Go vs. TypeScript)
- Different AWS services (DynamoDB vs. Step Functions)

```
✅ Good separation:
   - Task 1: DynamoDB table + indexes (OpenTofu)
   - Task 2: Lambda function (Go)
   - Task 3: API Gateway (OpenTofu)

❌ Bad separation:
   - Task 1: DynamoDB table
   - Task 2: DynamoDB indexes (same deployment!)
```

### Test Colocation
Keep tests with their implementation:
```
✅ Good: "Implement user service with tests"
   - user_service.go
   - user_service_test.go

❌ Bad: Separate "Write tests" task
   (Creates dependency on implementation task)
```

## Creating GitHub Issues

### Issue Format

Use the `@claude Task` issue template which includes:
- **Base Branch**: The group branch (e.g., `group-4/account-vending`)
- **Wave**: Wave number (1-5)
- **Dependencies**: List of issue numbers this task depends on
- **Task Description**: Detailed implementation requirements
- **Files to Create/Modify**: Explicit file list
- **Acceptance Criteria**: Testable checkboxes
- **Contracts**: Type definitions or API contracts to conform to

### Dependency Syntax

In the Dependencies section, reference issues by number:
```markdown
## Dependencies
- #42 (DynamoDB table)
- #43 (Lambda scaffold)
```

The `check-dependencies.yml` workflow parses this and manages labels automatically.

### Labels to Apply

| Label | When to Apply |
|-------|---------------|
| `claude-task` | Always - enables automation |
| `wave-N` | Wave assignment (1-5) |
| `{epic-name}` | Epic label for filtering |
| `infrastructure` | OpenTofu/IaC tasks |
| `backend` | Go service tasks |
| `api` | API Gateway tasks |

### Example Issue Creation

```bash
gh issue create \
  --title "[E1-4.1] Create DynamoDB table for account vending" \
  --label "claude-task,wave-1,platform-foundation,infrastructure" \
  --body "$(cat <<'EOF'
## Base Branch
group-4/account-vending-service

## Wave
1

## Dependencies
None

## Task Description
Create the DynamoDB table and indexes for the Account Vending Service.

**Table Structure:**
- Table name: `AccountVendingRequests`
- Partition key: `PK` (String)
- Sort key: `SK` (String)
- GSI: `StatusIndex` on `status` attribute

## Files to Create/Modify
- `infrastructure/modules/account-vending/dynamodb.tf` (create)
- `infrastructure/modules/account-vending/outputs.tf` (modify)

## Acceptance Criteria
- [ ] DynamoDB table created with correct schema
- [ ] GSI configured for status queries
- [ ] Point-in-time recovery enabled
- [ ] Table ARN exported as output
- [ ] `tofu validate` passes
- [ ] `tofu plan` shows expected resources

## Contracts
```hcl
# Expected outputs
output "table_arn" {
  value = aws_dynamodb_table.account_vending.arn
}
output "table_name" {
  value = aws_dynamodb_table.account_vending.name
}
```

---
@claude implement this task following TDD
EOF
)"
```

## Wave Execution Timeline

```
T=0: Create all Wave 1 issues
     └── @claude instances start on each Wave 1 task

T=1: Wave 1 tasks complete
     └── Dependency check workflow marks Wave 2 tasks as "ready"
     └── @claude instances start on Wave 2 tasks

T=2: Wave 2 tasks complete
     └── Wave coordinator creates integration notice
     └── Human reviews and merges Wave 1+2 to group branch
     └── Wave 3 tasks become "ready"

...continues until all waves complete...

T=N: All waves complete
     └── PR group branch to dev
     └── CI/CD runs full validation
```

## Troubleshooting

### Task stuck as "blocked"
1. Check the Dependencies section format
2. Ensure referenced issues exist
3. Verify dependent issues are closed (not just merged)

### Wave running sequentially instead of parallel
1. Verify all tasks in the wave have no inter-dependencies
2. Check that dependency issues are properly closed
3. Review workflow logs for errors

### Merge conflicts during integration
1. Review files modified by each task
2. Check for overlapping changes
3. Follow integration checklist (see `integration-checklist.md`)

## Best Practices

1. **Start small**: 2-3 waves for first epic, expand as team gains experience
2. **Buffer time**: Account for integration overhead between waves
3. **Clear contracts**: Define interfaces in Wave 0 that all waves can reference
4. **Limit wave size**: 3-5 parallel tasks per wave is manageable
5. **Integration checkpoints**: Don't let waves pile up - integrate after each wave
