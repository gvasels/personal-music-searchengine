# Epic Planning Template

Use this template when planning a new epic for wave-based execution with @claude.

---

## Epic: [Name]

**Label**: `epic-name` (create via `gh label create epic-name --color "#5319e7"`)
**Group Branch**: `group-N/epic-name`
**Base Branch**: `dev`

### Overview

[Brief description of what this epic delivers]

### Goals

1. [Primary goal]
2. [Secondary goal]
3. [Tertiary goal]

### Non-Goals

- [What this epic explicitly does NOT cover]

---

## Wave 0: Foundation (Local Only)

> Wave 0 tasks are done locally without GitHub issues. They establish contracts and interfaces that all subsequent waves depend on.

### Tasks

| Task | Description | Deliverables |
|------|-------------|--------------|
| 0.1 | Define data models | Type definitions, schemas |
| 0.2 | Define API contracts | OpenAPI spec, interface definitions |
| 0.3 | Create test fixtures | Mock data, test utilities |

### Deliverables Checklist
- [ ] Data models defined in `design.md`
- [ ] API contracts documented
- [ ] Shared types exported
- [ ] Test fixtures created

---

## Wave 1: Infrastructure

> Independent infrastructure tasks that can run in parallel.

### Tasks

| # | Title | Description | Dependencies | Files |
|---|-------|-------------|--------------|-------|
| 1.1 | [Task title] | [Description] | None | `path/to/file.tf` |
| 1.2 | [Task title] | [Description] | None | `path/to/file.tf` |
| 1.3 | [Task title] | [Description] | None | `path/to/file.tf` |

### Issue Creation Commands

```bash
# Create Wave 1 issues
gh issue create --title "[E{N}-1.1] {Title}" \
  --label "claude-task,wave-1,epic-name,infrastructure" \
  --body "$(cat <<'EOF'
## Base Branch
group-N/epic-name

## Wave
1

## Dependencies
None

## Task Description
{Detailed description}

## Files to Create/Modify
- `{path}` (create/modify)

## Acceptance Criteria
- [ ] {Criterion 1}
- [ ] {Criterion 2}
- [ ] `tofu validate` passes

---
@claude implement this task following TDD
EOF
)"

# Repeat for 1.2, 1.3, etc.
```

---

## Wave 2: Core Implementation

> Tasks that depend on Wave 1 infrastructure.

### Tasks

| # | Title | Description | Dependencies | Files |
|---|-------|-------------|--------------|-------|
| 2.1 | [Task title] | [Description] | #1.1 | `path/to/file.go` |
| 2.2 | [Task title] | [Description] | #1.2 | `path/to/file.go` |
| 2.3 | [Task title] | [Description] | #1.1, #1.3 | `path/to/file.go` |

### Issue Creation Commands

```bash
# Create Wave 2 issues (after Wave 1 issues exist)
gh issue create --title "[E{N}-2.1] {Title}" \
  --label "claude-task,wave-2,epic-name,backend" \
  --body "$(cat <<'EOF'
## Base Branch
group-N/epic-name

## Wave
2

## Dependencies
- #{issue_number} (Wave 1 infrastructure)

## Task Description
{Detailed description}

## Files to Create/Modify
- `{path}` (create/modify)

## Acceptance Criteria
- [ ] {Criterion 1}
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Contracts
```go
// From Wave 0
type {Interface} interface {
    {Method}(ctx context.Context) error
}
```

---
@claude implement this task following TDD
EOF
)"
```

---

## Wave 3: Integration

> Tasks that integrate Wave 2 components.

### Tasks

| # | Title | Description | Dependencies | Files |
|---|-------|-------------|--------------|-------|
| 3.1 | [Task title] | [Description] | #2.1, #2.2 | `path/to/file` |
| 3.2 | [Task title] | [Description] | #2.3 | `path/to/file` |

---

## Wave 4: Polish & Documentation

> Final tasks for documentation, testing, and polish.

### Tasks

| # | Title | Description | Dependencies | Files |
|---|-------|-------------|--------------|-------|
| 4.1 | Integration tests | End-to-end tests | #3.1, #3.2 | `*_test.go` |
| 4.2 | Documentation | README, CLAUDE.md | #3.1, #3.2 | `*.md` |

---

## Dependency Graph

```
Wave 0 (Local):
  └── Contracts & Types

Wave 1 (Parallel):
  ├── 1.1 ─────────┐
  ├── 1.2 ────────┐│
  └── 1.3 ───────┐││
                 │││
Wave 2 (Parallel):│││
  ├── 2.1 ◄──────┘││
  ├── 2.2 ◄───────┘│
  └── 2.3 ◄────────┘

Wave 3 (Parallel):
  ├── 3.1 ◄── 2.1, 2.2
  └── 3.2 ◄── 2.3

Wave 4 (After 3):
  └── 4.1, 4.2 ◄── 3.1, 3.2
```

---

## Execution Timeline

| Phase | Duration | Activities |
|-------|----------|------------|
| Wave 0 | 1-2 hours | Local foundation work |
| Wave 1 | ~30 min | Create issues, @claude parallel execution |
| Integration 1 | 15 min | Merge Wave 1 to group branch |
| Wave 2 | ~30 min | @claude parallel execution |
| Integration 2 | 15 min | Merge Wave 2 to group branch |
| Wave 3 | ~30 min | @claude parallel execution |
| Integration 3 | 15 min | Merge Wave 3 to group branch |
| Wave 4 | ~30 min | @claude parallel execution |
| Final | 30 min | PR to dev, review, merge |

**Estimated Total**: 3-4 hours (vs 8+ hours sequential)

---

## Checklist

### Planning
- [ ] Epic label created
- [ ] Group branch created from dev
- [ ] Wave 0 completed locally
- [ ] All tasks assigned to waves
- [ ] Dependencies documented
- [ ] Issue templates prepared

### Execution
- [ ] Wave 1 issues created and assigned
- [ ] Wave 1 complete, merged to group branch
- [ ] Wave 2 issues created and assigned
- [ ] Wave 2 complete, merged to group branch
- [ ] Wave 3 issues created and assigned
- [ ] Wave 3 complete, merged to group branch
- [ ] Wave 4 issues created and assigned
- [ ] Wave 4 complete, merged to group branch

### Completion
- [ ] All tests passing
- [ ] Documentation updated
- [ ] PR created to dev
- [ ] PR reviewed and merged
- [ ] Task branches cleaned up
