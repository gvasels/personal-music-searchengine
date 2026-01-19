# Claude Code Starter Project

A ready-to-use project template with a complete SDLC (Software Development Lifecycle) workflow powered by Claude Code automation.

## What's Included

This starter project provides:

- **7-Phase SDLC Workflow** - Structured development from spec to verified deployment
- **SDLC Plugins** - Orchestrate each development phase
- **Specialized Agents** - Subagents for implementation, testing, security, and docs
- **Reusable Skills** - Code review and documentation generation capabilities
- **Slash Commands** - Quick access to common workflows
- **MCP Server Configuration** - Pre-configured external tools

## Quick Start

### 1. Clone or Copy This Template

```bash
# Copy to your new project
cp -r claude-code-starter-project my-new-project
cd my-new-project

# Initialize git
git init
```

### 2. Customize CLAUDE.md

Open `CLAUDE.md` and update:
- Project name and description
- Technology stack tables
- Any project-specific guidelines

### 3. Set Up Environment Variables

This project uses [direnv](https://direnv.net/) to manage environment variables for MCP servers.

```bash
# Install direnv (macOS)
brew install direnv

# Add to your shell (add to ~/.zshrc or ~/.bashrc)
eval "$(direnv hook zsh)"  # or bash

# Copy the example file and fill in your values
cp .envrc.example .envrc

# Allow direnv to load the file
direnv allow
```

Edit `.envrc` with your credentials:

```bash
# Required
export GITHUB_TOKEN="your-github-personal-access-token"

# AWS (for AWS MCP servers)
export AWS_PROFILE="default"
export AWS_REGION="us-east-1"

# Optional (for Google Workspace integration)
export GOOGLE_OAUTH_CLIENT_ID="your-client-id"
export GOOGLE_OAUTH_CLIENT_SECRET="your-client-secret"
```

**Security Note**: `.envrc` is gitignored - never commit secrets to version control.

### 4. Set Up GitHub @claude Integration (Optional)

To enable `@claude` triggers in GitHub issues and PRs:

1. **Get a Claude Code OAuth Token**:
   - Sign up for [Claude Code](https://claude.com/claude-code) if you haven't
   - Generate an OAuth token from your account settings

2. **Add the secret to your repository**:
   - Go to your GitHub repository
   - Navigate to **Settings** → **Secrets and variables** → **Actions**
   - Click **New repository secret**
   - Name: `CLAUDE_CODE_OAUTH_TOKEN`
   - Value: Your Claude Code OAuth token

3. **Test the integration**:
   - Create a new issue with `@claude` in the title or body
   - The GitHub Action will trigger and Claude will respond

The workflows are pre-configured in `.github/workflows/`:
- `claude.yml` - Triggers on `@claude` mentions in issues/comments
- `claude-code-review.yml` - Automatic PR code reviews

### 5. Start Building

Run the SDLC workflow:

```bash
# In Claude Code
/sdlc
```

Or start with individual phases:
- "Run the spec phase for [feature]"
- "Write tests for this design"
- "Implement this feature"
- "Verify deployment"

## Project Structure

```
├── CLAUDE.md                    # Project documentation for AI/humans
├── README.md                    # This file
├── CHANGELOG.md                 # Project changelog
├── .mcp.json                    # MCP server configuration
├── .envrc.example               # Environment variable template
├── .gitignore                   # Git ignore patterns
├── .github/
│   └── workflows/
│       ├── claude.yml           # @claude trigger workflow
│       └── claude-code-review.yml # PR auto-review workflow
│
├── .claude/
│   ├── plugins/                 # SDLC workflow plugins
│   │   ├── spec-writer.md       # Phase 1: Requirements & design
│   │   ├── test-writer.md       # Phase 2: TDD test creation
│   │   ├── code-implementer.md  # Phase 3: Implementation
│   │   ├── builder.md           # Phase 4: Build verification
│   │   ├── security-checker.md  # Phase 5: Security audit
│   │   ├── docs-generator.md    # Phase 6: Documentation
│   │   └── deploy-verifier.md   # Phase 7: Post-deploy verification
│   │
│   ├── agents/                  # Specialized subagents
│   │   ├── implementation-agent.md
│   │   ├── test-engineer.md
│   │   ├── code-review.md
│   │   ├── security-auditor.md
│   │   └── documentation-generator.md
│   │
│   ├── skills/                  # Reusable capabilities
│   │   ├── code-reviewer/
│   │   │   ├── code-reviewer.md
│   │   │   ├── scripts/         # Analysis scripts
│   │   │   │   ├── analyze-metrics.py
│   │   │   │   └── compare-complexity.py
│   │   │   └── templates/       # Review templates
│   │   │       ├── finding-template.md
│   │   │       └── review-checklist.md
│   │   └── documentation-generator.md
│   │
│   ├── commands/                # Slash commands
│   │   ├── sdlc.md              # /sdlc - Full workflow
│   │   ├── update-claudemd.md   # /update-claudemd
│   │   ├── code-review.md       # /code-review
│   │   └── test-file.md         # /test-file
│   │
│   └── docs/                    # Lessons learned & patterns
│       ├── lessons-learned-template.md  # Template for adding lessons
│       └── common-patterns.md   # Cross-technology patterns
│
├── src/                         # Your source code (create as needed)
├── tests/                       # Your tests (create as needed)
└── docs/                        # Documentation (create as needed)
```

## SDLC Workflow

The workflow follows 7 phases, executed sequentially or individually:

```
┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         ┌─────────┐
│ 1.SPEC  │──►│ 2.TEST  │──►│ 3.CODE  │──►│ 4.BUILD │──►│5.SECURE │──►│ 6.DOCS  │──Deploy─►│7.VERIFY │
└─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘   └─────────┘         └─────────┘
```

### Phase 1: Specification
- Gather requirements and user stories
- Create technical design with data models
- Define API contracts
- Break down into implementation tasks

**Artifacts**: `requirements.md`, `design.md`, `tasks.md`

### Phase 2: Testing (TDD)
- Write failing unit tests from data models
- Write failing integration tests from API contracts
- Verify all tests fail (Red phase)

**Artifacts**: Test files in `tests/`

### Phase 3: Implementation
- Implement data models and validation
- Implement business logic
- Create API endpoints
- Make all tests pass (Green phase)

**Artifacts**: Source files in `src/`

### Phase 4: Build Verification
- Lint check (0 errors)
- Type check (0 errors)
- Test coverage (80%+)
- Build artifacts

### Phase 5: Security Audit
- Dependency vulnerability scan
- Secrets detection
- SAST analysis
- Infrastructure security (if applicable)

**Gate**: 0 critical/high vulnerabilities

### Phase 6: Documentation
- Generate/update OpenAPI spec
- Add code documentation (TSDoc/GoDoc)
- Create/update CLAUDE.md files
- Update CHANGELOG

### Phase 7: Post-Deployment Verification
- Health check endpoints
- Contract validation against OpenAPI spec
- Smoke tests on critical API paths
- Automatic rollback on failure

**Gate**: All smoke tests pass, schemas match spec

**On Failure**: Automatic rollback + incident notification

## Available Commands

| Command | Description |
|---------|-------------|
| `/sdlc` | Start full SDLC workflow |
| `/update-claudemd` | Update CLAUDE.md from git changes |
| `/code-review` | Run comprehensive code review |
| `/test-file <path>` | Generate tests for a file |

## Available Agents

| Agent | When to Use |
|-------|-------------|
| `implementation-agent` | Implementing features |
| `test-engineer` | Writing tests, TDD |
| `code-review` | Code quality analysis |
| `security-auditor` | Security vulnerability checks |
| `documentation-generator` | Generating docs |

Spawn agents with:
```
Use Task tool with subagent_type='agent-name'
```

## MCP Servers

Pre-configured servers in `.mcp.json`:

| Server | Purpose | Requires |
|--------|---------|----------|
| `spec-workflow` | Specification management and approvals | - |
| `context7` | Documentation lookup | - |
| `github` | GitHub integration | `GITHUB_TOKEN` |
| `docs-mcp-server` | Documentation indexing | - |
| `playwright` | Browser automation and testing | - |
| `godoc` | Go documentation lookup | `godoc-mcp` installed |

**AWS MCP Servers** (require `AWS_PROFILE` or AWS credentials):

| Server | Purpose |
|--------|---------|
| `awslabs.aws-documentation-mcp-server` | AWS documentation search |
| `aws-knowledge-mcp-server` | AWS knowledge base |
| `awslabs.terraform-mcp-server` | OpenTofu/Terraform for AWS |
| `awslabs.dynamodb-mcp-server` | DynamoDB data modeling |
| `awslabs.syntheticdata-mcp-server` | Synthetic data generation |
| `awslabs.aws-api-mcp-server` | AWS API access |
| `awslabs.finch-mcp-server` | Finch container operations |
| `well-architected-security-mcp-server` | AWS Well-Architected security |

**Optional Servers**:

| Server | Purpose | Requires |
|--------|---------|----------|
| `google_workspace` | Google Workspace integration | OAuth credentials |

## Customization

### Adding More MCP Servers

Update `.mcp.json` to add additional MCP servers for your stack. See the [MCP Registry](https://github.com/modelcontextprotocol/servers) for available servers.

### Adding Custom Commands

Create new commands in `.claude/commands/`:

```markdown
---
description: Your command description
---

# Command Name

Instructions for what the command should do.
```

### Adding Custom Agents

Create new agents in `.claude/agents/`:

```markdown
---
name: my-agent
description: What this agent does
tools: read, write, bash
---

# My Agent

## Purpose
...

## Workflow
...
```

## Best Practices

1. **Always start with a spec** for non-trivial features
2. **Write tests before code** (TDD)
3. **Update CLAUDE.md** when adding new modules
4. **Run security checks** before deployment
5. **Keep CHANGELOG updated** with changes

## Documentation Standards

### CLAUDE.md Files

Every major directory should have a `CLAUDE.md` with:
- Overview of the directory's purpose
- File descriptions
- Key functions with signatures
- Dependencies

### CHANGELOG

Follow [Keep a Changelog](https://keepachangelog.com/):

```markdown
## [Unreleased]

### Added
- New features

### Changed
- Changes to existing functionality

### Fixed
- Bug fixes
```

### Lessons Learned

Capture troubleshooting patterns in `.claude/docs/`:

- **Add technology-specific files**: `go-lessons.md`, `aws-lessons.md`, `typescript-lessons.md`
- **Use the template**: See `lessons-learned-template.md` for the entry format
- **Include actual error messages**: Helps Claude find relevant solutions
- **Add debugging steps**: Help others investigate similar issues

These files are read by Claude Code when troubleshooting, providing project-specific knowledge beyond general documentation.

## License

This starter template is provided as-is for use in your projects.

---

Ready to start building? Run `/sdlc` in Claude Code!
