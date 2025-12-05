// Configuration loader for frontend
// Reads from environment variables which can be set from configs/frontend.conf

interface FrontendConfig {
  apiBaseUrl: string;
  frontendBaseUrl: string;
}

function loadConfig(): FrontendConfig {
  const isServer = typeof window === 'undefined';

  // API Configuration (server inside Docker talks to backend service; client talks via exposed port)
  const apiProtocol = process.env.NEXT_PUBLIC_API_PROTOCOL || 'http';
  const apiPort = process.env.NEXT_PUBLIC_API_PORT || '8080';
  const apiVersion = process.env.NEXT_PUBLIC_API_VERSION || 'v1';

  const apiInternalProtocol = process.env.API_INTERNAL_PROTOCOL || apiProtocol;
  const apiInternalPort = process.env.API_INTERNAL_PORT || apiPort;

  const apiHost = isServer
    ? process.env.API_INTERNAL_HOST || process.env.NEXT_PUBLIC_API_HOST || 'localhost'
    : process.env.NEXT_PUBLIC_API_HOST || (typeof window !== 'undefined' ? window.location.hostname : 'localhost');

  const selectedApiProtocol = isServer ? apiInternalProtocol : apiProtocol;
  const selectedApiPort = isServer ? apiInternalPort : apiPort;

  const apiBaseUrl = `${selectedApiProtocol}://${apiHost}:${selectedApiPort}/api/${apiVersion}`;

  // Frontend Configuration
  const frontendProtocol = process.env.NEXT_PUBLIC_FRONTEND_PROTOCOL || 'http';
  const frontendHost = process.env.NEXT_PUBLIC_FRONTEND_HOST || (typeof window !== 'undefined' ? window.location.hostname : 'localhost');
  const frontendPort = process.env.NEXT_PUBLIC_FRONTEND_PORT || '3000';
  
  // For Vercel deployments, use VERCEL_URL if available
  const frontendBaseUrl = process.env.VERCEL_URL 
    ? `https://${process.env.VERCEL_URL}`
    : `${frontendProtocol}://${frontendHost}:${frontendPort}`;

  return {
    apiBaseUrl,
    frontendBaseUrl,
  };
}

export const config = loadConfig();
