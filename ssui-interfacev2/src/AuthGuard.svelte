<script>
  import { onMount, onDestroy } from 'svelte';
  import { authState, syncAuthState } from './lib/services/api';
  import Login from './Login.svelte';
  
  // Props
  export let checkAuth = true; // Set to false to disable auth checking
  
  // Local state
  let isChecking = true;
  let unsubscribe;
  
  onMount(async () => {
    if (checkAuth) {
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
{:else if !$authState.isAuthenticated && checkAuth}
  <Login />
{:else}
  <slot></slot>
{/if}

<style>
  .loading-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f5f5;
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
</style>