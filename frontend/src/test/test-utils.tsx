/* eslint-disable react-refresh/only-export-components */
import { ReactElement, ReactNode } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
}

function AllTheProviders({ children }: { children: ReactNode }) {
  const queryClient = createTestQueryClient();
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}

function customRender(ui: ReactElement, options?: Omit<RenderOptions, 'wrapper'>) {
  return {
    user: userEvent.setup(),
    ...render(ui, { wrapper: AllTheProviders, ...options }),
  };
}

export * from '@testing-library/react';
export { customRender as render, userEvent };

export const createMockAudioFile = (name = 'test.mp3', size = 1024 * 1024) => {
  const file = new File(['audio'], name, { type: 'audio/mpeg' });
  Object.defineProperty(file, 'size', { value: size });
  return file;
};
