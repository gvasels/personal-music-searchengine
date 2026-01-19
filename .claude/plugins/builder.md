---
name: builder
description: Build verification, linting, and type checking
phase: 4-build
skills: []
agents: []
mcp_servers: []
---

# Builder Plugin

Verifies code quality through builds, linting, type checking, and test execution.

## Phase Position

```
1. SPEC → 2. TEST → 3. CODE → [4. BUILD] → 5. SECURITY → 6. DOCS
                                   ▲
                                   YOU ARE HERE
```

## Prerequisites

From previous phases:
- Implementation code (from code-implementer)
- Test suites (from test-writer)

## Build Pipeline

```
┌─────────────────────────────────────────┐
│              LINT CHECK                  │
│   ESLint, golangci-lint, flake8         │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│             TYPE CHECK                   │
│      TypeScript tsc, Go build           │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│             TEST SUITE                   │
│   Unit + Integration + Coverage         │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│           BUILD ARTIFACTS                │
│    Compile, bundle, package             │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         BUILD VERIFICATION              │
│    Start server, smoke test             │
└─────────────────────────────────────────┘
```

## Workflow

### Step 1: Lint Check

```bash
# TypeScript/JavaScript
npm run lint
# or
npx eslint . --ext .ts,.tsx --max-warnings 0

# Go
golangci-lint run ./...

# Python
flake8 .
black --check .
```

**Fix lint errors before proceeding.**

### Step 2: Type Check

```bash
# TypeScript
npx tsc --noEmit

# Go (built into compiler)
go build ./...

# Python (optional but recommended)
mypy .
```

**Fix type errors before proceeding.**

### Step 3: Run Full Test Suite

```bash
# TypeScript/JavaScript
npm test -- --coverage --coverageThreshold='{"global":{"branches":80,"functions":80,"lines":80}}'

# Go
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out

# Python
pytest --cov=src --cov-fail-under=80
```

**Coverage Requirements:**

| Metric | Minimum | Target |
|--------|---------|--------|
| Lines | 80% | 90% |
| Branches | 80% | 85% |
| Functions | 80% | 90% |

### Step 4: Build Artifacts

```bash
# TypeScript/Vite
npm run build
# Outputs to dist/

# Go Lambda
GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
zip function.zip bootstrap

# Python Lambda
pip install -r requirements.txt -t package/
cd package && zip -r ../function.zip .
```

### Step 5: Build Verification

```bash
# Start locally and verify
npm run preview
# or
./bootstrap &
curl http://localhost:3000/health

# Smoke test endpoints
curl -X POST http://localhost:3000/api/resource \
  -H "Content-Type: application/json" \
  -d '{"name": "test"}'
```

## Quality Gates

| Check | Tool | Threshold |
|-------|------|-----------|
| Lint | ESLint/golangci-lint | 0 errors, 0 warnings |
| Types | TypeScript/Go | 0 errors |
| Unit Tests | Vitest/Testify | 100% pass |
| Coverage | c8/go cover | 80% minimum |
| Build | Vite/Go | Success |
| Bundle Size | Vite | < 500KB initial |

## Common Build Issues

### TypeScript

```typescript
// Issue: Implicit any
const handler = (req) => { }  // Error
const handler = (req: Request) => { }  // Fixed

// Issue: Unused variables
const unused = 'value';  // Error
const _unused = 'value';  // Prefix with _ if intentional
```

### Go

```go
// Issue: Unused import
import "fmt"  // Error if fmt not used

// Issue: Error not handled
result, _ := riskyFunction()  // Lint warning
result, err := riskyFunction()
if err != nil { return err }
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| Build output | `dist/` or `bootstrap` | Deployable artifacts |
| Coverage report | `coverage/` | Test coverage HTML |
| Lint report | Console | Code quality issues |
| Type errors | Console | Type safety issues |

## Failure Handling

If any check fails:

1. **Lint fails** → Return to `code-implementer`, fix style issues
2. **Type check fails** → Return to `code-implementer`, fix types
3. **Tests fail** → Return to `code-implementer`, fix implementation
4. **Coverage low** → Return to `test-writer`, add more tests
5. **Build fails** → Investigate dependency/config issues

## Handoff to Next Phase

After successful build:
1. All lints pass
2. All types check
3. All tests pass with coverage
4. Build artifacts generated
5. **NEXT**: Pass to `security-checker` plugin for security audit
