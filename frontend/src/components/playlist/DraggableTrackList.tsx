/**
 * Draggable Track List for playlist reordering
 * Wraps tracks with @dnd-kit sortable context
 */
import { useState, useCallback } from 'react';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  TouchSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
  DragOverlay,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { restrictToVerticalAxis, restrictToParentElement } from '@dnd-kit/modifiers';
import { DraggableTrackItem } from './DraggableTrackItem';
import { useSelectionOptional } from '@/hooks/useSelection';
import { usePlayerStore } from '@/lib/store/playerStore';
import type { Track } from '@/types';

interface DraggableTrackListProps {
  tracks: Track[];
  onReorder: (newOrder: string[]) => void;
  onRemoveTrack: (trackId: string) => void;
  isReordering?: boolean;
}

export function DraggableTrackList({
  tracks,
  onReorder,
  onRemoveTrack,
  isReordering = false,
}: DraggableTrackListProps) {
  const [activeId, setActiveId] = useState<string | null>(null);
  const selection = useSelectionOptional();
  const { setQueue, currentTrack, isPlaying } = usePlayerStore();

  // Configure sensors for both pointer and touch
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8, // Require 8px movement before drag starts
      },
    }),
    useSensor(TouchSensor, {
      activationConstraint: {
        delay: 250, // Long press to initiate drag on touch
        tolerance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;

      setActiveId(null);

      if (over && active.id !== over.id) {
        const oldIndex = tracks.findIndex((t) => t.id === active.id);
        const newIndex = tracks.findIndex((t) => t.id === over.id);

        // Check if we're dragging selected items
        const selectedIds = selection?.selectedIds;
        const isSelectedDrag = selectedIds && selectedIds.size > 1 && selectedIds.has(active.id as string);

        if (isSelectedDrag) {
          // Multi-item drag: move all selected items to new position
          const selectedTrackIds = tracks
            .filter((t) => selectedIds.has(t.id))
            .map((t) => t.id);

          const unselectedTracks = tracks.filter((t) => !selectedIds.has(t.id));
          const insertIndex = unselectedTracks.findIndex((t) => t.id === over.id);

          // Insert selected tracks at the drop position
          const beforeInsert = unselectedTracks.slice(0, insertIndex >= 0 ? insertIndex : unselectedTracks.length).map((t) => t.id);
          const afterInsert = unselectedTracks.slice(insertIndex >= 0 ? insertIndex : unselectedTracks.length).map((t) => t.id);
          const newOrder: string[] = [
            ...beforeInsert,
            ...selectedTrackIds,
            ...afterInsert,
          ];

          onReorder(newOrder);
        } else {
          // Single item drag
          const newTracks = arrayMove(tracks, oldIndex, newIndex);
          onReorder(newTracks.map((t) => t.id));
        }
      }
    },
    [tracks, selection?.selectedIds, onReorder]
  );

  const handlePlayTrack = useCallback(
    (index: number) => {
      setQueue(tracks, index);
    },
    [tracks, setQueue]
  );

  const handleSelect = useCallback(
    (trackId: string, e: React.MouseEvent) => {
      if (!selection) return;

      const trackIds = tracks.map((t) => t.id);

      if (e.shiftKey && selection.lastSelectedId) {
        // Range select
        selection.selectRange(selection.lastSelectedId, trackId, trackIds);
      } else if (e.metaKey || e.ctrlKey) {
        // Toggle single item
        selection.toggle(trackId);
      } else {
        // Single select
        selection.select(trackId);
      }
    },
    [selection, tracks]
  );

  const activeTrack = activeId ? tracks.find((t) => t.id === activeId) : null;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      modifiers={[restrictToVerticalAxis, restrictToParentElement]}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className="overflow-x-auto">
        <table className="table table-zebra">
          <thead>
            <tr>
              <th className="w-8"></th>
              <th className="w-12">#</th>
              <th>Title</th>
              <th>Artist</th>
              <th>Album</th>
              <th>Duration</th>
              <th className="w-12"></th>
            </tr>
          </thead>
          <tbody>
            <SortableContext
              items={tracks.map((t) => t.id)}
              strategy={verticalListSortingStrategy}
            >
              {tracks.map((track, index) => (
                <DraggableTrackItem
                  key={track.id}
                  track={track}
                  index={index}
                  isCurrentTrack={currentTrack?.id === track.id}
                  isPlaying={isPlaying}
                  isSelected={selection?.isSelected(track.id) ?? false}
                  onPlay={() => handlePlayTrack(index)}
                  onRemove={() => onRemoveTrack(track.id)}
                  onSelect={selection ? (e) => handleSelect(track.id, e) : undefined}
                />
              ))}
            </SortableContext>
          </tbody>
        </table>
      </div>

      {/* Drag overlay for visual feedback */}
      <DragOverlay>
        {activeTrack && (
          <div className="bg-base-200 p-3 rounded-lg shadow-xl border border-primary">
            <div className="flex items-center gap-3">
              {activeTrack.coverArtUrl && (
                <img
                  src={activeTrack.coverArtUrl}
                  alt=""
                  className="w-10 h-10 rounded object-cover"
                />
              )}
              <div>
                <div className="font-medium">{activeTrack.title}</div>
                <div className="text-sm text-base-content/60">{activeTrack.artist}</div>
              </div>
              {selection && selection.selectedIds.size > 1 && selection.isSelected(activeTrack.id) && (
                <div className="badge badge-primary">
                  +{selection.selectedIds.size - 1}
                </div>
              )}
            </div>
          </div>
        )}
      </DragOverlay>

      {/* Reordering indicator */}
      {isReordering && (
        <div className="fixed bottom-24 left-1/2 -translate-x-1/2 z-50">
          <div className="alert alert-info shadow-lg">
            <span className="loading loading-spinner loading-sm" />
            <span>Saving new order...</span>
          </div>
        </div>
      )}
    </DndContext>
  );
}

export default DraggableTrackList;
