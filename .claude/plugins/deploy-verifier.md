---
name: deploy-verifier
description: Post-deployment API verification with automatic rollback
phase: 7-verify
skills: []
agents: []
mcp_servers: []
---

# Deploy Verifier Plugin

Verifies API Gateway + Lambda deployments are functioning correctly and triggers rollback on failure.

## Phase Position

```
1. SPEC → 2. TEST → 3. CODE → 4. BUILD → 5. SECURITY → 6. DOCS → Deploy → [7. VERIFY]
                                                                              ▲
                                                                              YOU ARE HERE
```

## Prerequisites

From previous phases:
- Deployed API Gateway + Lambda infrastructure
- Test definitions from TDD phase (Phase 2)
- OpenAPI spec from docs phase (Phase 6)
- Deployment metadata (stack name, environment, version)

## Verification Pipeline

```
┌─────────────────────────────────────────┐
│         EXTRACT TEST CASES              │
│   From TDD specs + OpenAPI contract     │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         HEALTH CHECK                    │
│   Verify endpoints are reachable        │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         CONTRACT VALIDATION             │
│   Request/response schema validation    │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         SMOKE TESTS                     │
│   Critical path API calls               │
└─────────────────────────────────────────┘
                    │
                    ▼
        ┌───────────────────┐
        │   ALL PASSED?     │
        └───────────────────┘
           │           │
          YES          NO
           │           │
           ▼           ▼
     ┌─────────┐  ┌─────────────┐
     │ SUCCESS │  │  ROLLBACK   │
     └─────────┘  └─────────────┘
```

## Workflow

### Step 1: Extract Test Cases from TDD

Reuse existing test expectations from Phase 2:

```typescript
// From test files, extract API contract expectations
interface DeploymentTest {
  endpoint: string;           // e.g., "POST /api/v1/resources"
  description: string;        // From test description
  request: {
    method: string;
    path: string;
    headers?: Record<string, string>;
    body?: unknown;
  };
  expectedResponse: {
    status: number;
    bodySchema?: object;      // JSON Schema for validation
    bodyContains?: string[];  // Required fields/values
  };
}
```

**Extract from test files:**

```typescript
// Original TDD test (Phase 2)
describe('POST /api/v1/resources', () => {
  it('creates a resource with valid input', async () => {
    const input = { name: 'my-resource', type: 'standard' };
    const response = await api.post('/api/v1/resources', input);
    expect(response.status).toBe(201);
    expect(response.body).toHaveProperty('id');
  });
});

// Extracted deployment test
const deploymentTest: DeploymentTest = {
  endpoint: "POST /api/v1/resources",
  description: "creates a resource with valid input",
  request: {
    method: "POST",
    path: "/api/v1/resources",
    body: { name: "smoke-test-resource", type: "standard" }
  },
  expectedResponse: {
    status: 201,
    bodyContains: ["id"]
  }
};
```

### Step 2: Health Check

Verify API Gateway endpoints are reachable:

```bash
#!/bin/bash
# health-check.sh

API_BASE_URL="${API_BASE_URL:-https://api.example.com}"
HEALTH_ENDPOINT="${HEALTH_ENDPOINT:-/health}"
MAX_RETRIES=5
RETRY_DELAY=10

for i in $(seq 1 $MAX_RETRIES); do
  response=$(curl -s -o /dev/null -w "%{http_code}" "${API_BASE_URL}${HEALTH_ENDPOINT}")

  if [ "$response" = "200" ]; then
    echo "✅ Health check passed"
    exit 0
  fi

  echo "⏳ Attempt $i/$MAX_RETRIES - Status: $response"
  sleep $RETRY_DELAY
done

echo "❌ Health check failed after $MAX_RETRIES attempts"
exit 1
```

### Step 3: Contract Validation

Validate responses match OpenAPI spec:

```typescript
// verify-contract.ts
import Ajv from 'ajv';
import { OpenAPIV3 } from 'openapi-types';

async function verifyContract(
  apiBaseUrl: string,
  openApiSpec: OpenAPIV3.Document,
  test: DeploymentTest
): Promise<VerificationResult> {
  const ajv = new Ajv();

  // Make actual API call
  const response = await fetch(`${apiBaseUrl}${test.request.path}`, {
    method: test.request.method,
    headers: {
      'Content-Type': 'application/json',
      ...test.request.headers
    },
    body: test.request.body ? JSON.stringify(test.request.body) : undefined
  });

  // Verify status code
  if (response.status !== test.expectedResponse.status) {
    return {
      passed: false,
      error: `Expected status ${test.expectedResponse.status}, got ${response.status}`
    };
  }

  // Verify response body schema
  if (test.expectedResponse.bodySchema) {
    const body = await response.json();
    const valid = ajv.validate(test.expectedResponse.bodySchema, body);
    if (!valid) {
      return {
        passed: false,
        error: `Schema validation failed: ${ajv.errorsText()}`
      };
    }
  }

  return { passed: true };
}
```

### Step 4: Run Smoke Tests

Execute critical path tests:

```typescript
// smoke-test.ts
interface SmokeTestConfig {
  apiBaseUrl: string;
  tests: DeploymentTest[];
  cleanupAfter: boolean;
}

async function runSmokeTests(config: SmokeTestConfig): Promise<SmokeTestResult> {
  const results: TestResult[] = [];
  const createdResources: string[] = [];

  for (const test of config.tests) {
    console.log(`🧪 Testing: ${test.description}`);

    try {
      const result = await executeTest(config.apiBaseUrl, test);
      results.push(result);

      // Track created resources for cleanup
      if (result.passed && test.request.method === 'POST' && result.resourceId) {
        createdResources.push(result.resourceId);
      }

      if (result.passed) {
        console.log(`  ✅ Passed`);
      } else {
        console.log(`  ❌ Failed: ${result.error}`);
      }
    } catch (error) {
      results.push({
        test: test.endpoint,
        passed: false,
        error: error.message
      });
      console.log(`  ❌ Error: ${error.message}`);
    }
  }

  // Cleanup test data
  if (config.cleanupAfter) {
    await cleanupResources(config.apiBaseUrl, createdResources);
  }

  const passed = results.every(r => r.passed);
  return {
    passed,
    total: results.length,
    succeeded: results.filter(r => r.passed).length,
    failed: results.filter(r => !r.passed).length,
    results
  };
}
```

### Step 5: Rollback on Failure

Trigger automatic rollback if verification fails:

```yaml
# buildkite pipeline step
- label: "🔍 Post-Deploy Verification"
  command: |
    npm run verify:deployment
  env:
    API_BASE_URL: "${DEPLOYED_API_URL}"
    ENVIRONMENT: "${ENVIRONMENT}"

- label: "⏪ Rollback on Failure"
  command: |
    echo "Verification failed - initiating rollback"

    # OpenTofu rollback to previous state
    cd infrastructure
    tofu workspace select ${ENVIRONMENT}

    # Get previous version from state
    PREVIOUS_VERSION=$(aws lambda get-function \
      --function-name ${FUNCTION_NAME} \
      --query 'Configuration.Version' \
      --output text)

    # Update alias to previous version
    aws lambda update-alias \
      --function-name ${FUNCTION_NAME} \
      --name ${ENVIRONMENT} \
      --function-version ${PREVIOUS_VERSION}

    echo "✅ Rollback complete to version ${PREVIOUS_VERSION}"
  depends_on: "post-deploy-verification"
  if: build.state == "failed"
```

**AWS CodeBuild rollback:**

```yaml
# buildspec-verify.yml
version: 0.2

phases:
  build:
    commands:
      - echo "Running post-deployment verification..."
      - npm run verify:deployment

  post_build:
    commands:
      - |
        if [ "$CODEBUILD_BUILD_SUCCEEDING" = "0" ]; then
          echo "Verification failed - triggering rollback"
          aws codepipeline put-approval-result \
            --pipeline-name ${PIPELINE_NAME} \
            --stage-name Deploy \
            --action-name Approval \
            --result summary="Verification failed",status="Rejected"
        fi
```

## Verification Test File Format

Create `verify.config.json` in the service root:

```json
{
  "apiBaseUrl": "${API_BASE_URL}",
  "healthEndpoint": "/health",
  "healthTimeout": 60,
  "tests": [
    {
      "name": "Create Resource",
      "endpoint": "POST /api/v1/resources",
      "request": {
        "body": {
          "name": "smoke-test-${TIMESTAMP}",
          "type": "standard"
        }
      },
      "expected": {
        "status": 201,
        "bodyContains": ["id", "name", "createdAt"]
      },
      "cleanup": true
    },
    {
      "name": "List Resources",
      "endpoint": "GET /api/v1/resources",
      "expected": {
        "status": 200,
        "bodyContains": ["items", "total"]
      }
    },
    {
      "name": "Get Nonexistent Resource",
      "endpoint": "GET /api/v1/resources/nonexistent-id",
      "expected": {
        "status": 404
      }
    }
  ],
  "rollback": {
    "enabled": true,
    "strategy": "lambda-alias",
    "notifyOnRollback": ["platform-admin@oopo.io"]
  }
}
```

## Integration with CI/CD

### Buildkite Pipeline

```yaml
steps:
  # ... previous deployment steps ...

  - wait

  - label: "🔍 Verify Deployment"
    command: ".buildkite/scripts/verify-deployment.sh"
    env:
      ENVIRONMENT: "${ENVIRONMENT}"
    plugins:
      - artifacts#v1.9.0:
          upload: "verification-report.json"
    retry:
      automatic:
        - exit_status: 1
          limit: 2

  - label: "⏪ Rollback"
    command: ".buildkite/scripts/rollback.sh"
    if: build.state == "failed"
    env:
      ENVIRONMENT: "${ENVIRONMENT}"
```

### Environment-Specific Configuration

```bash
# Environment variables per stage
# dev.env
API_BASE_URL=https://dev-api.oopo.io
SKIP_CLEANUP=true
VERBOSE=true

# staging.env
API_BASE_URL=https://staging-api.oopo.io
SKIP_CLEANUP=false
VERBOSE=true

# prod.env
API_BASE_URL=https://api.oopo.io
SKIP_CLEANUP=false
VERBOSE=false
NOTIFY_ON_FAILURE=true
```

## Outputs

| Artifact | Location | Purpose |
|----------|----------|---------|
| Verification report | `verification-report.json` | Detailed test results |
| Console output | CI/CD logs | Real-time progress |
| Rollback log | `rollback.log` | Rollback audit trail |

## Rollback Strategies

| Strategy | Use Case | Speed |
|----------|----------|-------|
| **Lambda Alias** | Single function updates | Fast (~seconds) |
| **API Gateway Stage** | API configuration changes | Fast (~seconds) |
| **OpenTofu State** | Infrastructure changes | Medium (~minutes) |
| **Blue-Green** | Full stack deployment | Fast (DNS switch) |

## SDLC Complete with Verification

```
✅ 1. SPEC      - Requirements and design documented
✅ 2. TEST      - Tests written (TDD)
✅ 3. CODE      - Implementation complete
✅ 4. BUILD     - Quality gates passed
✅ 5. SECURITY  - Security audit passed
✅ 6. DOCS      - Documentation generated
✅ 7. VERIFY    - Post-deployment validation passed

🚀 Deployment verified and stable!
```

## Failure Handling

```
┌─────────────────────────────────────────┐
│         VERIFICATION FAILED             │
└─────────────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
  ┌───────────┐          ┌───────────┐
  │ ROLLBACK  │          │  NOTIFY   │
  │ Previous  │          │  Team     │
  │ Version   │          │           │
  └───────────┘          └───────────┘
        │                       │
        └───────────┬───────────┘
                    ▼
        ┌───────────────────┐
        │  CREATE INCIDENT  │
        │  GitHub Issue     │
        └───────────────────┘
```
