<script>
  import { onMount } from 'svelte';
  import AppSettings from './AppSettings.svelte';
  import RunfileSettings from './RunfileSettings.svelte';
  import BackendSettings from './BackendSettings.svelte';
  import DetectionManager from './DetectionManager.svelte';
  import ConfigManager from './ConfigManager.svelte';
  import FileView from './FileView.svelte';
  import AlternativeRunfileSettings from './AlternativeRunfileSettings.svelte';
  // State management
  let activeSidebarTab = $state('SSUI Settings'); // Default to General tab in sidebar
  
  // Handle sidebar tab selection
  function selectSidebarTab(tab) {
    activeSidebarTab = tab;
  }
</script>

<div class="settings-container">
  <div class="settings-sidebar">
    <button 
      class="settings-nav {activeSidebarTab === 'SSUI Settings' ? 'active' : ''}" 
      onclick={() => selectSidebarTab('SSUI Settings')}>SSUI Settings</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Runfile' ? 'active' : ''}" 
      onclick={() => selectSidebarTab('Runfile')}>Game Settings</button>
    <button
      class="settings-nav {activeSidebarTab === 'Alternative View: Runfile Settings' ? 'active' : ''}"
      onclick={() => selectSidebarTab('Alternative View: Runfile Settings')}>Alternative View: Game Settings</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Backends' ? 'active' : ''}" 
      onclick={() => selectSidebarTab('Backends')}>Backends</button>
    <button 
        class="settings-nav {activeSidebarTab === 'Files' ? 'active' : ''}" 
        onclick={() => selectSidebarTab('Files')}>Files</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Legacy Detection Manager' ? 'active' : ''}"
      onclick={() => selectSidebarTab('Legacy Detection Manager')}>Legacy Detection Manager</button>
    <button 
      class="settings-nav {activeSidebarTab === 'Legacy Config Manager' ? 'active' : ''}" 
      onclick={() => selectSidebarTab('Legacy Config Manager')}>Legacy Config Manager</button>
  </div>
  
  <div class="settings-content">
    {#if activeSidebarTab === 'SSUI Settings'}
      <AppSettings {activeSidebarTab} />
    {:else if activeSidebarTab === 'Runfile'}
      <RunfileSettings />
    {:else if activeSidebarTab === 'Alternative View: Game Settings'}
      <AlternativeRunfileSettings />
    {:else if activeSidebarTab === 'Backends'}
      <BackendSettings />
    {:else if activeSidebarTab === 'Files'}
      <FileView />
    {:else if activeSidebarTab === 'Legacy Detection Manager'}
      <DetectionManager />
    {:else if activeSidebarTab === 'Legacy Config Manager'}
      <ConfigManager />
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