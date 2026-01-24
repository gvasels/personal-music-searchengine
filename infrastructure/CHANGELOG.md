# Changelog

All notable changes to the infrastructure will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- FFmpeg Lambda layer for audio processing
- Analyzer Lambda processor infrastructure
- S3 Tables configuration (disabled pending AWS provider support)

### Changed
- Applied OpenTofu formatting to all configuration files
- Updated CI workflow with fetch-depth for security scanning

### Fixed
- Placeholder for FFmpeg layer in CI validation
- Gitleaks security scan fetch-depth configuration
