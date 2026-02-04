import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UserSearchForm } from '../UserSearchForm';

describe('UserSearchForm', () => {
  const mockOnSearch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('rendering', () => {
    it('renders search input', () => {
      render(<UserSearchForm onSearch={mockOnSearch} />);
      expect(screen.getByRole('textbox', { name: /search users/i })).toBeInTheDocument();
    });

    it('renders with default placeholder', () => {
      render(<UserSearchForm onSearch={mockOnSearch} />);
      expect(screen.getByPlaceholderText(/search by email or name/i)).toBeInTheDocument();
    });

    it('renders with custom placeholder', () => {
      render(<UserSearchForm onSearch={mockOnSearch} placeholder="Find users..." />);
      expect(screen.getByPlaceholderText('Find users...')).toBeInTheDocument();
    });

    it('renders helper text', () => {
      render(<UserSearchForm onSearch={mockOnSearch} />);
      expect(screen.getByText(/minimum 1 character/i)).toBeInTheDocument();
    });
  });

  describe('search behavior', () => {
    it('calls onSearch when typing', async () => {
      const user = userEvent.setup();
      render(<UserSearchForm onSearch={mockOnSearch} />);

      await user.type(screen.getByRole('textbox'), 'test');
      expect(mockOnSearch).toHaveBeenCalledWith('t');
      expect(mockOnSearch).toHaveBeenCalledWith('te');
      expect(mockOnSearch).toHaveBeenCalledWith('tes');
      expect(mockOnSearch).toHaveBeenCalledWith('test');
    });

    it('shows clear button when input has value', async () => {
      const user = userEvent.setup();
      render(<UserSearchForm onSearch={mockOnSearch} />);

      expect(screen.queryByRole('button', { name: /clear/i })).not.toBeInTheDocument();

      await user.type(screen.getByRole('textbox'), 'query');
      expect(screen.getByRole('button', { name: /clear/i })).toBeInTheDocument();
    });

    it('clears input and calls onSearch when clear button clicked', async () => {
      const user = userEvent.setup();
      render(<UserSearchForm onSearch={mockOnSearch} />);

      await user.type(screen.getByRole('textbox'), 'query');
      mockOnSearch.mockClear();

      await user.click(screen.getByRole('button', { name: /clear/i }));

      expect(screen.getByRole('textbox')).toHaveValue('');
      expect(mockOnSearch).toHaveBeenCalledWith('');
    });
  });

  describe('loading state', () => {
    it('shows loading spinner when isLoading is true', () => {
      render(<UserSearchForm onSearch={mockOnSearch} isLoading={true} />);
      expect(screen.getByRole('textbox').parentElement?.querySelector('.loading')).toBeInTheDocument();
    });

    it('hides clear button when loading', async () => {
      const user = userEvent.setup();
      const { rerender } = render(<UserSearchForm onSearch={mockOnSearch} />);

      await user.type(screen.getByRole('textbox'), 'test');
      expect(screen.getByRole('button', { name: /clear/i })).toBeInTheDocument();

      rerender(<UserSearchForm onSearch={mockOnSearch} isLoading={true} />);
      expect(screen.queryByRole('button', { name: /clear/i })).not.toBeInTheDocument();
    });
  });

  describe('form submission', () => {
    it('calls onSearch on form submit', async () => {
      const user = userEvent.setup();
      render(<UserSearchForm onSearch={mockOnSearch} />);

      await user.type(screen.getByRole('textbox'), 'search term');
      mockOnSearch.mockClear();

      await user.keyboard('{Enter}');
      expect(mockOnSearch).toHaveBeenCalledWith('search term');
    });
  });
});
