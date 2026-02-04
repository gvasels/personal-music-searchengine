/**
 * Keyboard Shortcuts Modal
 * Displays all available keyboard shortcuts grouped by category
 */
import { useEffect, useRef } from 'react';
import { getShortcutsByCategory, categoryLabels, formatKeyForDisplay, type ShortcutConfig } from '@/lib/shortcuts/shortcuts';
import { useIsMac } from '@/hooks/useKeyboardShortcuts';

interface ShortcutsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

interface ShortcutRowProps {
  shortcut: ShortcutConfig;
  isMac: boolean;
}

function ShortcutRow({ shortcut, isMac }: ShortcutRowProps) {
  const keyDisplay = formatKeyForDisplay(shortcut, isMac);

  return (
    <div className="flex justify-between items-center py-1.5">
      <span className="text-base-content/70 text-sm">{shortcut.description}</span>
      <kbd className="kbd kbd-sm bg-base-200 text-base-content font-mono">
        {keyDisplay}
      </kbd>
    </div>
  );
}

export function ShortcutsModal({ isOpen, onClose }: ShortcutsModalProps) {
  const modalRef = useRef<HTMLDialogElement>(null);
  const isMac = useIsMac();
  const shortcutsByCategory = getShortcutsByCategory();

  // Focus trap and modal control
  useEffect(() => {
    const modal = modalRef.current;
    if (!modal) return;

    if (isOpen) {
      modal.showModal();
    } else {
      modal.close();
    }
  }, [isOpen]);

  // Handle Escape key to close
  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape' && isOpen) {
        onClose();
      }
    }

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  // Handle click outside to close
  function handleDialogClick(event: React.MouseEvent<HTMLDialogElement>) {
    const dialog = modalRef.current;
    if (!dialog) return;

    const rect = dialog.getBoundingClientRect();
    const isInDialog =
      event.clientX >= rect.left &&
      event.clientX <= rect.right &&
      event.clientY >= rect.top &&
      event.clientY <= rect.bottom;

    if (!isInDialog) {
      onClose();
    }
  }

  const categories = Object.entries(shortcutsByCategory).filter(
    ([, shortcuts]) => shortcuts.length > 0
  );

  return (
    <dialog
      ref={modalRef}
      className="modal modal-bottom sm:modal-middle"
      onClick={handleDialogClick}
    >
      <div className="modal-box max-w-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg">Keyboard Shortcuts</h3>
          <button
            className="btn btn-sm btn-circle btn-ghost"
            onClick={onClose}
            aria-label="Close"
          >
            ✕
          </button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          {categories.map(([category, shortcuts]) => (
            <div key={category}>
              <h4 className="font-semibold text-primary mb-2">
                {categoryLabels[category as keyof typeof categoryLabels]}
              </h4>
              <div className="space-y-1 border-l-2 border-base-300 pl-3">
                {shortcuts.map((shortcut) => (
                  <ShortcutRow
                    key={shortcut.action}
                    shortcut={shortcut}
                    isMac={isMac}
                  />
                ))}
              </div>
            </div>
          ))}
        </div>

        <div className="modal-action">
          <p className="text-xs text-base-content/50">
            Press <kbd className="kbd kbd-xs">?</kbd> anytime to show this help
          </p>
        </div>
      </div>

      {/* Backdrop */}
      <form method="dialog" className="modal-backdrop">
        <button onClick={onClose}>close</button>
      </form>
    </dialog>
  );
}

export default ShortcutsModal;
