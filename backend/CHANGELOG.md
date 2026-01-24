# Changelog

All notable changes to the backend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Audio analysis package (`internal/analysis/`) with BPM detection using multi-segment autocorrelation algorithm
- Camelot wheel mapping for harmonic mixing support (24 key mappings including enharmonic equivalents)
- Analyzer Lambda processor (`cmd/processor/analyzer/`) for Step Functions integration
- Migration service (`internal/service/migration.go`) for string-to-entity artist migration
- Playlist reorder endpoint for track position management
- Matching service for DJ-style track compatibility scoring
- FFmpeg input validation to prevent command injection

### Changed
- Updated CI coverage threshold from 19% to 24%
- Added golangci-lint job to CI workflow

### Fixed
- CORS handling for playlist reorder endpoint
- 404 error on playlist reorder route
