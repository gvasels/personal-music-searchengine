/**
 * Environment Configuration
 * Supports both production AWS and LocalStack development environments
 */

export interface AppConfig {
  apiUrl: string;
  cognito: {
    userPoolId: string;
    clientId: string;
    region: string;
    endpoint?: string;
  };
  isLocalStack: boolean;
  isDevelopment: boolean;
}

/**
 * Detect if running in LocalStack mode.
 */
export function isLocalStackMode(): boolean {
  return import.meta.env.VITE_LOCAL_STACK === 'true';
}

/**
 * Get application configuration based on environment.
 */
export function getConfig(): AppConfig {
  const isLocalStack = isLocalStackMode();
  const isDevelopment = import.meta.env.DEV;

  return {
    apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
    cognito: {
      userPoolId: import.meta.env.VITE_COGNITO_USER_POOL_ID || '',
      clientId: import.meta.env.VITE_COGNITO_CLIENT_ID || '',
      region: import.meta.env.VITE_COGNITO_REGION || 'us-east-1',
      endpoint: isLocalStack ? import.meta.env.VITE_COGNITO_ENDPOINT : undefined,
    },
    isLocalStack,
    isDevelopment,
  };
}

/**
 * Singleton config instance.
 */
export const config = getConfig();

/**
 * Test user credentials for LocalStack development.
 * Only available when VITE_LOCAL_STACK=true
 */
export const localTestUsers = {
  admin: {
    email: 'admin@local.test',
    password: 'LocalTest123!',
  },
  subscriber: {
    email: 'subscriber@local.test',
    password: 'LocalTest123!',
  },
  artist: {
    email: 'artist@local.test',
    password: 'LocalTest123!',
  },
};

/**
 * Log configuration on startup (non-sensitive values only).
 */
export function logConfig(): void {
  if (config.isDevelopment) {
    console.log('[Config] Environment:', {
      isLocalStack: config.isLocalStack,
      apiUrl: config.apiUrl,
      cognitoRegion: config.cognito.region,
      cognitoEndpoint: config.cognito.endpoint || 'AWS (default)',
      hasUserPoolId: !!config.cognito.userPoolId,
      hasClientId: !!config.cognito.clientId,
    });

    if (config.isLocalStack) {
      console.log('[Config] LocalStack mode enabled. Test users available:');
      console.log('  - admin@local.test (admin)');
      console.log('  - subscriber@local.test (subscriber)');
      console.log('  - artist@local.test (artist)');
      console.log('  Password: LocalTest123!');
    }
  }
}
