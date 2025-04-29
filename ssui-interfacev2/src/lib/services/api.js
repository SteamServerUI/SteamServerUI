// api.js

import { writable, get } from 'svelte/store';

// Store for backend configuration
export const backendConfig = writable({
  active: 'default', // Currently active backend
  backends: {
    default: {
      url: 'https://localhost:8443', // Default backend URL
      cookie: null // Authentication cookie
    }
  }
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
  return getCurrentBackend().url;
}

// Helper to get the current authentication cookie
export function getCurrentAuthCookie() {
  return getCurrentBackend().cookie;
}

// Add or update a backend
export function setBackend(id, url, cookie = null) {
  backendConfig.update(config => {
    config.backends[id] = { url, cookie };
    return config;
  });
}

// Set the active backend and ensure cookie is properly set
export function setActiveBackend(id) {
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.active = id;
    }
    return config;
  });
}

// Update the cookie for a backend and persist it
export function updateCookie(id, cookie) {
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.backends[id].cookie = cookie;
      
      // If this is the active backend, ensure the cookie is immediately available
      if (config.active === id) {
        config.backends[config.active].cookie = cookie;
      }
    }
    return config;
  });
}

/**
 * Fetch wrapper that automatically adds the backend URL and handles authentication
 * @param {string} endpoint - The API endpoint (e.g., "/api/v2/whatever")
 * @param {Object} options - Fetch options
 * @returns {Promise} - The fetch promise
 */
export async function apiFetch(endpoint, options = {}) {
  // Get the current backend configuration
  const backend = getCurrentBackend();
  
  // Ensure endpoint starts with "/" and handle backendUrl that might end with "/"
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const normalizedBackendUrl = backend.url.endsWith('/') ? backend.url.slice(0, -1) : backend.url;
  
  // Construct the full URL
  const url = `${normalizedBackendUrl}${normalizedEndpoint}`;
  
  // Set up headers if not provided
  options.headers = options.headers || {};
  
  // If we have an auth cookie, add it to the request
  if (backend.cookie) {
    // For non-SSE requests, we'll manually set the Cookie header
    options.headers['Cookie'] = backend.cookie;
  }
  
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
  // Get the current backend configuration
  const backend = getCurrentBackend();
  
  // Ensure endpoint starts with "/" and handle backendUrl that might end with "/"
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const normalizedBackendUrl = backend.url.endsWith('/') ? backend.url.slice(0, -1) : backend.url;
  
  // Construct the full URL
  const url = `${normalizedBackendUrl}${normalizedEndpoint}`;
  
  let eventSource = null;
  let isActive = true;
  let currentBackendId = get(backendConfig).active;
  
  // Create EventSource for SSE with credentials if cookie exists
  const eventSourceInit = backend.cookie ? { withCredentials: true } : {};
  eventSource = new EventSource(url, eventSourceInit);
  
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
            url: 'https://localhost:8443',
            cookie: parsed.backends?.default?.cookie || null
          }
        }
      };

      // Merge all backends from storage
      if (parsed.backends) {
        for (const [id, backend] of Object.entries(parsed.backends)) {
          if (id !== 'default') {
            validatedConfig.backends[id] = {
              url: backend.url,
              cookie: backend.cookie || null
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
    
    // Return cleanup function (though this is a service, so it won't typically be cleaned up)
    return unsubscribe;
  } catch (error) {
    console.error('Error initializing API service:', error);
    isInitialized = true; // Prevent repeated failed initializations
  }
}

export async function syncAuthState() {
  const backend = getCurrentBackend();
  if (!backend.cookie) return;

  try {
    // Make a simple request to verify the cookie
    const response = await apiFetch('/api/v2/auth/check', {
      method: 'GET',
      headers: {
        'Accept': 'application/json'
      }
    });

    if (!response.ok) {
      // Cookie is invalid, clear it
      updateCookie(get(backendConfig).active, null);
    }
  } catch (error) {
    console.warn('Auth check failed:', error);
  }
}

// Automatically sync auth state when the module loads
syncAuthState().catch(console.error);