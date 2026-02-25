import { useMutation } from '@tanstack/react-query';
import { reprocessTrack } from '@/lib/api/tracks';
import toast from 'react-hot-toast';

interface ReprocessButtonProps {
  trackId: string;
  onSuccess?: () => void;
}

export function ReprocessButton({ trackId, onSuccess }: ReprocessButtonProps) {
  const mutation = useMutation({
    mutationFn: () => reprocessTrack(trackId),
    onSuccess: () => {
      toast.success('Reprocess started');
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Reprocess failed: ${error.message}`);
    },
  });

  return (
    <button
      type="button"
      className={`btn btn-ghost btn-xs${mutation.isPending ? ' loading' : ''}`}
      onClick={(e) => {
        e.stopPropagation();
        mutation.mutate();
      }}
      disabled={mutation.isPending}
      aria-label="Reprocess track"
      title="Reprocess AI analysis"
    >
      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
    </button>
  );
}
