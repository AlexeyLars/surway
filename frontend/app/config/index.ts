// Configuration loader for frontend
// Reads from environment variables which can be set from configs/frontend.conf

interface FrontendConfig {
  apiBaseUrl: string;
  frontendBaseUrl: string;
}

function loadConfig(): FrontendConfig {
  // API Configuration
  const apiProtocol = process.env.NEXT_PUBLIC_API_PROTOCOL || 'http';
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || 'localhost';
  const apiPort = process.env.NEXT_PUBLIC_API_PORT || '8080';
  const apiVersion = process.env.NEXT_PUBLIC_API_VERSION || 'v1';
  
  const apiBaseUrl = `${apiProtocol}://${apiHost}:${apiPort}/api/${apiVersion}`;

  // Frontend Configuration
  const frontendProtocol = process.env.NEXT_PUBLIC_FRONTEND_PROTOCOL || 'http';
  const frontendHost = process.env.NEXT_PUBLIC_FRONTEND_HOST || 'localhost';
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
