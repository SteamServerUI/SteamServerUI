<script>
  import { onMount, onDestroy } from 'svelte';
  import { authState, syncAuthState } from './lib/services/api';
  import Login from './Login.svelte';
  
  
  /**
   * @typedef {Object} Props
   * @property {boolean} [checkAuth] - Props - Set to false to disable auth checking
   * @property {string} [serverStatus] - 'online', 'offline', 'error'
   * @property {any} [serverError]
   * @property {import('svelte').Snippet} [children]
   */

  /** @type {Props} */
  let {
    checkAuth = true,
    serverStatus = 'online',
    serverError = null,
    children
  } = $props();
  
  // Local state
  let isChecking = $state(true);
  let unsubscribe;
  
  onMount(async () => {
    if (checkAuth && serverStatus === 'online') {
      // Subscribe to auth state changes
      unsubscribe = authState.subscribe(state => {
        isChecking = state.isAuthenticating;
      });
      
      // Sync auth state on mount
      try {
        await syncAuthState();
      } catch (err) {
        console.error('Auth check failed:', err);
        // Even on error, we're done checking
        isChecking = false;
      }
    } else {
      isChecking = false;
    }
  });
  
  onDestroy(() => {
    if (unsubscribe) unsubscribe();
  });
</script>

{#if isChecking}
  <div class="loading-container">
    <div class="loading-spinner"></div>
    <p>Checking authentication...</p>
  </div>
{:else if serverStatus === 'offline' || serverStatus === 'error'}
  <div class="server-error-container">
    <div class="status-icon {serverStatus}">
      {#if serverStatus === 'offline'}
        🔴
      {:else}
        ⚠️
      {/if}
    </div>
    <h2>Server Unavailable</h2>
    <p class="error-message">
      {#if serverStatus === 'offline'}
        Cannot connect to the server. Please check your connection or try again later.
      {:else}
        There was an error connecting to the server.
      {/if}
    </p>
    {#if serverError}
      <div class="error-details">
        {serverError}
      </div>
    {/if}
    <button class="retry-button" onclick={() => window.location.reload()}>
      Retry Connection
    </button>
  </div>
{:else if !$authState.isAuthenticated && checkAuth}
  <Login />
{:else}
  {@render children?.()}
{/if}

<style>
  .loading-container, .server-error-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f5f5;
    padding: 2rem;
    text-align: center;
  }
  
  .loading-spinner {
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-radius: 50%;
    border-top: 4px solid #4a90e2;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .status-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
  }
  
  .status-icon.offline {
    color: #f44336;
  }
  
  .status-icon.error {
    color: #ff9800;
  }
  
  h2 {
    margin-bottom: 1rem;
    color: #333;
  }

</style>