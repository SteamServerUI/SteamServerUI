<script>
  import { getCurrentBackendUrl, apiFetch } from '../../services/api';
  import { onMount } from 'svelte';
  
  let iframeRef;
  let loading = true;
  let iframeSrc = '';
  let error = false;

  async function constructIframeSrc() {
    try {
      // Fetch the working directory
      const response = await apiFetch('/api/v2/getwd');
      if (!response.ok) {
        throw new Error('Failed to fetch working directory');
      }
      const data = await response.json();
      if (data.status !== 'OK' || !data.WorkingDir) {
        throw new Error('Invalid response from /api/v2/getwd');
      }

      const backendUrl = getCurrentBackendUrl();
      // Construct the full URL with the folder query parameter
      iframeSrc = `${backendUrl}/api/v2/codeserver?folder=${encodeURIComponent(data.WorkingDir)}`;
    } catch (error) {
      console.error('Error constructing iframe src:', error);
      // Fallback to URL without folder parameter
      const backendUrl = getCurrentBackendUrl();
      iframeSrc = `${backendUrl}/api/v2/codeserver`;
    }
  }

  function handleIframeLoad() {
    loading = false;
    error = false;
  }

  function handleIframeError() {
    loading = false;
    error = true;
    console.error('Iframe failed to load');
  }

  // Run on component mount
  onMount(() => {
    constructIframeSrc();
  });
</script>

{#if loading && iframeSrc}
  <div class="loading">
    <div class="spinner"></div>
    <p>Loading Code Server...</p>
  </div>
{/if}

{#if error}
  <div class="error">
    <p>Failed to load Code Server</p>
  </div>
{/if}

<div class="iframe-container" class:hidden={loading || !iframeSrc}>
  {#if iframeSrc}
    <iframe
      bind:this={iframeRef}
      src={iframeSrc}
      title="Code Server"
      on:load={handleIframeLoad}
      on:error={handleIframeError}
    ></iframe>
  {/if}
</div>

<style>
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
    min-height: 200px;
  }

  .loading p {
    margin: 1rem 0 0 0;
    color: #666;
    font-size: 14px;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #007acc;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .error {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
    min-height: 200px;
    color: #d73a49;
  }

  .error p {
    margin: 0;
    font-size: 14px;
  }
  
  .iframe-container {
    width: 100%;
    height: 100%;
    border: none;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .iframe-container.hidden {
    display: none;
  }
  
  iframe {
    width: 100%;
    height: 100%;
    border: none;
    background: white;
  }
</style>