# Infrastructure Backend - CLAUDE.md

## Overview

Backend infrastructure: API Gateway with Cognito authorizer, Lambda functions, Step Functions for upload processing, MediaConvert for HLS transcoding, CloudFront for media streaming, and EventBridge for async events. Pure serverless architecture with no VPC.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.tf` | Provider configuration, remote state references, outputs |
| `step-functions.tf` | Upload processor state machine with transcode step |
| `api-gateway.tf` | HTTP API with Cognito authorizer |
| `lambda-api.tf` | Main API Lambda function |
| `lambda-processors.tf` | Step Functions processor Lambdas |
| `lambda-nixiesearch.tf` | Nixiesearch search engine Lambda (container image) |
| `mediaconvert.tf` | MediaConvert queue, IAM, and transcode Lambdas |
| `cloudfront.tf` | CloudFront distribution with signed URLs |
| `eventbridge.tf` | EventBridge rules for MediaConvert and scheduled tasks |

## Resources Created

### Step Functions (`step-functions.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_sfn_state_machine` | `music-library-prod-upload-processor` | Upload processing workflow |
| `aws_iam_role` | `music-library-prod-step-functions` | Step Functions execution role |
| `aws_cloudwatch_log_group` | `/aws/vendedlogs/states/...` | Execution logs |

### API Gateway (`api-gateway.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_apigatewayv2_api` | `music-library-prod` | HTTP API |
| `aws_apigatewayv2_authorizer` | `cognito` | Cognito JWT authorizer |
| `aws_apigatewayv2_stage` | `$default` | Default stage |

### Lambda Functions
| Lambda | File | Purpose |
|--------|------|---------|
| `api` | `lambda-api.tf` | Main API handler (Echo) |
| `metadata-extractor` | `lambda-processors.tf` | Extract audio metadata |
| `cover-art-processor` | `lambda-processors.tf` | Extract and store cover art |
| `track-creator` | `lambda-processors.tf` | Create track in DynamoDB |
| `file-mover` | `lambda-processors.tf` | Move file to media storage |
| `search-indexer` | `lambda-processors.tf` | Index track in Nixiesearch |
| `upload-status-updater` | `lambda-processors.tf` | Update upload status |
| `nixiesearch` | `lambda-nixiesearch.tf` | Embedded search engine (container) |
| `transcode-start` | `mediaconvert.tf` | Start MediaConvert HLS job |
| `transcode-complete` | `mediaconvert.tf` | Handle transcode completion |
| `index-rebuild` | `eventbridge.tf` | Daily search index rebuild |

### MediaConvert (`mediaconvert.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_media_convert_queue` | `music-library-prod-transcoding` | On-demand transcoding queue |
| `aws_iam_role` | `music-library-prod-mediaconvert` | MediaConvert job role |

### CloudFront (`cloudfront.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_cloudfront_distribution` | `music-library-prod-media-cdn` | Media streaming CDN |
| `aws_cloudfront_origin_access_control` | `music-library-prod-media-oac` | S3 origin security |
| `aws_cloudfront_public_key` | `music-library-prod-signing-key` | URL signing key |
| `aws_cloudfront_key_group` | `music-library-prod-signing-group` | Key group for signed URLs |
| `aws_secretsmanager_secret` | `music-library-prod/cloudfront-signing-key` | Private key storage |

### EventBridge (`eventbridge.tf`)
| Resource | Purpose |
|----------|---------|
| `mediaconvert-complete` | Trigger Lambda on transcode success |
| `mediaconvert-error` | Trigger Lambda on transcode failure |
| `daily-index-rebuild` | 3 AM UTC daily index rebuild |

## Step Functions Workflow

```
StartAt: ExtractMetadata
         ↓
ExtractMetadata → ProcessCoverArt → CreateTrackRecord → MoveToMediaStorage
         ↓              ↓                  ↓                    ↓
    [On Error] → MarkUploadFailed     [On Error]          [On Error]
                                           ↓                    ↓
                                    StartTranscode → IndexForSearch → MarkUploadCompleted
                                           ↓                  ↓
                                    [On Error] →         [On Error] → Continue
                                    Continue             (non-critical)
```

**Note**: StartTranscode is async - the actual transcode completion is handled by EventBridge triggering `transcode-complete` Lambda.

## Outputs

| Output | Description |
|--------|-------------|
| `step_functions_arn` | Upload processor state machine ARN |
| `api_lambda_arn` | Main API Lambda ARN |
| `api_gateway_url` | API Gateway invoke URL |
| `nixiesearch_lambda_arn` | Nixiesearch Lambda ARN |
| `cloudfront_domain_name` | CloudFront domain for media streaming |
| `cloudfront_key_pair_id` | Key pair ID for signed URLs |
| `mediaconvert_queue_arn` | MediaConvert queue ARN |

## Deployment

```bash
cd infrastructure/backend

# Ensure global and shared layers are deployed first
tofu init
tofu plan
tofu apply
```

## Architecture Notes

### Pure Serverless (No VPC)
- All Lambdas run in AWS public network
- Nixiesearch uses S3 for index storage (no EFS)
- No NAT Gateway or VPC endpoints required
- Reduces operational complexity and costs

### CloudFront Signed URLs
- All media content (HLS, downloads, cover art) requires signed URLs
- Private key stored in Secrets Manager
- 24-hour URL expiration for convenience
- API Lambda reads key on demand

### HLS Adaptive Streaming
- MediaConvert creates 3 quality levels (96k, 192k, 320k AAC)
- HLS master playlist at `/hls/{userId}/{trackId}/master.m3u8`
- Fallback to original file if HLS not ready

## Remote State References

```hcl
locals {
  dynamodb_table_name        = data.terraform_remote_state.shared.outputs.dynamodb_table_name
  media_bucket_name          = data.terraform_remote_state.shared.outputs.media_bucket_name
  search_indexes_bucket_name = data.terraform_remote_state.shared.outputs.search_indexes_bucket_name
  cognito_user_pool_id       = data.terraform_remote_state.shared.outputs.cognito_user_pool_id
  lambda_role_arn            = data.terraform_remote_state.global.outputs.lambda_execution_role_arn
}
```
