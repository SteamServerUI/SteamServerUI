<script>
    import { onMount } from 'svelte';
    import { authState, backendConfig, login, getCurrentBackendUrl } from './lib/services/api';
    
    // Form data
    let username = '';
    let password = '';
    let rememberMe = false;
    let errorMessage = '';
    let isSubmitting = false;
    
    // Subscribe to auth state
    let unsubscribe;
    
    onMount(() => {
      unsubscribe = authState.subscribe(state => {
        if (state.authError) {
          errorMessage = state.authError;
        }
        isSubmitting = state.isAuthenticating;
      });
      
      return () => {
        if (unsubscribe) unsubscribe();
      };
    });
    
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
    
    // Get current backend URL for display
    $: currentBackend = $backendConfig.backends[$backendConfig.active];
    $: backendUrl = currentBackend?.url || '/';
    $: displayUrl = backendUrl === '/' ? window.location.origin : backendUrl;
  </script>
  
  <div class="login-container">
    <div class="login-card">
      <h2>Login Required</h2>
      <p class="server-url">Server: {displayUrl}</p>
      
      {#if errorMessage}
        <div class="error-message">
          {errorMessage}
        </div>
      {/if}
      
      <form on:submit|preventDefault={handleSubmit}>
        <div class="form-group">
          <label for="username">Username</label>
          <input 
            type="text" 
            id="username" 
            bind:value={username} 
            disabled={isSubmitting} 
            placeholder="Enter username"
            autocomplete="username"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            type="password" 
            id="password" 
            bind:value={password} 
            disabled={isSubmitting} 
            placeholder="Enter password"
            autocomplete="current-password"
          />
        </div>
        
        <div class="form-group checkbox">
          <label>
            <input type="checkbox" bind:checked={rememberMe} disabled={isSubmitting} />
            <span>Remember me</span>
          </label>
        </div>
        
        <button type="submit" class="login-button" disabled={isSubmitting}>
          {isSubmitting ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  </div>
  
  <style>
    .login-container {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
      background-color: #f5f5f5;
    }
    
    .login-card {
      background: white;
      padding: 2rem;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      border-radius: 8px;
      width: 100%;
      max-width: 400px;
    }
    
    h2 {
      text-align: center;
      margin-bottom: 1.5rem;
      color: #333;
    }
    
    .server-url {
      text-align: center;
      margin-bottom: 1.5rem;
      color: #666;
      font-size: 0.9rem;
    }
    
    .form-group {
      margin-bottom: 1rem;
    }
    
    label {
      display: block;
      margin-bottom: 0.5rem;
      color: #333;
      font-weight: 500;
    }
    
    .checkbox {
      display: flex;
      align-items: center;
    }
    
    .checkbox label {
      display: flex;
      align-items: center;
      margin-bottom: 0;
    }
    
    .checkbox span {
      margin-left: 0.5rem;
    }
    
    input[type="text"],
    input[type="password"] {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 1rem;
    }
    
    input[type="text"]:focus,
    input[type="password"]:focus {
      outline: none;
      border-color: #4a90e2;
      box-shadow: 0 0 0 2px rgba(74, 144, 226, 0.2);
    }
    
    .login-button {
      width: 100%;
      padding: 0.75rem;
      background-color: #4a90e2;
      color: white;
      border: none;
      border-radius: 4px;
      font-size: 1rem;
      cursor: pointer;
      transition: background-color 0.2s;
    }
    
    .login-button:hover {
      background-color: #3a80d2;
    }
    
    .login-button:disabled {
      background-color: #a0c1e2;
      cursor: not-allowed;
    }
    
    .error-message {
      padding: 0.75rem;
      background-color: #ffebee;
      color: #d32f2f;
      border-radius: 4px;
      margin-bottom: 1rem;
      font-size: 0.9rem;
    }
  </style>