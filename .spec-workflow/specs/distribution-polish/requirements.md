# Requirements Document - Distribution & Polish (Epic 6)

## Introduction

Epic 6 focuses on deploying the frontend application to production, establishing comprehensive CI/CD pipelines, and finalizing project documentation. This epic completes the Personal Music Search Engine by making it accessible to users via a public URL and ensuring maintainability through automated testing and documentation.

## Alignment with Product Vision

This epic delivers the final piece of the user experience by:
- Enabling users to access the application from any device via a URL (US-6.1)
- Ensuring code quality through automated testing and deployment pipelines (US-6.2)
- Maintaining project sustainability through comprehensive documentation (US-6.3)

## Requirements

### Requirement 1: Frontend Hosting (US-6.1)

**User Story:** As a User, I want to access the application via a URL, so that I can use it from any device.

#### Acceptance Criteria

1. WHEN user navigates to the application URL THEN CloudFront SHALL serve the React SPA from S3
2. WHEN user requests any SPA route (e.g., /tracks, /albums) THEN CloudFront SHALL return index.html with 200 status (SPA routing support)
3. WHEN user requests static assets (JS, CSS, images) THEN CloudFront SHALL serve them with appropriate cache headers (max-age=31536000 for hashed assets)
4. WHEN user requests index.html THEN CloudFront SHALL serve it with no-cache to enable app updates
5. IF user requests a non-existent path THEN CloudFront SHALL return index.html for client-side routing
6. WHEN assets are requested THEN CloudFront SHALL compress responses with gzip/brotli for supported clients
7. IF API requests are made from the frontend THEN CORS SHALL be properly configured on the API Gateway

### Requirement 2: Automated Testing & CI/CD (US-6.2)

**User Story:** As a Developer, I want automated tests and deployment pipelines, so that code quality is maintained and deployments are reliable.

#### Acceptance Criteria

1. WHEN a PR is created THEN GitHub Actions SHALL run Go unit tests with minimum 80% coverage
2. WHEN a PR is created THEN GitHub Actions SHALL run frontend tests (Vitest) with coverage reporting
3. WHEN a PR is created THEN GitHub Actions SHALL run linting (golangci-lint, ESLint)
4. WHEN a PR is created THEN GitHub Actions SHALL run type checking (go build, tsc --noEmit)
5. WHEN a PR is created THEN GitHub Actions SHALL run OpenTofu validation (tofu validate)
6. WHEN a PR is merged to main THEN GitHub Actions SHALL deploy infrastructure changes via OpenTofu
7. WHEN a PR is merged to main THEN GitHub Actions SHALL build and deploy the frontend to S3
8. WHEN a PR is merged to main THEN GitHub Actions SHALL invalidate CloudFront cache
9. IF any CI check fails THEN the PR SHALL be blocked from merging
10. WHEN backend Lambda code changes THEN CI SHALL build Docker images and push to ECR
11. WHEN tests pass THEN CI SHALL generate and publish coverage reports

### Requirement 3: Documentation (US-6.3)

**User Story:** As a Developer, I want comprehensive documentation, so that the project can be maintained and extended.

#### Acceptance Criteria

1. WHEN accessing the project THEN root CLAUDE.md SHALL contain complete project overview, tech stack, and quick start guide
2. WHEN a directory contains code THEN it SHALL have a CLAUDE.md with file descriptions and key function signatures
3. WHEN reviewing the codebase THEN all public Go functions SHALL have GoDoc comments
4. WHEN reviewing the codebase THEN all public TypeScript functions/components SHALL have TSDoc comments
5. WHEN deploying the application THEN deployment documentation SHALL include step-by-step instructions
6. WHEN using the API THEN API documentation SHALL be available (OpenAPI spec is already defined)
7. WHEN reviewing changes THEN CHANGELOG.md SHALL be maintained with version history

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each OpenTofu module handles one infrastructure concern (frontend hosting)
- **Modular Design**: CI/CD workflows are separated by concern (test, deploy, review)
- **Dependency Management**: Infrastructure modules use proper state management and outputs
- **Clear Interfaces**: Workflows use standardized inputs/outputs for composability

### Performance
- **CloudFront Edge Caching**: Static assets cached at edge locations for <100ms latency globally
- **Gzip/Brotli Compression**: All text assets compressed for reduced transfer size
- **Cache Headers**: Proper cache-control headers for optimal browser caching
- **CI Pipeline Speed**: Full CI run should complete in <10 minutes

### Security
- **HTTPS Only**: All traffic served over HTTPS via CloudFront
- **S3 Origin Access Control**: S3 bucket not publicly accessible, only via CloudFront OAC
- **Secrets Management**: GitHub Secrets for API keys, AWS credentials
- **Dependency Scanning**: Automated vulnerability scanning in CI
- **IAM Least Privilege**: GitHub Actions OIDC role with minimal permissions

### Reliability
- **Multi-AZ S3**: S3 Standard storage class with 99.99% availability
- **CloudFront High Availability**: Built-in redundancy across edge locations
- **Deployment Rollback**: Ability to rollback frontend deployments
- **CI Retry Logic**: Transient failures handled with retry mechanisms

### Usability
- **SPA Deep Linking**: All routes work when directly accessed or refreshed
- **Deployment Notifications**: Slack/email notifications for deployment status
- **Clear Error Pages**: Custom 404/error pages for better user experience

## Dependencies

| Dependency | Type | Purpose |
|------------|------|---------|
| Epic 5 (Frontend) | Internal | Frontend build artifacts to deploy |
| Epic 1 (Global Infrastructure) | Internal | S3 state bucket, ECR, base IAM |
| Epic 2 (Backend API) | Internal | API Gateway endpoint for CORS configuration |
| GitHub Actions | External | CI/CD platform |
| AWS CloudFront | External | CDN for frontend distribution |
| AWS S3 | External | Frontend asset storage |

## Out of Scope

- Custom domain name configuration (can be added later)
- WAF/DDoS protection (can be added later)
- Blue-green deployments (future enhancement)
- Performance monitoring/APM integration
- Feature flags system
