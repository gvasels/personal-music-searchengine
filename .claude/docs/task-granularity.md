# Task Granularity Guidelines

This document defines rules for splitting work into tasks for wave-based @claude execution.

## Core Principle

**A task should be the smallest unit of work that can be independently implemented, tested, and merged.**

## Granularity Rules

### Rule 1: Bundle Sequential Code Changes

When code changes must happen in a specific order within the same module, they belong in ONE task.

**Good (Single Task):**
```
Task: "Implement user service"
Files:
  - internal/models/user.go        (types)
  - internal/repository/user.go    (data access)
  - internal/service/user.go       (business logic)
  - internal/service/user_test.go  (tests)
```

**Bad (Multiple Tasks):**
```
Task 1: "Create user model"           ← These create artificial
Task 2: "Create user repository"      ← dependencies and
Task 3: "Create user service"         ← prevent parallelization
Task 4: "Write user tests"
```

### Rule 2: Separate by Deployment Boundary

Different deployment targets = different tasks.

**Good (Separate Tasks):**
```
Task 1: "Create DynamoDB table"       → OpenTofu deployment
Task 2: "Implement Lambda function"   → Go build + Lambda deployment
Task 3: "Configure API Gateway"       → OpenTofu deployment
```

**Why?** Each can be validated and deployed independently.

### Rule 3: Separate by Runtime/Language

Different languages or runtimes = different tasks.

**Good (Separate Tasks):**
```
Task 1: "Backend API endpoints"       → Go
Task 2: "Frontend components"         → TypeScript/React
Task 3: "Infrastructure"              → OpenTofu (HCL)
```

### Rule 4: Keep Tests with Implementation

Tests should be in the same task as the code they test.

**Good:**
```
Task: "Implement order service"
Files:
  - internal/service/order.go
  - internal/service/order_test.go
```

**Bad:**
```
Task 1: "Implement order service"
Task 2: "Write order service tests"   ← Creates unnecessary dependency
```

### Rule 5: Separate by AWS Service

Different AWS services that don't share resources = different tasks.

**Good (Separate Tasks):**
```
Task 1: "DynamoDB tables and indexes"
Task 2: "SQS queues and DLQ"
Task 3: "EventBridge rules"
Task 4: "IAM roles and policies"
```

### Rule 6: Group Related Schema Changes

Database schema changes that work together = one task.

**Good (Single Task):**
```
Task: "Create account vending schema"
- Table definition
- GSIs for access patterns
- LSIs if needed
- IAM policy for table access
```

**Bad (Separate Tasks):**
```
Task 1: "Create DynamoDB table"
Task 2: "Add GSI"                     ← Must be done together anyway
Task 3: "Add IAM policy"
```

## Size Guidelines

### Minimum Task Size
- At least 50+ lines of meaningful code
- At least one testable unit
- Clear acceptance criteria

### Maximum Task Size
- No more than 500 lines of code (excluding tests)
- Completable in under 30 minutes by @claude
- Reviewable in under 15 minutes by human

### Goldilocks Zone
- 100-300 lines of implementation code
- 50-150 lines of tests
- 3-5 files modified/created
- 15-20 minutes @claude execution time

## Decision Matrix

| Scenario | Single Task | Separate Tasks |
|----------|-------------|----------------|
| Same module, sequential files | ✅ | |
| Same module, parallel features | | ✅ |
| Different AWS services | | ✅ |
| Different languages | | ✅ |
| Tests for implementation | ✅ | |
| Integration tests across modules | | ✅ |
| Related schema changes | ✅ | |
| Unrelated infrastructure | | ✅ |

## Examples

### Example 1: API Endpoint

**Single Task - Correct:**
```
Task: "Implement GET /accounts endpoint"
- handlers/accounts.go (handler)
- handlers/accounts_test.go (tests)
- models/account.go (request/response types)
```

### Example 2: Multi-Service Feature

**Multiple Tasks - Correct:**
```
Task 1: "DynamoDB account table" (OpenTofu)
Task 2: "Account Lambda function" (Go)
Task 3: "API Gateway routes" (OpenTofu)
Task 4: "Integration tests" (Go)

Dependencies: 2→1, 3→2, 4→[1,2,3]
```

### Example 3: CRUD Operations

**Can go either way depending on complexity:**

Simple CRUD (Single Task):
```
Task: "Implement account CRUD operations"
- All handlers in one file
- Shared validation logic
- Related tests
```

Complex CRUD (Multiple Tasks):
```
Task 1: "Account creation with validation"
Task 2: "Account retrieval with caching"
Task 3: "Account update with audit logging"
Task 4: "Account deletion with soft-delete"
```

## Anti-Patterns

### Anti-Pattern 1: Over-Splitting
```
❌ Task 1: "Create user struct"
❌ Task 2: "Add validation to user struct"
❌ Task 3: "Add JSON tags to user struct"
```
**Fix:** Combine into "Implement user model with validation"

### Anti-Pattern 2: Under-Splitting
```
❌ Task: "Implement entire account vending service"
   - DynamoDB table
   - Lambda function
   - API Gateway
   - Step Functions
   - All handlers
   - All tests
```
**Fix:** Split by deployment boundary and AWS service

### Anti-Pattern 3: Test Separation
```
❌ Task 1: "Implement feature X"
❌ Task 2: "Write tests for feature X"
```
**Fix:** Include tests in the implementation task

### Anti-Pattern 4: Artificial Layering
```
❌ Wave 1: "Create all types"
❌ Wave 2: "Create all repositories"
❌ Wave 3: "Create all services"
❌ Wave 4: "Create all handlers"
```
**Fix:** Group by feature/module, not by layer

## Quick Reference

```
BUNDLE INTO SINGLE TASK:
├── Sequential code in same module
├── Types + implementation + tests
├── Related schema changes
└── Tightly coupled components

SEPARATE INTO MULTIPLE TASKS:
├── Different deployment targets
├── Different languages/runtimes
├── Different AWS services
├── Independent features
└── Cross-module integrations
```
