# Infrastructure Shared - CLAUDE.md

## Overview

Shared services used by multiple components: authentication (Cognito), database (DynamoDB), and media storage (S3). These resources are referenced by backend and frontend layers.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.tf` | Provider configuration, remote state references |
| `cognito.tf` | Cognito User Pool and App Client |
| `dynamodb.tf` | Single-table DynamoDB with GSIs |
| `s3.tf` | Media bucket with Intelligent-Tiering |

## Resources Created

### Authentication (`cognito.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_cognito_user_pool` | `music-library-prod` | User authentication |
| `aws_cognito_user_pool_client` | `music-library-prod-web` | Frontend app client |
| `aws_cognito_identity_pool` | `music-library-prod` | Federated identity (optional) |

### Database (`dynamodb.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_dynamodb_table` | `MusicLibrary` | Single-table design |

**Table Design:**
- Primary Key: `PK` (Partition), `SK` (Sort)
- GSI1: `GSI1PK`, `GSI1SK` (for artist/tag queries)
- Billing: Pay-per-request
- Encryption: Server-side (AES-256)
- Point-in-time recovery: Enabled
- TTL: `ExpiresAt` attribute

### Storage (`s3.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_s3_bucket` | `music-library-prod-media` | Audio files and cover art |
| `aws_s3_bucket` | `music-library-prod-search-index` | Nixiesearch indexes |

**S3 Configuration:**
- Storage class: Intelligent-Tiering
- Lifecycle: Archive after 90 days, Deep Archive after 180 days
- CORS: Configured for frontend uploads
- Versioning: Disabled (cost optimization)

## Outputs

| Output | Description |
|--------|-------------|
| `cognito_user_pool_id` | Cognito User Pool ID |
| `cognito_user_pool_arn` | Cognito User Pool ARN |
| `cognito_client_id` | App Client ID for frontend |
| `dynamodb_table_name` | DynamoDB table name |
| `dynamodb_table_arn` | DynamoDB table ARN |
| `media_bucket_name` | S3 media bucket name |
| `media_bucket_arn` | S3 media bucket ARN |
| `search_index_bucket_name` | S3 search index bucket name |

## Deployment

```bash
cd infrastructure/shared

# Ensure global layer is deployed first
tofu init
tofu plan
tofu apply
```

## Remote State Reference

Backend layer references these outputs:
```hcl
data "terraform_remote_state" "shared" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "shared/terraform.tfstate"
    region = "us-east-1"
  }
}

# Usage
locals {
  dynamodb_table_name = data.terraform_remote_state.shared.outputs.dynamodb_table_name
  cognito_user_pool_id = data.terraform_remote_state.shared.outputs.cognito_user_pool_id
}
```

## DynamoDB Access Patterns

See `dynamodb.tf` comments for detailed access patterns:

1. Get user profile: `PK = USER#{userId}, SK = PROFILE`
2. List tracks: `PK = USER#{userId}, SK begins_with TRACK#`
3. List by artist: GSI1 `GSI1PK = USER#{userId}#ARTIST#{artist}`
4. List by tag: GSI1 `GSI1PK = USER#{userId}#TAG#{tagName}`
5. List uploads by status: GSI1 `GSI1PK = UPLOAD#STATUS#{status}`
