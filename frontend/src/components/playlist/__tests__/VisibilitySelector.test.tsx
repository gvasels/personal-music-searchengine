/**
 * VisibilitySelector Components Tests - Global User Type Feature
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import {
  VisibilitySelector,
  VisibilityBadge,
  VisibilityRadioGroup,
} from '../VisibilitySelector';
import type { PlaylistVisibility } from '../../../types';

describe('VisibilitySelector', () => {
  it('should render select with all visibility options', () => {
    const onChange = vi.fn();
    render(<VisibilitySelector value="private" onChange={onChange} />);

    const select = screen.getByRole('combobox');
    expect(select).toBeInTheDocument();

    // Check all options are present
    expect(screen.getByText(/Private/)).toBeInTheDocument();
    expect(screen.getByText(/Unlisted/)).toBeInTheDocument();
    expect(screen.getByText(/Public/)).toBeInTheDocument();
  });

  it('should show selected value', () => {
    const onChange = vi.fn();
    render(<VisibilitySelector value="public" onChange={onChange} />);

    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('public');
  });

  it('should call onChange when selection changes', () => {
    const onChange = vi.fn();
    render(<VisibilitySelector value="private" onChange={onChange} />);

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: 'public' } });

    expect(onChange).toHaveBeenCalledWith('public');
  });

  it('should be disabled when disabled prop is true', () => {
    const onChange = vi.fn();
    render(<VisibilitySelector value="private" onChange={onChange} disabled />);

    const select = screen.getByRole('combobox');
    expect(select).toBeDisabled();
  });

  it('should apply size classes', () => {
    const onChange = vi.fn();
    const { rerender } = render(
      <VisibilitySelector value="private" onChange={onChange} size="sm" />
    );
    expect(screen.getByRole('combobox')).toHaveClass('select-sm');

    rerender(<VisibilitySelector value="private" onChange={onChange} size="lg" />);
    expect(screen.getByRole('combobox')).toHaveClass('select-lg');
  });
});

describe('VisibilityBadge', () => {
  it('should render badge for private visibility', () => {
    render(<VisibilityBadge visibility="private" />);
    expect(screen.getByText(/Private/)).toBeInTheDocument();
    expect(screen.getByText(/Private/).closest('span')).toHaveClass('badge-ghost');
  });

  it('should render badge for unlisted visibility', () => {
    render(<VisibilityBadge visibility="unlisted" />);
    expect(screen.getByText(/Unlisted/)).toBeInTheDocument();
    expect(screen.getByText(/Unlisted/).closest('span')).toHaveClass('badge-warning');
  });

  it('should render badge for public visibility', () => {
    render(<VisibilityBadge visibility="public" />);
    expect(screen.getByText(/Public/)).toBeInTheDocument();
    expect(screen.getByText(/Public/).closest('span')).toHaveClass('badge-success');
  });

  it('should apply size classes', () => {
    const { rerender } = render(<VisibilityBadge visibility="public" size="sm" />);
    expect(screen.getByText(/Public/).closest('span')).toHaveClass('badge-sm');

    rerender(<VisibilityBadge visibility="public" size="lg" />);
    expect(screen.getByText(/Public/).closest('span')).toHaveClass('badge-lg');
  });
});

describe('VisibilityRadioGroup', () => {
  it('should render all radio options', () => {
    const onChange = vi.fn();
    render(<VisibilityRadioGroup value="private" onChange={onChange} />);

    const radios = screen.getAllByRole('radio');
    expect(radios).toHaveLength(3);

    expect(screen.getByText('Only you can see this playlist')).toBeInTheDocument();
    expect(screen.getByText('Anyone with the link can see')).toBeInTheDocument();
    expect(screen.getByText('Visible to everyone')).toBeInTheDocument();
  });

  it('should show selected option as checked', () => {
    const onChange = vi.fn();
    render(<VisibilityRadioGroup value="unlisted" onChange={onChange} />);

    const radios = screen.getAllByRole('radio') as HTMLInputElement[];
    const unlistedRadio = radios.find((r) => r.closest('label')?.textContent?.includes('Unlisted'));
    expect(unlistedRadio?.checked).toBe(true);
  });

  it('should call onChange when option is selected', () => {
    const onChange = vi.fn();
    render(<VisibilityRadioGroup value="private" onChange={onChange} />);

    const radios = screen.getAllByRole('radio');
    const publicRadio = radios[2]; // Public is the third option
    fireEvent.click(publicRadio);

    expect(onChange).toHaveBeenCalledWith('public');
  });

  it('should disable all radios when disabled prop is true', () => {
    const onChange = vi.fn();
    render(<VisibilityRadioGroup value="private" onChange={onChange} disabled />);

    const radios = screen.getAllByRole('radio');
    radios.forEach((radio) => {
      expect(radio).toBeDisabled();
    });
  });
});
