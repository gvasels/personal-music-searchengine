# AWS Lessons Learned

Troubleshooting patterns and solutions for AWS services.

## API Gateway HTTP API

### Problem: 404 errors with named stages

**Symptom**: Requests to `https://api-id.execute-api.region.amazonaws.com/prod/health` return 404.

**Root Cause**: API Gateway HTTP API with named stages includes the stage name in the path sent to Lambda. The request path becomes `/prod/health`, not `/health`.

**Solution**: Handle stage prefix in your application routing (see go-lessons.md).

**Alternative**: Use the `$default` stage which doesn't add a prefix.

---

### Problem: CORS errors from browser

**Symptom**: Browser shows CORS error even though Lambda sets headers.

**Solution**: Configure CORS at API Gateway level, not just Lambda:

```hcl
resource "aws_apigatewayv2_api" "api" {
  name          = "my-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["https://myapp.com"]
    allow_methods = ["GET", "POST", "OPTIONS"]
    allow_headers = ["Content-Type", "Authorization"]
    max_age       = 300
  }
}
```

---

## DynamoDB Single-Table Design

### Problem: Query returns empty results despite data existing

**Debugging steps**:

1. **Verify partition key format**:
```bash
aws dynamodb get-item \
  --table-name PlatformManifests \
  --key '{"PK": {"S": "ENV#prod"}, "SK": {"S": "ROUTE#product#service#mfe"}}'
```

2. **Scan to see all items** (dev only):
```bash
aws dynamodb scan --table-name PlatformManifests --limit 10
```

3. **Check attribute names** - DynamoDB is case-sensitive:
```
PK: "ENV#prod"     # Correct
PK: "env#prod"     # Wrong - won't match
```

---

### Problem: Understanding PK/SK patterns

**Common patterns used**:

| Record Type | PK | SK |
|-------------|----|----|
| Route | `ENV#{env}` | `ROUTE#{product}#{service}#{mfe}` |
| Deployment | `SERVICE#{service}#{mfe}` | `DEPLOY#{env}#{timestamp}` |
| User Override | `USER#{userId}` | `OVERRIDE#{product}#{service}#{feature}` |
| Global Config | `GLOBAL` | `GLOBAL#{env}#{type}` |

**Query by environment**:
```go
input := &dynamodb.QueryInput{
    TableName:              aws.String("PlatformManifests"),
    KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk":       &types.AttributeValueMemberS{Value: "ENV#prod"},
        ":skPrefix": &types.AttributeValueMemberS{Value: "ROUTE#"},
    },
}
```

---

## Lambda Debugging

### Problem: Lambda function errors with no clear message

**Solution**: Check CloudWatch Logs:

```bash
# Tail logs in real-time
aws logs tail /aws/lambda/manifest-api-prod --follow

# Get recent logs
aws logs tail /aws/lambda/manifest-api-prod --since 1h
```

**Common log patterns**:
```
# Init error - check handler configuration
Init Duration: 234.56 ms

# Cold start
INIT_START Runtime Version: provided:al2023.v19

# Request/Response
{"level":"info","uri":"/prod/api/manifest","method":"GET"}
```

---

### Problem: Lambda timeout

**Symptoms**: Function times out at 3 seconds (default).

**Solutions**:
1. Increase timeout in Lambda configuration
2. Check for blocking operations (missing context cancellation)
3. Verify DynamoDB/external service connectivity
4. Check VPC configuration if Lambda is in a VPC

```hcl
resource "aws_lambda_function" "api" {
  timeout = 30  # seconds

  # If in VPC, ensure NAT gateway for internet access
  vpc_config {
    subnet_ids         = var.private_subnet_ids
    security_group_ids = [aws_security_group.lambda.id]
  }
}
```

---

## Lambda Deployment

### Problem: Lambda not picking up new code

**Solution**: Ensure you're updating the function code:

```bash
# Build
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap cmd/api/main.go

# Package
zip deployment.zip bootstrap

# Deploy
aws lambda update-function-code \
  --function-name manifest-api-prod \
  --zip-file fileb://deployment.zip
```

**Verify deployment**:
```bash
aws lambda get-function --function-name manifest-api-prod \
  --query 'Configuration.LastModified'
```

---

## Cross-Account Access

### Problem: Access denied when accessing resources in another account

**Solution**: Use assume role with proper trust policy:

```hcl
# In target account - trust policy
data "aws_iam_policy_document" "trust" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::SOURCE_ACCOUNT:root"]
    }
  }
}

# In source account - assume role
provider "aws" {
  alias  = "target"
  region = "us-east-1"

  assume_role {
    role_arn = "arn:aws:iam::TARGET_ACCOUNT:role/DeploymentRole"
  }
}
```

---

## S3 State Backend

### Problem: State lock errors

**Symptom**:
```
Error acquiring the state lock
Lock Info: ID: xxx-xxx-xxx
```

**Solutions**:

1. **Wait** - another operation may be in progress
2. **Force unlock** (use with caution):
```bash
tofu force-unlock LOCK_ID
```

3. **Check DynamoDB locks table**:
```bash
aws dynamodb scan --table-name oopo-terraform-locks
```

---

## CloudWatch Logs Insights

### Problem: Need to search across many log events

**Solution**: Use CloudWatch Logs Insights:

```sql
-- Find errors
fields @timestamp, @message
| filter @message like /error/i
| sort @timestamp desc
| limit 100

-- Find slow requests
fields @timestamp, @duration
| filter @duration > 1000
| sort @duration desc

-- Count by status code
fields @message
| parse @message '"status":*,' as status
| stats count() by status
```

---

## IAM Debugging

### Problem: Access denied but policy looks correct

**Debugging steps**:

1. **Check effective permissions**:
```bash
aws iam simulate-principal-policy \
  --policy-source-arn arn:aws:iam::ACCOUNT:role/MyRole \
  --action-names dynamodb:Query \
  --resource-arns arn:aws:dynamodb:us-east-1:ACCOUNT:table/MyTable
```

2. **Check CloudTrail** for denied events:
```bash
aws cloudtrail lookup-events \
  --lookup-attributes AttributeKey=EventName,AttributeValue=Query \
  --start-time 2024-01-01T00:00:00Z
```

3. **Common issues**:
   - Resource ARN mismatch (table vs table/*)
   - Missing region in ARN
   - Condition key not met
   - SCP blocking the action
