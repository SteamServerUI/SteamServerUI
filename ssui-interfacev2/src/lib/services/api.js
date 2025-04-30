// api.js

import { writable, get } from 'svelte/store';

// Store for backend configuration
export const backendConfig = writable({
  active: 'default', // Currently active backend
  backends: {
    default: {
      url: '/', // Default backend URL is the current host
      token: null // Authentication token (JWT)
    }
  }
});

// Store for authentication state
export const authState = writable({
  isAuthenticated: false,
  isAuthenticating: false,
  authError: null
});

// Track initialization state
let isInitialized = false;

// Helper to get the current backend configuration
export function getCurrentBackend() {
  const config = get(backendConfig);
  return config.backends[config.active] || config.backends.default;
}

// Helper to get the current backend URL
export function getCurrentBackendUrl() {
  const backend = getCurrentBackend();
  return backend.url === '/' ? '' : backend.url;
}

// Helper to get the current authentication token
export function getCurrentAuthToken() {
  return getCurrentBackend().token;
}

// Add or update a backend
export function setBackend(id, url) {
  backendConfig.update(config => {
    // If the backend already exists, preserve its token
    const existingToken = config.backends[id]?.token || null;
    config.backends[id] = { url, token: existingToken };
    return config;
  });
}

// Set the active backend and verify authentication
export async function setActiveBackend(id) {
  let success = false;
  
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.active = id;
    }
    return config;
  });
  
  // After changing backend, check authentication status
  try {
    await syncAuthState();
    success = true;
  } catch (error) {
    console.error('Error syncing auth state after backend change:', error);
  }
  
  return success;
}

// Update the token for a backend and persist it
export function updateAuthToken(id, token) {
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.backends[id].token = token;
    }
    return config;
  });
  
  // Update the auth state
  if (id === get(backendConfig).active) {
    authState.update(state => ({
      ...state,
      isAuthenticated: !!token,
      authError: null
    }));
  }
}

// Clear authentication for the current backend
export function clearAuthentication() {
  const currentBackendId = get(backendConfig).active;
  updateAuthToken(currentBackendId, null);
}

/**
 * Fetch wrapper that automatically adds the backend URL and handles authentication
 * @param {string} endpoint - The API endpoint (e.g., "/api/v2/whatever")
 * @param {Object} options - Fetch options
 * @returns {Promise} - The fetch promise
 */
export async function apiFetch(endpoint, options = {}) {
  // Get the current backend configuration
  const backendUrl = getCurrentBackendUrl();
  
  // Ensure endpoint starts with "/" if it's not an empty string
  const normalizedEndpoint = endpoint.startsWith('/') || endpoint === '' ? endpoint : `/${endpoint}`;
  
  // Construct the full URL
  const url = `${backendUrl}${normalizedEndpoint}`;
  
  // Set up headers if not provided
  options.headers = options.headers || {};
  
  // Always include credentials for CORS requests
  options.credentials = 'include';
  
  // No need to add the Authorization header as we're using cookies now
  
  // Perform the fetch
  return fetch(url, options);
}

/**
 * Helper for JSON API calls
 * @param {string} endpoint - The API endpoint
 * @param {Object} options - Fetch options
 * @returns {Promise<Object>} - Parsed JSON response
 */
export async function apiJson(endpoint, options = {}) {
  // Default to JSON content type if not specified
  options.headers = options.headers || {};
  if (!options.headers['Content-Type'] && options.method && options.method !== 'GET') {
    options.headers['Content-Type'] = 'application/json';
  }
  
  const response = await apiFetch(endpoint, options);
  
  // Check for auth errors
  if (response.status === 401) {
    // Update auth state
    authState.update(state => ({
      ...state,
      isAuthenticated: false,
      authError: 'Unauthorized'
    }));
    throw new Error('Authentication required');
  }
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  
  return response.json();
}

/**
 * Helper for text API calls
 * @param {string} endpoint - The API endpoint
 * @param {Object} options - Fetch options
 * @returns {Promise<string>} - Text response
 */
export async function apiText(endpoint, options = {}) {
  const response = await apiFetch(endpoint, options);
  
  // Check for auth errors
  if (response.status === 401) {
    // Update auth state
    authState.update(state => ({
      ...state,
      isAuthenticated: false,
      authError: 'Unauthorized'
    }));
    throw new Error('Authentication required');
  }
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  
  return response.text();
}

/**
 * Helper for Server-Sent Events (SSE)
 * @param {string} endpoint - The SSE endpoint
 * @param {function} onMessage - Callback for each message
 * @param {function} onError - Error callback
 * @returns {Object} - Control object with a close() method
 */
export function apiSSE(endpoint, onMessage, onError = console.error) {
  return
  // Get the current backend URL
  const backendUrl = getCurrentBackendUrl();
  
  // Ensure endpoint starts with "/" if it's not an empty string
  const normalizedEndpoint = endpoint.startsWith('/') || endpoint === '' ? endpoint : `/${endpoint}`;
  
  // Construct the full URL - no need to add token as query param as we use cookies
  const baseUrl = backendUrl || window.location.origin;
  const url = new URL(`${baseUrl}${normalizedEndpoint}`);
  
  // No need to add token as query parameter - cookies will be sent automatically
  
  let eventSource = null;
  let isActive = true;
  let currentBackendId = get(backendConfig).active;
  
  // Create EventSource for SSE with withCredentials to send cookies
  const eventSourceOptions = { withCredentials: true };
  eventSource = new EventSource(url.toString(), eventSourceOptions);
  
  // Set up event handlers
  eventSource.onmessage = event => {
    try {
      // Try to parse as JSON first
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (e) {
      // If not JSON, pass the raw string
      onMessage(event.data);
    }
  };
  
  eventSource.onerror = error => {
    // Check if the error might be an authentication issue
    if (eventSource.readyState === EventSource.CLOSED) {
      // Update auth state if we suspect auth issues
      syncAuthState().catch(console.error);
    }
    onError(error);
  };
  
  // Subscribe to backend config changes to close this connection when backend changes
  const unsubscribe = backendConfig.subscribe(config => {
    if (isActive && config.active !== currentBackendId) {
      // Backend has changed, clean up this connection
      if (eventSource) {
        eventSource.close();
        eventSource = null;
      }
      isActive = false;
      unsubscribe();
    }
  });
  
  // Return control object with enhanced close method
  return {
    close: () => {
      if (eventSource) {
        eventSource.close();
        eventSource = null;
      }
      isActive = false;
      unsubscribe();
    }
  };
}

/**
 * Login to the current backend
 * @param {string} username - Username
 * @param {string} password - Password
 * @returns {Promise<boolean>} - Success status
 */
export async function login(username, password) {
  // Set authenticating state
  authState.update(state => ({
    ...state,
    isAuthenticating: true,
    authError: null
  }));
  
  try {
    const response = await apiFetch('/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include', // Ensure cookies are stored
      body: JSON.stringify({ username, password })
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Authentication failed');
    }
    
    const data = await response.json();
    
    // Save the token for future reference - the actual auth will use cookies
    updateAuthToken(get(backendConfig).active, data.token);
    
    // Update auth state
    authState.update(state => ({
      ...state,
      isAuthenticated: true,
      isAuthenticating: false,
      authError: null
    }));
    
    return true;
  } catch (error) {
    // Update auth state with error
    authState.update(state => ({
      ...state,
      isAuthenticated: false,
      isAuthenticating: false,
      authError: error.message
    }));
    
    return false;
  }
}

// Check if the current backend requires authentication
export async function syncAuthState() {
  const currentBackendId = get(backendConfig).active;
  const backend = getCurrentBackend();
  
  // Update auth state to checking
  authState.update(state => ({
    ...state,
    isAuthenticating: true
  }));
  
  try {
    // Make a simple request to verify authentication
    const response = await apiFetch('/api/v2/auth/check', {
      method: 'GET',
      headers: {
        'Accept': 'application/json'
      },
      credentials: 'include' // Ensure cookies are sent
    });
    
    if (response.status === 401) {
      // Authentication required but we're not authenticated
      authState.update(state => ({
        ...state,
        isAuthenticated: false,
        isAuthenticating: false,
        authError: 'Authentication required'
      }));
      return false;
    } else if (!response.ok) {
      // Some other error
      authState.update(state => ({
        ...state,
        isAuthenticating: false,
        authError: `API error: ${response.status} ${response.statusText}`
      }));
      return false;
    }
    
    // Successfully authenticated
    authState.update(state => ({
      ...state,
      isAuthenticated: true,
      isAuthenticating: false,
      authError: null
    }));
    return true;
  } catch (error) {
    // Connection error or other issue
    console.warn('Auth check failed:', error);
    authState.update(state => ({
      ...state,
      isAuthenticating: false,
      authError: 'Connection error'
    }));
    return false;
  }
}

// Initial setup function to load saved backend configurations
export function initializeApiService() {
  if (isInitialized) return;
  
  try {
    // Try to load saved config from localStorage
    const savedConfig = localStorage.getItem('ssui-backend-config');
    if (savedConfig) {
      const parsed = JSON.parse(savedConfig);
      
      // Validate and normalize the loaded configuration
      const validatedConfig = {
        active: parsed.active && parsed.backends?.[parsed.active] ? parsed.active : 'default',
        backends: {
          default: {
            url: '/',
            token: parsed.backends?.default?.token || null
          }
        }
      };

      // Merge all backends from storage
      if (parsed.backends) {
        for (const [id, backend] of Object.entries(parsed.backends)) {
          if (id !== 'default') {
            validatedConfig.backends[id] = {
              url: backend.url,
              token: backend.token || null
            };
          }
        }
      }

      // Apply the validated config
      backendConfig.set(validatedConfig);
    }

    // Subscribe to changes and save to localStorage
    const unsubscribe = backendConfig.subscribe(value => {
      try {
        localStorage.setItem('ssui-backend-config', JSON.stringify(value));
      } catch (error) {
        console.error('Error saving backend config:', error);
      }
    });

    isInitialized = true;
    
    // Check authentication status
    syncAuthState().catch(console.error);
    
    // Return cleanup function (though this is a service, so it won't typically be cleaned up)
    return unsubscribe;
  } catch (error) {
    console.error('Error initializing API service:', error);
    isInitialized = true; // Prevent repeated failed initializations
  }
}