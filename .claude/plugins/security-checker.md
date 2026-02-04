---
name: security-checker
description: Security vulnerability scanning and audit
phase: 5-security
skills: [code-reviewer]
agents: [security-auditor]
mcp_servers: []
---

# Security Checker Plugin

Performs comprehensive security audits before code is deployed, catching vulnerabilities early.

## Phase Position

```
1. SPEC → 2. TEST → 3. CODE → 4. BUILD → [5. SECURITY] → 6. DOCS
                                              ▲
                                              YOU ARE HERE
```

## Prerequisites

From previous phases:
- Built and tested code (from builder)
- All quality gates passed

## Security Scan Pipeline

```
┌─────────────────────────────────────────┐
│         DEPENDENCY SCAN                  │
│   npm audit, govulncheck, safety        │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│          SECRETS SCAN                    │
│   gitleaks, trufflehog                  │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│           SAST SCAN                      │
│   semgrep, CodeQL, gosec                │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│       INFRASTRUCTURE SCAN                │
│   checkov, tfsec (for OpenTofu)         │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         MANUAL REVIEW                    │
│   security-auditor agent                │
└─────────────────────────────────────────┘
```

## Workflow

### Step 1: Dependency Vulnerability Scan

```bash
# TypeScript/JavaScript
npm audit --audit-level=high
# or
npx better-npm-audit audit

# Go
govulncheck ./...

# Python
pip-audit
# or
safety check -r requirements.txt
```

**Severity Thresholds:**

| Severity | Action |
|----------|--------|
| Critical | Block deployment, fix immediately |
| High | Block deployment, fix required |
| Medium | Warning, fix before next release |
| Low | Informational, track in backlog |

### Step 2: Secrets Detection

```bash
# Scan for hardcoded secrets
gitleaks detect --source . --verbose

# Alternative
trufflehog filesystem . --only-verified
```

**Common Patterns Detected:**

| Pattern | Example |
|---------|---------|
| AWS Keys | `AKIA[0-9A-Z]{16}` |
| Private Keys | `-----BEGIN RSA PRIVATE KEY-----` |
| API Tokens | `ghp_`, `sk-`, `Bearer` |
| Passwords | `password\s*=\s*['"]` |

### Step 3: Static Application Security Testing (SAST)

```bash
# Multi-language SAST
semgrep --config=auto .

# Go specific
gosec ./...

# TypeScript/JavaScript specific
npx eslint . --ext .ts,.tsx --config .eslintrc.security.js
```

**OWASP Top 10 Checks:**

| # | Vulnerability | Check |
|---|---------------|-------|
| A01 | Broken Access Control | Authorization checks |
| A02 | Cryptographic Failures | Encryption usage |
| A03 | Injection | Input validation |
| A04 | Insecure Design | Architecture review |
| A05 | Security Misconfiguration | Config audit |
| A06 | Vulnerable Components | Dependency scan |
| A07 | Auth Failures | Session management |
| A08 | Data Integrity Failures | Signature verification |
| A09 | Logging Failures | Audit logging |
| A10 | SSRF | URL validation |

### Step 4: Infrastructure Security Scan

```bash
# OpenTofu/Terraform
checkov -d infrastructure/ --framework terraform

# Alternative
tfsec infrastructure/
```

**AWS Security Checks:**

| Resource | Check |
|----------|-------|
| S3 | Encryption, public access block |
| IAM | Least privilege, no wildcards |
| Lambda | VPC config, env var encryption |
| API Gateway | Authorization, throttling |
| DynamoDB | Encryption, backup enabled |

### Step 5: Manual Security Review

Spawn `security-auditor` agent for deep analysis:

```
Use Task tool with subagent_type='security-auditor'
Provide:
- List of files changed
- Authentication/authorization flows
- Data handling patterns
- External integrations
```

## Security Findings Format

```markdown
### Finding: SQL Injection Vulnerability

- **Severity**: Critical
- **CWE**: CWE-89
- **Location**: src/repositories/user.ts:45
- **Finding**: User input directly concatenated into query
- **Risk**: Database compromise, data exfiltration
- **Remediation**:
  ```typescript
  // Before (vulnerable)
  const query = `SELECT * FROM users WHERE id = '${userId}'`;

  // After (secure)
  const query = 'SELECT * FROM users WHERE id = ?';
  const result = await db.execute(query, [userId]);
  ```
- **References**:
  - https://owasp.org/Top10/A03_2021-Injection/
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| Dependency report | `security/dependencies.json` | CVE findings |
| Secrets report | `security/secrets.json` | Leaked credentials |
| SAST report | `security/sast.json` | Code vulnerabilities |
| IaC report | `security/infrastructure.json` | Config issues |
| Audit summary | `security/AUDIT.md` | Human-readable summary |

## Quality Gates

| Check | Threshold |
|-------|-----------|
| Critical vulnerabilities | 0 |
| High vulnerabilities | 0 |
| Hardcoded secrets | 0 |
| SAST critical findings | 0 |
| IaC critical misconfigs | 0 |

## Failure Handling

If security scan fails:

1. **Critical/High CVE** → Update dependency, return to builder
2. **Secrets detected** → Remove secret, rotate credential, update history
3. **SAST finding** → Return to code-implementer, fix vulnerability
4. **IaC misconfiguration** → Update infrastructure code

## Handoff to Next Phase

After security approval:
1. No critical/high vulnerabilities
2. No secrets in codebase
3. SAST findings addressed
4. Infrastructure secure
5. **NEXT**: Pass to `docs-generator` plugin for documentation
