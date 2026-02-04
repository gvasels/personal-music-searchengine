# Integration Cadence Guidelines

This document defines when and how to integrate completed waves into the group branch.

## Core Principle

**Integrate early, integrate often.** Don't let waves pile up - merge each wave to the group branch as soon as it completes.

## Integration Timing

### When to Integrate

| Trigger | Action |
|---------|--------|
| All Wave N tasks closed | Integrate Wave N immediately |
| Wave completion notice created | Begin integration within 1 hour |
| 2+ waves pending integration | Stop new work, integrate first |
| End of work session | Integrate any completed waves |

### When NOT to Integrate

| Condition | Action |
|-----------|--------|
| Wave has open tasks | Wait for completion |
| Build failures in task branches | Fix before integrating |
| Unresolved merge conflicts | Resolve in task branch first |
| Blocking review comments | Address before merging |

## Integration Workflow

### 1. Immediate Integration (Recommended)

```
Wave 1 Complete
     │
     ▼
Integrate to Group (15 min)
     │
     ▼
Start Wave 2
     │
     ▼
Wave 2 Complete
     │
     ▼
Integrate to Group (15 min)
     │
     ▼
Continue...
```

**Benefits:**
- Catches integration issues early
- Smaller merge conflicts
- Clear progress tracking
- Easier rollback if needed

### 2. Batched Integration (Use Sparingly)

```
Wave 1 Complete ─┐
                 │
Wave 2 Complete ─┼─► Batch Integrate (30-45 min)
                 │
Wave 3 Complete ─┘
```

**When acceptable:**
- Very small waves (1-2 tasks each)
- High confidence in compatibility
- Time-constrained session

**Risks:**
- Larger merge conflicts
- Harder to isolate issues
- More complex rollback

## Integration SLAs

| Scenario | Target Time | Maximum Time |
|----------|-------------|--------------|
| Single wave | 15 minutes | 30 minutes |
| Two waves batched | 30 minutes | 1 hour |
| Full epic (all waves) | 1 hour | 2 hours |

## Integration Checklist

### Pre-Integration (5 minutes)
```
□ All wave tasks show as closed in GitHub
□ Wave completion notice exists (or create manually)
□ No "blocked" tasks remaining in the wave
□ Task branches exist and are pushed
```

### Integration (5-10 minutes)
```
□ Checkout group branch
□ Pull latest from origin
□ Merge each task branch (use merge-wave.sh)
□ Resolve any conflicts
```

### Validation (5 minutes)
```
□ go build ./... passes
□ go test ./... passes (or npm test, etc.)
□ tofu validate passes (if infrastructure)
□ No new lint errors
```

### Post-Integration (2 minutes)
```
□ Push group branch to origin
□ Comment on wave completion issue
□ Close wave completion issue
□ Delete merged task branches
```

## Conflict Resolution Priority

When merge conflicts occur, resolve using this priority:

1. **Contracts/Interfaces** - Use the version from the earlier wave
2. **Implementation** - Use the version with more complete functionality
3. **Tests** - Include ALL tests from both branches
4. **Configuration** - Merge both configurations (don't drop settings)
5. **Documentation** - Combine content from both sources

## Integration Frequency by Epic Size

| Epic Size | Waves | Integration Frequency |
|-----------|-------|----------------------|
| Small (5-10 tasks) | 2-3 | After each wave |
| Medium (10-20 tasks) | 3-4 | After each wave |
| Large (20+ tasks) | 4-5 | After each wave |

**Rule:** Always integrate after each wave, regardless of epic size.

## Handling Integration Failures

### Build Failure
```
1. Identify failing component
2. Check which task introduced the issue
3. Fix in the task branch
4. Re-merge task branch
5. Verify build passes
```

### Test Failure
```
1. Identify failing test
2. Determine if test or implementation is wrong
3. Fix in appropriate task branch
4. Re-merge and verify
```

### Conflict That Can't Be Resolved
```
1. Document the conflict
2. Create a new "integration fix" task
3. Assign to Wave N+1 or handle manually
4. Don't block other integrations
```

## Monitoring Integration Health

### Healthy Signs
- ✅ Waves integrating within 15 minutes
- ✅ Zero or minimal merge conflicts
- ✅ All tests passing after integration
- ✅ Clear progress through waves

### Warning Signs
- ⚠️ Multiple waves waiting for integration
- ⚠️ Repeated merge conflicts in same files
- ⚠️ Test failures after integration
- ⚠️ Integration taking > 30 minutes

### Critical Issues
- ❌ Build broken after integration
- ❌ Unable to resolve merge conflicts
- ❌ Dependent waves starting before integration
- ❌ More than 3 waves pending

## Best Practices

### Do
- Integrate immediately after wave completion
- Run full validation after each integration
- Keep integration sessions short and focused
- Document any issues encountered
- Clean up task branches after successful integration

### Don't
- Start new waves before integrating completed ones
- Skip validation steps to save time
- Force-push to group branches
- Ignore failing tests "to fix later"
- Let waves accumulate without integration

## Quick Commands

```bash
# Check wave status
./scripts/wave-status.sh platform-foundation

# Merge a completed wave
./scripts/merge-wave.sh group-4/account-vending 1

# Verify integration
go build ./... && go test ./...

# Push integrated changes
git push origin group-4/account-vending

# Clean up task branches (after verification)
git push origin --delete task-42-dynamodb
```

## Escalation Path

If integration is blocked:

1. **Self-resolve** (15 min) - Try standard conflict resolution
2. **Peer review** (15 min) - Get another developer's perspective
3. **Skip and document** - Move forward, create follow-up task
4. **Rollback** - Reset to pre-integration state if critical
