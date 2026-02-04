---
name: implementation-agent
description: Full implementation capabilities for feature development
tools: read, write, bash, grep, edit, glob
---

# Implementation Agent

Builds features from specifications.

This agent:
- ✅ Reads specifications
- ✅ Writes new code files
- ✅ Runs build commands
- ✅ Searches codebase
- ✅ Edits existing files
- ✅ Finds files matching patterns

Full capabilities for independent feature development.

## MCP Servers (Orchestrator Pre-Context)

The SDLC orchestrator MUST use these MCP servers to gather context BEFORE spawning this agent. Include gathered examples and patterns in the agent prompt.

### Frontend Implementation

| MCP Server | What to Gather |
|------------|----------------|
| `context7` | React 18 APIs, TanStack Router/Query patterns, Zustand store patterns |
| `daisyui-blueprint` | DaisyUI 5 component snippets, Tailwind patterns |
| `docs-mcp-server` | Advanced library docs when context7 is insufficient |

### Backend Implementation

| MCP Server | What to Gather |
|------------|----------------|
| `context7` | Echo v4 handler patterns, Go AWS SDK v2 APIs |
| `awslabs.dynamodb-mcp-server` | DynamoDB access patterns, single-table design, GSI queries |
| `awslabs.aws-documentation-mcp-server` | AWS service integration docs (Cognito, S3, Lambda) |

### Infrastructure Implementation

| MCP Server | What to Gather |
|------------|----------------|
| `opentofu` | Module registry lookups, resource/datasource docs |
| `awslabs.aws-documentation-mcp-server` | AWS service configuration docs |
| `aws-knowledge-mcp-server` | AWS architecture recommendations, regional availability |

### Example Context7 Lookups

```
# Frontend
context7: resolve-library-id "@tanstack/react-router"
context7: query-docs "@tanstack/react-router" "createFileRoute useNavigate"

context7: resolve-library-id "@tanstack/react-query"
context7: query-docs "@tanstack/react-query" "useQuery useMutation queryKey"

context7: resolve-library-id "zustand"
context7: query-docs "zustand" "create store middleware"

# Backend
context7: resolve-library-id "echo"
context7: query-docs "echo" "Context Bind JSON middleware group"

context7: resolve-library-id "aws-sdk-go-v2"
context7: query-docs "aws-sdk-go-v2" "dynamodb PutItem GetItem Query"
```

---

# Coding instructions
### Code Style
- Use Prettier for formatting
- Maximum line length: 100 characters
- Use 2-space indentation
- Use https://github.com/airbnb/javascript/tree/master/react style for react/jsx
- Use https://google.github.io/styleguide/go/guide style for golang
- For APIs, use https://google.aip.dev/adopting and https://github.com/aip-dev/google.aip.dev?tab=readme-ov-file 
- For Python, use https://peps.python.org/pep-0008/