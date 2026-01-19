# Infrastructure Global - CLAUDE.md

## Overview

Foundation infrastructure resources that must be deployed first. Creates the OpenTofu state backend, ECR repositories, and base IAM roles used by other layers.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.tf` | Provider configuration and variables |
| `state.tf` | S3 bucket and DynamoDB table for Tofu state |
| `ecr.tf` | ECR repositories for Lambda container images |
| `iam.tf` | Base IAM roles for Lambda execution |

## Resources Created

### State Management (`state.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_s3_bucket` | `music-library-prod-tofu-state` | Terraform state storage |
| `aws_dynamodb_table` | `music-library-prod-tofu-lock` | State locking |

### Container Registry (`ecr.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_ecr_repository` | `music-library-prod-api` | API Lambda image |
| `aws_ecr_repository` | `music-library-prod-processor` | Processor Lambda image |
| `aws_ecr_repository` | `music-library-prod-indexer` | Indexer Lambda image |

### IAM (`iam.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_iam_role` | `music-library-prod-lambda-execution` | Base Lambda execution role |
| `aws_iam_policy` | `music-library-prod-lambda-base` | Basic Lambda permissions |

## Outputs

| Output | Description |
|--------|-------------|
| `state_bucket_name` | S3 bucket for Terraform state |
| `lock_table_name` | DynamoDB table for state locking |
| `ecr_api_repository_url` | ECR URL for API image |
| `ecr_processor_repository_url` | ECR URL for processor image |
| `ecr_indexer_repository_url` | ECR URL for indexer image |
| `lambda_execution_role_arn` | Base Lambda execution role ARN |

## Deployment

```bash
cd infrastructure/global

# First deployment (no remote state yet)
tofu init
tofu apply

# After state bucket exists, add backend config to main.tf:
# terraform {
#   backend "s3" {
#     bucket         = "music-library-prod-tofu-state"
#     key            = "global/terraform.tfstate"
#     region         = "us-east-1"
#     dynamodb_table = "music-library-prod-tofu-lock"
#   }
# }

# Then migrate state
tofu init -migrate-state
```

## Important Notes

- Deploy this layer FIRST before any other infrastructure
- The state bucket has versioning enabled for recovery
- ECR repositories have lifecycle policies to limit image count
- Lambda execution role is referenced by other layers via remote state
