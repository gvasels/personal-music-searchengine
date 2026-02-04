# Infrastructure Frontend - CLAUDE.md

## Overview

OpenTofu infrastructure for hosting the React SPA frontend via S3 and CloudFront. The frontend is served globally with edge caching, HTTPS enforcement, and proper SPA routing support.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.tf` | Provider configuration, remote state references, variables, outputs |
| `s3.tf` | S3 bucket for static assets with versioning, encryption, OAC policy |
| `cloudfront.tf` | CloudFront distribution with SPA routing, caching, security headers |

## Resources Created

### S3 (`s3.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_s3_bucket` | `music-library-prod-frontend` | Static asset storage |
| `aws_s3_bucket_versioning` | - | Enable versioning for rollback |
| `aws_s3_bucket_server_side_encryption_configuration` | - | AES-256 encryption |
| `aws_s3_bucket_public_access_block` | - | Block all public access |
| `aws_s3_bucket_policy` | - | Allow CloudFront OAC access |

### CloudFront (`cloudfront.tf`)
| Resource | Name | Purpose |
|----------|------|---------|
| `aws_cloudfront_distribution` | `music-library-prod-frontend-cdn` | CDN distribution |
| `aws_cloudfront_origin_access_control` | `music-library-prod-frontend-oac` | S3 origin security |
| `aws_cloudfront_response_headers_policy` | `music-library-prod-frontend-security-headers` | Security headers |

## Outputs

| Output | Description |
|--------|-------------|
| `frontend_bucket_name` | S3 bucket name for deployment |
| `frontend_bucket_arn` | S3 bucket ARN |
| `frontend_cloudfront_distribution_id` | CloudFront ID for cache invalidation |
| `frontend_cloudfront_domain_name` | CloudFront URL (*.cloudfront.net) |

## SPA Routing

CloudFront handles SPA routing with custom error responses:
- 403 Forbidden → `/index.html` (200 OK)
- 404 Not Found → `/index.html` (200 OK)

This ensures client-side routing works for all routes.

## Cache Behaviors

| Path | TTL | Purpose |
|------|-----|---------|
| `/assets/*` | 1 week - 1 year | Hashed static assets (long cache) |
| `/*` (default) | 0 | index.html (no cache for updates) |

## Security Headers

| Header | Value |
|--------|-------|
| X-Content-Type-Options | nosniff |
| X-Frame-Options | DENY |
| Strict-Transport-Security | max-age=31536000; includeSubDomains |
| X-XSS-Protection | 1; mode=block |

## Deployment

```bash
cd infrastructure/frontend

# Initialize
tofu init

# Plan changes
tofu plan

# Apply changes
tofu apply
```

## Cache Invalidation

After deploying new frontend code:

```bash
# Get distribution ID
DIST_ID=$(tofu output -raw frontend_cloudfront_distribution_id)

# Invalidate all paths
aws cloudfront create-invalidation \
  --distribution-id $DIST_ID \
  --paths "/*"
```

## Remote State References

```hcl
data "terraform_remote_state" "shared" {
  # Cognito configuration
}

data "terraform_remote_state" "backend" {
  # API Gateway URL
}
```
