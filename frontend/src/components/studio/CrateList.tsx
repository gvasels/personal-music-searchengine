/**
 * CrateList Component
 * Displays user's DJ crates with create/manage functionality
 */
import { useState } from 'react';
import { useCrates, useCreateCrate, useDeleteCrate } from '@/hooks/useCrates';
import type { Crate } from '@/types';

interface CrateListProps {
  onSelectCrate?: (crate: Crate) => void;
  selectedCrateId?: string;
}

export function CrateList({ onSelectCrate, selectedCrateId }: CrateListProps) {
  const { data: crates, isLoading, isLocked, isFeatureEnabled } = useCrates();
  const createCrate = useCreateCrate();
  const deleteCrate = useDeleteCrate();
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newCrateName, setNewCrateName] = useState('');

  if (!isFeatureEnabled && isLocked) {
    return (
      <div className="p-4 text-center">
        <p className="text-base-content/70 mb-4">
          Crates are available for Artist users.
        </p>
        <span className="badge badge-outline">Artist role required</span>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="space-y-2 p-2">
        {[1, 2, 3].map((i) => (
          <div key={i} className="h-12 bg-base-300 animate-pulse rounded" />
        ))}
      </div>
    );
  }

  const handleCreateCrate = async () => {
    if (!newCrateName.trim()) return;
    await createCrate.mutateAsync({ name: newCrateName.trim() });
    setNewCrateName('');
    setShowCreateModal(false);
  };

  const handleDeleteCrate = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm('Are you sure you want to delete this crate?')) {
      await deleteCrate.mutateAsync(id);
    }
  };

  return (
    <div className="p-2">
      {/* Header */}
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-semibold text-sm">Crates</h3>
        <button
          className="btn btn-ghost btn-xs"
          onClick={() => setShowCreateModal(true)}
        >
          + New
        </button>
      </div>

      {/* Crate list */}
      <div className="space-y-1">
        {crates?.length === 0 ? (
          <p className="text-xs text-base-content/50 text-center py-4">
            No crates yet. Create one to organize your tracks.
          </p>
        ) : (
          crates?.map((crate) => (
            <div
              key={crate.id}
              className={`flex items-center gap-2 p-2 rounded cursor-pointer hover:bg-base-200 ${
                selectedCrateId === crate.id ? 'bg-base-200' : ''
              }`}
              onClick={() => onSelectCrate?.(crate)}
            >
              <span
                className="w-3 h-3 rounded-full flex-shrink-0"
                style={{ backgroundColor: crate.color || '#888' }}
              />
              <span className="flex-1 truncate text-sm">{crate.name}</span>
              <span className="text-xs text-base-content/50">
                {crate.trackCount}
              </span>
              <button
                className="btn btn-ghost btn-xs opacity-0 group-hover:opacity-100"
                onClick={(e) => handleDeleteCrate(crate.id, e)}
              >
                X
              </button>
            </div>
          ))
        )}
      </div>

      {/* Create modal */}
      {showCreateModal && (
        <div className="modal modal-open">
          <div className="modal-box max-w-xs">
            <h3 className="font-bold text-lg">New Crate</h3>
            <div className="py-4">
              <input
                type="text"
                placeholder="Crate name"
                className="input input-bordered w-full"
                value={newCrateName}
                onChange={(e) => setNewCrateName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleCreateCrate()}
                autoFocus
              />
            </div>
            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => setShowCreateModal(false)}
              >
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={handleCreateCrate}
                disabled={!newCrateName.trim() || createCrate.isPending}
              >
                {createCrate.isPending ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
          <div
            className="modal-backdrop"
            onClick={() => setShowCreateModal(false)}
          />
        </div>
      )}
    </div>
  );
}
