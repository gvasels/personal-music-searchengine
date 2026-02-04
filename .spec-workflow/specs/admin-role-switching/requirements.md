# Requirements: Admin Role Switching

## Introduction

This feature enables administrators to temporarily "impersonate" or view the application as a user with a different role (guest, subscriber, or artist). This allows admins to verify that role-based access controls work correctly, test the user experience for different roles, and troubleshoot issues reported by users without needing separate test accounts.

The feature is purely client-side for viewing purposes - it does not change the actual JWT or backend permissions. The admin sees what a user of that role would see in the UI.

## Alignment with Product Vision

This feature supports the Global User Type system by providing a way to:
- Verify role-based UI restrictions work correctly
- Test feature flags and permissions for different roles
- Debug user-reported issues by seeing their perspective
- Ensure quality assurance across all user types

---

## Requirements

### REQ-1: Role Switching UI

**User Story:** As an admin, I want to switch my view to see the application as a different user role so that I can verify the user experience for each role type.

#### Acceptance Criteria

1. WHEN an admin user is logged in THEN the system SHALL display a role switcher control in the header or sidebar
2. WHEN the admin clicks the role switcher THEN the system SHALL show options: Admin (current), Artist, Subscriber, Guest
3. WHEN the admin selects a different role THEN the system SHALL immediately update the UI to reflect that role's permissions
4. IF the admin is viewing as a non-admin role THEN the system SHALL display a visible indicator showing "Viewing as [role]"
5. WHEN the admin clicks "Exit" or selects "Admin" THEN the system SHALL restore the full admin view
6. IF the user is not an admin THEN the system SHALL NOT display the role switcher control

---

### REQ-2: Permission Simulation

**User Story:** As an admin viewing as another role, I want the UI to accurately reflect what that role can see and do so that I can verify the permission system works correctly.

#### Acceptance Criteria

1. WHEN viewing as Guest THEN the system SHALL:
   - Hide navigation items that require authentication
   - Show only public content (public playlists, artist profiles)
   - Hide upload, library management, and playlist creation options
   - Display login/signup prompts where appropriate

2. WHEN viewing as Subscriber THEN the system SHALL:
   - Show subscriber-level navigation (Library, Playlists, Upload)
   - Hide artist-specific features (Artist Profile management)
   - Hide admin features (User Management)
   - Allow viewing of the subscriber's own content context

3. WHEN viewing as Artist THEN the system SHALL:
   - Show artist-level navigation (includes Artist Profile)
   - Allow viewing artist profile management UI
   - Hide admin features (User Management)

4. WHEN viewing as any simulated role THEN the system SHALL:
   - NOT allow actual data modifications (read-only simulation)
   - Show a banner indicating simulation mode
   - Preserve actual admin session for security

---

### REQ-3: Persistence and State

**User Story:** As an admin, I want my role switching preference to persist during my session so that I don't have to reselect it when navigating between pages.

#### Acceptance Criteria

1. WHEN the admin selects a simulated role THEN the system SHALL persist that selection in session state
2. WHEN the admin navigates to a different page THEN the system SHALL maintain the simulated role view
3. WHEN the admin refreshes the page THEN the system SHALL restore the simulated role (optional: can reset to admin)
4. WHEN the admin logs out THEN the system SHALL clear the simulated role state
5. WHEN the admin's session expires THEN the system SHALL NOT persist simulated role to next session

---

### REQ-4: Visual Distinction

**User Story:** As an admin in role simulation mode, I want clear visual indicators so that I always know I'm viewing a simulated experience.

#### Acceptance Criteria

1. WHEN viewing as a simulated role THEN the system SHALL display a persistent banner at the top of the page
2. The banner SHALL include:
   - The text "Viewing as [Role Name]"
   - A colored border/background distinct from normal UI (e.g., orange/amber)
   - An "Exit Simulation" button
3. IF the admin tries to perform a write action THEN the system SHALL show a toast: "Action blocked in simulation mode"
4. WHEN in simulation mode THEN the system SHALL NOT allow navigation to admin-only pages

---

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility**: Role simulation logic isolated in a dedicated hook (`useRoleSimulation`)
- **Modular Design**: Simulation banner as reusable component
- **Dependency Management**: Minimal changes to existing auth/permission system
- **Clear Interfaces**: Simulation state separate from actual auth state

### Performance
- Role switching SHALL be instantaneous (no API calls required)
- UI updates SHALL complete within 100ms of role selection
- Simulation state stored in memory/sessionStorage (no server round-trips)

### Security
- Simulation is CLIENT-SIDE ONLY - backend permissions remain unchanged
- API calls still use admin's actual JWT token
- Backend continues to enforce real permissions
- No sensitive data exposure through simulation
- Simulation cannot escalate privileges (only demote view)

### Reliability
- If simulation state becomes corrupted, system defaults to actual admin role
- Page refresh recovers gracefully
- Session expiry clears simulation state completely

### Usability
- Role switcher easily accessible (header dropdown or sidebar)
- Clear visual feedback when in simulation mode
- One-click exit from simulation
- Keyboard accessible controls

---

## Out of Scope

- **Backend impersonation**: This is UI-only; API calls use real admin credentials
- **Specific user impersonation**: Cannot impersonate a specific user, only role types
- **Data modification**: Write operations blocked in simulation mode
- **Audit logging**: Logging of simulation actions (future enhancement)
- **Cross-session persistence**: Simulation resets on logout/session expiry

---

## Dependencies

- Existing `useAuth` hook with `isAdmin` property
- Existing role-based feature flags (`useFeatureFlags`)
- Global User Type permission system
- React Context or Zustand for state management
