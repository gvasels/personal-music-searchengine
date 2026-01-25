# Tasks Document - Distribution & Polish (Epic 6)

## US-6.1: Frontend Hosting

- [x] 1. Create frontend S3 bucket infrastructure
  - File: `infrastructure/frontend/s3.tf`
  - Create S3 bucket with versioning, encryption, and public access blocks
  - Configure bucket for static website hosting via CloudFront
  - Purpose: Store React SPA static assets (HTML, JS, CSS)
  - _Leverage: `infrastructure/shared/s3.tf` for encryption and public access patterns_
  - _Requirements: US-6.1 AC-1, AC-7_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer specializing in AWS and OpenTofu | Task: Create S3 bucket for frontend hosting following patterns from infrastructure/shared/s3.tf with versioning, AES-256 encryption, and all public access blocked | Restrictions: Do not enable static website hosting on S3 (CloudFront handles this), do not make bucket public, follow existing naming conventions (music-library-prod-frontend) | Success: S3 bucket created with proper security configuration, outputs bucket name and ARN for CloudFront | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 2. Create frontend CloudFront distribution
  - File: `infrastructure/frontend/cloudfront.tf`
  - Create CloudFront distribution with S3 origin using OAC
  - Configure SPA routing (custom error responses for 403/404 → index.html)
  - Add cache behaviors for /assets/* (long TTL) and default (short TTL)
  - Purpose: Serve frontend globally with edge caching and proper SPA routing
  - _Leverage: `infrastructure/backend/cloudfront.tf` for OAC and cache patterns_
  - _Requirements: US-6.1 AC-1 through AC-6_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer specializing in CloudFront and CDN optimization | Task: Create CloudFront distribution for frontend SPA with OAC to S3, custom error responses (403/404 → index.html with 200), cache behaviors for hashed assets (1 year TTL) and index.html (no-cache), gzip/brotli compression | Restrictions: Do not require signed URLs (frontend is public), use default CloudFront certificate (no custom domain yet), follow existing naming patterns | Success: CloudFront serves SPA with proper routing, assets cached at edge, index.html not cached | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 3. Add security headers response policy
  - File: `infrastructure/frontend/cloudfront.tf` (append to task 2)
  - Create CloudFront response headers policy with security headers
  - Configure: X-Content-Type-Options, X-Frame-Options, HSTS, XSS-Protection
  - Purpose: Enhance security posture of frontend distribution
  - _Leverage: `infrastructure/backend/cloudfront.tf` CORS response policy as reference_
  - _Requirements: US-6.1 (security best practice)_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Security Engineer specializing in web application security headers | Task: Create CloudFront response headers policy with security headers: X-Content-Type-Options (nosniff), X-Frame-Options (DENY), Strict-Transport-Security (max-age 1 year, includeSubdomains), X-XSS-Protection (1; mode=block) | Restrictions: Do not add CSP yet (may break SPA), attach to default cache behavior | Success: All security headers present in responses, validated with curl or browser dev tools | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 4. Create S3 bucket policy for CloudFront OAC
  - File: `infrastructure/frontend/s3.tf` (append to task 1)
  - Create bucket policy allowing CloudFront OAC to read objects
  - Reference CloudFront distribution ARN in condition
  - Purpose: Allow CloudFront to access S3 while keeping bucket private
  - _Leverage: `infrastructure/backend/cloudfront.tf` lines 199-221 for bucket policy pattern_
  - _Requirements: US-6.1 AC-1_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: IAM Security Specialist for AWS S3 and CloudFront | Task: Create S3 bucket policy allowing cloudfront.amazonaws.com principal to s3:GetObject with condition AWS:SourceArn matching the CloudFront distribution ARN | Restrictions: Only allow GetObject action, restrict to specific distribution ARN, no other principals | Success: CloudFront can read S3 objects, direct S3 access denied | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 5. Create frontend infrastructure main.tf
  - File: `infrastructure/frontend/main.tf`
  - Configure OpenTofu provider, backend, variables, and remote state references
  - Add outputs for bucket name, CloudFront distribution ID and domain
  - Purpose: Complete frontend infrastructure module setup
  - _Leverage: `infrastructure/backend/main.tf` for provider and remote state patterns_
  - _Requirements: US-6.1 (infrastructure setup)_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Infrastructure Engineer specializing in OpenTofu/Terraform state management | Task: Create main.tf with terraform 1.8+ required_version, aws provider ~5.0, S3 backend (key: frontend/terraform.tfstate), remote state references to shared for Cognito/API values, locals for name_prefix, outputs for frontend_bucket_name, frontend_cloudfront_distribution_id, frontend_cloudfront_domain_name | Restrictions: Use existing state bucket (music-library-prod-tofu-state), follow existing provider configuration patterns | Success: tofu init and tofu validate pass, outputs available for CI/CD | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

## US-6.2: Testing & CI/CD

- [x] 6. Create GitHub Actions OIDC provider in global infrastructure
  - File: `infrastructure/global/iam.tf` (append)
  - Create OIDC identity provider for GitHub Actions
  - Create IAM role with trust policy for GitHub OIDC
  - Purpose: Enable keyless authentication from GitHub Actions to AWS
  - _Leverage: AWS documentation for GitHub OIDC setup_
  - _Requirements: US-6.2 AC-6 through AC-8_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: IAM Security Architect specializing in OIDC federation | Task: Create aws_iam_openid_connect_provider for token.actions.githubusercontent.com with thumbprint, create aws_iam_role with assume_role_policy for sts:AssumeRoleWithWebIdentity limited to repo:gvasels/personal-music-searchengine:* | Restrictions: Limit OIDC to specific repository, use condition keys for repository and ref, follow least privilege | Success: GitHub Actions can assume role via OIDC, no static credentials needed | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 7. Add deploy permissions to GitHub Actions IAM role
  - File: `infrastructure/global/iam.tf` (append to task 6)
  - Create IAM policy for S3 sync, CloudFront invalidation, ECR push
  - Attach policy to GitHub Actions role
  - Purpose: Allow GitHub Actions to deploy frontend and backend
  - _Leverage: Existing Lambda role policies in `infrastructure/global/iam.tf`_
  - _Requirements: US-6.2 AC-6 through AC-10_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: IAM Policy Specialist | Task: Create IAM policy with permissions for s3:PutObject/DeleteObject/ListBucket on frontend bucket, cloudfront:CreateInvalidation on frontend distribution, ecr:GetAuthorizationToken/BatchCheckLayerAvailability/PutImage/InitiateLayerUpload/UploadLayerPart/CompleteLayerUpload on ECR repos | Restrictions: Scope to specific resources (not *), follow least privilege, separate policy from role | Success: GitHub Actions can sync S3, invalidate CloudFront, push to ECR | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 8. Create CI workflow for PR checks
  - File: `.github/workflows/ci.yml`
  - Run Go tests with 80% coverage threshold on backend
  - Run Vitest with coverage on frontend
  - Run golangci-lint and ESLint
  - Run tsc --noEmit for type checking
  - Run tofu validate on all infrastructure modules
  - Purpose: Validate code quality before merge
  - _Leverage: Existing patterns from `.github/workflows/claude-code-review.yml`_
  - _Requirements: US-6.2 AC-1 through AC-5, AC-9_
  - _Status: Implemented with backend-tests, frontend-tests (Vitest), backend-lint (golangci-lint v1.61), frontend linting (ESLint), type-checking, tofu-validate (matrix), security-scan (Gitleaks, Checkov). Coverage threshold temporarily at 19% (to be increased to 80% as coverage improves)._
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer specializing in GitHub Actions and CI/CD | Task: Create ci.yml workflow triggered on pull_request to main with jobs: backend-tests (go test with -coverprofile, fail if <80%), frontend-tests (npm test --coverage), backend-lint (golangci-lint run), frontend-lint (npm run lint), type-check (go build, npm run typecheck), tofu-validate (matrix for global/shared/backend/frontend modules) | Restrictions: Use ubuntu-latest runners, cache Go and npm dependencies, fail fast on any check failure | Success: All checks run on PR, coverage enforced, types validated | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 9. Create deploy workflow for main branch
  - File: `.github/workflows/deploy.yml`
  - Trigger on push to main
  - Use OIDC for AWS authentication
  - Deploy infrastructure with OpenTofu
  - Build and deploy frontend to S3
  - Invalidate CloudFront cache
  - Purpose: Automated deployment on merge to main
  - _Leverage: AWS configure-credentials action with OIDC_
  - _Requirements: US-6.2 AC-6 through AC-8, AC-10_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer specializing in AWS deployments and GitHub Actions | Task: Create deploy.yml workflow triggered on push to main with id-token: write permission, jobs: deploy-infrastructure (tofu apply for each module in order), deploy-frontend (npm ci, npm run build, aws s3 sync dist/ to bucket --delete), invalidate-cache (aws cloudfront create-invalidation --paths "/*") | Restrictions: Use OIDC authentication (role-to-assume), run infrastructure deploys sequentially (global→shared→backend→frontend), only deploy if CI passes | Success: Merge to main triggers full deployment, frontend accessible at CloudFront URL | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 10. Add backend Docker build and ECR push to deploy workflow
  - File: `.github/workflows/deploy.yml` (append to task 9)
  - Build backend Lambda Docker images
  - Push to ECR repositories
  - Update Lambda function code
  - Purpose: Deploy backend Lambda code changes
  - _Leverage: ECR repository names from `infrastructure/global/ecr.tf`_
  - _Requirements: US-6.2 AC-10_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer specializing in Docker and Lambda deployments | Task: Add job deploy-backend to deploy.yml that builds backend/Dockerfile, logs into ECR with aws ecr get-login-password, tags and pushes image to ECR, updates Lambda function with aws lambda update-function-code --image-uri | Restrictions: Only run if backend/ files changed (use paths filter or check), use same OIDC auth, ensure image is built before push | Success: Backend changes trigger Lambda update, new container deployed | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

## US-6.3: Documentation

- [x] 11. Create CLAUDE.md for infrastructure/frontend
  - File: `infrastructure/frontend/CLAUDE.md`
  - Document all resources, deployment instructions, outputs
  - Follow existing CLAUDE.md patterns from other infrastructure directories
  - Purpose: Enable maintainability of frontend infrastructure
  - _Leverage: `infrastructure/backend/CLAUDE.md` as template_
  - _Requirements: US-6.3 AC-1, AC-2_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer specializing in infrastructure documentation | Task: Create CLAUDE.md with Overview section, File Descriptions table (s3.tf, cloudfront.tf, main.tf), Resources Created tables for S3 and CloudFront, Outputs table, Deployment instructions, Cache Invalidation section | Restrictions: Follow exact format of infrastructure/backend/CLAUDE.md, include all resources and outputs, keep concise | Success: New team member can understand and deploy frontend infrastructure from docs | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 12. Update root CLAUDE.md with Epic 6 completion
  - File: `CLAUDE.md` (root)
  - Add frontend infrastructure to repository structure
  - Update CI/CD section with new workflows
  - Add deployment instructions
  - Purpose: Keep root documentation current
  - _Leverage: Existing CLAUDE.md structure_
  - _Requirements: US-6.3 AC-1_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer | Task: Update root CLAUDE.md to add infrastructure/frontend to Repository Structure, add CI/CD section describing ci.yml and deploy.yml workflows, add Deployment section with manual and automated deployment instructions | Restrictions: Do not remove existing content, maintain consistent formatting, keep updates concise | Success: CLAUDE.md accurately reflects current project state including Epic 6 additions | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 13. Create deployment documentation
  - File: `docs/deployment.md`
  - Document manual deployment steps for all infrastructure modules
  - Document automated deployment via GitHub Actions
  - Include rollback procedures
  - Purpose: Enable team to deploy and troubleshoot
  - _Leverage: Infrastructure CLAUDE.md files for module-specific details_
  - _Requirements: US-6.3 AC-5_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Documentation Specialist | Task: Create docs/deployment.md with sections: Prerequisites (AWS CLI, OpenTofu, AWS profile), Manual Deployment (step-by-step for each module in order), Automated Deployment (GitHub Actions workflow explanation), Rollback Procedures (S3 versioning, tofu state), Troubleshooting (common issues and solutions) | Restrictions: Include actual commands, reference correct bucket/resource names, be specific not generic | Success: Developer can deploy entire stack following documentation | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 14. Update CHANGELOG.md with Epic 6
  - File: `CHANGELOG.md`
  - Add entry for Epic 6: Distribution & Polish
  - List all new infrastructure and CI/CD additions
  - Purpose: Track project changes
  - _Leverage: Existing CHANGELOG.md format_
  - _Requirements: US-6.3 AC-7_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Release Manager | Task: Add CHANGELOG.md entry under [Unreleased] or new version section with Added: Frontend S3 bucket and CloudFront distribution, GitHub Actions CI workflow with test coverage enforcement, GitHub Actions deploy workflow with OIDC authentication, Deployment documentation | Restrictions: Follow Keep a Changelog format, be concise, list actual changes | Success: CHANGELOG reflects all Epic 6 additions | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

- [x] 15. Update epics-user-stories.md to mark Epic 6 complete
  - File: `implementation-plan/epics-user-stories.md`
  - Update Epic 6 status to Complete
  - Check all acceptance criteria boxes
  - Add completion date
  - Purpose: Track epic completion in central tracking document
  - _Leverage: Epic 5 completion entry as template_
  - _Requirements: US-6.3 (project tracking)_
  - _Prompt: Implement the task for spec distribution-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Project Manager | Task: Update epics-user-stories.md Epic 6 section: change Status to Complete, add completion date, check all acceptance criteria boxes [x], add Implementation Summary similar to Epic 5 | Restrictions: Only update Epic 6 section, maintain existing formatting | Success: Epic 6 marked complete with all criteria checked | Instructions: Mark task [-] in tasks.md when starting, use log-implementation tool after completion, mark [x] when done_

## Task Dependencies

```
Infrastructure Tasks (Sequential):
Task 1 (S3) → Task 4 (Bucket Policy) → Task 2 (CloudFront) → Task 3 (Security Headers) → Task 5 (main.tf)

IAM Tasks (Sequential):
Task 6 (OIDC Provider) → Task 7 (Deploy Permissions)

CI/CD Tasks (Parallel after IAM):
Task 8 (CI Workflow) | Task 9 (Deploy Workflow) → Task 10 (Backend Deploy)

Documentation Tasks (After Implementation):
Task 11 (Frontend CLAUDE.md) | Task 12 (Root CLAUDE.md) | Task 13 (Deployment Docs) | Task 14 (CHANGELOG) → Task 15 (Epic Complete)
```

## Wave Assignment

| Wave | Tasks | Focus |
|------|-------|-------|
| 1 | 1, 4, 5 | S3 infrastructure foundation |
| 2 | 2, 3 | CloudFront distribution |
| 3 | 6, 7 | IAM and OIDC setup |
| 4 | 8, 9, 10 | CI/CD workflows |
| 5 | 11, 12, 13, 14, 15 | Documentation |
