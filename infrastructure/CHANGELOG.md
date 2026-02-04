# Changelog

All notable changes to the infrastructure will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Cognito admin IAM permissions for Lambda (`backend/iam-cognito.tf`)
  - Allows user management operations (ListUsers, AdminGetUser, AdminAddUserToGroup, etc.)
  - Attached to API Lambda execution role for admin panel functionality
- FFmpeg Lambda layer for audio processing
- Analyzer Lambda processor infrastructure
- S3 Tables configuration (disabled pending AWS provider support)
- `frontend_cloudfront_domain` variable for CORS configuration

### Changed
- Applied OpenTofu formatting to all configuration files
- Updated CI workflow with fetch-depth for security scanning

### Security
- Restricted Bedrock Gateway CORS to localhost development origins and configurable frontend domain
- Added explicit header allowlist (Authorization, Content-Type, X-Request-ID) instead of wildcard
- Added documentation about API key validation in Lambda

### Fixed
- Placeholder for FFmpeg layer in CI validation
- Gitleaks security scan fetch-depth configuration
