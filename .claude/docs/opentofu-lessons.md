# OpenTofu/Terraform Lessons Learned

Troubleshooting patterns and solutions for Infrastructure as Code.

## State Backend Configuration

### Problem: "Backend configuration changed" error

**Symptom**:
```
Error: Backend configuration changed
```

**Solution**: Run `tofu init -reconfigure`:

```bash
tofu init -reconfigure
```

Or migrate state:
```bash
tofu init -migrate-state
```

---

### Problem: State lock stuck

**Symptom**:
```
Error acquiring the state lock
```

**Solutions**:

1. Wait for the other operation to complete
2. Check if another process is running
3. Force unlock (use with caution):
```bash
tofu force-unlock LOCK_ID
```

4. Check DynamoDB for stale locks:
```bash
aws dynamodb scan --table-name terraform-locks \
  --filter-expression "attribute_exists(LockID)"
```

---

## Module Development

### Problem: Changes to module not reflected

**Solution**: Re-initialize to pick up module changes:

```bash
tofu init -upgrade
```

For local modules, changes are picked up automatically, but you may need to run:
```bash
tofu get -update
```

---

### Problem: Module output not available

**Symptom**:
```
Error: Unsupported attribute
module.example.some_output
```

**Solution**: Ensure the output is defined in the module's `outputs.tf`:

```hcl
# modules/example/outputs.tf
output "some_output" {
  description = "Description of this output"
  value       = aws_resource.example.id
}
```

---

## Provider Configuration

### Problem: Using wrong AWS account

**Debugging**:
```bash
# Check current identity
aws sts get-caller-identity

# Check which profile is being used
echo $AWS_PROFILE
```

**Solution**: Explicitly set the profile or assume role:

```hcl
provider "aws" {
  region  = "us-east-1"
  profile = "my-profile"  # Or use assume_role

  # Or assume a role
  assume_role {
    role_arn = "arn:aws:iam::ACCOUNT:role/DeploymentRole"
  }
}
```

---

### Problem: Provider version conflicts

**Symptom**:
```
Error: Failed to query available provider packages
```

**Solution**: Pin provider versions in `versions.tf`:

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

---

## Resource Naming

### Problem: Resource name conflicts

**Solution**: Use consistent naming with project/environment prefixes:

```hcl
locals {
  name_prefix = "${var.project}-${var.environment}"
}

resource "aws_s3_bucket" "example" {
  bucket = "${local.name_prefix}-${var.bucket_suffix}"
}
```

---

### Problem: Lambda function name too long

**Constraint**: Lambda function names max 64 characters.

**Solution**: Use shorter prefixes and abbreviations:

```hcl
# Instead of: oopo-dev-platform-manifest-api-production
# Use: manifest-api-prod
resource "aws_lambda_function" "api" {
  function_name = "${var.service_name}-${var.environment}"
}
```

---

## Data Sources vs Resources

### Problem: Need to reference existing infrastructure

**Solution**: Use data sources for existing resources:

```hcl
# Reference existing VPC
data "aws_vpc" "main" {
  tags = {
    Name = "main-vpc"
  }
}

# Reference existing certificate
data "aws_acm_certificate" "api" {
  domain   = "api.example.com"
  statuses = ["ISSUED"]
}

# Use in resources
resource "aws_lb" "example" {
  subnets = data.aws_subnets.private.ids
}
```

---

## Conditional Resources

### Problem: Need to create resource only in certain conditions

**Solution**: Use `count` or `for_each`:

```hcl
# Create only if enabled
resource "aws_route53_record" "custom_domain" {
  count = var.enable_custom_domain ? 1 : 0

  zone_id = var.hosted_zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name    = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].target_domain_name
    zone_id = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].hosted_zone_id
  }
}

# Reference conditional resource
output "custom_domain" {
  value = var.enable_custom_domain ? aws_route53_record.custom_domain[0].fqdn : null
}
```

---

## Import Existing Resources

### Problem: Resource exists in AWS but not in state

**Solution**: Import the resource:

```bash
# Import format
tofu import aws_dynamodb_table.example table-name

# Import with module
tofu import module.db.aws_dynamodb_table.main PlatformManifests
```

**Generate configuration** (OpenTofu 1.6+):
```bash
tofu plan -generate-config-out=generated.tf
```

---

## Debugging Plan Output

### Problem: Unexpected changes in plan

**Debugging**:

1. **Show detailed diff**:
```bash
tofu plan -out=tfplan
tofu show -json tfplan | jq '.resource_changes[]'
```

2. **Check state**:
```bash
tofu state show aws_lambda_function.api
```

3. **Refresh state** to sync with AWS:
```bash
tofu refresh
```

---

## Handling Secrets

### Problem: Secrets in state file

**Solution**: Use AWS Secrets Manager or SSM Parameter Store:

```hcl
# Store secret in Secrets Manager
resource "aws_secretsmanager_secret" "db_password" {
  name = "${var.project}/db-password"
}

# Reference in Lambda
resource "aws_lambda_function" "api" {
  environment {
    variables = {
      DB_SECRET_ARN = aws_secretsmanager_secret.db_password.arn
    }
  }
}
```

**Never** put secrets directly in `.tf` files or variables.

---

## Module Versioning

### Problem: Breaking changes in modules

**Solution**: Use semantic versioning for modules:

```hcl
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"  # Any 5.x version
}

# For internal modules, use git tags
module "internal" {
  source = "git::https://github.com/org/modules.git//vpc?ref=v1.2.0"
}
```

---

## Workspace vs Directory Structure

### Problem: Managing multiple environments

**Recommendation**: Use directory structure over workspaces for clarity:

```
infrastructure/
├── modules/           # Reusable modules
├── accounts/
│   ├── dev/          # Dev environment
│   │   ├── main.tf
│   │   └── terraform.tfvars
│   ├── staging/      # Staging environment
│   └── prod/         # Production environment
```

Each environment has its own state file and can be deployed independently.

---

## Common Gotchas

### Lambda + API Gateway Integration

Always use `aws_lambda_permission` to allow API Gateway to invoke Lambda:

```hcl
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}
```

### DynamoDB Table Changes

Some DynamoDB changes require replacement (destroy + create):
- Changing partition key or sort key
- Changing GSI key schema

Plan carefully and consider data migration.

### S3 Bucket Naming

S3 bucket names are globally unique. Include account ID or random suffix:

```hcl
resource "aws_s3_bucket" "state" {
  bucket = "${var.project}-terraform-state-${data.aws_caller_identity.current.account_id}"
}
```
