<script>
  import { getCurrentBackendUrl, apiFetch } from '../../services/api';
  import { onMount } from 'svelte';
  
  let iframeRef;
  let loading = true;
  let iframeSrc = '';

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
  }

  function handleIframeError() {
    loading = false;
    console.error('Iframe failed to load');
  }

  // Run on component mount
  onMount(() => {
    constructIframeSrc();
  });
</script>

{#if loading}
  <div class="loading">
    <p>Loading Code Server...</p>
  </div>
{/if}

<div class="iframe-container" class:hidden={loading}>
  <iframe
    bind:this={iframeRef}
    src={iframeSrc}
    title="Code Server"
    on:load={handleIframeLoad}
    on:error={handleIframeError}
  ></iframe>
</div>

<style>
  .loading {
    padding: 2rem;
    text-align: center;
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
    background: none;
  }
</style>