<script>
  import { onMount } from 'svelte';
  import { authState, backendConfig, login, getCurrentBackendUrl, setActiveBackend, setBackend } from './lib/services/api';
  import { get } from 'svelte/store';
  
  // Form data
  let username = '';
  let password = '';
  let errorMessage = '';
  let isSubmitting = false;
  let backends = [];
  let activeBackend = '';
  let backendStatuses = {};
  let showBackendSelector = false;
  let testingBackend = false;
  
  // New backend form data
  let showNewBackendForm = false;
  let newBackendId = '';
  let newBackendUrl = '';
  
  // Subscribe to auth state
  let unsubscribe;
  let unsubscribeBackend;
  
  onMount(() => {
    unsubscribe = authState.subscribe(state => {
      if (state.authError) {
        errorMessage = state.authError;
        
        // Show new backend form if endpoint not found
        if (state.authError === 'endpoint not found') {
          showNewBackendForm = true;
        }
      }
      isSubmitting = state.isAuthenticating;
    });
    
    unsubscribeBackend = backendConfig.subscribe(config => {
      backends = Object.keys(config.backends);
      activeBackend = config.active;
      
      // Initialize backend statuses
      backends.forEach(id => {
        if (!(id in backendStatuses)) {
          backendStatuses[id] = { status: 'unknown', lastChecked: null, error: null };
        }
      });
    });
    
    // Check active backend status
    checkBackendStatus(activeBackend);
    
    return () => {
      if (unsubscribe) unsubscribe();
      if (unsubscribeBackend) unsubscribeBackend();
    };
  });
  
// Test backend connection
async function checkBackendStatus(id) {
  testingBackend = true;
  backendStatuses[id] = { status: 'checking', lastChecked: new Date(), error: null };
  backendStatuses = { ...backendStatuses };
  
  try {
    const config = get(backendConfig);
    const backendUrl = config.backends[id].url === '/' ? '' : config.backends[id].url;
    
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000);
    
    const response = await fetch(`${backendUrl}/api/v2/server/status`, {
      method: 'GET',
      signal: controller.signal
    });
    
    clearTimeout(timeoutId);
    
    // Handle specific status codes
    if (response.ok) {
      backendStatuses[id] = {
        status: 'online',
        lastChecked: new Date(),
        error: null
      };
    } else if (response.status === 401) {
      // Treat 401 Unauthorized as "ok" to show login box
      backendStatuses[id] = {
        status: 'online',
        lastChecked: new Date(),
        authRequired: true,
        error: null
      };
    } else {
      backendStatuses[id] = {
        status: 'offline',
        lastChecked: new Date(),
        error: `Server responded with ${response.status} ${response.statusText}`
      };
    }
  } catch (error) {
    let status = 'error';
    let errorMsg = 'Failed to connect to the server';
    let certificateHint = false;
    const isHttps = backendUrl.startsWith('https://');

    if (error.name === 'AbortError') {
      errorMsg = 'Connection timed out. The server may be slow or unreachable.';
    } else if (isHttps) {
      status = 'cert-error';
      certificateHint = true;
      errorMsg = 'SSL Certificate Error. Please visit the server URL to accept the certificate.';
    } else {
      errorMsg = 'Server not found. The server may be down or the URL may be incorrect.';
    }

    backendStatuses[id] = {
      status,
      lastChecked: new Date(),
      error: errorMsg,
      certificateHint,
      backendUrl // Store the backend URL for the link
    };
  } finally {
    testingBackend = false;
    backendStatuses = { ...backendStatuses };
  }
}
  
  // Handle form submission
  async function handleSubmit() {
    errorMessage = '';
    isSubmitting = true;
    
    if (!username || !password) {
      errorMessage = 'Username and password are required';
      isSubmitting = false;
      return;
    }
    
    try {
      const success = await login(username, password);
      
      if (!success) {
        // Error message will be set via the authState subscription
      }
    } catch (error) {
      errorMessage = error.message || 'Login failed';
    } finally {
      isSubmitting = false;
    }
  }
  
  // Change active backend
  async function changeBackend() {
    await setActiveBackend(activeBackend);
    checkBackendStatus(activeBackend);
    showBackendSelector = false;
    // Reset error and new backend form when changing backend
    errorMessage = '';
    showNewBackendForm = false;
  }
  
  // Toggle backend selector
  function toggleBackendSelector() {
    showBackendSelector = !showBackendSelector;
  }
  
  // Add new backend from login page
  async function addNewBackend() {
    if (!newBackendId || !newBackendUrl) {
      return;
    }
    
    // Add the new backend
    setBackend(newBackendId, newBackendUrl);
    
    // Set it as active
    activeBackend = newBackendId;
    await setActiveBackend(newBackendId);
    
    // Check its status
    checkBackendStatus(newBackendId);
    
    // Reset form and hide it
    showNewBackendForm = false;
    newBackendId = '';
    newBackendUrl = '';
    errorMessage = '';
  }
  
  // Get current backend URL for display
  $: currentBackend = $backendConfig.backends[$backendConfig.active];
  $: backendUrl = currentBackend?.url || '/';
  $: displayUrl = backendUrl === '/' ? window.location.origin : backendUrl;
  $: activeStatus = backendStatuses[activeBackend]?.status || 'unknown';
  $: activeError = backendStatuses[activeBackend]?.error;
</script>

<div class="login-page">
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h2>Login Required</h2>
        
        <div class="server-selector">
          <div class="server-info" on:click={toggleBackendSelector}>
            <span class="server-label">Server:</span>
            <span class="server-url">{displayUrl}</span>
            
            <span class="status-indicator {activeStatus}">
              {#if activeStatus === 'online'}
                🟢
              {:else if activeStatus === 'offline'}
                🔴
              {:else if activeStatus === 'error' || activeStatus === 'unreachable'}
                ⚠️
              {:else if activeStatus === 'cert-error'}
                🔒
              {:else}
                ⚪
              {/if}
            </span>
            
            <button class="change-server-btn" title="Change server">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
            </button>
          </div>
          
          {#if showBackendSelector}
            <div class="backend-dropdown">
              <div class="dropdown-header">
                <h3>Select Backend</h3>
                <button class="close-btn" on:click={() => showBackendSelector = false}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>
              
              <div class="backend-list">
                {#each backends as backendId}
                  <div class="backend-item">
                    <label class="backend-option">
                      <input 
                        type="radio" 
                        name="backend" 
                        value={backendId} 
                        bind:group={activeBackend}
                      />
                      <div class="backend-details">
                        <span class="backend-name">{backendId}</span>
                        <span class="backend-url">{$backendConfig.backends[backendId].url}</span>
                        {#if backendStatuses[backendId]?.status === 'cert-error'}
                          <p class="backend-error">
                            {backendStatuses[backendId].error}
                            <a href={backendStatuses[backendId].backendUrl} target="_blank" rel="noopener">
                              Click here to accept the certificate
                            </a>
                          </p>
                        {:else if backendStatuses[backendId]?.error}
                          <p class="backend-error">{backendStatuses[backendId].error}</p>
                        {/if}
                      </div>
                      
                      <span class="status-indicator {backendStatuses[backendId]?.status || 'unknown'}">
                        {#if backendStatuses[backendId]?.status === 'online'}
                          🟢
                        {:else if backendStatuses[backendId]?.status === 'offline'}
                          🔴
                        {:else if backendStatuses[backendId]?.status === 'error' || backendStatuses[backendId]?.status === 'unreachable'}
                          ⚠️
                        {:else if backendStatuses[backendId]?.status === 'cert-error'}
                          🔒
                        {:else}
                          ⚪
                        {/if}
                      </span>
                    </label>
                  </div>
                {/each}
              </div>
              
              <div class="dropdown-actions">
                <button 
                  class="test-btn" 
                  on:click={() => checkBackendStatus(activeBackend)}
                  disabled={testingBackend}
                >
                  {testingBackend ? 'Testing...' : 'Test Connection'}
                </button>
                
                <button 
                  class="switch-btn"
                  on:click={changeBackend}
                  disabled={activeBackend === $backendConfig.active}
                >
                  Switch Server
                </button>
              </div>
            </div>
          {/if}
        </div>
      </div>
      {#if errorMessage || activeError}
      <div class="error-message">
        {#if errorMessage === 'endpoint not found'}
          The selected server doesn't have the expected login endpoint.
        {:else if activeStatus === 'cert-error'}
          {activeError}
          <a href={backendStatuses[activeBackend].backendUrl} target="_blank" rel="noopener">
            Click here to accept the certificate
          </a>
        {:else if activeError}
          {activeError}
        {:else}
          {errorMessage}
        {/if}
      </div>
      {/if}
      
      {#if showNewBackendForm}
        <div class="new-backend-form">
          <h3>Add New Server</h3>
          <p>The current server doesn't have the expected login endpoint or is unreachable. Add a new server with the correct authentication endpoint.</p>
          
          <div class="form-group">
            <label for="new-backend-id">Server Name</label>
            <input 
              type="text" 
              id="new-backend-id" 
              bind:value={newBackendId} 
              placeholder="e.g., production"
              disabled={isSubmitting}
            />
          </div>
          
          <div class="form-group">
            <label for="new-backend-url">Server URL</label>
            <input 
              type="text" 
              id="new-backend-url" 
              bind:value={newBackendUrl} 
              placeholder="e.g., https://api.example.com"
              disabled={isSubmitting}
            />
          </div>
          
          <div class="backend-form-actions">
            <button 
              class="cancel-btn" 
              on:click={() => showNewBackendForm = false}
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button 
              class="add-backend-btn" 
              on:click={addNewBackend}
              disabled={isSubmitting || !newBackendId || !newBackendUrl}
            >
              Add Server
            </button>
          </div>
        </div>
      {:else}
        <form on:submit|preventDefault={handleSubmit}>
          <div class="form-group">
            <label for="username">Username</label>
            <div class="input-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="input-icon">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              <input 
                type="text" 
                id="username" 
                bind:value={username} 
                disabled={isSubmitting} 
                placeholder="Enter username"
                autocomplete="username"
              />
            </div>
          </div>
          
          <div class="form-group">
            <label for="password">Password</label>
            <div class="input-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="input-icon">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              <input 
                type="password" 
                id="password" 
                bind:value={password} 
                disabled={isSubmitting} 
                placeholder="Enter password"
                autocomplete="current-password"
              />
            </div>
          </div>
          
          <button type="submit" class="login-button" disabled={isSubmitting || activeStatus === 'unreachable' || activeStatus === 'cert-error'}>
            {#if isSubmitting}
              <svg class="loading-spinner" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
              </svg>
              Logging in...
            {:else}
              Login
            {/if}
          </button>
        </form>
      {/if}
    </div>
  </div>
</div>

<style>
  .login-page {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f0f2f5;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }
  
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    max-width: 460px;
    padding: 1rem;
  }
  
  .login-card {
    background: white;
    padding: 2.5rem;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
    width: 100%;
    transition: transform 0.3s ease;
  }
  
  .login-card:hover {
    transform: translateY(-5px);
  }
  
  .login-header {
    margin-bottom: 1.5rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
    color: #333;
    font-weight: 600;
    font-size: 1.75rem;
  }
  
  h3 {
    color: #333;
    font-weight: 500;
    font-size: 1.2rem;
    margin-top: 0;
    margin-bottom: 1rem;
  }
  
  .server-selector {
    position: relative;
    margin-bottom: 1rem;
  }
  
  .server-info {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 0.75rem 1rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .server-info:hover {
    background-color: #eef1f6;
  }
  
  .server-label {
    font-weight: 500;
    color: #666;
    margin-right: 0.5rem;
  }
  
  .server-url {
    flex: 1;
    color: #333;
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .status-indicator {
    margin: 0 0.5rem;
  }
  
  .status-indicator.cert-error {
    color: #ff9800;
  }
  
  .change-server-btn {
    background: none;
    border: none;
    color: #666;
    cursor: pointer;
    padding: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
  }
  
  .change-server-btn:hover {
    color: #4a90e2;
  }
  
  .backend-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border-radius: 8px;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
    margin-top: 0.5rem;
    z-index: 10;
    overflow: hidden;
    animation: slideDown 0.2s ease-out;
  }
  
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .dropdown-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;
  }
  
  .dropdown-header h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 500;
  }
  
  .close-btn {
    background: none;
    border: none;
    color: #666;
    cursor: pointer;
    padding: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
  }
  
  .close-btn:hover {
    color: #4a90e2;
  }
  
  .backend-list {
    max-height: 250px;
    overflow-y: auto;
  }
  
  .backend-item {
    padding: 0.25rem 0.5rem;
  }
  
  .backend-option {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .backend-option:hover {
    background-color: #f5f7fa;
  }
  
  .backend-option input {
    margin-right: 0.75rem;
  }
  
  .backend-details {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: hidden;
  }
  
  .backend-name {
    font-weight: 500;
    color: #333;
  }
  
  .backend-url {
    font-size: 0.85rem;
    color: #666;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .backend-error {
    font-size: 0.8rem;
    color: #d32f2f;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .backend-error a {
  color: #4a90e2;
  text-decoration: underline;
  margin-left: 0.5rem;
  }

  .backend-error a:hover {
    color: #3a80d2;
  }
  
  .dropdown-actions {
    display: flex;
    justify-content: space-between;
    padding: 1rem;
    border-top: 1px solid #eee;
  }
  
  .test-btn {
    padding: 0.5rem 1rem;
    background-color: #f5f7fa;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    color: #333;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .test-btn:hover {
    background-color: #eef1f6;
    border-color: #ccc;
  }
  
  .switch-btn {
    padding: 0.5rem 1rem;
    background-color: #4a90e2;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    color: white;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .switch-btn:hover {
    background-color: #3a80d2;
  }
  
  .switch-btn:disabled {
    background-color: #a0c1e2;
    cursor: not-allowed;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #333;
    font-weight: 500;
    font-size: 0.95rem;
  }
  
  .input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }
  
  .input-icon {
    position: absolute;
    left: 1rem;
    color: #999;
  }
  
  input[type="text"],
  input[type="password"] {
    width: 100%;
    padding: 0.85rem 1rem 0.85rem 2.75rem;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
    transition: all 0.2s;
  }
  
  .new-backend-form input[type="text"] {
    padding: 0.85rem 1rem;
  }
  
  input[type="text"]:focus,
  input[type="password"]:focus {
    outline: none;
    border-color: #4a90e2;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.15);
  }
  
  .checkbox {
    display: flex;
    align-items: center;
  }
  
  .checkbox label {
    display: flex;
    align-items: center;
    margin-bottom: 0;
    cursor: pointer;
  }
  
  .checkbox span {
    margin-left: 0.5rem;
    font-weight: normal;
  }
  
  .login-button {
    width: 100%;
    padding: 0.85rem;
    background-color: #4a90e2;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .login-button:hover {
    background-color: #3a80d2;
  }
  
  .login-button:disabled {
    background-color: #a0c1e2;
    cursor: not-allowed;
  }
  
  .loading-spinner {
    animation: spin 1s linear infinite;
    margin-right: 0.5rem;
  }
  
  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
  
  .error-message {
    padding: 0.85rem;
    background-color: #ffebee;
    color: #d32f2f;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    font-size: 0.95rem;
    display: flex;
    align-items: center;
    animation: fadeIn 0.3s ease-out;
  }
  
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  
  /* Status colors */
  .status-indicator.online {
    color: #4caf50;
  }
  
  .status-indicator.offline {
    color: #f44336;
  }
  
  .status-indicator.error, .status-indicator.unreachable {
    color: #ff9800;
  }
  
  .status-indicator.unknown {
    color: #9e9e9e;
  }
  
  .new-backend-form {
    background-color: #f8fafc;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
    border: 1px solid #e0e7ff;
    animation: fadeIn 0.3s ease-out;
  }
  
  .new-backend-form p {
    color: #555;
    font-size: 0.95rem;
    margin-bottom: 1.25rem;
  }
  
  .backend-form-actions {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
  }
  
  .cancel-btn {
    flex: 1;
    padding: 0.75rem;
    background-color: #f5f7fa;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 0.95rem;
    font-weight: 500;
    color: #555;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .cancel-btn:hover {
    background-color: #eef1f6;
    border-color: #ccc;
  }
  
  .add-backend-btn {
    flex: 2;
    padding: 0.75rem;
    background-color: #4a90e2;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .add-backend-btn:hover {
    background-color: #3a80d2;
  }
  
  .add-backend-btn:disabled {
    background-color: #a0c1e2;
    cursor: not-allowed;
  }
  
  @media (max-width: 480px) {
    .login-container {
      padding: 0.5rem;
    }
    
    .login-card {
      padding: 1.5rem;
    }
    
    h2 {
      font-size: 1.5rem;
    }
    
    .dropdown-actions {
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .test-btn, .switch-btn {
      width: 100%;
    }
    
    .backend-form-actions {
      flex-direction: column;
    }
    
    .cancel-btn, .add-backend-btn {
      width: 100%;
    }
  }
</style>