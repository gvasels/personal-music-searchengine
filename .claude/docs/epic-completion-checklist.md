# Epic Completion Checklist

This checklist ensures all aspects of an epic are properly completed, documented, and deployed before marking the epic as complete.

## When to Use This Checklist

Use this checklist when:
- All tasks in an epic's `tasks.md` are marked as completed
- All acceptance criteria in `epics-user-stories.md` are met
- You're ready to mark the epic as complete

## Completion Steps

### 1. Code & Implementation ✅

- [ ] **All tasks completed** - Verify all tasks in `.spec-workflow/specs/{spec}/tasks.md` are marked `[x] Completed`
- [ ] **All tests passing** - Run full test suite:
  ```bash
  # Go services
  cd services/{service} && go test ./...

  # OpenTofu validation
  cd infrastructure/accounts/{account} && terragrunt validate
  ```
- [ ] **Code reviewed** - All PRs reviewed and merged
- [ ] **Breaking changes documented** - If any breaking changes exist, document in PR description

### 2. Documentation Updates 📝

#### 2.1 Epic Status Documentation

- [ ] **Update epics-user-stories.md**
  - Mark epic status as `✅ **COMPLETE** (YYYY-MM-DD)`
  - Check off all acceptance criteria: `- [x]`
  - Add implementation details to each criterion
  - Add epic to "Completed Epics" section in Implementation Notes

- [ ] **Update root CLAUDE.md**
  - Update project status line to reflect completed epics
  - Example: `Status: Active development - Epics 1-3 complete. Epic 4 in progress.`

- [ ] **Update CHANGELOG.md**
  - Add comprehensive epic entry with date
  - Include all deliverables (modules, services, features)
  - Document security enhancements
  - Include deployment table if infrastructure was deployed
  - Add architecture decisions
  - Document breaking changes and migration notes

#### 2.2 Planning Document Sync

- [ ] **Review planning-docs/ for references**
  - Search for fictional/planned account names that should now be actual account IDs
  - Update architecture diagrams if infrastructure changed
  - Verify service names match actual implementation
  - Check that technology stack references are current (no EC2 if using Lambda, etc.)

- [ ] **Validate account references**
  - Cross-reference all account names against `infrastructure/docs/aws-organizations.md`
  - Replace fictional accounts with actual deployed accounts + IDs
  - Mark future accounts explicitly with "future:" prefix or 🔮 emoji
  - Ensure email formats use `platform-admin+{service}{environment}@oopo.io`

- [ ] **Update architecture diagrams**
  - If epic changed infrastructure architecture, update diagrams in planning docs
  - Ensure diagrams reflect current state, not just vision
  - Add deployment status indicators (✅ Deployed, 🔮 Planned)

#### 2.3 Code Documentation (MANDATORY - DO NOT SKIP)

**⚠️ DO NOT PROCEED to Section 3 until all CLAUDE.md files are verified.**

- [ ] **Discover ALL changed directories**:
  ```bash
  git diff --name-only $(git merge-base HEAD main) HEAD | xargs -I {} dirname {} | sort -u | grep -v '^\.$'
  ```

- [ ] **Verify each changed directory has CLAUDE.md**:
  ```bash
  MISSING=""
  for dir in $(git diff --name-only $(git merge-base HEAD main) HEAD | xargs -I {} dirname {} | sort -u | grep -v '^\.$'); do
    if [ -d "$dir" ] && [ ! -f "$dir/CLAUDE.md" ] && [[ "$dir" != .claude* ]] && [[ "$dir" != .spec-workflow* ]] && [[ "$dir" != .github* ]]; then
      MISSING="$MISSING\n  MISSING: $dir/CLAUDE.md"
    fi
  done
  if [ -n "$MISSING" ]; then
    echo "⚠️ STOP - Missing CLAUDE.md files:$MISSING"
  else
    echo "✅ All changed directories have CLAUDE.md"
  fi
  ```

- [ ] **⚠️ STOP if any CLAUDE.md missing** - Create them before proceeding

- [ ] **Spawn documentation-generator for missing/outdated CLAUDE.md**:
  - Run `documentation-generator` agent with list of directories needing docs
  - Each CLAUDE.md MUST include: Overview, Files, Key Functions/Exports, Dependencies, Usage Examples

- [ ] **Re-verify - zero missing CLAUDE.md files**:
  - Re-run the verification check above
  - Output MUST be empty

- [ ] **Update API documentation**
  - Generate/update OpenAPI specs for new APIs
  - Update API endpoint documentation
  - Document request/response formats

### 3. Deployment & Infrastructure 🚀

- [ ] **Infrastructure deployed** (if applicable)
  - All OpenTofu/Terragrunt configurations applied
  - Resources created in correct AWS accounts
  - Outputs documented in account CLAUDE.md files

- [ ] **Services deployed** (if applicable)
  - Lambda functions deployed with latest code
  - API Gateways configured
  - DynamoDB tables created and seeded if needed

- [ ] **DNS & Certificates configured** (if applicable)
  - Route 53 records created
  - ACM certificates validated
  - Custom domains mapped to APIs

- [ ] **Update deployment documentation**
  - Add deployment instructions to relevant account CLAUDE.md
  - Document any manual steps required
  - Update `infrastructure/docs/aws-organizations.md` if new accounts created

### 4. Validation & Testing ✅

- [ ] **Verify CLAUDE.md completeness (MANDATORY)**
  - ALL changed directories have CLAUDE.md (verified in 2.3)
  - ALL CLAUDE.md files reflect current code (not stale)
  - **⚠️ DO NOT PROCEED if any changed directory is missing CLAUDE.md**

- [ ] **Run automated consistency checks**
  - Execute `doc-consistency-checker` agent
  - Fix any flagged inconsistencies
  - Verify account names, epic statuses, architecture references
  - Verify CLAUDE.md coverage (should flag missing files as HIGH severity)

- [ ] **Integration testing** (if applicable)
  - Run integration tests with `INTEGRATION_TEST=true`
  - Verify cross-account role assumption works
  - Test API endpoints end-to-end

- [ ] **Smoke testing** (if applicable)
  - Verify deployed services are reachable
  - Test critical user flows
  - Verify monitoring and logging configured

### 5. Handoff & Communication 📢

- [ ] **PR Description complete**
  - Include epic summary
  - List all completed tasks
  - Document breaking changes
  - Include deployment notes
  - Add test results summary

- [ ] **Update project board** (if using GitHub Projects)
  - Move epic card to "Complete"
  - Close related issues
  - Update milestone progress

## Quick Reference: Files That Must Be Updated

When completing ANY epic, these files typically need updates:

| File | Update Required | Example |
|------|----------------|---------|
| `implementation-plan/epics-user-stories.md` | Mark epic complete, check criteria | `✅ **COMPLETE** (2024-12-22)` |
| `CLAUDE.md` (root) | Update project status | `Epics 1-3 complete. Epic 4 in progress.` |
| `CHANGELOG.md` | Add epic entry | `## [2024-12-22] - Epic N: Feature Name` |
| `.spec-workflow/specs/{spec}/tasks.md` | All tasks marked `[x]` | Verify completion dates |

Additional files (if applicable):

| File | When to Update | Example |
|------|---------------|---------|
| `infrastructure/docs/aws-organizations.md` | New accounts created | Add account to inventory table |
| `planning-docs/*.md` | Architecture changed | Update diagrams, account references |
| Account-level `CLAUDE.md` | Infrastructure deployed | Document deployed resources |
| Module `CLAUDE.md` | New module created | Document module purpose, exports |

## Common Mistakes to Avoid

❌ **Don't:**
- Skip CLAUDE.md for changed directories (most common violation)
- Skip CHANGELOG.md updates (critical for tracking architectural changes)
- Leave fictional account names in planning docs after accounts are created
- Forget to update root CLAUDE.md project status
- Mark epic complete without checking all acceptance criteria
- Leave "future:" markers on accounts that are now deployed
- Create PR before verifying all CLAUDE.md files exist

✅ **Do:**
- Run CLAUDE.md discovery check FIRST before any other doc step
- Use this checklist every time you complete an epic
- Run the doc-consistency-checker agent before marking complete
- Update planning docs to reflect actual implementation vs. original vision
- Document breaking changes and migration paths
- Cross-reference account names with aws-organizations.md
- Verify zero missing CLAUDE.md files before creating PR

## Automation Support

The following agents can help with epic completion:

- **documentation-generator**: Creates/updates CLAUDE.md files
- **doc-consistency-checker**: Validates documentation consistency
- **code-reviewer**: Reviews code quality before marking complete

## Example Epic Completion Session

```bash
# 1. Verify all tasks complete
cat .spec-workflow/specs/cross-account-iam/tasks.md | grep "^\[x\]" | wc -l

# 2. Run tests
cd services/account-vending && go test ./...

# 3. Update documentation
# ... manually edit epics-user-stories.md, CLAUDE.md, CHANGELOG.md

# 4. Run consistency checks
claude --agent doc-consistency-checker

# 5. Commit documentation updates
git add implementation-plan/epics-user-stories.md CLAUDE.md CHANGELOG.md
git commit -m "docs: Mark Epic 3 (Cross-Account IAM) as complete"

# 6. Push and create/update PR
git push origin group-N/{feature}
gh pr create --base dev --head group-N/{feature}
```

## Related Documentation

- `.claude/docs/tdd-workflow.md` - Test-driven development workflow
- `.claude/docs/wave-assignment.md` - Wave-based execution for parallel work
- `.claude/commands/sdlc.md` - Full SDLC workflow documentation
