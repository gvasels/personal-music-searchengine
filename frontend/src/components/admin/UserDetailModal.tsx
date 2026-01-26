/**
 * UserDetailModal - Admin Panel Component
 * Modal for viewing and editing user details (role and status)
 */
import { useState, useEffect } from 'react';
import {
  useUserDetails,
  useUpdateUserRole,
  useUpdateUserStatus,
} from '../../hooks/useAdmin';
import type { UserRole } from '../../types';

interface UserDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  userId: string | null;
}

// Role options for the selector
const ROLE_OPTIONS: { value: UserRole; label: string; description: string }[] = [
  { value: 'guest', label: 'Guest', description: 'Limited access, browse only' },
  { value: 'subscriber', label: 'Subscriber', description: 'Can create playlists and follow artists' },
  { value: 'artist', label: 'Artist', description: 'Can upload tracks and manage profile' },
  { value: 'admin', label: 'Admin', description: 'Full administrative access' },
];

// Role badge colors
const ROLE_BADGE_COLORS: Record<UserRole, string> = {
  guest: 'badge-ghost',
  subscriber: 'badge-info',
  artist: 'badge-secondary',
  admin: 'badge-primary',
};

export function UserDetailModal({ isOpen, onClose, userId }: UserDetailModalProps) {
  const { data: user, isLoading, error } = useUserDetails(userId ?? undefined);
  const updateRoleMutation = useUpdateUserRole();
  const updateStatusMutation = useUpdateUserStatus();

  const [selectedRole, setSelectedRole] = useState<UserRole | ''>('');
  const [showRoleConfirm, setShowRoleConfirm] = useState(false);
  const [showStatusConfirm, setShowStatusConfirm] = useState(false);
  const [actionError, setActionError] = useState('');

  // Reset state when modal opens or user changes
  useEffect(() => {
    if (user) {
      setSelectedRole(user.role);
    }
    setShowRoleConfirm(false);
    setShowStatusConfirm(false);
    setActionError('');
  }, [user, isOpen]);

  const handleRoleChange = (newRole: UserRole) => {
    if (newRole !== user?.role) {
      setSelectedRole(newRole);
      setShowRoleConfirm(true);
    }
  };

  const confirmRoleChange = async () => {
    if (!userId || !selectedRole) return;

    setActionError('');
    try {
      await updateRoleMutation.mutateAsync({ userId, role: selectedRole });
      setShowRoleConfirm(false);
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to update role');
    }
  };

  const handleToggleStatus = () => {
    setShowStatusConfirm(true);
  };

  const confirmStatusChange = async () => {
    if (!userId || !user) return;

    setActionError('');
    try {
      await updateStatusMutation.mutateAsync({ userId, disabled: !user.disabled });
      setShowStatusConfirm(false);
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to update status');
    }
  };

  if (!isOpen) return null;

  const isUpdating = updateRoleMutation.isPending || updateStatusMutation.isPending;

  return (
    <dialog className="modal modal-open">
      <div className="modal-box max-w-2xl">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h3 className="font-bold text-lg">User Details</h3>
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-circle"
            onClick={onClose}
            aria-label="Close"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Content */}
        {isLoading ? (
          <div className="flex justify-center py-8">
            <span className="loading loading-spinner loading-lg" />
          </div>
        ) : error ? (
          <div className="alert alert-error">
            <span>Failed to load user details</span>
          </div>
        ) : user ? (
          <>
            {/* Error alert */}
            {actionError && (
              <div className="alert alert-error mb-4">
                <span>{actionError}</span>
              </div>
            )}

            {/* User Info Section */}
            <div className="card bg-base-200 p-4 mb-6">
              <div className="flex items-start gap-4">
                {/* Avatar placeholder */}
                <div className="avatar placeholder">
                  <div className="bg-neutral text-neutral-content rounded-full w-16">
                    <span className="text-2xl">
                      {user.displayName?.charAt(0).toUpperCase() || user.email.charAt(0).toUpperCase()}
                    </span>
                  </div>
                </div>

                <div className="flex-1">
                  <h4 className="font-semibold text-lg">
                    {user.displayName || 'No display name'}
                  </h4>
                  <p className="text-base-content/70">{user.email}</p>
                  <div className="flex items-center gap-2 mt-2">
                    <span className={`badge ${ROLE_BADGE_COLORS[user.role]}`}>
                      {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                    </span>
                    {user.disabled && (
                      <span className="badge badge-error badge-outline">Disabled</span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-6">
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Tracks</div>
                <div className="stat-value text-lg">{user.trackCount}</div>
              </div>
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Playlists</div>
                <div className="stat-value text-lg">{user.playlistCount}</div>
              </div>
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Albums</div>
                <div className="stat-value text-lg">{user.albumCount}</div>
              </div>
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Followers</div>
                <div className="stat-value text-lg">{user.followerCount}</div>
              </div>
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Following</div>
                <div className="stat-value text-lg">{user.followingCount}</div>
              </div>
              <div className="stat bg-base-200 rounded-lg p-3">
                <div className="stat-title text-xs">Storage</div>
                <div className="stat-value text-lg">
                  {(user.storageUsed / (1024 * 1024)).toFixed(1)} MB
                </div>
              </div>
            </div>

            {/* Dates */}
            <div className="text-sm text-base-content/60 mb-6">
              <p>
                Joined: {new Date(user.createdAt).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </p>
              {user.lastLoginAt && (
                <p>
                  Last login: {new Date(user.lastLoginAt).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </p>
              )}
            </div>

            {/* Role Selector */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text font-medium">User Role</span>
              </label>
              <select
                className="select select-bordered w-full"
                value={selectedRole}
                onChange={(e) => handleRoleChange(e.target.value as UserRole)}
                disabled={isUpdating}
              >
                {ROLE_OPTIONS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label} - {option.description}
                  </option>
                ))}
              </select>
            </div>

            {/* Role Change Confirmation */}
            {showRoleConfirm && (
              <div className="alert alert-warning mb-4">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div>
                  <p>Change role from <strong>{user.role}</strong> to <strong>{selectedRole}</strong>?</p>
                  <p className="text-sm opacity-80">This will update both the database and Cognito user groups.</p>
                </div>
                <div className="flex gap-2">
                  <button
                    className="btn btn-sm"
                    onClick={() => {
                      setSelectedRole(user.role);
                      setShowRoleConfirm(false);
                    }}
                    disabled={isUpdating}
                  >
                    Cancel
                  </button>
                  <button
                    className="btn btn-sm btn-warning"
                    onClick={confirmRoleChange}
                    disabled={isUpdating}
                  >
                    {updateRoleMutation.isPending ? (
                      <span className="loading loading-spinner loading-xs" />
                    ) : (
                      'Confirm'
                    )}
                  </button>
                </div>
              </div>
            )}

            {/* Status Toggle */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text font-medium">Account Status</span>
              </label>
              <div className="flex items-center gap-4">
                <label className="label cursor-pointer gap-3">
                  <span className="label-text">
                    {user.disabled ? 'Account Disabled' : 'Account Active'}
                  </span>
                  <input
                    type="checkbox"
                    className="toggle toggle-success"
                    checked={!user.disabled}
                    onChange={handleToggleStatus}
                    disabled={isUpdating}
                  />
                </label>
              </div>
            </div>

            {/* Status Change Confirmation */}
            {showStatusConfirm && (
              <div className={`alert ${user.disabled ? 'alert-success' : 'alert-error'} mb-4`}>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div>
                  <p>
                    {user.disabled
                      ? 'Re-enable this user account?'
                      : 'Disable this user account?'}
                  </p>
                  <p className="text-sm opacity-80">
                    {user.disabled
                      ? 'The user will be able to sign in again.'
                      : 'The user will not be able to sign in until re-enabled.'}
                  </p>
                </div>
                <div className="flex gap-2">
                  <button
                    className="btn btn-sm"
                    onClick={() => setShowStatusConfirm(false)}
                    disabled={isUpdating}
                  >
                    Cancel
                  </button>
                  <button
                    className={`btn btn-sm ${user.disabled ? 'btn-success' : 'btn-error'}`}
                    onClick={confirmStatusChange}
                    disabled={isUpdating}
                  >
                    {updateStatusMutation.isPending ? (
                      <span className="loading loading-spinner loading-xs" />
                    ) : user.disabled ? (
                      'Enable'
                    ) : (
                      'Disable'
                    )}
                  </button>
                </div>
              </div>
            )}
          </>
        ) : (
          <div className="py-8 text-center text-base-content/60">
            Select a user to view details
          </div>
        )}

        {/* Footer */}
        <div className="modal-action">
          <button type="button" className="btn" onClick={onClose}>
            Close
          </button>
        </div>
      </div>

      {/* Backdrop */}
      <form method="dialog" className="modal-backdrop">
        <button onClick={onClose}>close</button>
      </form>
    </dialog>
  );
}
