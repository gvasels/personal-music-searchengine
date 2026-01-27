/**
 * SimulationBanner Component Tests - Admin Role Switching Feature
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SimulationBanner } from '../SimulationBanner';

// Mock useRoleSimulation hook
const mockStopSimulation = vi.fn();
const mockUseRoleSimulation = vi.fn();

vi.mock('../../../hooks/useRoleSimulation', () => ({
  useRoleSimulation: () => mockUseRoleSimulation(),
}));

describe('SimulationBanner', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseRoleSimulation.mockReturnValue({
      isSimulating: true,
      effectiveRole: 'subscriber',
      stopSimulation: mockStopSimulation,
    });
  });

  describe('rendering', () => {
    it('should render when simulating', () => {
      render(<SimulationBanner />);

      expect(screen.getByText(/viewing as/i)).toBeInTheDocument();
    });

    it('should not render when not simulating', () => {
      mockUseRoleSimulation.mockReturnValue({
        isSimulating: false,
        effectiveRole: 'admin',
        stopSimulation: mockStopSimulation,
      });

      const { container } = render(<SimulationBanner />);

      expect(container).toBeEmptyDOMElement();
    });
  });

  describe('role labels', () => {
    it('should display "Subscriber" for subscriber role', () => {
      mockUseRoleSimulation.mockReturnValue({
        isSimulating: true,
        effectiveRole: 'subscriber',
        stopSimulation: mockStopSimulation,
      });

      render(<SimulationBanner />);

      expect(screen.getByText('Viewing as Subscriber')).toBeInTheDocument();
    });

    it('should display "Guest" for guest role', () => {
      mockUseRoleSimulation.mockReturnValue({
        isSimulating: true,
        effectiveRole: 'guest',
        stopSimulation: mockStopSimulation,
      });

      render(<SimulationBanner />);

      expect(screen.getByText('Viewing as Guest')).toBeInTheDocument();
    });

    it('should display "Artist" for artist role', () => {
      mockUseRoleSimulation.mockReturnValue({
        isSimulating: true,
        effectiveRole: 'artist',
        stopSimulation: mockStopSimulation,
      });

      render(<SimulationBanner />);

      expect(screen.getByText('Viewing as Artist')).toBeInTheDocument();
    });

    it('should handle unknown roles gracefully', () => {
      mockUseRoleSimulation.mockReturnValue({
        isSimulating: true,
        // Test unknown role fallback behavior
        effectiveRole: 'unknown' as 'guest' | 'subscriber' | 'artist' | 'admin',
        stopSimulation: mockStopSimulation,
      });

      render(<SimulationBanner />);

      // Should fall back to the raw role value
      expect(screen.getByText('Viewing as unknown')).toBeInTheDocument();
    });
  });

  describe('exit button', () => {
    it('should render exit button', () => {
      render(<SimulationBanner />);

      expect(screen.getByRole('button', { name: /exit simulation/i })).toBeInTheDocument();
    });

    it('should call stopSimulation when clicked', () => {
      render(<SimulationBanner />);

      const exitButton = screen.getByRole('button', { name: /exit simulation/i });
      fireEvent.click(exitButton);

      expect(mockStopSimulation).toHaveBeenCalledTimes(1);
    });
  });

  describe('styling', () => {
    it('should have warning background color', () => {
      const { container } = render(<SimulationBanner />);

      const banner = container.firstChild as HTMLElement;
      expect(banner?.className).toContain('bg-warning');
    });

    it('should be sticky at top', () => {
      const { container } = render(<SimulationBanner />);

      const banner = container.firstChild as HTMLElement;
      expect(banner?.className).toContain('sticky');
      expect(banner?.className).toContain('top-0');
    });

    it('should have high z-index', () => {
      const { container } = render(<SimulationBanner />);

      const banner = container.firstChild as HTMLElement;
      expect(banner?.className).toContain('z-50');
    });
  });

  describe('icon', () => {
    it('should display eye icon', () => {
      render(<SimulationBanner />);

      const banner = screen.getByText(/viewing as/i).closest('div');
      const svg = banner?.querySelector('svg');
      expect(svg).toBeInTheDocument();
    });
  });
});
