# Lessons Learned Template

This directory contains troubleshooting patterns and solutions discovered during development. These files help AI assistants and future developers avoid repeating mistakes and quickly find solutions to common problems.

## Purpose

- **Knowledge capture**: Document solutions to non-obvious problems
- **AI context**: Provide Claude with project-specific troubleshooting knowledge
- **Onboarding**: Help new team members learn from past experiences

## File Naming Convention

Create files by technology or domain:

| File | Content |
|------|---------|
| `go-lessons.md` | Go language patterns, struct tags, testing |
| `typescript-lessons.md` | TypeScript patterns, type issues, build configs |
| `aws-lessons.md` | AWS service gotchas, IAM debugging, SDK patterns |
| `opentofu-lessons.md` | IaC patterns, state management, provider issues |
| `react-lessons.md` | React patterns, hooks, state management |
| `database-lessons.md` | Database-specific patterns and queries |

## Entry Format

Each entry should follow this structure:

```markdown
---

### Problem: [Brief description]

**Symptom**:
```
[Error message or unexpected behavior]
```

**Root Cause**: [Why this happens]

**Solution**: [How to fix it]

```code
[Example code or commands]
```

**Debugging**: [How to investigate similar issues]

---
```

## Example Entry

```markdown
---

### Problem: Empty fields when unmarshaling JSON

**Symptom**:
```
Struct fields are nil/empty after json.Unmarshal() even though data exists
```

**Root Cause**: JSON field names don't match Go struct tags (case sensitivity).

**Solution**: Use json tags that match the JSON field names exactly:

```go
// WRONG - JSON has lowercase "name"
type User struct {
    Name string `json:"Name"`
}

// CORRECT - matches JSON exactly
type User struct {
    Name string `json:"name"`
}
```

**Debugging**: Print raw JSON and compare with struct tags:
```go
fmt.Printf("Raw JSON: %s\n", jsonBytes)
```

---
```

## Best Practices

1. **Be specific**: Include actual error messages, not paraphrased versions
2. **Include context**: Explain why the problem occurs, not just how to fix it
3. **Add debugging steps**: Help others investigate similar issues
4. **Keep entries focused**: One problem per entry, use sections for variations
5. **Update when learning**: Add new entries as you discover solutions
6. **Reference docs**: Link to official documentation when relevant

## When to Add Entries

Add a lessons learned entry when:

- You spend significant time debugging an issue
- The solution isn't obvious from error messages
- The problem is likely to recur
- A workaround or pattern is project-specific
- Official documentation doesn't cover the scenario

## Integration with Claude

Claude Code reads these files when troubleshooting. To maximize effectiveness:

- Use clear problem descriptions (Claude searches by keywords)
- Include exact error messages (helps with pattern matching)
- Document solutions with working code examples
- Group related issues in the same file
