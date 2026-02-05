/**
 * SearchInput Component Tests - Hello World Local Dev Feature
 *
 * Tests for the controlled search input component.
 * MUST FAIL until SearchInput component is implemented.
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';

describe('SearchInput Component', () => {
  it('renders an input element', async () => {
    const { SearchInput } = await import('../SearchInput');
    const onChange = vi.fn();
    render(<SearchInput value="" onChange={onChange} />);

    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('displays the current value', async () => {
    const { SearchInput } = await import('../SearchInput');
    const onChange = vi.fn();
    render(<SearchInput value="jazz" onChange={onChange} />);

    expect(screen.getByRole('textbox')).toHaveValue('jazz');
  });

  it('calls onChange when user types', async () => {
    const { SearchInput } = await import('../SearchInput');
    const onChange = vi.fn();
    const { user } = render(<SearchInput value="" onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.type(input, 'rock');

    expect(onChange).toHaveBeenCalled();
  });

  it('shows placeholder text', async () => {
    const { SearchInput } = await import('../SearchInput');
    const onChange = vi.fn();
    render(
      <SearchInput value="" onChange={onChange} placeholder="Search tracks..." />
    );

    expect(screen.getByPlaceholderText('Search tracks...')).toBeInTheDocument();
  });

  it('shows default placeholder when none provided', async () => {
    const { SearchInput } = await import('../SearchInput');
    const onChange = vi.fn();
    render(<SearchInput value="" onChange={onChange} />);

    // Default placeholder should exist
    const input = screen.getByRole('textbox');
    expect(input).toHaveAttribute('placeholder');
  });
});
