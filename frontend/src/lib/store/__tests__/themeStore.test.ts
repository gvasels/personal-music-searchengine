/**
 * Theme Store Tests - REQ-5.2
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { act, renderHook } from '@testing-library/react';
import { useThemeStore } from '@/lib/store/themeStore';

describe('Theme Store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset store to initial state
    useThemeStore.setState({ theme: 'dark' });
  });

  describe('Initial State', () => {
    it('REQ-5.2: should default to dark theme', () => {
      const { result } = renderHook(() => useThemeStore());
      expect(result.current.theme).toBe('dark');
    });
  });

  describe('toggleTheme', () => {
    it('REQ-5.2: should toggle from dark to light', () => {
      const { result } = renderHook(() => useThemeStore());

      act(() => {
        result.current.toggleTheme();
      });

      expect(result.current.theme).toBe('light');
    });

    it('should toggle from light to dark', () => {
      const { result } = renderHook(() => useThemeStore());

      act(() => {
        result.current.toggleTheme();
        result.current.toggleTheme();
      });

      expect(result.current.theme).toBe('dark');
    });

    it('should apply theme to document element', () => {
      const { result } = renderHook(() => useThemeStore());

      act(() => {
        result.current.toggleTheme();
      });

      expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    });
  });
});
