# Deployment Guide

## Prerequisites

- **AWS CLI** v2.x configured with profile `gvasels-muza`
- **OpenTofu** 1.8.0+ (`tofu --version`)
- **Node.js** 20.x (`node --version`)
- **Go** 1.22+ (`go version`)
- **Docker** for backend Lambda builds

### AWS Profile Setup

```bash
# Configure AWS profile
aws configure --profile gvasels-muza

# Set as default
export AWS_PROFILE=gvasels-muza

# Verify access
aws sts get-caller-identity
```

## Infrastructure Deployment Order

Infrastructure must be deployed in order due to dependencies:

```
1. global → 2. shared → 3. backend → 4. frontend
```

### Manual Deployment

#### 1. Global Module (State, ECR, IAM)

```bash
cd infrastructure/global
tofu init
tofu plan -out=tfplan
tofu apply tfplan
```

**Creates:**
- S3 state bucket
- DynamoDB lock table
- ECR repositories
- Base IAM roles
- GitHub OIDC provider

#### 2. Shared Module (Cognito, DynamoDB, S3)

```bash
cd infrastructure/shared
tofu init
tofu plan -out=tfplan
tofu apply tfplan
```

**Creates:**
- Cognito User Pool
- DynamoDB table
- S3 media bucket

#### 3. Backend Module (API, Lambda, Step Functions)

```bash
cd infrastructure/backend
tofu init
tofu plan -out=tfplan
tofu apply tfplan
```

**Creates:**
- API Gateway
- Lambda functions
- Step Functions state machine
- CloudFront (media)

#### 4. Frontend Module (S3, CloudFront)

```bash
cd infrastructure/frontend
tofu init
tofu plan -out=tfplan
tofu apply tfplan
```

**Creates:**
- S3 frontend bucket
- CloudFront distribution

## Automated Deployment (GitHub Actions)

### Setup

1. Add repository secrets in GitHub:
   - `AWS_OIDC_ROLE_ARN`: Get from `tofu output github_actions_role_arn` in global module
   - `VITE_API_URL`: API Gateway URL from backend module
   - `VITE_COGNITO_USER_POOL_ID`: From shared module
   - `VITE_COGNITO_CLIENT_ID`: From shared module

2. Merge to `main` branch triggers:
   - Infrastructure deployment (all modules)
   - Frontend build and S3 sync
   - CloudFront cache invalidation

### Workflow Files

- `.github/workflows/ci.yml` - Runs on PRs (tests, lint, validation)
- `.github/workflows/deploy.yml` - Runs on merge to main (deploy)

## Frontend Deployment

### Manual Frontend Deploy

```bash
cd frontend

# Install dependencies
npm ci

# Build for production
npm run build

# Sync to S3 (assets with long cache)
aws s3 sync dist/ s3://music-library-prod-frontend \
  --delete \
  --cache-control "public, max-age=31536000, immutable" \
  --exclude "index.html" \
  --exclude "*.json"

# Upload index.html with no-cache
aws s3 cp dist/index.html s3://music-library-prod-frontend/index.html \
  --cache-control "no-cache, no-store, must-revalidate"

# Invalidate CloudFront
DIST_ID=$(cd ../infrastructure/frontend && tofu output -raw frontend_cloudfront_distribution_id)
aws cloudfront create-invalidation --distribution-id $DIST_ID --paths "/*"
```

## Backend Deployment

### Manual Backend Deploy

```bash
cd backend

# Build Docker image
docker build -t music-library-prod-api:latest -f cmd/api/Dockerfile .

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 887395463840.dkr.ecr.us-east-1.amazonaws.com

# Tag and push
docker tag music-library-prod-api:latest \
  887395463840.dkr.ecr.us-east-1.amazonaws.com/music-library-prod-api:latest
docker push 887395463840.dkr.ecr.us-east-1.amazonaws.com/music-library-prod-api:latest

# Update Lambda
aws lambda update-function-code \
  --function-name music-library-prod-api \
  --image-uri 887395463840.dkr.ecr.us-east-1.amazonaws.com/music-library-prod-api:latest
```

## Rollback Procedures

### Frontend Rollback (S3 Versioning)

```bash
# List object versions
aws s3api list-object-versions \
  --bucket music-library-prod-frontend \
  --prefix index.html

# Restore previous version
aws s3api copy-object \
  --bucket music-library-prod-frontend \
  --copy-source music-library-prod-frontend/index.html?versionId=PREVIOUS_VERSION_ID \
  --key index.html

# Invalidate cache
aws cloudfront create-invalidation \
  --distribution-id DIST_ID \
  --paths "/index.html"
```

### Infrastructure Rollback (OpenTofu State)

```bash
cd infrastructure/frontend

# List state versions
aws s3api list-object-versions \
  --bucket music-library-prod-tofu-state \
  --prefix frontend/terraform.tfstate

# Download previous state
aws s3api get-object \
  --bucket music-library-prod-tofu-state \
  --key frontend/terraform.tfstate \
  --version-id PREVIOUS_VERSION_ID \
  terraform.tfstate.backup

# Review and apply
tofu plan -state=terraform.tfstate.backup
```

### Lambda Rollback

```bash
# List Lambda versions/aliases
aws lambda list-versions-by-function \
  --function-name music-library-prod-api

# Update to previous version
aws lambda update-function-code \
  --function-name music-library-prod-api \
  --image-uri 887395463840.dkr.ecr.us-east-1.amazonaws.com/music-library-prod-api:PREVIOUS_SHA
```

## Troubleshooting

### Common Issues

#### OIDC Authentication Fails

```
Error: Could not assume role with OIDC
```

**Solution:** Verify GitHub repository matches IAM trust policy:
```bash
# Check trust policy
aws iam get-role --role-name music-library-prod-github-actions \
  --query 'Role.AssumeRolePolicyDocument'
```

#### S3 Access Denied

```
Error: Access Denied when syncing to S3
```

**Solution:** Verify bucket policy and IAM permissions:
```bash
aws s3api get-bucket-policy --bucket music-library-prod-frontend
```

#### CloudFront 403 Error

```
Error: CloudFront returns 403 Forbidden
```

**Solution:** Check OAC configuration:
```bash
# Verify OAC exists
aws cloudfront list-origin-access-controls

# Check distribution origin
aws cloudfront get-distribution --id DIST_ID \
  --query 'Distribution.DistributionConfig.Origins'
```

#### OpenTofu State Lock

```
Error: Error acquiring the state lock
```

**Solution:** Check DynamoDB lock and force unlock if needed:
```bash
# List locks
aws dynamodb scan --table-name music-library-prod-tofu-lock

# Force unlock (use with caution)
tofu force-unlock LOCK_ID
```

## Monitoring

### CloudFront Metrics

```bash
# View distribution metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/CloudFront \
  --metric-name Requests \
  --dimensions Name=DistributionId,Value=DIST_ID \
  --start-time $(date -u -v-1H +%Y-%m-%dT%H:%M:%SZ) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --period 300 \
  --statistics Sum
```

### Lambda Logs

```bash
# View recent logs
aws logs tail /aws/lambda/music-library-prod-api --follow
```
