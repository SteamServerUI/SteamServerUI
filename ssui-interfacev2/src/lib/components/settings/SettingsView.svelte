<script>
  import { onMount } from 'svelte';
  import AppSettings from './AppSettings.svelte';
  import RunfileSettings from './RunfileSettings.svelte';
  import DiscordSettings from './DiscordSettings.svelte';
  
  // State management
  let activeSidebarTab = 'General'; // Default to General tab in sidebar
  
  // Handle sidebar tab selection
  function selectSidebarTab(tab) {
    activeSidebarTab = tab;
  }
</script>

<div class="settings-container">
  <div class="settings-sidebar">
    <button 
      class="settings-nav {activeSidebarTab === 'General' ? 'active' : ''}" 
      on:click={() => selectSidebarTab('General')}>General</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Runfile' ? 'active' : ''}" 
      on:click={() => selectSidebarTab('Runfile')}>Runfile Settings</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Advanced' ? 'active' : ''}" 
      on:click={() => selectSidebarTab('Advanced')}>Advanced</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Discord' ? 'active' : ''}" 
      on:click={() => selectSidebarTab('Discord')}>Discord</button>
  </div>
  
  <div class="settings-content">
    {#if ['General', 'Appearance', 'SteamCMD', 'Notifications', 'Advanced'].includes(activeSidebarTab)}
      <AppSettings {activeSidebarTab} />
    {:else if activeSidebarTab === 'Runfile'}
      <RunfileSettings />
    {:else if activeSidebarTab === 'Discord'}
      <DiscordSettings />
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