<script>
  import { onMount } from 'svelte';
  import { initializeApiService, syncAuthState } from './lib/services/api';
  
  let isInitialized = false;
  
  onMount(() => {
    // Initialize the API service
    initializeApiService();
    // Start auth check immediately
    syncAuthState().catch(console.error);
    isInitialized = true;
  });
</script>

{#if isInitialized}
  <slot />
{:else}
  <div class="initializing">
    <div class="loading-spinner"></div>
    <p>Initializing application...</p>
  </div>
{/if}

<style>
  .initializing {
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