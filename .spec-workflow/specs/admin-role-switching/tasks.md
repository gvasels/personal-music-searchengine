# Tasks: Admin Role Switching

## Overview

This document breaks down the Admin Role Switching feature into implementable tasks following the TDD workflow.

---

## Task Groups

### Group 1: State Management (Frontend)

#### Task 1.1: Create Role Simulation Zustand Store
- **Description**: Create Zustand store for persisting simulation state
- **Files**:
  - `frontend/src/lib/store/roleSimulationStore.ts`
  - `frontend/src/lib/store/roleSimulationStore.test.ts`
- **Acceptance Criteria**:
  - Store has `simulatedRole: UserRole | null` state
  - Store has `setSimulatedRole(role)` action
  - Store has `clearSimulation()` action
  - State persists in sessionStorage
  - Tests verify all state transitions

#### Task 1.2: Create useRoleSimulation Hook
- **Description**: Create hook that combines auth state with simulation state
- **Files**:
  - `frontend/src/hooks/useRoleSimulation.ts`
  - `frontend/src/hooks/__tests__/useRoleSimulation.test.tsx`
- **Acceptance Criteria**:
  - Hook returns `effectiveRole` (simulated or actual)
  - Hook returns `isSimulating` boolean
  - Hook returns `canSimulate` (true if admin)
  - Hook provides `startSimulation(role)` and `stopSimulation()` methods
  - Only admins can start simulation
  - Tests cover all scenarios

---

### Group 2: UI Components (Frontend)

#### Task 2.1: Create RoleSwitcher Dropdown Component
- **Description**: Dropdown in header for selecting simulation role
- **Files**:
  - `frontend/src/components/admin/RoleSwitcher.tsx`
  - `frontend/src/components/admin/__tests__/RoleSwitcher.test.tsx`
- **Acceptance Criteria**:
  - Shows current effective role
  - Dropdown lists: Admin, Artist, Subscriber, Guest
  - Only renders for admin users
  - Calls `startSimulation` on role selection
  - Shows checkmark on current selection
  - Uses DaisyUI dropdown component

#### Task 2.2: Create SimulationBanner Component
- **Description**: Persistent banner showing simulation status
- **Files**:
  - `frontend/src/components/admin/SimulationBanner.tsx`
  - `frontend/src/components/admin/__tests__/SimulationBanner.test.tsx`
- **Acceptance Criteria**:
  - Only renders when `isSimulating` is true
  - Shows "Viewing as [Role]" text
  - Has amber/orange background
  - Has "Exit Simulation" button
  - Button calls `stopSimulation()`
  - Fixed position below header

---

### Group 3: Integration (Frontend)

#### Task 3.1: Integrate RoleSwitcher into Header
- **Description**: Add RoleSwitcher to main header component
- **Files**:
  - `frontend/src/components/layout/Header.tsx` (modify)
- **Acceptance Criteria**:
  - RoleSwitcher appears before user avatar
  - Only visible when user is admin
  - Responsive on mobile (collapses appropriately)

#### Task 3.2: Integrate SimulationBanner into Layout
- **Description**: Add SimulationBanner to app layout
- **Files**:
  - `frontend/src/components/layout/AppLayout.tsx` (modify)
- **Acceptance Criteria**:
  - Banner appears below header when simulating
  - Does not scroll with page content
  - Takes full width

#### Task 3.3: Update Feature Flags to Use Effective Role
- **Description**: Modify useFeatureFlags to respect simulated role
- **Files**:
  - `frontend/src/hooks/useFeatureFlags.ts` (modify)
  - `frontend/src/hooks/__tests__/useFeatureFlags.test.tsx` (modify)
- **Acceptance Criteria**:
  - Uses `effectiveRole` from useRoleSimulation
  - Feature flags change when simulation starts
  - Tests verify role-based feature changes

#### Task 3.4: Update Sidebar Navigation for Simulated Role
- **Description**: Hide/show nav items based on effective role
- **Files**:
  - `frontend/src/components/layout/Sidebar.tsx` (modify)
- **Acceptance Criteria**:
  - Admin section hidden when simulating non-admin
  - Artist features hidden when simulating subscriber/guest
  - Upload hidden when simulating guest

#### Task 3.5: Block Write Actions in Simulation Mode
- **Description**: Show toast and prevent mutations when simulating
- **Files**:
  - `frontend/src/lib/api/client.ts` (modify)
  - `frontend/src/hooks/useRoleSimulation.ts` (add `isWriteBlocked`)
- **Acceptance Criteria**:
  - POST/PUT/DELETE blocked when simulating non-admin
  - Toast shows "Action blocked in simulation mode"
  - Read operations (GET) still work

---

## Implementation Order

```
1.1 → 1.2 → 2.1 → 2.2 → 3.1 → 3.2 → 3.3 → 3.4 → 3.5
```

All tasks follow TDD: write failing tests first, then implement.

---

## Estimated Effort

| Group | Tasks | Complexity |
|-------|-------|------------|
| State Management | 2 | Medium |
| UI Components | 2 | Low |
| Integration | 5 | Medium |
| **Total** | **9 tasks** | |

---

## Dependencies

- Existing `useAuth` hook
- Existing `useFeatureFlags` hook
- DaisyUI components
- Zustand store patterns from existing codebase
