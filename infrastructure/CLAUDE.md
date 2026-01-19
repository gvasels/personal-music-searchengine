# Infrastructure - CLAUDE.md

## Overview

OpenTofu (Terraform-compatible) infrastructure as code for the Personal Music Search Engine. Organized into layers for separation of concerns and independent deployment.

## Directory Structure

```
infrastructure/
├── global/         # Foundation: state bucket, ECR, base IAM
├── shared/         # Shared services: Cognito, DynamoDB, S3 media
├── backend/        # Backend: API Gateway, Lambdas, Step Functions
└── frontend/       # Frontend: S3 static hosting, CloudFront
```

## Layer Dependencies

```
global (Wave 1) ──► shared (Wave 1) ──► backend (Wave 2) ──► frontend (Wave 5)
     │                    │                   │
     └─ State bucket      └─ Cognito          └─ API Gateway
     └─ ECR repos         └─ DynamoDB         └─ Lambda functions
     └─ Base IAM          └─ S3 media         └─ Step Functions
```

## Deployment Order

1. **global/** - Must be deployed first (creates state bucket)
2. **shared/** - Deploy after global (uses remote state)
3. **backend/** - Deploy after shared (references Cognito, DynamoDB, S3)
4. **frontend/** - Deploy last (references API Gateway, CloudFront for media)

## Common Commands

```bash
# Initialize (first time only)
cd infrastructure/global && tofu init

# Plan changes
tofu plan

# Apply changes
tofu apply

# Destroy (be careful!)
tofu destroy
```

## State Management

- **Backend**: S3 with DynamoDB locking
- **Bucket**: `music-library-prod-tofu-state`
- **Lock Table**: `music-library-prod-tofu-lock`
- **Region**: us-east-1

Each layer has its own state file:
- `global/terraform.tfstate`
- `shared/terraform.tfstate`
- `backend/terraform.tfstate`
- `frontend/terraform.tfstate`

## AWS Profile

All infrastructure uses AWS profile `gvasels-muza`:

```bash
export AWS_PROFILE=gvasels-muza
```

Or configure in provider:
```hcl
provider "aws" {
  region  = "us-east-1"
  profile = "gvasels-muza"
}
```

## Naming Convention

Resources follow the pattern: `{project}-{environment}-{resource}`

Example: `music-library-prod-api-gateway`

## Tags

All resources are tagged with:
```hcl
default_tags {
  tags = {
    Project     = "personal-music-searchengine"
    Environment = var.environment
    ManagedBy   = "opentofu"
  }
}
```
