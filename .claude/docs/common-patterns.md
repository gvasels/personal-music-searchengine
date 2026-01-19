# Common Patterns and Lessons

Cross-technology troubleshooting patterns and solutions.

## Git & Version Control

### Problem: Merge conflicts in generated files

**Symptom**:
```
CONFLICT (content): Merge conflict in package-lock.json
```

**Solution**: Regenerate the file instead of manual merge:

```bash
# For package-lock.json
git checkout --ours package-lock.json
npm install

# For go.sum
git checkout --ours go.sum
go mod tidy
```

---

### Problem: Accidentally committed sensitive data

**Symptom**: Secrets, API keys, or credentials in commit history.

**Solution**: Use git-filter-repo (NOT filter-branch):

```bash
# Install
pip3 install git-filter-repo

# Remove file from all history
git filter-repo --invert-paths --path secrets.json

# Force push (coordinate with team)
git push --force-with-lease
```

**Prevention**:
- Use `.gitignore` for sensitive files
- Add pre-commit hooks with gitleaks
- Use environment variables for secrets

---

## Environment & Configuration

### Problem: "Command not found" after installation

**Symptom**:
```
command not found: tofu
```

**Solution**: Check PATH and shell configuration:

```bash
# Check if installed
which tofu || ls /usr/local/bin/tofu

# If exists but not in PATH, add to shell config
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

---

### Problem: Wrong environment variables loaded

**Symptom**: Application uses wrong configuration (staging vs prod).

**Debugging**:
```bash
# Check current environment
printenv | grep -E "(ENV|NODE_ENV|AWS_PROFILE)"

# Check .env file loading order
cat .env .env.local .env.development 2>/dev/null | grep KEY_NAME
```

**Solution**: Ensure proper .env file precedence and verify with logging.

---

## API & HTTP

### Problem: CORS errors from browser

**Symptom**:
```
Access to fetch blocked by CORS policy: No 'Access-Control-Allow-Origin' header
```

**Root Cause**: Backend doesn't include CORS headers, or headers are incorrect.

**Solution**: Configure CORS at the API/gateway level:

```typescript
// Express
app.use(cors({
  origin: 'https://myapp.com',
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  credentials: true
}));
```

```hcl
# API Gateway (Terraform)
cors_configuration {
  allow_origins = ["https://myapp.com"]
  allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allow_headers = ["Content-Type", "Authorization"]
}
```

**Note**: For preflight requests, ensure OPTIONS method returns 200 with headers.

---

### Problem: API returns 200 but response body is empty

**Symptom**: Successful HTTP status but no data in response.

**Debugging**:
```bash
# Check full response including headers
curl -v https://api.example.com/endpoint

# Check Content-Type header
curl -I https://api.example.com/endpoint
```

**Common causes**:
- Response not serialized (missing `return` or `json()`)
- Content-Type mismatch (expecting JSON, getting HTML)
- Empty query result returned as empty array/object

---

## Debugging

### Problem: Log messages not appearing

**Symptom**: Print/log statements don't show in output.

**Common causes and solutions**:

1. **Buffered output**:
   ```python
   # Python - flush buffer
   print("debug", flush=True)
   ```

2. **Wrong log level**:
   ```javascript
   // Check log level
   console.log("Level:", process.env.LOG_LEVEL);
   ```

3. **Logging to wrong stream**:
   ```bash
   # Check both stdout and stderr
   command 2>&1 | tee output.log
   ```

---

### Problem: Tests pass locally but fail in CI

**Symptom**: Green locally, red in CI pipeline.

**Debugging checklist**:

1. **Environment differences**:
   ```yaml
   # CI config - match local versions
   node-version: '20.x'
   go-version: '1.22'
   ```

2. **Time-dependent tests**:
   ```javascript
   // Use fixed time in tests
   jest.useFakeTimers().setSystemTime(new Date('2024-01-01'));
   ```

3. **Race conditions**: Tests may execute in different order
4. **Missing dependencies**: Check CI installs all required packages
5. **File system**: Case sensitivity differs between OS

---

## Performance

### Problem: Build times increasing significantly

**Symptom**: CI/local builds taking much longer than before.

**Investigation**:
```bash
# Measure build phases
time npm run lint
time npm run build
time npm test
```

**Common solutions**:
- Enable incremental builds (`tsconfig.json`: `"incremental": true`)
- Cache node_modules in CI
- Use turbo/nx for monorepo caching
- Parallelize independent steps

---

## Security

### Problem: Dependency vulnerability alerts

**Symptom**: GitHub Dependabot or npm audit warnings.

**Triage approach**:

1. **Check exploitability**:
   - Is the vulnerable code path actually used?
   - Is it only in devDependencies?

2. **Resolution options**:
   ```bash
   # Update to fixed version
   npm update vulnerable-package

   # If transitive dependency
   npm update parent-package

   # Override if needed (package.json)
   "overrides": {
     "vulnerable-package": "^2.0.0"
   }
   ```

3. **Accept risk** (document decision) if:
   - No fix available
   - Impact is minimal
   - Upgrade breaks functionality

---

## Documentation

### Problem: Documentation out of sync with code

**Prevention**:
- Add CLAUDE.md checks to PR workflow
- Auto-generate API docs from code comments
- Include doc updates in Definition of Done
- Use documentation-generator agent after code changes

---

## Add Your Lessons

As you encounter and solve problems, add entries following the template in `lessons-learned-template.md`.
