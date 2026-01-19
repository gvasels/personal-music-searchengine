# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- Initial project setup
- SDLC workflow plugins (6 phases)
- Specialized agents for development tasks
- Reusable skills for code review and documentation
- Slash commands for common operations
- MCP server configuration
- Project documentation (CLAUDE.md, README.md)

---

## Template Information

This project was created from the Claude Code Starter Project template.

### Included Components

**Plugins (SDLC Workflow)**
- `spec-writer` - Requirements and technical design
- `test-writer` - TDD test creation
- `code-implementer` - Implementation
- `builder` - Build verification
- `security-checker` - Security audit
- `docs-generator` - Documentation generation

**Agents**
- `implementation-agent` - Feature implementation
- `test-engineer` - Test strategy and TDD
- `code-review` - Code quality analysis
- `security-auditor` - Security vulnerability analysis
- `documentation-generator` - Documentation generation

**Skills**
- `code-reviewer` - Code review capabilities
- `documentation-generator` - Documentation capabilities

**Commands**
- `/sdlc` - Full SDLC workflow
- `/update-claudemd` - Update CLAUDE.md from git
- `/code-review` - Code review
- `/test-file` - Generate tests for a file
