---
name: docs-generator
description: API and code documentation generation
phase: 6-documentation
skills: [documentation-generator]
agents: [documentation-generator]
mcp_servers: []
---

# Docs Generator Plugin

Generates comprehensive documentation as the final phase before deployment.

## Phase Position

```
1. SPEC → 2. TEST → 3. CODE → 4. BUILD → 5. SECURITY → [6. DOCS]
                                                           ▲
                                                           YOU ARE HERE
```

## Prerequisites

From previous phases:
- Implemented code (from code-implementer)
- Security-approved (from security-checker)
- Design specifications (from spec-writer)

## Documentation Pipeline

```
┌─────────────────────────────────────────┐
│          API DOCUMENTATION               │
│   OpenAPI spec, endpoint docs           │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         CODE DOCUMENTATION               │
│   TSDoc, GoDoc, inline comments         │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│        CLAUDE.md GENERATION              │
│   Directory-level documentation         │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│       LESSONS LEARNED UPDATE             │
│   Capture troubleshooting patterns      │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         TASKS.MD UPDATE                  │
│   Mark completed tasks in spec          │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│      ⚠️ CHANGELOG UPDATE (REQUIRED)      │
│   MUST document changes in CHANGELOG.md │
└─────────────────────────────────────────┘
```

## CHANGELOG.md Requirements (MANDATORY)

**Every feature or fix MUST have a CHANGELOG.md entry before the PR is complete.**

### Location
- Root `CHANGELOG.md` for all changes

### Format (Keep a Changelog)
```markdown
## [Unreleased]

### Added
- New features

### Changed
- Changes to existing functionality

### Fixed
- Bug fixes

### Security
- Security-related changes
```

### Entry Template
```markdown
### Added
- **[Feature Name]** - Brief description of what was added
  - Sub-point with details
  - Link to issue: #123
```

### When to Update
| Change Type | Section | Example |
|-------------|---------|---------|
| New feature | Added | New API endpoint, new module |
| Enhancement | Changed | Performance improvement, UX update |
| Bug fix | Fixed | Error correction, edge case handling |
| Vulnerability fix | Security | Dependency update, auth fix |
| Removal | Removed | Deprecated feature removal |

## Workflow

### Step 1: Generate API Documentation

From implemented endpoints, create OpenAPI spec:

```yaml
# docs/openapi.yaml
openapi: 3.1.0
info:
  title: Feature API
  version: 1.0.0
  description: API for managing resources

paths:
  /api/v1/resources:
    post:
      operationId: createResource
      summary: Create a new resource
      tags: [Resources]
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateResourceRequest'
            example:
              name: "my-resource"
              type: "standard"
      responses:
        '201':
          description: Resource created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Resource'
        '400':
          $ref: '#/components/responses/ValidationError'
        '409':
          $ref: '#/components/responses/ConflictError'
```

### Step 2: Generate Code Documentation

**TypeScript (TSDoc):**

```typescript
/**
 * Creates a new resource in the system.
 *
 * @param input - The resource creation parameters
 * @returns The created resource with generated ID
 * @throws {ValidationError} When input fails validation
 * @throws {ConflictError} When resource name already exists
 *
 * @example
 * ```typescript
 * const resource = await createResource({
 *   name: 'my-resource',
 *   type: 'standard'
 * });
 * console.log(resource.id); // 'res_abc123'
 * ```
 */
export async function createResource(
  input: CreateResourceInput
): Promise<Resource> {
  // Implementation
}
```

**Go (GoDoc):**

```go
// CreateResource creates a new resource in the system.
//
// It validates the input, checks for duplicates, and persists
// the resource to the database.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines
//   - input: Resource creation parameters
//
// Returns the created resource or an error if:
//   - Validation fails (ValidationError)
//   - Resource name exists (ConflictError)
//
// Example:
//
//	resource, err := svc.CreateResource(ctx, CreateResourceInput{
//	    Name: "my-resource",
//	    Type: "standard",
//	})
func (s *Service) CreateResource(
    ctx context.Context,
    input CreateResourceInput,
) (*Resource, error) {
    // Implementation
}
```

### Step 3: Generate CLAUDE.md Files

For each directory with new/modified code:

```markdown
# src/services/resources/CLAUDE.md

## Overview
Resource management service handling CRUD operations for the resources feature.

## Files

| File | Description |
|------|-------------|
| `service.ts` | Main service class with business logic |
| `repository.ts` | Database operations for resources |
| `types.ts` | TypeScript interfaces and schemas |
| `validation.ts` | Input validation logic |

## Key Functions

### `ResourceService.create(input: CreateResourceInput): Promise<Resource>`
Creates a new resource after validation.
- **Input**: `CreateResourceInput` - name, type, optional metadata
- **Output**: `Resource` - created resource with ID
- **Throws**: `ValidationError`, `ConflictError`

### `ResourceService.findById(id: string): Promise<Resource | null>`
Retrieves a resource by ID.
- **Input**: `id` - UUID of the resource
- **Output**: `Resource` or `null` if not found

## Dependencies

- `@/lib/database` - Database client
- `@/lib/validation` - Zod schemas
- `@/lib/errors` - Custom error types
```

### Step 4: Update Lessons Learned

**Capture troubleshooting patterns discovered during implementation.**

After completing a feature, review what was learned and add entries to the appropriate lessons file in `.claude/docs/`.

#### When to Add Entries

Add a lessons learned entry when:
- You spent significant time debugging an issue
- The solution wasn't obvious from error messages
- A workaround or pattern is project-specific
- Official documentation didn't cover the scenario

#### File Selection by Technology

| Technology | File |
|------------|------|
| Go (Lambda, APIs, testing) | `.claude/docs/go-lessons.md` |
| AWS (API Gateway, DynamoDB, Lambda, IAM) | `.claude/docs/aws-lessons.md` |
| OpenTofu/Terraform (IaC, state, modules) | `.claude/docs/opentofu-lessons.md` |
| TypeScript/React | `.claude/docs/typescript-lessons.md` |
| General/Cross-cutting | `.claude/docs/common-patterns.md` |

#### Entry Format

```markdown
---

### Problem: [Brief description]

**Symptom**:
```
[Exact error message or unexpected behavior]
```

**Root Cause**: [Why this happens]

**Solution**: [How to fix it]

```code
[Working code example]
```

**Debugging**: [Steps to investigate similar issues]

---
```

#### Example: Adding a Go Lesson

If during implementation you discovered that DynamoDB struct tags must be lowercase:

```markdown
---

### Problem: Empty fields when unmarshaling DynamoDB items

**Symptom**:
```
Struct fields are empty after attributevalue.UnmarshalMap() even though data exists in DynamoDB.
```

**Root Cause**: Field names in DynamoDB don't match Go struct tags.

**Solution**: Use lowercase `dynamodbav` tags that match DynamoDB attribute names:

```go
// WRONG - uses uppercase field names in DynamoDB
type RouteRecord struct {
    Product string `dynamodbav:"Product"`  // DynamoDB has "product"
}

// CORRECT - matches DynamoDB attribute names exactly
type RouteRecord struct {
    Product string `dynamodbav:"product"`
}
```

**Debugging**: Print the raw DynamoDB item to see actual attribute names:
```go
item, _ := client.GetItem(ctx, &dynamodb.GetItemInput{...})
fmt.Printf("Raw item: %+v\n", item.Item)
```

---
```

#### Example: Adding an AWS Lesson

If you discovered API Gateway stage prefix behavior:

```markdown
---

### Problem: 404 errors with named stages

**Symptom**: Requests to `https://api-id.execute-api.region.amazonaws.com/prod/health` return 404.

**Root Cause**: API Gateway HTTP API with named stages includes the stage name in the path sent to Lambda.

**Solution**: Handle stage prefix in your application routing:

```go
// Direct paths (for local testing and $default stage)
e.GET("/health", handlers.HealthCheck)

// Stage-prefixed paths (for named stages like "prod")
e.GET("/:stage/health", handlers.HealthCheck)
```

---
```

#### Workflow Integration

1. **During implementation**: Note any non-obvious issues encountered
2. **After tests pass**: Review notes and identify reusable patterns
3. **Before PR**: Add entries to appropriate lessons file
4. **Include in commit**: Stage lessons files with other documentation

### Step 5: Update tasks.md (Spec Completion Tracking)

**CRITICAL**: Mark completed tasks in the spec's `tasks.md` file.

1. **Locate the spec**: Find `.spec-workflow/specs/{feature}/tasks.md`
2. **Identify completed tasks**: Match implemented work to task IDs
3. **Update task status**: Change `[ ]` to `[x]` for completed tasks

**Task Status Markers:**

| Marker | Status | Meaning |
|--------|--------|---------|
| `[ ]` | Pending | Not started |
| `[-]` | In Progress | Currently being worked on |
| `[x]` | Completed | Implementation finished |

**Example Update:**

```markdown
## Group 2: State Management

### Task 2.1: Set up S3 bucket for state storage
- [x] **Status**: Completed
- **Files**: `infrastructure/modules/oopo-state-backend/main.tf`
- **Completed**: 2024-12-13

### Task 2.2: Set up DynamoDB table for state locking
- [x] **Status**: Completed
- **Files**: `infrastructure/modules/oopo-state-backend/dynamodb.tf`
- **Completed**: 2024-12-13
```

**Using spec-workflow MCP:**

```javascript
// Use log-implementation tool to record completion
await mcp.spec_workflow.log_implementation({
  specName: 'feature-name',
  taskId: '2.1',
  summary: 'Implemented S3 bucket with versioning and encryption',
  filesModified: ['infrastructure/modules/oopo-state-backend/main.tf'],
  filesCreated: [],
  statistics: { linesAdded: 150, linesRemoved: 0 },
  artifacts: {
    // Document what was created
  }
});
```

**Manual Update Process:**

1. Read the current `tasks.md` file
2. Find the task section that was implemented
3. Change `[ ]` to `[x]` in the status line
4. Add completion date
5. Verify file paths match implementation

### Step 6: Update CHANGELOG

Add entry for the feature:

```markdown
## [Unreleased]

### Added
- **Resource Management API** - CRUD operations for resources
  - `POST /api/v1/resources` - Create resource
  - `GET /api/v1/resources/:id` - Get resource by ID
  - `PUT /api/v1/resources/:id` - Update resource
  - `DELETE /api/v1/resources/:id` - Delete resource
- Resource validation with Zod schemas
- DynamoDB repository with single-table design
```

## Subagent Delegation

Spawn `documentation-generator` agent:

```
Use Task tool with subagent_type='documentation-generator'
Provide:
- Source files to document
- API endpoints implemented
- Design specs for reference
```

## Documentation Standards

| Type | Format | Location |
|------|--------|----------|
| API Spec | OpenAPI 3.1 | `docs/openapi.yaml` |
| Code Docs | TSDoc/GoDoc | Inline in source |
| Directory Docs | Markdown | `CLAUDE.md` per directory |
| Change Log | Keep a Changelog | `CHANGELOG.md` |

## Quality Checks

```bash
# Validate OpenAPI spec
npx @redocly/cli lint docs/openapi.yaml

# Check TypeScript docs
npx typedoc --validation

# Generate Go docs
go doc -all ./...
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| OpenAPI spec | `docs/openapi.yaml` | API documentation |
| API reference | `docs/api/` | Generated HTML docs |
| CLAUDE.md files | Per directory | AI/dev context |
| Lessons learned | `.claude/docs/*-lessons.md` | Troubleshooting patterns |
| CHANGELOG entry | `CHANGELOG.md` | Release notes |

## CLAUDE.md Update Requirements

**CRITICAL**: CLAUDE.md updates are a mandatory part of the documentation phase.

| Scope | Action | When |
|-------|--------|------|
| **Directory-level** | Create `CLAUDE.md` in new directories | New modules/features added |
| **Directory-level** | Update existing `CLAUDE.md` | Functions/files added or modified |
| **Project-level** | Update root `CLAUDE.md` | New agents, skills, workflows, or architecture changes |

Use `/update-claudemd` command for git-based analysis of project-level changes.

## Handoff to Deployment

After documentation complete:
1. OpenAPI spec generated and valid
2. Code documented with TSDoc/GoDoc
3. **CLAUDE.md files created/updated** for all affected directories
4. Root CLAUDE.md updated if project structure changed
5. **Lessons learned captured** in appropriate `.claude/docs/*-lessons.md` files
6. CHANGELOG updated
7. **READY FOR DEPLOYMENT** - Feature complete!

## SDLC Complete

```
✅ 1. SPEC      - Requirements and design documented
✅ 2. TEST      - Tests written (TDD)
✅ 3. CODE      - Implementation complete
✅ 4. BUILD     - Quality gates passed
✅ 5. SECURITY  - Security audit passed
✅ 6. DOCS      - Documentation generated

🚀 Ready for PR and deployment!
```
