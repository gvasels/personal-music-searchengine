# Changelog

All notable changes to the frontend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Admin Role Simulation Feature**
  - `RoleSwitcher` component - Dropdown for admin to simulate different user roles
  - `SimulationBanner` component - Alert banner shown during role simulation
  - `useRoleSimulation` hook - Combines auth state with simulation state
  - `roleSimulationStore` - Zustand store with localStorage persistence
  - `useFeatureFlags` integration - Uses simulated role for `hasRole()` checks
  - Visual indicators: btn-warning styling during simulation, checkmark for current role
  - Accessibility: ARIA role="alert" on banner
  - Comprehensive test coverage (77 tests for simulation feature)

- **Admin Panel & Track Visibility Feature**
  - Admin users page (`routes/admin/users.tsx`) for user management
  - Admin components: `UserSearchForm`, `UserCard`, `UserDetailModal`
  - Admin API client (`lib/api/admin.ts`) for user search and management
  - Admin hooks (`hooks/useAdmin.ts`) with debounced search
  - Track visibility selector on track detail page
  - "Uploaded By" column in TrackList (admin-only, hidden by default)
  - Settings toggle for showing/hiding "Uploaded By" column
  - Admin navigation section in Sidebar
  - `TrackVisibility` type (private/unlisted/public)
  - `showUploadedBy` preference in preferencesStore
  - Component tests for admin components
