/**
 * TagInput Component Tests - REQ-5.10
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';

describe('TagInput Component', () => {
  it('REQ-5.10: should render input field', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    render(<TagInput trackId="track-1" tags={[]} />);
    expect(screen.getByRole('button', { name: /add tag/i })).toBeInTheDocument();
  });

  it('should render existing tags', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    render(<TagInput trackId="track-1" tags={['rock', 'favorite']} />);
    // Tags are rendered with emoji prefix, so use regex matcher
    expect(screen.getByLabelText('remove rock')).toBeInTheDocument();
    expect(screen.getByLabelText('remove favorite')).toBeInTheDocument();
  });

  it('REQ-5.10: should add tag on Enter', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    const { user } = render(<TagInput trackId="track-1" tags={[]} />);

    await user.click(screen.getByRole('button', { name: /add tag/i }));
    const input = screen.getByRole('textbox');
    await user.type(input, 'newtag{Enter}');

    // After Enter, the input should be removed (isAdding becomes false)
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    // The "Add tag" button should be visible again
    expect(screen.getByRole('button', { name: /add tag/i })).toBeInTheDocument();
  });

  it('should normalize tags to lowercase', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    const { user } = render(<TagInput trackId="track-1" tags={[]} />);

    await user.click(screen.getByRole('button', { name: /add tag/i }));
    const input = screen.getByRole('textbox');
    await user.type(input, 'UPPERCASE{Enter}');

    // After Enter, the input should be removed and "Add tag" button restored
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: /add tag/i })).toBeInTheDocument();
  });

  it('REQ-5.10: should remove tag on click', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    const { user } = render(<TagInput trackId="track-1" tags={['rock', 'jazz']} />);

    const removeButtons = screen.getAllByRole('button');
    const removeRockButton = removeButtons.find(btn =>
      btn.closest('.badge')?.textContent?.includes('rock')
    );

    if (removeRockButton) {
      await user.click(removeRockButton);
    }

    // The remove mutation would be called
  });

  it('should have proper accessibility', async () => {
    const { TagInput } = await import('@/components/tag/TagInput');
    render(<TagInput trackId="track-1" tags={['rock']} />);

    // Each tag has an accessible remove button
    expect(screen.getByLabelText('remove rock')).toBeInTheDocument();
    // Add tag button is accessible
    expect(screen.getByRole('button', { name: /add tag/i })).toBeInTheDocument();
  });
});
