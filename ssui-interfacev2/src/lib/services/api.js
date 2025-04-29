// api.js - A global API service for managing backend URLs and requests

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

// Helper to get the current backend URL
export function getCurrentBackendUrl() {
  const config = get(backendConfig);
  return config.backends[config.active]?.url || '';
}

// Helper to get the current authentication cookie
export function getCurrentAuthCookie() {
  const config = get(backendConfig);
  return config.backends[config.active]?.cookie || null;
}

// Add or update a backend
export function setBackend(id, url, cookie = null) {
  backendConfig.update(config => {
    config.backends[id] = { url, cookie };
    return config;
  });
}

// Set the active backend
export function setActiveBackend(id) {
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.active = id;
    }
    return config;
  });
}

// Update the cookie for a backend
export function updateCookie(id, cookie) {
  backendConfig.update(config => {
    if (config.backends[id]) {
      config.backends[id].cookie = cookie;
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
  // Get the current backend URL
  const backendUrl = getCurrentBackendUrl();
  
  // Ensure endpoint starts with "/" and handle backendUrl that might end with "/"
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const normalizedBackendUrl = backendUrl.endsWith('/') ? backendUrl.slice(0, -1) : backendUrl;
  
  // Construct the full URL
  const url = `${normalizedBackendUrl}${normalizedEndpoint}`;
  
  // Set up headers if not provided
  if (!options.headers) {
    options.headers = {};
  }
  
  // If we have an auth cookie, add it to the request
  const authCookie = getCurrentAuthCookie();
  if (authCookie) {
    // In a real implementation, you might want to handle this differently
    // depending on how your auth system works
    options.credentials = 'include';
    // You could also manually set the cookie if needed
    // options.headers['Cookie'] = authCookie;
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
  if (!options.headers) {
    options.headers = {};
  }
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
  // Get the current backend URL
  const backendUrl = getCurrentBackendUrl();
  
  // Ensure endpoint starts with "/" and handle backendUrl that might end with "/"
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const normalizedBackendUrl = backendUrl.endsWith('/') ? backendUrl.slice(0, -1) : backendUrl;
  
  // Construct the full URL
  const url = `${normalizedBackendUrl}${normalizedEndpoint}`;
  
  let eventSource = null;
  let isActive = true;
  let currentBackendId = get(backendConfig).active;
  
  // Create EventSource for SSE
  eventSource = new EventSource(url);
  
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
  try {
    // Try to load saved config from localStorage
    const savedConfig = localStorage.getItem('ssui-backend-config');
    if (savedConfig) {
      const parsed = JSON.parse(savedConfig);
      backendConfig.set(parsed);
    }
    
    // Subscribe to changes and save to localStorage
    backendConfig.subscribe(value => {
      localStorage.setItem('ssui-backend-config', JSON.stringify(value));
    });
  } catch (error) {
    console.error('Error initializing API service:', error);
  }
}