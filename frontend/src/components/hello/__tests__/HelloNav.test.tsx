import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '../../../test/test-utils';
import { HelloNav } from '../HelloNav';

describe('HelloNav', () => {
  it('renders navbar with app name', () => {
    render(<HelloNav searchValue="" onSearchChange={vi.fn()} />);
    expect(screen.getByText('MusicSearch')).toBeInTheDocument();
  });

  it('renders search input in navbar', () => {
    render(<HelloNav searchValue="" onSearchChange={vi.fn()} />);
    expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
  });

  it('renders theme toggle button', () => {
    const { container } = render(
      <HelloNav searchValue="" onSearchChange={vi.fn()} />
    );
    const themeController = container.querySelector('.theme-controller');
    expect(themeController).toBeInTheDocument();
  });

  it('renders dock for mobile viewport', () => {
    const { container } = render(
      <HelloNav searchValue="" onSearchChange={vi.fn()} />
    );
    const dock = container.querySelector('.dock');
    expect(dock).toBeInTheDocument();
  });
});
