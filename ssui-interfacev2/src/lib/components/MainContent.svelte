<script>
  import DashboardView from './Dashboard/DashboardView.svelte';
  import SettingsView from './settings/SettingsView.svelte';
  import LogsView from './views/LogsView.svelte';
  import ConsoleView from './views/ConsoleView.svelte';

  /**
   * @typedef {Object} Props
   * @property {string} [activeView]
   */

  /** @type {Props} */
  let { activeView = 'dashboard' } = $props();
  
  // View metadata for headers
  const viewContent = {
    dashboard: {
      title: 'Dashboard',
      description: 'Server overview and statistics'
    },
    servers: {
      title: 'Servers',
      description: 'Manage your game servers'
    },
    settings: {
      title: 'Settings',
      description: 'Configure application settings'
    },
    logs: {
      title: 'Logs',
      description: 'View server logs and events'
    },
    console: {
      title: 'Console',
      description: 'View server console output'
    }
  };
</script>

<main class="main-content">
    <div class="view-header">
      <h1 class:hide={activeView === 'dashboard'}>{viewContent[activeView].title}</h1>
      <p class="description" class:hide={activeView === 'dashboard'}>{viewContent[activeView].description}</p>
    </div>
  
  <div class="view-content">
    {#if activeView === 'dashboard'}
      <DashboardView />
    {:else if activeView === 'settings'}
      <SettingsView />
    {:else if activeView === 'logs'}
      <LogsView />
    {:else if activeView === 'console'}
      <ConsoleView />
    {/if}
  </div>
</main>

<style>
  .main-content {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  
  .view-header {
    margin-bottom: 1.5rem;
  }
  
  .view-header h1 {
    margin: 0 0 0.5rem 0;
    font-size: 1.8rem;
    font-weight: 500;
  }
  
  .view-header .hide {
  display: none;
  }
  
  .description {
    color: var(--text-secondary);
    margin: 0;
  }
  
  .view-content {
    flex: 1;
  }
</style>