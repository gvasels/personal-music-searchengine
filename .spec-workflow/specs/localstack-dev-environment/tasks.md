# Tasks Document: LocalStack Development Environment

## Phase 1: Docker & LocalStack Configuration

- [x] 1.1 Update docker-compose.yml to add Cognito service
  - File: `docker/docker-compose.yml`
  - Add `cognito-idp` to SERVICES environment variable
  - Verify LocalStack 3.4+ supports required services
  - Purpose: Enable local Cognito user pool for authentication
  - _Leverage: Existing docker-compose.yml configuration_
  - _Requirements: 1.1, 4.1_

- [x] 1.2 Create Cognito initialization script
  - File: `docker/localstack-init/init-cognito.sh`
  - Create local Cognito user pool and app client
  - Create test user groups (admin, subscriber, artist)
  - Create test users with known credentials
  - Purpose: Provide local authentication for development
  - _Leverage: Existing init-aws.sh pattern_
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [x] 1.3 Create wait-for-localstack.sh script
  - File: `scripts/wait-for-localstack.sh`
  - Poll LocalStack health endpoint until ready
  - Check specific services (dynamodb, s3, cognito-idp)
  - Timeout after configurable duration
  - Purpose: Reliable startup sequencing
  - _Leverage: Health check pattern from existing init scripts_
  - _Requirements: 7.1_

## Phase 2: Integration Test Framework

- [x] 2.1 Create LocalStack test utilities
  - File: `backend/internal/testutil/localstack.go`
  - SetupLocalStack function with cleanup
  - TestContext struct with AWS clients
  - Skip logic when LocalStack unavailable
  - Purpose: Foundation for integration tests
  - _Leverage: Existing AWS client setup in cmd/api/main.go:75-109_
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 2.2 Create test fixtures
  - File: `backend/internal/testutil/fixtures.go`
  - TestUsers map with credentials
  - CreateTestTrack helper function
  - CreateTestUser helper function
  - Purpose: Consistent test data across tests
  - _Leverage: models package for data structures_
  - _Requirements: 6.3_

- [x] 2.3 Create cleanup utilities
  - File: `backend/internal/testutil/cleanup.go`
  - CleanupUser function
  - CleanupTrack function
  - CleanupAll function for test teardown
  - Purpose: Isolated test data between tests
  - _Leverage: DynamoDB DeleteItem operations_
  - _Requirements: 6.3_

- [x] 2.4 Write sample integration test
  - File: `backend/internal/service/track_integration_test.go`
  - Use `//go:build integration` tag
  - Test TrackService.CreateTrack against LocalStack
  - Verify data persists in DynamoDB
  - Purpose: Prove integration test framework works
  - _Leverage: testutil package, service package_
  - _Requirements: 6.1, 6.2_

## Phase 3: Frontend Local Mode

- [x] 3.1 Create frontend local environment template
  - File: `frontend/.env.local.example`
  - Document all local configuration variables
  - Include comments explaining each variable
  - Purpose: Easy local setup for developers
  - _Leverage: Existing .env patterns_
  - _Requirements: 5.1, 5.2_

- [x] 3.2 Update Amplify configuration for LocalStack
  - File: `frontend/src/lib/config.ts`
  - Detect VITE_LOCAL_STACK environment variable
  - Configure for LocalStack Cognito endpoint
  - Purpose: Frontend auth works with local Cognito
  - _Leverage: Existing Amplify configuration_
  - _Requirements: 5.2, 5.3_

- [x] 3.3 Add npm script for local development
  - File: `frontend/package.json`
  - Add `dev:local` script that loads .env.local
  - Purpose: One command to start frontend in local mode
  - _Leverage: Existing dev script_
  - _Requirements: 5.4_

## Phase 4: One-Command Setup

- [x] 4.1 Create root Makefile
  - File: `Makefile` (project root)
  - `make local` - start full environment
  - `make local-stop` - stop environment
  - `make test-integration` - run integration tests
  - Purpose: Simple commands for common operations
  - _Leverage: Scripts created in earlier tasks_
  - _Requirements: 7.1, 7.2, 7.3_

- [x] 4.2 Create shell script alternative
  - File: `scripts/local-dev.sh`
  - Same functionality as Makefile for non-make users
  - Support start, stop, test subcommands
  - Purpose: Flexibility for different developer preferences
  - _Leverage: Scripts and Makefile logic_
  - _Requirements: 7.1, 7.2_

## Phase 5: Documentation

- [x] 5.1 Update docker/CLAUDE.md
  - File: `docker/CLAUDE.md`
  - Document Cognito setup
  - Update quick start instructions
  - Add troubleshooting for new services
  - Purpose: Keep documentation current
  - _Leverage: Existing CLAUDE.md structure_
  - _Requirements: All_

- [x] 5.2 Update root README or create LOCAL_DEV.md
  - File: `LOCAL_DEV.md` (project root)
  - Complete local development guide
  - Prerequisites, setup, troubleshooting
  - Purpose: Onboarding documentation
  - _Leverage: CLAUDE.md files_
  - _Requirements: All_

- [x] 5.3 Create backend testutil CLAUDE.md
  - File: `backend/internal/testutil/CLAUDE.md`
  - Document test utilities and usage
  - Include examples
  - Purpose: Enable developers to write integration tests
  - _Leverage: Code created in Phase 2_
  - _Requirements: 6.1, 6.2, 6.3_
