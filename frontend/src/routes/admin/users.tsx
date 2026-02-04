/**
 * Admin Users Page - Admin Panel & Track Visibility Feature
 * User management page for administrators
 */
import { useState, useEffect } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { toast } from 'react-hot-toast';
import { useAuth } from '../../hooks/useAuth';
import { useSearchUsers } from '../../hooks/useAdmin';
import { UserSearchForm, UserCard, UserDetailModal } from '../../components/admin';

export default function AdminUsersPage() {
  const navigate = useNavigate();
  const { isAdmin, isLoading: authLoading } = useAuth();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const { data: searchResults, isLoading: searchLoading } = useSearchUsers(searchQuery);

  // Redirect non-admins
  useEffect(() => {
    if (!authLoading && !isAdmin) {
      toast.error('You do not have permission to access this page');
      navigate({ to: '/' });
    }
  }, [authLoading, isAdmin, navigate]);

  // Show loading while checking auth
  if (authLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg" />
      </div>
    );
  }

  // Don't render if not admin (will redirect)
  if (!isAdmin) {
    return null;
  }

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const handleSelectUser = (userId: string) => {
    setSelectedUserId(userId);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    // Keep selected user for context, clear after modal animation
    setTimeout(() => {
      if (!isModalOpen) {
        setSelectedUserId(null);
      }
    }, 300);
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">User Management</h1>
        <p className="text-base-content/60 mt-1">
          Search and manage user accounts, roles, and status
        </p>
      </div>

      {/* Search Form */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title text-lg mb-4">Search Users</h2>
          <UserSearchForm onSearch={handleSearch} isLoading={searchLoading} />
        </div>
      </div>

      {/* Search Results */}
      {searchQuery.length > 0 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">
              Search Results
              {searchResults?.items && (
                <span className="text-base-content/60 font-normal ml-2">
                  ({searchResults.items.length} found)
                </span>
              )}
            </h2>
          </div>

          {searchLoading ? (
            <div className="flex justify-center py-8">
              <span className="loading loading-spinner loading-lg" />
            </div>
          ) : searchResults?.items && searchResults.items.length > 0 ? (
            <div className="grid gap-3">
              {searchResults.items.map((user) => (
                <UserCard
                  key={user.id}
                  user={user}
                  onSelect={handleSelectUser}
                  isSelected={user.id === selectedUserId}
                />
              ))}
            </div>
          ) : (
            <div className="card bg-base-100 shadow-sm">
              <div className="card-body text-center">
                <p className="text-base-content/60">
                  No users found matching "{searchQuery}"
                </p>
              </div>
            </div>
          )}

          {/* Pagination hint */}
          {searchResults?.nextCursor && (
            <div className="text-center">
              <p className="text-sm text-base-content/60">
                More results available. Refine your search to see specific users.
              </p>
            </div>
          )}
        </div>
      )}

      {/* Empty state when no search */}
      {searchQuery.length === 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body text-center py-12">
            <svg
              className="w-16 h-16 mx-auto text-base-content/30 mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
            <h3 className="text-lg font-semibold mb-2">Search for Users</h3>
            <p className="text-base-content/60 max-w-md mx-auto">
              Enter an email address or display name to find users. You can then
              view their details, change their role, or manage their account status.
            </p>
          </div>
        </div>
      )}

      {/* User Detail Modal */}
      <UserDetailModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        userId={selectedUserId}
      />
    </div>
  );
}
