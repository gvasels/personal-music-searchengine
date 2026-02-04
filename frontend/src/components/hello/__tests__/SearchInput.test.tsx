import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '../../../test/test-utils';
import { SearchInput } from '../SearchInput';

describe('SearchInput', () => {
  it('renders input element with placeholder text', () => {
    render(<SearchInput value="" onChange={vi.fn()} />);
    expect(screen.getByPlaceholderText('Search tracks, artists...')).toBeInTheDocument();
  });

  it('calls onChange when user types', async () => {
    const onChange = vi.fn();
    const { user } = render(<SearchInput value="" onChange={onChange} />);

    const input = screen.getByPlaceholderText('Search tracks, artists...');
    await user.type(input, 'a');

    expect(onChange).toHaveBeenCalled();
  });

  it('calls onSubmit on Enter key press', async () => {
    const onSubmit = vi.fn();
    const { user } = render(
      <SearchInput value="test" onChange={vi.fn()} onSubmit={onSubmit} />
    );

    const input = screen.getByPlaceholderText('Search tracks, artists...');
    await user.type(input, '{Enter}');

    expect(onSubmit).toHaveBeenCalled();
  });

  it('clears input on Escape key press', async () => {
    const onChange = vi.fn();
    const { user } = render(<SearchInput value="test" onChange={onChange} />);

    const input = screen.getByPlaceholderText('Search tracks, artists...');
    await user.type(input, '{Escape}');

    expect(onChange).toHaveBeenCalledWith('');
  });

  it('autoFocus prop focuses input on mount', () => {
    render(<SearchInput value="" onChange={vi.fn()} autoFocus />);
    const input = screen.getByPlaceholderText('Search tracks, artists...');
    expect(input).toHaveFocus();
  });
});
