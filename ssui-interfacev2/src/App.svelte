<script>
  import TopNav from './lib/components/TopNav.svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import MainContent from './lib/components/MainContent.svelte';
  import BackendInitializer from './BackendInitializer.svelte';
  import AuthGuard from './AuthGuard.svelte';
  import './lib/theme.css';

  // Track active view
  let activeView = 'dashboard';
  
  // Views available in the app
  const views = [
    { id: 'dashboard', name: 'Dashboard', icon: 'grid' },
    { id: 'servers', name: 'Servers', icon: 'server' },
    { id: 'settings', name: 'Settings', icon: 'settings' },
    { id: 'logs', name: 'Logs', icon: 'file-text' },
    { id: 'console', name: 'Console', icon: 'terminal' }
  ];

  // Set active view function
  function setActiveView(viewId) {
    activeView = viewId;
  }

  // Server status state
  let serverStatus = 'checking';
  let serverError = null;
  
  function handleStatusChange(status) {
    serverStatus = status.status;
    serverError = status.error;
  }
</script>

<BackendInitializer onStatusChange={handleStatusChange}>
  <AuthGuard serverStatus={serverStatus} serverError={serverError}>
    <div class="app-container">
      <TopNav {views} {activeView} {setActiveView} />
      
      <div class="main-container">
        <Sidebar {views} {activeView} {setActiveView} />
        <MainContent {activeView} />
      </div>
    </div>
  </AuthGuard>
</BackendInitializer>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
  }
  
  .main-container {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
</style>