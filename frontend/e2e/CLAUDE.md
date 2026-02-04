# E2E Tests - CLAUDE.md

## Overview
Playwright end-to-end tests for frontend routes.

## Files
- `hello.spec.ts` - HelloPage validation test

## Running Tests
```bash
npx playwright test
```

## Test Pattern

E2E tests follow TDD workflow:
1. Create test file in `frontend/e2e/{feature}.spec.ts`
2. Run and verify test FAILS (Red phase)
3. Implement feature code
4. Run and verify test PASSES (Green phase)
