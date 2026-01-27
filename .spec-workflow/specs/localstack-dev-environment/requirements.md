# Requirements Document: LocalStack Development Environment

## Introduction

This feature establishes a local AWS development environment using LocalStack and Docker to enable realistic integration testing without deploying to AWS. The goal is to move away from mock-heavy unit tests toward integration tests that exercise real AWS service behaviors locally, catching issues earlier and improving test reliability.

Currently, the project relies heavily on mocked AWS services in tests, which:
- Miss integration issues between components
- Don't validate actual AWS API contracts
- Can't test IAM permissions or service configurations
- Make refactoring risky due to false confidence from passing mock tests

## Alignment with Product Vision

This feature supports the development velocity and quality goals by:
- Reducing deployment cycles for testing (no AWS round-trips)
- Catching integration bugs before they reach staging/production
- Enabling developers to run the full stack locally
- Providing consistent test environments across team members

## Requirements

### Requirement 1: Docker-based LocalStack Infrastructure

**User Story:** As a developer, I want to run AWS services locally via Docker, so that I can test against realistic AWS APIs without incurring costs or deployment delays.

#### Acceptance Criteria

1. WHEN developer runs `docker-compose up` THEN the system SHALL start LocalStack with DynamoDB, Lambda, API Gateway, S3, Cognito, and Step Functions services
2. WHEN LocalStack starts THEN the system SHALL automatically create the DynamoDB table with the production schema
3. WHEN LocalStack starts THEN the system SHALL create the S3 buckets (media, search-indexes)
4. IF LocalStack container is already running THEN `docker-compose up` SHALL be idempotent (no errors, no duplicate resources)
5. WHEN developer runs `docker-compose down` THEN all containers SHALL stop and local state SHALL be cleared

### Requirement 2: Local Lambda Execution

**User Story:** As a developer, I want to run the Go Lambda functions locally against LocalStack, so that I can test API handlers with real DynamoDB and S3 interactions.

#### Acceptance Criteria

1. WHEN the local Lambda is invoked THEN it SHALL connect to LocalStack DynamoDB instead of AWS DynamoDB
2. WHEN the local Lambda performs S3 operations THEN it SHALL use LocalStack S3
3. WHEN environment variable `LOCAL_STACK=true` is set THEN the Lambda SHALL use LocalStack endpoints
4. WHEN running locally THEN Lambda logs SHALL be visible in the terminal
5. IF LocalStack is not running THEN the Lambda SHALL fail with a clear error message

### Requirement 3: Local API Gateway

**User Story:** As a developer, I want a local API Gateway that routes to local Lambdas, so that the frontend can make requests to the same endpoints as production.

#### Acceptance Criteria

1. WHEN the local API Gateway starts THEN it SHALL expose the same routes as production (`/tracks`, `/albums`, `/playlists`, etc.)
2. WHEN frontend makes requests to `localhost:4566` THEN requests SHALL be routed to local Lambda functions
3. WHEN API Gateway receives requests THEN it SHALL validate JWT tokens (or bypass auth in dev mode)
4. IF a route doesn't exist THEN API Gateway SHALL return 404 matching production behavior

### Requirement 4: Local Cognito Authentication

**User Story:** As a developer, I want local Cognito user pools, so that I can test authentication flows without hitting AWS Cognito.

#### Acceptance Criteria

1. WHEN LocalStack starts THEN a local Cognito user pool SHALL be created
2. WHEN developer runs a seed script THEN test users (admin, subscriber, artist) SHALL be created
3. WHEN frontend authenticates against local Cognito THEN it SHALL receive valid JWT tokens
4. WHEN local JWT tokens are used THEN API Gateway SHALL accept them for authorization

### Requirement 5: Frontend Local Development Mode

**User Story:** As a developer, I want the React frontend to connect to LocalStack services, so that I can develop and test the full stack locally.

#### Acceptance Criteria

1. WHEN `VITE_LOCAL_STACK=true` is set THEN frontend SHALL use LocalStack API endpoint
2. WHEN `VITE_LOCAL_STACK=true` is set THEN frontend SHALL use local Cognito configuration
3. WHEN running `npm run dev:local` THEN frontend SHALL start with LocalStack configuration
4. WHEN media URLs are requested THEN they SHALL point to LocalStack S3

### Requirement 6: Integration Test Framework

**User Story:** As a developer, I want an integration test framework that runs against LocalStack, so that I can write tests that exercise real AWS interactions.

#### Acceptance Criteria

1. WHEN running `go test -tags=integration ./...` THEN tests SHALL run against LocalStack
2. WHEN integration tests start THEN they SHALL wait for LocalStack to be healthy
3. WHEN each integration test runs THEN it SHALL have isolated data (cleanup between tests)
4. WHEN integration tests complete THEN they SHALL report coverage metrics
5. IF LocalStack is not running THEN integration tests SHALL skip with a clear message

### Requirement 7: One-Command Setup

**User Story:** As a developer, I want a single command to start the entire local environment, so that onboarding and daily development are frictionless.

#### Acceptance Criteria

1. WHEN developer runs `make local` or `./scripts/local-dev.sh` THEN the system SHALL:
   - Start LocalStack via Docker Compose
   - Wait for services to be healthy
   - Run database migrations/seeds
   - Build and deploy local Lambda
   - Start the frontend dev server
2. WHEN any component fails to start THEN the system SHALL display clear error messages and cleanup
3. WHEN developer runs `make local-stop` THEN all local services SHALL stop cleanly

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: LocalStack setup scripts separate from test utilities
- **Modular Design**: AWS client factory pattern to switch between real AWS and LocalStack
- **Dependency Management**: Use environment variables for endpoint configuration, not code changes
- **Clear Interfaces**: Same service interfaces work for both local and production

### Performance
- LocalStack containers SHALL start within 60 seconds
- Integration tests SHALL complete within 5 minutes for the full suite
- Local Lambda cold start SHALL be under 3 seconds

### Security
- Local development SHALL NOT use real AWS credentials
- Test data SHALL NOT contain real user information
- LocalStack SHALL be accessible only from localhost

### Reliability
- LocalStack configuration SHALL be version-controlled
- Docker Compose SHALL use health checks for service readiness
- Scripts SHALL be idempotent (safe to run multiple times)

### Usability
- Setup instructions SHALL be documented in README
- Error messages SHALL suggest remediation steps
- Logs SHALL be aggregated and easily viewable

## Out of Scope

- CloudFront CDN emulation (not supported by LocalStack free tier)
- MediaConvert emulation (not supported by LocalStack)
- Production deployment changes
- CI/CD pipeline integration (future spec)
