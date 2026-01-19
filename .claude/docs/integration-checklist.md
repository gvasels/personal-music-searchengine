# Wave Integration Checklist

Use this checklist when integrating completed wave tasks into the group branch.

## Pre-Integration

### 1. Verify Wave Completion
- [ ] All tasks in the wave have closed issues
- [ ] Each task has a merged PR to its task branch
- [ ] Wave completion notice created (automatic via `wave-coordinator.yml`)
- [ ] No "blocked" issues remaining in the wave

### 2. Review Task Branches
For each task branch in the wave:
- [ ] Tests pass locally
- [ ] No uncommitted changes
- [ ] Clean git history (squash if needed)

```bash
# List task branches for the wave
git branch -r | grep "task-" | grep "origin"

# Check status of each task branch
for branch in $(git branch -r | grep "task-" | awk '{print $1}'); do
  echo "=== $branch ==="
  git log --oneline -3 $branch
done
```

## Integration Steps

### 3. Update Group Branch
```bash
# Ensure group branch is up to date
git checkout group-N/epic-name
git pull origin group-N/epic-name
```

### 4. Merge Task Branches
Merge each task branch from the wave:

```bash
# For each task in the wave (in dependency order if any)
git merge origin/task-XX-description --no-edit

# If conflicts occur, resolve and continue
git add .
git commit -m "Resolve merge conflicts from task-XX"
```

### 5. Resolve Merge Conflicts
Common conflict patterns and resolutions:

| Conflict Type | Resolution |
|--------------|------------|
| Import paths | Use the most complete import set |
| Type definitions | Prefer the version with more fields |
| Function signatures | Keep backward-compatible version |
| Configuration files | Merge both configurations |
| Test fixtures | Include all test data |

**Conflict Resolution Process:**
```bash
# See conflicting files
git diff --name-only --diff-filter=U

# For each file, resolve conflicts
# Look for <<<<<<< ======= >>>>>>> markers

# After resolving
git add <resolved-file>
git commit -m "Resolve: <description of resolution>"
```

### 6. Verify Integration
Run full validation suite:

```bash
# Go code
go build ./...
go test ./... -v
golangci-lint run

# OpenTofu
cd infrastructure/modules/<module>
tofu init
tofu validate
tofu plan

# TypeScript (if applicable)
npm run lint
npm test
npm run build
```

## Post-Integration

### 7. Update Documentation
- [ ] CLAUDE.md files updated for new directories
- [ ] README updated if public APIs changed
- [ ] CHANGELOG.md entry added

### 8. Push and Verify CI
```bash
git push origin group-N/epic-name

# Monitor CI status
gh run list --branch group-N/epic-name
```

### 9. Clean Up Task Branches
After successful integration:

```bash
# Delete merged task branches (remote)
git push origin --delete task-XX-description

# Delete local task branches
git branch -d task-XX-description
```

### 10. Notify and Proceed
- [ ] Comment on wave completion issue that integration is done
- [ ] Close the wave completion issue
- [ ] Verify next wave tasks are marked "ready"

## Integration Checklist by Wave Size

### Small Wave (2-3 tasks)
```
[ ] Review all task PRs
[ ] Merge to group branch
[ ] Run tests
[ ] Push
```

### Medium Wave (4-6 tasks)
```
[ ] Review all task PRs
[ ] Create integration branch if needed
[ ] Merge in dependency order
[ ] Resolve conflicts iteratively
[ ] Run full test suite
[ ] Run security scan
[ ] Push
```

### Large Wave (7+ tasks)
```
[ ] Split into sub-waves if possible
[ ] Assign integration lead
[ ] Create integration branch
[ ] Merge and test incrementally
[ ] Full regression test
[ ] Performance test if applicable
[ ] Security scan
[ ] Code review integration
[ ] Push
```

## Rollback Procedure

If integration fails and cannot be fixed quickly:

```bash
# Reset group branch to pre-integration state
git checkout group-N/epic-name
git reset --hard origin/group-N/epic-name

# Or reset to specific commit
git reset --hard <commit-hash-before-integration>

# Force push (if already pushed broken state)
git push origin group-N/epic-name --force
```

## Troubleshooting

### Build Failures After Merge
1. Check for missing imports from merged code
2. Verify dependency versions match
3. Look for renamed types or functions

### Test Failures After Merge
1. Check for conflicting test data
2. Verify mock configurations
3. Look for order-dependent tests

### CI Failures
1. Check for environment-specific issues
2. Verify secrets and configurations
3. Review CI logs for specific error

## Quick Reference

```bash
# Full integration flow (copy-paste ready)
WAVE=1
EPIC="account-vending"
GROUP="group-4"

git checkout ${GROUP}/${EPIC}
git pull origin ${GROUP}/${EPIC}

# List and merge all task branches for the wave
# (Adjust task numbers as needed)
git merge origin/task-42-dynamodb --no-edit
git merge origin/task-43-lambda-scaffold --no-edit
git merge origin/task-44-iam-roles --no-edit

# Verify
go build ./...
go test ./...
tofu validate

# Push
git push origin ${GROUP}/${EPIC}

# Clean up (after CI passes)
git push origin --delete task-42-dynamodb
git push origin --delete task-43-lambda-scaffold
git push origin --delete task-44-iam-roles

echo "Wave ${WAVE} integration complete!"
```
