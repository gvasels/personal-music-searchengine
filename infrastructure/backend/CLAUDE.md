# Infrastructure Backend - CLAUDE.md

## Overview

Backend infrastructure: API Gateway with Cognito authorizer, Lambda functions, and Step Functions for upload processing. References shared resources via remote state.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.tf` | Provider configuration, remote state references, outputs |
| `step-functions.tf` | Upload processor state machine and IAM |
| `api-gateway.tf` | HTTP API with Cognito authorizer (TODO) |
| `lambda-api.tf` | Main API Lambda function (TODO) |
| `lambda-processors.tf` | Step Functions processor Lambdas (TODO) |

## Resources Created

### Step Functions (`step-functions.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_sfn_state_machine` | `music-library-prod-upload-processor` | Upload processing workflow |
| `aws_iam_role` | `music-library-prod-step-functions` | Step Functions execution role |
| `aws_cloudwatch_log_group` | `/aws/vendedlogs/states/...` | Execution logs |

### API Gateway (`api-gateway.tf`) - TODO
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_apigatewayv2_api` | `music-library-prod` | HTTP API |
| `aws_apigatewayv2_authorizer` | `cognito` | Cognito JWT authorizer |
| `aws_apigatewayv2_stage` | `$default` | Default stage |

### Lambda Functions (TODO)
| Lambda | Purpose |
|--------|---------|
| `api` | Main API handler (Echo) |
| `metadata-extractor` | Extract audio metadata |
| `cover-art-processor` | Extract and store cover art |
| `track-creator` | Create track in DynamoDB |
| `file-mover` | Move file to media storage |
| `search-indexer` | Index track in Nixiesearch |
| `upload-status-updater` | Update upload status |

## Step Functions Workflow

```
StartAt: ExtractMetadata
         ↓
ExtractMetadata → ProcessCoverArt → CreateTrackRecord → MoveToMediaStorage
         ↓              ↓                  ↓                    ↓
    [On Error] → MarkUploadFailed     [On Error]          [On Error]
                                           ↓                    ↓
                                    IndexForSearch → MarkUploadCompleted
                                           ↓
                                    [On Error] → Continue (non-critical)
```

## Outputs

| Output | Description |
|--------|-------------|
| `step_functions_arn` | Upload processor state machine ARN |
| `api_lambda_arn` | Main API Lambda ARN |
| `api_gateway_url` | API Gateway invoke URL |

## Deployment

```bash
cd infrastructure/backend

# Ensure shared layer is deployed first
tofu init
tofu plan
tofu apply
```

## API Gateway Cognito Authorizer

**IMPORTANT**: Use API Gateway's native Cognito JWT authorizer, NOT custom middleware.

```hcl
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito"

  jwt_configuration {
    audience = [data.terraform_remote_state.shared.outputs.cognito_client_id]
    issuer   = "https://cognito-idp.${var.aws_region}.amazonaws.com/${data.terraform_remote_state.shared.outputs.cognito_user_pool_id}"
  }
}
```

## Remote State References

```hcl
data "terraform_remote_state" "shared" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "shared/terraform.tfstate"
    region = "us-east-1"
  }
}

data "terraform_remote_state" "global" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "global/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  dynamodb_table_name  = data.terraform_remote_state.shared.outputs.dynamodb_table_name
  media_bucket_name    = data.terraform_remote_state.shared.outputs.media_bucket_name
  cognito_user_pool_id = data.terraform_remote_state.shared.outputs.cognito_user_pool_id
  lambda_role_arn      = data.terraform_remote_state.global.outputs.lambda_execution_role_arn
}
```
