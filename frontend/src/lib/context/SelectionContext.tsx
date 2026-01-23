/**
 * Selection Context for batch operations
 * Manages multi-select state for tracks and other items
 */
import { createContext, useContext, useMemo, useState, useCallback, type ReactNode } from 'react';

export interface SelectionContextValue {
  /** Set of selected item IDs */
  selectedIds: Set<string>;
  /** Number of selected items */
  selectionCount: number;
  /** Check if an item is selected */
  isSelected: (id: string) => boolean;
  /** Select a single item (replaces selection) */
  select: (id: string) => void;
  /** Toggle selection of an item */
  toggle: (id: string) => void;
  /** Deselect an item */
  deselect: (id: string) => void;
  /** Select a range of items (for shift-click) */
  selectRange: (fromId: string, toId: string, orderedIds: string[]) => void;
  /** Select all items from the provided list */
  selectAll: (ids: string[]) => void;
  /** Clear all selection */
  clearSelection: () => void;
  /** Last selected item ID (for range selection anchor) */
  lastSelectedId: string | null;
}

const SelectionContext = createContext<SelectionContextValue | null>(null);

interface SelectionProviderProps {
  children: ReactNode;
}

export function SelectionProvider({ children }: SelectionProviderProps) {
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [lastSelectedId, setLastSelectedId] = useState<string | null>(null);

  const isSelected = useCallback(
    (id: string) => selectedIds.has(id),
    [selectedIds]
  );

  const select = useCallback((id: string) => {
    setSelectedIds(new Set([id]));
    setLastSelectedId(id);
  }, []);

  const toggle = useCallback((id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
    setLastSelectedId(id);
  }, []);

  const deselect = useCallback((id: string) => {
    setSelectedIds((prev) => {
      if (!prev.has(id)) return prev;
      const next = new Set(prev);
      next.delete(id);
      return next;
    });
  }, []);

  const selectRange = useCallback((fromId: string, toId: string, orderedIds: string[]) => {
    const fromIndex = orderedIds.indexOf(fromId);
    const toIndex = orderedIds.indexOf(toId);

    if (fromIndex === -1 || toIndex === -1) return;

    const startIndex = Math.min(fromIndex, toIndex);
    const endIndex = Math.max(fromIndex, toIndex);

    const rangeIds = orderedIds.slice(startIndex, endIndex + 1);
    setSelectedIds((prev) => {
      const next = new Set(prev);
      for (const id of rangeIds) {
        next.add(id);
      }
      return next;
    });
    setLastSelectedId(toId);
  }, []);

  const selectAll = useCallback((ids: string[]) => {
    setSelectedIds(new Set(ids));
    if (ids.length > 0) {
      setLastSelectedId(ids[ids.length - 1]);
    }
  }, []);

  const clearSelection = useCallback(() => {
    setSelectedIds(new Set());
    setLastSelectedId(null);
  }, []);

  const value = useMemo<SelectionContextValue>(
    () => ({
      selectedIds,
      selectionCount: selectedIds.size,
      isSelected,
      select,
      toggle,
      deselect,
      selectRange,
      selectAll,
      clearSelection,
      lastSelectedId,
    }),
    [selectedIds, isSelected, select, toggle, deselect, selectRange, selectAll, clearSelection, lastSelectedId]
  );

  return (
    <SelectionContext.Provider value={value}>
      {children}
    </SelectionContext.Provider>
  );
}

/**
 * Hook to access selection context
 * @throws Error if used outside SelectionProvider
 */
export function useSelection(): SelectionContextValue {
  const context = useContext(SelectionContext);
  if (!context) {
    throw new Error('useSelection must be used within a SelectionProvider');
  }
  return context;
}

/**
 * Optional hook that returns null if not in SelectionProvider
 */
export function useSelectionOptional(): SelectionContextValue | null {
  return useContext(SelectionContext);
}

export default SelectionContext;
