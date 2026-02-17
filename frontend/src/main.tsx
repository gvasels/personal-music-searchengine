/* eslint-disable react-refresh/only-export-components */
import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import {
  RouterProvider,
  createRouter,
  createRootRoute,
  createRoute,
  Outlet,
  redirect,
} from '@tanstack/react-router';
import { Toaster } from 'react-hot-toast';
import './index.css';

// Configure Amplify Auth
import { configureAuth } from './lib/auth';

// Import page components
import HomePage from './routes/index';
import LoginPage from './routes/login';
import SignupPage from './routes/signup';
import SearchPage from './routes/search';
import UploadPage from './routes/upload';
import TracksPage from './routes/tracks/index';
import TrackDetailPage from './routes/tracks/$trackId';
import AlbumsPage from './routes/albums/index';
import AlbumDetailPage from './routes/albums/$albumId';
import ArtistsPage from './routes/artists/index';
import ArtistDetailPage from './routes/artists/$artistName';
import PlaylistsPage from './routes/playlists/index';
import PlaylistDetailPage from './routes/playlists/$playlistId';
import TagsPage from './routes/tags/index';
import TagDetailPage from './routes/tags/$tagName';
import SettingsPage from './routes/settings';
import AdminUsersPage from './routes/admin/users';
import PermissionDeniedPage from './routes/permission-denied';
import MusicIndexPage from './routes/music/index';
import VideosPage from './routes/videos/index';
import GamingPage from './routes/gaming/index';
// Import layout components
import { Layout } from './components/layout';
import { AuthGuard } from './components/auth';
import MusicLayout from './components/layout/MusicLayout';

// Helper to wrap protected pages with AuthGuard
function withAuthGuard<P extends object>(Component: React.ComponentType<P>) {
  return function ProtectedComponent(props: P) {
    return (
      <AuthGuard>
        <Component {...props} />
      </AuthGuard>
    );
  };
}

// Configure auth
const cognitoConfig = {
  userPoolId: import.meta.env.VITE_COGNITO_USER_POOL_ID || '',
  userPoolClientId: import.meta.env.VITE_COGNITO_CLIENT_ID || '',
};

if (cognitoConfig.userPoolId && cognitoConfig.userPoolClientId) {
  configureAuth(cognitoConfig);
}

// Create query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
});

// Create root route with layout
const rootRoute = createRootRoute({
  component: () => (
    <Layout>
      <Outlet />
    </Layout>
  ),
});

// Public routes
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
});

const signupRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/signup',
  component: SignupPage,
});

const permissionDeniedRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/permission-denied',
  component: PermissionDeniedPage,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
});

// Protected routes
const searchRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/search',
  component: withAuthGuard(SearchPage),
});

const uploadRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/upload',
  component: withAuthGuard(UploadPage),
});

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/settings',
  component: withAuthGuard(SettingsPage),
});

const adminUsersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/users',
  component: withAuthGuard(AdminUsersPage),
});

// Music layout route — wraps all /music/* routes
const musicLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/music',
  component: withAuthGuard(MusicLayout),
});

// Music child routes
const musicIndexRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/',
  component: MusicIndexPage,
});

const musicTracksRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/tracks',
  component: TracksPage,
});

const musicTrackDetailRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/tracks/$trackId',
  component: TrackDetailPage,
});

const musicAlbumsRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/albums',
  component: AlbumsPage,
});

const musicAlbumDetailRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/albums/$albumId',
  component: AlbumDetailPage,
});

const musicArtistsRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/artists',
  component: ArtistsPage,
});

const musicArtistDetailRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/artists/$artistName',
  component: ArtistDetailPage,
});

const musicPlaylistsRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/playlists',
  component: PlaylistsPage,
});

const musicPlaylistDetailRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/playlists/$playlistId',
  component: PlaylistDetailPage,
});

const musicTagsRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/tags',
  component: TagsPage,
});

const musicTagDetailRoute = createRoute({
  getParentRoute: () => musicLayoutRoute,
  path: '/tags/$tagName',
  component: TagDetailPage,
});

// New section routes
const videosRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/videos',
  component: withAuthGuard(VideosPage),
});

const gamingRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/gaming',
  component: withAuthGuard(GamingPage),
});

// Legacy redirect routes
const legacyTracksRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tracks',
  beforeLoad: () => { throw redirect({ to: '/music/tracks' }); },
  component: () => null,
});

const legacyTrackDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tracks/$trackId',
  beforeLoad: ({ params }) => { throw redirect({ to: '/music/tracks/$trackId', params }); },
  component: () => null,
});

const legacyAlbumsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/albums',
  beforeLoad: () => { throw redirect({ to: '/music/albums' }); },
  component: () => null,
});

const legacyAlbumDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/albums/$albumId',
  beforeLoad: ({ params }) => { throw redirect({ to: '/music/albums/$albumId', params }); },
  component: () => null,
});

const legacyArtistsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/artists',
  beforeLoad: () => { throw redirect({ to: '/music/artists' }); },
  component: () => null,
});

const legacyArtistDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/artists/$artistName',
  beforeLoad: ({ params }) => { throw redirect({ to: '/music/artists/$artistName', params }); },
  component: () => null,
});

const legacyPlaylistsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/playlists',
  beforeLoad: () => { throw redirect({ to: '/music/playlists' }); },
  component: () => null,
});

const legacyPlaylistDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/playlists/$playlistId',
  beforeLoad: ({ params }) => { throw redirect({ to: '/music/playlists/$playlistId', params }); },
  component: () => null,
});

const legacyTagsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tags',
  beforeLoad: () => { throw redirect({ to: '/music/tags' }); },
  component: () => null,
});

const legacyTagDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tags/$tagName',
  beforeLoad: ({ params }) => { throw redirect({ to: '/music/tags/$tagName', params }); },
  component: () => null,
});

// Build route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  signupRoute,
  permissionDeniedRoute,
  searchRoute,
  uploadRoute,
  settingsRoute,
  adminUsersRoute,
  musicLayoutRoute.addChildren([
    musicIndexRoute,
    musicTracksRoute,
    musicTrackDetailRoute,
    musicAlbumsRoute,
    musicAlbumDetailRoute,
    musicArtistsRoute,
    musicArtistDetailRoute,
    musicPlaylistsRoute,
    musicPlaylistDetailRoute,
    musicTagsRoute,
    musicTagDetailRoute,
  ]),
  videosRoute,
  gamingRoute,
  // Legacy redirects
  legacyTracksRoute,
  legacyTrackDetailRoute,
  legacyAlbumsRoute,
  legacyAlbumDetailRoute,
  legacyArtistsRoute,
  legacyArtistDetailRoute,
  legacyPlaylistsRoute,
  legacyPlaylistDetailRoute,
  legacyTagsRoute,
  legacyTagDetailRoute,
]);

// Create router
const router = createRouter({ routeTree });

// Type registration for router
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

// App component
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster position="bottom-right" />
    </QueryClientProvider>
  );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
