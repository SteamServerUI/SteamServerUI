<script>
  
  let iframeRef;
  let loading = true;
  
  function handleIframeLoad() {
    loading = false;
  }
  
  function handleIframeError() {
    loading = false;
  }
</script>

{#if loading}
  <div class="loading">
    <p>Loading Code Server...</p>
  </div>
{/if}

<div class="iframe-container" class:hidden={loading}>
  <iframe 
    bind:this={iframeRef}
    src="/api/v2/codeserver"
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
    border-radius: 4px;
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