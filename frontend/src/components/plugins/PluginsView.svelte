<script>
  import { onMount } from 'svelte';
  import PluginContainer from './PluginContainer.svelte';
  import { pluginsList } from '../../services/plugins';
  
  let activeSidebarTab = $state('');
  let plugins = $derived($pluginsList);

  function extractPluginName(pluginPath) {
    // Extract plugin name from path like "/plugins/ExamplePlugin/"
    const name = pluginPath.split('/').filter(Boolean).pop();
    // Insert spaces before capital letters: ExamplePlugin -> Example Plugin
    return name.replace(/([A-Z])/g, ' $1').trim();
  }

  $effect(() => {
    // Set the first plugin as active if available
    if (plugins.length > 0 && !activeSidebarTab) {
      activeSidebarTab = plugins[0];
    }
  });
</script>

<div class="settings-container">
  <div class="settings-sidebar">
    {#each plugins as pluginPath}
      <button 
        class="settings-nav {activeSidebarTab === pluginPath ? 'active' : ''}" 
        onclick={() => { activeSidebarTab = pluginPath; }}>
        {extractPluginName(pluginPath)}
      </button>
    {/each}
  </div>
  
  <div class="settings-content">
    {#if activeSidebarTab}
      <PluginContainer pluginPath={activeSidebarTab} />
    {/if}
  </div>
</div>

<style>
  .settings-container {
    display: flex;
    gap: 2rem;
    height: 100%;
  }
  
  .settings-sidebar {
    width: 180px;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .settings-nav {
    text-align: left;
    padding: 0.75rem 1rem;
    background-color: transparent;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }
  
  .settings-nav:hover {
    background-color: var(--bg-hover);
  }
  
  .settings-nav.active {
    background-color: var(--bg-active);
    color: var(--accent-primary);
  }
  
  .settings-content {
    flex: 1;
    background-color: var(--bg-secondary);
    border-radius: 4px;
    padding: 1.5rem;
    overflow-y: auto;
  }
  
  @media (max-width: 768px) {
    .settings-container {
      flex-direction: column;
    }
    
    .settings-sidebar {
      width: 100%;
      flex-direction: row;
      overflow-x: auto;
      padding-bottom: 0.5rem;
    }
  }
</style>