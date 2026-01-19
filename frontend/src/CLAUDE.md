# Frontend Source - CLAUDE.md

## Overview

Source code for the React frontend application. Contains components, utilities, routes, and hooks.

## Directory Structure

```
src/
├── components/     # Reusable React components
├── hooks/          # Custom React hooks
├── lib/            # Utilities (API, auth, state)
├── pages/          # Page components
├── routes/         # TanStack Router file-based routes
├── main.tsx        # Application entry point
├── index.css       # Global styles
└── routeTree.gen.ts # Generated route tree (auto-generated)
```

## Entry Point (`main.tsx`)

```typescript
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { routeTree } from './routeTree.gen'
import { configureAuth } from './lib/auth'

// Configure Cognito auth
configureAuth()

const queryClient = new QueryClient()
const router = createRouter({ routeTree })

ReactDOM.createRoot(document.getElementById('root')!).render(
  <QueryClientProvider client={queryClient}>
    <RouterProvider router={router} />
  </QueryClientProvider>
)
```

## Subdirectories

### components/
Reusable UI components. See `components/CLAUDE.md` for details.

### hooks/
Custom React hooks for shared logic.

### lib/
Utilities and configurations:
- `api.ts` - API client and types
- `auth.ts` - Cognito authentication
- `store.ts` - Zustand state stores

### routes/
TanStack Router file-based routes:
- `__root.tsx` - Root layout
- `index.tsx` - Home page

## Testing

All components should have corresponding test files:
- `Component.tsx` → `Component.test.tsx`

Run tests: `npm run test`
