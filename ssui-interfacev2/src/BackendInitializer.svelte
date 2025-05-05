<script>
  import { onMount } from 'svelte';
  import { initializeApiService } from './lib/services/api';
  import InitializingView from './lib/components/resuables/InitializingView.svelte';
  
  /**
   * @typedef {Object} Props
   * @property {any} [onStatusChange] - Default to a no-op function if no handler is provided
   * @property {import('svelte').Snippet} [children]
   */

  /** @type {Props} */
  let { onStatusChange = (statusObj) => {}, children } = $props();
  let isInitialized = $state(false);
  let serverStatus = $state('initializing'); // 'checking', 'online', 'offline', 'error', 'cert-error', 'unreachable'
  
  // Add AbortSignal polyfill for older browsers
  if (!AbortSignal.timeout) {
    AbortSignal.timeout = function timeout(ms) {
      const controller = new AbortController();
      setTimeout(() => controller.abort(), ms);
      return controller.signal;
    };
  }
  
  onMount(() => {
    // Initialize the API service
    isInitialized = initializeApiService();
  });
</script>

{#if isInitialized}
  {@render children?.()}
{:else}
serverStatus
  <InitializingView {serverStatus} />
{/if}