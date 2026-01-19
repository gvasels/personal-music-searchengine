---
name: spec-writer
description: Requirements gathering and technical design specification
phase: 1-specification
skills: []
agents: []
mcp_servers: [spec-workflow]
---

# Spec Writer Plugin

Orchestrates the specification phase of the SDLC, producing requirements and design documents before any code is written.

## Phase Position

```
[1. SPEC] → 2. TEST → 3. CODE → 4. BUILD → 5. SECURITY → 6. DOCS
   ▲
   YOU ARE HERE
```

## Workflow

### Step 1: Requirements Gathering

Use `spec-workflow` MCP to create requirements document:

```
1. Analyze user request / GitHub issue
2. Identify stakeholders and acceptance criteria
3. Define user stories (As a... I want... So that...)
4. Document non-functional requirements (performance, security, scalability)
5. Create requirements.md via spec-workflow
```

### Step 2: Technical Design

Produce design artifacts:

```
1. Data models (TypeScript interfaces, Go structs)
2. API contracts (OpenAPI spec with request/response schemas)
3. Architecture decisions (component diagram, data flow)
4. Integration points (external services, databases)
5. Create design.md via spec-workflow
```

### Step 3: Task Breakdown

Generate implementation tasks:

```
1. Break design into discrete, testable tasks
2. Identify dependencies between tasks
3. Estimate complexity (not time)
4. Assign to SDLC phases (which tasks need tests first, etc.)
5. Create tasks.md via spec-workflow
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| `requirements.md` | `.specs/{feature}/` | User stories, acceptance criteria |
| `design.md` | `.specs/{feature}/` | Data models, API contracts |
| `tasks.md` | `.specs/{feature}/` | Implementation breakdown |

## Data Model Template

```typescript
// Define before writing any code
interface EntityName {
  id: string;                    // UUID v4
  requiredField: string;         // Description of field
  optionalField?: string;        // Why optional
  status: 'active' | 'inactive'; // Enum values
  createdAt: ISO8601Timestamp;
  updatedAt: ISO8601Timestamp;
}
```

## API Contract Template

```yaml
POST /api/{resource}:
  summary: Brief description
  requestBody:
    required: true
    content:
      application/json:
        schema:
          type: object
          required: [field1, field2]
          properties:
            field1:
              type: string
              minLength: 3
              maxLength: 50
  responses:
    201:
      description: Created successfully
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/EntityName'
    400:
      description: Validation error
    409:
      description: Conflict (duplicate)
```

## Handoff to Next Phase

After spec approval:
1. Requirements reviewed and approved
2. Design documents complete with data models
3. Tasks broken down and ready for test-first development
4. **NEXT**: Pass to `test-writer` plugin for TDD test creation

## MCP Integration

```javascript
// Use spec-workflow MCP server
await mcp.spec_workflow.create_spec({
  name: 'feature-name',
  phase: 'requirements'
});

await mcp.spec_workflow.approvals({
  action: 'request',
  category: 'spec',
  categoryName: 'feature-name',
  filePath: '.specs/feature-name/requirements.md',
  title: 'Requirements for Feature Name',
  type: 'document'
});
```
