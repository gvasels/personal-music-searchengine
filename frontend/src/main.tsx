import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import {
  RouterProvider,
  createRouter,
  createRootRoute,
  createRoute,
  Outlet,
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

// Import layout components
import { Layout } from './components/layout';
import { AuthGuard } from './components/auth';

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

// Create routes - Public routes (no auth required)
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

// Create routes - Protected routes (auth required)
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: withAuthGuard(HomePage),
});

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

const tracksRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tracks',
  component: withAuthGuard(TracksPage),
});

const trackDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tracks/$trackId',
  component: withAuthGuard(TrackDetailPage),
});

const albumsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/albums',
  component: withAuthGuard(AlbumsPage),
});

const albumDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/albums/$albumId',
  component: withAuthGuard(AlbumDetailPage),
});

const artistsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/artists',
  component: withAuthGuard(ArtistsPage),
});

const artistDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/artists/$artistName',
  component: withAuthGuard(ArtistDetailPage),
});

const playlistsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/playlists',
  component: withAuthGuard(PlaylistsPage),
});

const playlistDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/playlists/$playlistId',
  component: withAuthGuard(PlaylistDetailPage),
});

const tagsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tags',
  component: withAuthGuard(TagsPage),
});

const tagDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/tags/$tagName',
  component: withAuthGuard(TagDetailPage),
});

// Build route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  signupRoute,
  searchRoute,
  uploadRoute,
  tracksRoute,
  trackDetailRoute,
  albumsRoute,
  albumDetailRoute,
  artistsRoute,
  artistDetailRoute,
  playlistsRoute,
  playlistDetailRoute,
  tagsRoute,
  tagDetailRoute,
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
