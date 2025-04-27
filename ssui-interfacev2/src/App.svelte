<script>
  import TopNav from './lib/components/TopNav.svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import MainContent from './lib/components/MainContent.svelte';
  import './lib/theme.css';
  
  // Track active view
  let activeView = 'dashboard';
  
  // Views available in the app
  const views = [
    { id: 'dashboard', name: 'Dashboard', icon: 'grid' },
    { id: 'servers', name: 'Servers', icon: 'server' },
    { id: 'settings', name: 'Settings', icon: 'settings' },
    { id: 'logs', name: 'Logs', icon: 'file-text' }
  ];
  
  // Set active view function
  function setActiveView(viewId) {
    activeView = viewId;
  }
</script>

<div class="app-container">
  <TopNav {views} {activeView} {setActiveView} />
  
  <div class="main-container">
    <Sidebar {views} {activeView} {setActiveView} />
    <MainContent {activeView} />
  </div>
</div>

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