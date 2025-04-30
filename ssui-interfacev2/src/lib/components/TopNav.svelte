<script>
  import { onMount, onDestroy } from 'svelte';
  import { backendConfig, setActiveBackend } from '../services/api';
  
  export let views = [];
  export let activeView = 'dashboard';
  export let setActiveView;
  
  let currentTime = new Date();
  let showBackendDropdown = false;
  let showUserMenu = false;
  let backends = [];
  let activeBackend = '';
  let backendStatus = {};
  let timeoutId;
  let clickOutsideHandler;
  
  // Subscribe to backend config
  const unsubscribe = backendConfig.subscribe(value => {
    backends = Object.keys(value.backends);
    activeBackend = value.active;
    
    // Initialize status objects
    backends.forEach(id => {
      if (!(id in backendStatus)) {
        backendStatus[id] = { status: 'unknown', lastChecked: null };
      }
    });
  });
  
  // Handle clicks outside dropdowns
  function setupClickOutsideHandler() {
    clickOutsideHandler = (event) => {
      const backendElement = document.querySelector('.backend-selector');
      const userMenuElement = document.querySelector('.user-menu-container');
      
      if (backendElement && !backendElement.contains(event.target)) {
        showBackendDropdown = false;
      }
      
      if (userMenuElement && !userMenuElement.contains(event.target)) {
        showUserMenu = false;
      }
    };
    
    document.addEventListener('click', clickOutsideHandler);
  }
  
  // Update the time every minute
  onMount(() => {
    setupClickOutsideHandler();
    
    const interval = setInterval(() => {
      currentTime = new Date();
    }, 60000);
    
    return () => {
      clearInterval(interval);
      unsubscribe();
      document.removeEventListener('click', clickOutsideHandler);
    };
  });
  
  function toggleBackendDropdown(event) {
    event.stopPropagation();
    showBackendDropdown = !showBackendDropdown;
    if (showBackendDropdown) {
      showUserMenu = false;
    }
  }
  
  function toggleUserMenu(event) {
    event.stopPropagation();
    showUserMenu = !showUserMenu;
    if (showUserMenu) {
      showBackendDropdown = false;
    }
  }
  
  function changeActiveBackend(id) {
    setActiveBackend(id);
    
    // Show feedback
    backendStatus[id] = { ...backendStatus[id], status: 'connecting' };
    
    // Simulate connection process (would be replaced with actual async operation)
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
      backendStatus[id] = { status: 'connected', lastChecked: new Date() };
    }, 800);
    
    showBackendDropdown = false;
  }

  function getStatusIndicator(id) {
    const status = backendStatus[id]?.status || 'unknown';
    switch(status) {
      case 'connected': return '🟢';
      case 'connecting': return '🟠';
      case 'error': return '🔴';
      default: return '⚪';
    }
  }
  
  function slide(node, { duration = 300 }) {
    return {
      duration,
      css: t => `
        transform: translateY(${(1 - t) * -10}px);
        opacity: ${t};
      `
    };
  }

  $: formattedTime = currentTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  $: formattedDate = currentTime.toLocaleDateString([], { month: 'short', day: 'numeric' });
</script>

<nav class="top-nav">
  <div class="nav-left">
    <div class="logo" on:click={() => setActiveView('dashboard')}>
      <span class="logo-icon">⚙️</span>
      <span class="logo-text">SSUI</span>
    </div>
    
    <div class="nav-buttons">
      {#each views as view}
        <button 
          class={activeView === view.id ? 'active' : ''} 
          on:click={() => setActiveView(view.id)}
        >
          {view.name}
        </button>
      {/each}
    </div>
  </div>

  <div class="nav-right">
    <div class="backend-selector" on:click|stopPropagation>
      <button class="backend-toggle" on:click={toggleBackendDropdown}>
        <span class="status-indicator">{getStatusIndicator(activeBackend)}</span>
        <span class="backend-label">{activeBackend}</span>
        <span class="dropdown-arrow">{showBackendDropdown ? '▲' : '▼'}</span>
      </button>
      
      {#if showBackendDropdown}
        <div class="backend-dropdown" in:slide={{ duration: 150 }} out:slide={{ duration: 150 }}>
          <div class="dropdown-header">
            <h3>Select Backend</h3>
          </div>
          <div class="dropdown-content">
            {#each backends as backendId}
              <div 
                class="dropdown-item {backendId === activeBackend ? 'active' : ''}"
                on:click={() => changeActiveBackend(backendId)}
              >
                <div class="backend-info">
                  <span class="status-dot" style="background-color: {backendStatus[backendId]?.status === 'connected' ? 'var(--accent-primary)' : backendStatus[backendId]?.status === 'connecting' ? 'var(--text-warning)' : 'var(--text-secondary)'}"></span>
                  <span class="backend-name">{backendId}</span>
                </div>
                {#if backendId === activeBackend}
                  <span class="active-marker">✓</span>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
    
    <div class="datetime">
      <span class="date">{formattedDate}</span>
      <span class="time">{formattedTime}</span>
    </div>
    
    <div class="user-menu-container" on:click|stopPropagation>
      <button class="user-button" on:click={toggleUserMenu}>
        <span class="user-avatar">SA</span>
      </button>
      
      {#if showUserMenu}
        <div class="user-dropdown" in:slide={{ duration: 150 }} out:slide={{ duration: 150 }}>
          <div class="user-dropdown-header">
            <div class="user-info">
              <div class="user-avatar large">SA</div>
              <div class="user-details">
                <div class="user-name">Admin</div>
                <div class="user-email">Superadmin</div>
              </div>
            </div>
          </div>
          <div class="dropdown-content">
            <div class="dropdown-item">
              <span class="item-icon">🌙</span>
              <span>Theme</span>
            </div>
            <div class="divider"></div>
            <div class="dropdown-item logout">
              <span class="item-icon">🚪</span>
              <span>Logout</span>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</nav>

<style>
  .top-nav {
    height: var(--top-nav-height);
    background-color: var(--bg-secondary);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1rem;
    box-shadow: var(--shadow-light);
    z-index: 100;
    border-bottom: 1px solid var(--border-color);
  }
  
  .nav-left, .nav-right {
    display: flex;
    align-items: center;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 600;
    font-size: 1.1rem;
    margin-right: 2rem;
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    cursor: pointer;
    transition: background-color var(--transition-speed) ease;
  }
  
  .logo:hover {
    background-color: var(--bg-hover);
  }
  
  .logo-icon {
    font-size: 1.3rem;
  }
  
  .logo-text {
    color: var(--accent-primary);
  }
  
  .nav-buttons {
    display: flex;
    gap: 0.25rem;
  }
  
  .nav-buttons button {
    background: transparent;
    color: var(--text-secondary);
    border: none;
    border-radius: 6px;
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    transition: all var(--transition-speed) ease;
    position: relative;
  }
  
  .nav-buttons button:hover {
    background-color: var(--bg-hover);
    color: var(--text-primary);
  }
  
  .nav-buttons button.active {
    color: var(--accent-primary);
    background-color: var(--bg-active);
    font-weight: 500;
  }
  
  .nav-buttons button.active::after {
    content: "";
    position: absolute;
    bottom: -1px;
    left: 25%;
    width: 50%;
    height: 2px;
    background-color: var(--accent-primary);
    border-radius: 2px;
  }
  
  /* Backend selector styles */
  .backend-selector {
    position: relative;
    margin-right: 1rem;
  }
  
  .backend-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.4rem 0.75rem;
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-size: 0.85rem;
    color: var(--text-primary);
    cursor: pointer;
    transition: all var(--transition-speed) ease;
  }
  
  .backend-toggle:hover {
    background-color: var(--bg-hover);
    border-color: var(--accent-tertiary);
  }
  
  .status-indicator {
    font-size: 0.7rem;
  }
  
  .dropdown-arrow {
    font-size: 0.7rem;
    color: var(--text-secondary);
    margin-left: 0.25rem;
  }
  
  .backend-dropdown, .user-dropdown {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    min-width: 220px;
    background: var(--bg-secondary);
    border-radius: 8px;
    box-shadow: var(--shadow-medium);
    z-index: 20;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }
  
  @keyframes slide {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  .dropdown-header, .user-dropdown-header {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border-color);
  }
  
  .dropdown-header h3 {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--text-primary);
  }
  
  .dropdown-content {
    max-height: 300px;
    overflow-y: auto;
  }
  
  .dropdown-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.65rem 1rem;
    cursor: pointer;
    transition: background-color var(--transition-speed) ease;
    font-size: 0.85rem;
  }
  
  .dropdown-item:hover {
    background-color: var(--bg-hover);
  }
  
  .dropdown-item.active {
    background-color: rgba(106, 153, 85, 0.1);
    font-weight: 500;
  }
  
  .backend-info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
  }
  
  .backend-name {
    color: var(--text-primary);
  }
  
  .active-marker {
    color: var(--accent-primary);
    font-weight: bold;
  }
  
  /* DateTime styles */
  .datetime {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    margin-right: 1rem;
  }
  
  .date {
    font-size: 0.7rem;
    color: var(--text-secondary);
    line-height: 1;
  }
  
  .time {
    font-size: 0.9rem;
    color: var(--text-primary);
    font-weight: 500;
  }
  
  /* User menu styles */
  .user-menu-container {
    position: relative;
  }
  
  .user-button {
    padding: 0;
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    border: none;
    background: transparent;
    cursor: pointer;
  }
  
  .user-avatar {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    background-color: var(--accent-secondary);
    color: var(--text-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.8rem;
    font-weight: 500;
    transition: all var(--transition-speed) ease;
  }
  
  .user-button:hover .user-avatar {
    background-color: var(--accent-primary);
  }
  
  .user-dropdown {
    right: 0;
    width: 240px;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 0;
  }
  
  .user-avatar.large {
    width: 2.5rem;
    height: 2.5rem;
    font-size: 1rem;
  }
  
  .user-details {
    display: flex;
    flex-direction: column;
  }
  
  .user-name {
    font-weight: 500;
    color: var(--text-primary);
    font-size: 0.9rem;
  }
  
  .user-email {
    color: var(--text-secondary);
    font-size: 0.75rem;
  }
  
  .item-icon {
    margin-right: 0.5rem;
    font-size: 0.9rem;
    min-width: 1rem;
    display: inline-flex;
    justify-content: center;
  }
  
  .divider {
    height: 1px;
    background-color: var(--border-color);
    margin: 0.25rem 0;
  }
  
  .dropdown-item.logout {
    color: var(--text-warning);
  }
  
  /* For animations */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  /* Slide animation for dropdowns */
  .slide-enter {
    opacity: 0;
    transform: translateY(-10px);
  }
  
  .slide-enter-active {
    opacity: 1;
    transform: translateY(0);
    transition: opacity 150ms, transform 150ms;
  }
  
  .slide-exit {
    opacity: 1;
    transform: translateY(0);
  }
  
  .slide-exit-active {
    opacity: 0;
    transform: translateY(-10px);
    transition: opacity 150ms, transform 150ms;
  }

  /* Additional function for transitions */
</style>