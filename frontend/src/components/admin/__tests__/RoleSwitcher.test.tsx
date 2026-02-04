/**
 * RoleSwitcher Component Tests - Admin Role Switching Feature
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { RoleSwitcher } from '../RoleSwitcher';

// Mock useRoleSimulation hook
const mockStartSimulation = vi.fn();
const mockUseRoleSimulation = vi.fn();

vi.mock('../../../hooks/useRoleSimulation', () => ({
  useRoleSimulation: () => mockUseRoleSimulation(),
}));

describe('RoleSwitcher', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseRoleSimulation.mockReturnValue({
      canSimulate: true,
      effectiveRole: 'admin',
      isSimulating: false,
      startSimulation: mockStartSimulation,
    });
  });

  describe('rendering', () => {
    it('should render when user can simulate', () => {
      const { container } = render(<RoleSwitcher />);

      // The trigger is a label element, not a button
      expect(container.querySelector('label.btn')).toBeInTheDocument();
    });

    it('should not render when user cannot simulate', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: false,
        effectiveRole: 'subscriber',
        isSimulating: false,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      expect(container).toBeEmptyDOMElement();
    });

    it('should display the current role icon', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'admin',
        isSimulating: false,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      // Admin icon is crown (first occurrence is in the trigger)
      const trigger = container.querySelector('label.btn');
      expect(trigger?.textContent).toContain('👑');
    });

    it('should display different icons for different roles', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'subscriber',
        isSimulating: true,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      // Subscriber icon in the trigger
      const trigger = container.querySelector('label.btn');
      expect(trigger?.textContent).toContain('👤');
    });
  });

  describe('styling', () => {
    it('should use btn-ghost when not simulating', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'admin',
        isSimulating: false,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      const trigger = container.querySelector('label.btn');
      expect(trigger?.className).toContain('btn-ghost');
      expect(trigger?.className).not.toContain('btn-warning');
    });

    it('should use btn-warning when simulating', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'subscriber',
        isSimulating: true,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      const trigger = container.querySelector('label.btn');
      expect(trigger?.className).toContain('btn-warning');
      expect(trigger?.className).not.toContain('btn-ghost');
    });
  });

  describe('dropdown options', () => {
    it('should display all role options', () => {
      const { container } = render(<RoleSwitcher />);

      // Find the dropdown menu
      const dropdownMenu = container.querySelector('ul.dropdown-content');

      // Check that all role options are in the dropdown
      expect(dropdownMenu?.textContent).toContain('Admin');
      expect(dropdownMenu?.textContent).toContain('Artist');
      expect(dropdownMenu?.textContent).toContain('Subscriber');
      expect(dropdownMenu?.textContent).toContain('Guest');
    });

    it('should display role descriptions', () => {
      render(<RoleSwitcher />);

      expect(screen.getByText('Full access (actual role)')).toBeInTheDocument();
      expect(screen.getByText('Artist features enabled')).toBeInTheDocument();
      expect(screen.getByText('Standard user access')).toBeInTheDocument();
      expect(screen.getByText('Unauthenticated view')).toBeInTheDocument();
    });

    it('should display role icons in menu', () => {
      const { container } = render(<RoleSwitcher />);

      const dropdownMenu = container.querySelector('ul.dropdown-content');

      // All icons should be present in the dropdown
      expect(dropdownMenu?.textContent).toContain('👑'); // Admin
      expect(dropdownMenu?.textContent).toContain('🎨'); // Artist
      expect(dropdownMenu?.textContent).toContain('👤'); // Subscriber
      expect(dropdownMenu?.textContent).toContain('👻'); // Guest
    });
  });

  describe('interactions', () => {
    it('should call startSimulation when role is selected', () => {
      const { container } = render(<RoleSwitcher />);

      // Find buttons in the dropdown menu
      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const buttons = dropdownMenu?.querySelectorAll('button');
      // Find the subscriber button (index 2: admin=0, artist=1, subscriber=2, guest=3)
      const subscriberButton = buttons?.[2];
      fireEvent.click(subscriberButton!);

      expect(mockStartSimulation).toHaveBeenCalledWith('subscriber');
    });

    it('should call startSimulation with guest when guest is selected', () => {
      const { container } = render(<RoleSwitcher />);

      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const buttons = dropdownMenu?.querySelectorAll('button');
      // Guest is index 3
      const guestButton = buttons?.[3];
      fireEvent.click(guestButton!);

      expect(mockStartSimulation).toHaveBeenCalledWith('guest');
    });

    it('should call startSimulation with artist when artist is selected', () => {
      const { container } = render(<RoleSwitcher />);

      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const buttons = dropdownMenu?.querySelectorAll('button');
      // Artist is index 1
      const artistButton = buttons?.[1];
      fireEvent.click(artistButton!);

      expect(mockStartSimulation).toHaveBeenCalledWith('artist');
    });

    it('should call startSimulation with admin when admin is selected', () => {
      const { container } = render(<RoleSwitcher />);

      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const buttons = dropdownMenu?.querySelectorAll('button');
      // Admin is index 0
      const adminButton = buttons?.[0];
      fireEvent.click(adminButton!);

      expect(mockStartSimulation).toHaveBeenCalledWith('admin');
    });
  });

  describe('active state indicator', () => {
    it('should mark the current effective role as active', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'subscriber',
        isSimulating: true,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      // Find the dropdown menu and look for active button
      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const activeButton = dropdownMenu?.querySelector('button.active');
      expect(activeButton).toBeInTheDocument();
      expect(activeButton?.textContent).toContain('Subscriber');
    });

    it('should show checkmark for current role', () => {
      mockUseRoleSimulation.mockReturnValue({
        canSimulate: true,
        effectiveRole: 'artist',
        isSimulating: true,
        startSimulation: mockStartSimulation,
      });

      const { container } = render(<RoleSwitcher />);

      // Find the active button and check for the checkmark SVG
      const dropdownMenu = container.querySelector('ul.dropdown-content');
      const activeButton = dropdownMenu?.querySelector('button.active');
      const checkmark = activeButton?.querySelector('svg[fill="currentColor"]');
      expect(checkmark).toBeInTheDocument();
    });
  });
});
