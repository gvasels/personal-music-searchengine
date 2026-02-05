/**
 * SearchInput Component Tests - Hello World Feature (Red Phase)
 *
 * Tests for a controlled text input component that accepts value, onChange, and optional placeholder.
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { SearchInput } from '../SearchInput';

describe('SearchInput', () => {
  it('renders an input element', () => {
    render(<SearchInput value="" onChange={() => {}} />);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('displays the current value', () => {
    render(<SearchInput value="jazz" onChange={() => {}} />);
    expect(screen.getByDisplayValue('jazz')).toBeInTheDocument();
  });

  it('calls onChange when user types', async () => {
    const onChange = vi.fn();
    const { user } = render(<SearchInput value="" onChange={onChange} />);
    await user.type(screen.getByRole('textbox'), 'a');
    expect(onChange).toHaveBeenCalled();
  });

  it('shows placeholder text', () => {
    render(
      <SearchInput value="" onChange={() => {}} placeholder="Search tracks..." />
    );
    expect(screen.getByPlaceholderText('Search tracks...')).toBeInTheDocument();
  });
});
