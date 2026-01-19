# Infrastructure Deployer Skill

Initialize and manage OpenTofu deployments with proper S3 backend configuration.

## When to Use This Skill

**Invoke this skill when:**
1. **First-time deployment** - Setting up infrastructure in a new directory
2. **State backend setup** - Configuring S3 remote state for OpenTofu
3. **Bootstrap procedure** - Following the two-step bootstrap process for new deployments

**Do NOT use when:**
- Running routine `tofu plan` or `tofu apply` (backend already configured)
- Modifying existing infrastructure (just run tofu commands directly)
- Working with modules that don't have their own state (they inherit from root)

## Platform State Architecture

The project uses **centralized state management** in a single AWS account:

| Component | Resource | Description |
|-----------|----------|-------------|
| State Bucket | `music-library-prod-tofu-state` | S3 bucket with versioning |
| Lock Table | `music-library-prod-tofu-lock` | DynamoDB table for state locking |
| Region | `us-east-1` | Primary deployment region |
| Encryption | AES-256 | Server-side encryption enabled |

## Infrastructure Directories

| Directory | Purpose |
|-----------|---------|
| `infrastructure/global/` | State bucket, ECR repos, base IAM roles |
| `infrastructure/shared/` | Cognito, DynamoDB, S3 media bucket |
| `infrastructure/backend/` | API Gateway, Lambda, Step Functions |
| `infrastructure/frontend/` | CloudFront, frontend S3 bucket |

## Usage Modes

### 1. Standard Initialization

For directories that already have remote state configured:

```bash
cd infrastructure/shared
tofu init
tofu plan
tofu apply
```

### 2. Bootstrap Mode (Two-Step Process)

For **first-time deployments** that need to create the state bucket itself:

**Step 1: Initialize locally and create state resources**
```bash
cd infrastructure/global
tofu init  # Uses local state initially
tofu apply  # Creates the state bucket and lock table
```

**Step 2: Migrate to S3 backend**
```bash
# Add backend configuration to main.tf
tofu init -migrate-state
```

### 3. Configuring Remote State

After the global infrastructure is deployed, add this backend block to other directories:

```hcl
terraform {
  backend "s3" {
    bucket         = "music-library-prod-tofu-state"
    key            = "infrastructure/{module}/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "music-library-prod-tofu-lock"
    encrypt        = true
  }
}
```

Replace `{module}` with the directory name (e.g., `shared`, `backend`, `frontend`).

## State Key Convention

State keys follow the pattern:
```
infrastructure/{directory}/terraform.tfstate
```

Examples:
- `infrastructure/global/terraform.tfstate`
- `infrastructure/shared/terraform.tfstate`
- `infrastructure/backend/terraform.tfstate`
- `infrastructure/frontend/terraform.tfstate`

## Deployment Order

Infrastructure must be deployed in dependency order:

```
1. global/     ──► Creates state bucket, ECR, base IAM
2. shared/     ──► Creates Cognito, DynamoDB, S3 media
3. backend/    ──► Creates API Gateway, Lambda, Step Functions
4. frontend/   ──► Creates CloudFront, S3 frontend bucket
```

## Common Commands

```bash
# Validate configuration
tofu validate

# Preview changes
tofu plan

# Apply changes
tofu apply

# Destroy resources (use with caution)
tofu destroy

# Format files
tofu fmt

# Show current state
tofu show
```

## Troubleshooting

### "Error: Failed to get existing workspaces"
- State bucket doesn't exist yet -> Use bootstrap mode (Step 1)

### "Error: Error acquiring the state lock"
- Another process holds the lock -> Wait or check DynamoDB for stale locks
- Force unlock (dangerous): `tofu force-unlock <lock-id>`

### "Error: Access Denied"
- Check AWS credentials: `aws sts get-caller-identity`
- Verify AWS_PROFILE environment variable
- Check IAM permissions for S3 and DynamoDB

### "Provider version conflict"
- Update provider lock file: `tofu init -upgrade`

## Environment Variables

```bash
# AWS credentials
export AWS_PROFILE=gvasels-muza
export AWS_REGION=us-east-1

# Optional: Enable detailed logging
export TF_LOG=DEBUG
```

## Related Documentation

- `infrastructure/CLAUDE.md` - Infrastructure overview
- `infrastructure/global/CLAUDE.md` - Global resources documentation
- `infrastructure/shared/CLAUDE.md` - Shared services documentation
