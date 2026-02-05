/**
 * SearchInput Component Tests - TDD Red Phase
 *
 * Tests for the hello search input component.
 * Source module (SearchInput.tsx) does NOT exist yet - these tests MUST fail.
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SearchInput } from '../SearchInput';

describe('SearchInput', () => {
  it('renders input with placeholder text', () => {
    render(<SearchInput value="" onChange={vi.fn()} />);

    const input = screen.getByPlaceholderText(/search/i);
    expect(input).toBeInTheDocument();
  });

  it('calls onChange when user types', () => {
    const handleChange = vi.fn();
    render(<SearchInput value="" onChange={handleChange} />);

    const input = screen.getByPlaceholderText(/search/i);
    fireEvent.change(input, { target: { value: 'hello world' } });

    expect(handleChange).toHaveBeenCalledWith('hello world');
  });

  it('displays current value', () => {
    render(<SearchInput value="current search" onChange={vi.fn()} />);

    const input = screen.getByDisplayValue('current search');
    expect(input).toBeInTheDocument();
  });
});
