# Current vs. Future Distinction Markers

This guide defines visual markers to distinguish between **currently deployed** resources and **planned/future** resources in planning documentation.

## Why This Matters

Documentation often mixes:
- **Current state**: What exists in production right now
- **Vision**: What we plan to build in the future

Without clear distinction, readers can't tell what's real vs. aspirational, leading to:
- Confusion about which accounts exist
- Attempts to deploy to non-existent accounts
- Misleading architecture diagrams

## Marker Conventions

### For AWS Accounts

| Status | Format | Example |
|--------|--------|---------|
| **Deployed** | Name + ID | `oopo-cicd-prod (634945387634)` |
| **Planned** | 🔮 + Name + "future" tag | `🔮 oopo-shell-prod (future)` |

### For Infrastructure Components

| Status | Marker | Example |
|--------|--------|---------|
| **Deployed** | ✅ | `✅ Route 53 Hosted Zones` |
| **In Progress** | 🔄 | `🔄 CodeBuild Projects` |
| **Planned** | 🔮 | `🔮 EKS Clusters (backlogged for post-MVP)` |
| **Deprecated** | ⚠️ | `⚠️ EC2 Agents (replaced by CodeBuild)` |

### For Features

| Status | Marker | Example |
|--------|--------|---------|
| **Live** | ✅ | `✅ Manifest API v1` |
| **Beta** | 🧪 | `🧪 Canary Deployments` |
| **Planned** | 🔮 | `🔮 Multi-region failover` |

## Application Guidelines

### In Planning Documents

**Account Lists:**
```markdown
### AWS Organizations Structure

Root
├── deployments/
│   ├── oopo-cicd-dev (471544433440)
│   ├── oopo-cicd-staging (543613944458)
│   └── oopo-cicd-prod (634945387634)
├── infrastructure/
│   └── oopo-infrastructure-dns (510931056307)
└── services/
    ├── 🔮 oopo-identity-dev (future)
    ├── 🔮 oopo-identity-staging (future)
    └── 🔮 oopo-identity-prod (future)
```

**Architecture Diagrams:**
```markdown
## Current Architecture (✅ Deployed)

[Diagram showing only deployed components]

## Future Vision (🔮 Planned)

[Diagram showing planned enhancements]
```

**Feature Lists:**
```markdown
## Platform Features

### ✅ Currently Available
- Cross-account IAM roles
- Manifest API for version management
- Buildkite + CodeBuild CI/CD

### 🔄 In Progress (Epic N)
- Account vending service
- Permission boundaries

### 🔮 Planned (Backlog)
- Multi-region failover
- EKS-based services
- Canary deployments
```

### In Source of Truth Files

**infrastructure/docs/aws-organizations.md:**
- Only list deployed accounts (with IDs)
- Do not include future accounts
- This is the single source of truth for what exists

**implementation-plan/epics-user-stories.md:**
- Use markers for epic status: `✅ COMPLETE`, `🔄 IN PROGRESS`, `🔮 PLANNED`

## Doc-Consistency-Checker Integration

The doc-consistency-checker agent validates:
- ✅ Deployed accounts have IDs
- ✅ Future accounts are marked with 🔮 emoji or "future:" prefix
- ✅ Planning docs distinguish current vs. vision
- ❌ Future accounts listed without markers
- ❌ Deployed accounts listed without IDs

## Examples

### ✅ Good: Clear Distinction

```markdown
---
doc_type: planning
last_reviewed: 2025-12-22
---

# Platform Architecture

> **Note**: This document shows both **current deployed state** (✅) and
> **future planned features** (🔮). See `infrastructure/docs/aws-organizations.md`
> for definitive current state.

## Deployment Accounts

### ✅ Currently Deployed
- oopo-cicd-dev (471544433440)
- oopo-cicd-staging (543613944458)
- oopo-cicd-prod (634945387634)

### 🔮 Planned for Q1 2026
- oopo-shell-dev (future)
- oopo-shell-staging (future)
- oopo-shell-prod (future)
```

### ❌ Bad: No Distinction

```markdown
# Platform Architecture

## Deployment Accounts
- oopo-cicd-dev
- oopo-cicd-staging
- oopo-cicd-prod
- oopo-shell-dev
- oopo-shell-staging
- oopo-shell-prod
```

**Problems:**
- No account IDs for deployed accounts
- Can't tell which accounts exist vs. planned
- Could lead to deployment attempts on non-existent accounts

## Related Documentation

- `.claude/agents/doc-consistency-checker.md` - Automated validation
- `.claude/docs/epic-completion-checklist.md` - Epic completion requirements
- `infrastructure/docs/aws-organizations.md` - Single source of truth for accounts
