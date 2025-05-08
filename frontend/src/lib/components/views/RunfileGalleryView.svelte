<script lang="js">
  import { apiFetch } from '../../services/api';
  import RunfileCard from '../Gallery/RunfileCard.svelte';

  // State
  let runfiles = $state([]);
  let loading = $state(true);
  let error = $state(null);

  // Fetch runfiles on mount
  async function fetchRunfiles(forceUpdate = false) {
    loading = true;
    error = null;
    
    try {
      // Get the response from apiFetch
      const response = await apiFetch(`/api/v2/gallery${forceUpdate ? '?forceUpdate=true' : ''}`);
      
      // If your apiFetch doesn't parse JSON automatically, do this:
      let data;
      if (response instanceof Response) {
        // This means apiFetch returns the raw Response object
        console.log('API Response status:', response.status);
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        data = await response.json();
      } else {
        // apiFetch already parsed the JSON for us
        data = response;
      }
      
      console.log('Parsed data:', data);
      
      // Handle different response formats
      if (data && data.data && Array.isArray(data.data)) {
        // Format: { data: [...runfiles] }
        runfiles = data.data;
      } else if (Array.isArray(data)) {
        // Format: [...runfiles]
        runfiles = data;
      } else {
        console.warn('Unexpected response format:', data);
        runfiles = [];
      }
      
      // Add missing identifiers if needed
      runfiles = runfiles.map(runfile => {
        if (!runfile.identifier) {
          runfile.identifier = runfile.filename || runfile.name;
        }
        return runfile;
      });
      
      console.log('Processed runfiles:', runfiles);
      
    } catch (err) {
      error = err.message || 'Failed to fetch runfiles';
      console.error('Error fetching runfiles:', err);
      runfiles = []; // Ensure runfiles is always an array
    } finally {
      loading = false;
    }
  }

  // Refresh gallery
  function refreshGallery() {
    fetchRunfiles(true);
  }

  // Fetch runfiles on mount
  $effect(() => {
    fetchRunfiles();
  });
</script>

<div class="gallery-page">
  <div class="header">
    <button 
      class="refresh-button" 
      onclick={refreshGallery} 
      disabled={loading}
      aria-label="Refresh runfile gallery"
    >
      {loading ? 'Loading...' : 'Refresh Gallery'}
    </button>
  </div>
  
  {#if error}
    <div class="error-banner">
      <p>{error}</p>
    </div>
  {/if}
  
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading runfiles...</p>
    </div>
  {:else if !runfiles || runfiles.length === 0}
    <div class="empty-state">
      <p>No runfiles available. Try refreshing the gallery.</p>
    </div>
  {:else}
    <div class="gallery-grid">
      {#each runfiles as runfile (runfile.identifier)}
        <RunfileCard {runfile} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .gallery-page {
    padding: 1rem;
    background-color: var(--bg-primary, #f0f0f0);
    min-height: 100vh;
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  h1 {
    color: var(--text-primary, #333333);
    margin: 0;
  }
  
  .refresh-button {
    background-color: var(--accent-primary, #2C6E49);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s ease;
  }
  
  .refresh-button:hover:not(:disabled) {
    transform: scale(1.05);
  }
  
  .refresh-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .error-banner {
    background-color: var(--bg-error, #ffeded);
    color: var(--text-warning, #FF6B6B);
    padding: 1rem;
    border-radius: 6px;
    margin-bottom: 1rem;
  }
  
  .error-banner p {
    margin: 0;
  }
  
  .loading, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--bg-secondary, #eeeeee);
    border-top: 4px solid var(--accent-primary, #2C6E49);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }
  
  .gallery-grid {
    display: ruby;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  @media (max-width: 600px) {
    .gallery-grid {
      grid-template-columns: 1fr;
    }
    
    .header {
      flex-direction: column;
      gap: 1rem;
      align-items: flex-start;
    }
  }
</style>