/**
 * EditArtistProfileModal Component - Create or edit artist profile
 */
import { useState, useEffect } from 'react';
import { useCreateArtistProfile, useUpdateArtistProfile } from '../../hooks/useArtistProfiles';
import type { ArtistProfile } from '../../types';

interface EditArtistProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  profile?: ArtistProfile | null;
  onSuccess?: (profile: ArtistProfile) => void;
}

export function EditArtistProfileModal({
  isOpen,
  onClose,
  profile,
  onSuccess,
}: EditArtistProfileModalProps) {
  const [displayName, setDisplayName] = useState('');
  const [bio, setBio] = useState('');
  const [location, setLocation] = useState('');
  const [website, setWebsite] = useState('');
  const [error, setError] = useState('');

  const createMutation = useCreateArtistProfile();
  const updateMutation = useUpdateArtistProfile(profile?.userId || '');

  const isEditing = !!profile;
  const isLoading = createMutation.isPending || updateMutation.isPending;

  useEffect(() => {
    if (profile) {
      setDisplayName(profile.displayName || '');
      setBio(profile.bio || '');
      setLocation(profile.location || '');
      setWebsite(profile.website || '');
    } else {
      setDisplayName('');
      setBio('');
      setLocation('');
      setWebsite('');
    }
    setError('');
  }, [profile, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!displayName.trim()) {
      setError('Display name is required');
      return;
    }

    try {
      let result: ArtistProfile;

      if (isEditing) {
        result = await updateMutation.mutateAsync({
          displayName: displayName.trim(),
          bio: bio.trim() || undefined,
          location: location.trim() || undefined,
          website: website.trim() || undefined,
        });
      } else {
        result = await createMutation.mutateAsync({
          displayName: displayName.trim(),
          bio: bio.trim() || undefined,
          location: location.trim() || undefined,
          website: website.trim() || undefined,
        });
      }

      onSuccess?.(result);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save profile');
    }
  };

  if (!isOpen) return null;

  return (
    <dialog className="modal modal-open">
      <div className="modal-box">
        <h3 className="font-bold text-lg mb-4">
          {isEditing ? 'Edit Artist Profile' : 'Create Artist Profile'}
        </h3>

        <form onSubmit={handleSubmit}>
          {error && (
            <div className="alert alert-error mb-4">
              <span>{error}</span>
            </div>
          )}

          <div className="form-control mb-4">
            <label className="label">
              <span className="label-text">Display Name *</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              placeholder="Your artist name"
              maxLength={100}
              required
            />
          </div>

          <div className="form-control mb-4">
            <label className="label">
              <span className="label-text">Bio</span>
            </label>
            <textarea
              className="textarea textarea-bordered h-24"
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Tell us about yourself..."
              maxLength={500}
            />
            <label className="label">
              <span className="label-text-alt">{bio.length}/500</span>
            </label>
          </div>

          <div className="form-control mb-4">
            <label className="label">
              <span className="label-text">Location</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              placeholder="City, Country"
              maxLength={100}
            />
          </div>

          <div className="form-control mb-4">
            <label className="label">
              <span className="label-text">Website</span>
            </label>
            <input
              type="url"
              className="input input-bordered"
              value={website}
              onChange={(e) => setWebsite(e.target.value)}
              placeholder="https://your-website.com"
            />
          </div>

          <div className="modal-action">
            <button type="button" className="btn" onClick={onClose} disabled={isLoading}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary" disabled={isLoading}>
              {isLoading ? (
                <span className="loading loading-spinner loading-sm"></span>
              ) : isEditing ? (
                'Save Changes'
              ) : (
                'Create Profile'
              )}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button onClick={onClose}>close</button>
      </form>
    </dialog>
  );
}
