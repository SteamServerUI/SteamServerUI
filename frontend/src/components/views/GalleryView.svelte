<script lang="js">
  import { fade } from 'svelte/transition';
  import RunfileGallery from './RunfileGalleryView.svelte';
  import PluginGallery from './PluginGalleryView.svelte';

  // Reactive state for current view
  let currentView = $state('runfile'); // Default to Runfile Gallery

  // Switch view and update URL
  function switchView(view) {
    currentView = view;
  }

  // Handle browser back/forward navigation
  $effect(() => {
    function handlePopstate() {
      const path = window.location.pathname.split('/').pop() || 'runfile';
      currentView = path === 'plugin' ? 'plugin' : 'runfile';
    }
    window.addEventListener('popstate', handlePopstate);
    return () => window.removeEventListener('popstate', handlePopstate);
  });
</script>

<div class="gallery-container">
  <!-- Toggle Buttons -->
  <div class="gallery-toggle">
    <button
      class="toggle-button"
      class:active={currentView === 'runfile'}
      onclick={() => switchView('runfile')}
      aria-label="Show Runfile Gallery"
    >
      Runfile Gallery
    </button>
    <button
      class="toggle-button"
      class:active={currentView === 'plugin'}
      onclick={() => switchView('plugin')}
      aria-label="Show Plugin Gallery"
    >
      Plugin Gallery
    </button>
  </div>

  <!-- Conditional Rendering of Galleries with Transitions -->
  {#if currentView === 'runfile'}
    <div transition:fade={{ duration: 250 }} class="gallery-wrapper">
      <RunfileGallery />
    </div>
  {:else if currentView === 'plugin'}
    <div transition:fade={{ duration: 250 }} class="gallery-wrapper">
      <PluginGallery />
    </div>
  {/if}
</div>

<style>
  .gallery-container {
    max-width: 100%;
    margin: 0 auto;
    padding: 1rem;
  }
  .gallery-toggle {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    justify-content: center;
  }
  .toggle-button {
    background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
    color: white;
    border: none;
    border-radius: 12px;
    padding: 0.75rem 1.5rem;
    cursor: pointer;
    font-weight: 600;
    font-size: 1rem;
    transition: all var(--transition-speed, 0.3s) ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    box-shadow: var(--shadow-light, 0 4px 8px rgba(0,0,0,0.1));
  }
  .toggle-button:hover:not(.active) {
    transform: translateY(-2px);
    box-shadow: var(--shadow-medium, 0 8px 16px rgba(0,0,0,0.15));
  }
  .toggle-button.active {
    background: linear-gradient(135deg, var(--accent-primary, #1F4A33), var(--accent-secondary, #3267B2));
    transform: translateY(0);
    box-shadow: var(--shadow-medium, 0 8px 16px rgba(0,0,0,0.15));
    cursor: default;
  }
  .toggle-button:not(.active) {
    opacity: 0.85;
  }
  .gallery-wrapper {
    position: relative;
  }
</style>