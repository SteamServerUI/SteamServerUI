<script>
    export let views = [];
    export let activeView = 'dashboard';
    export let setActiveView;
    
    let isHovered = false;
    let isPinned = false;
    
    function togglePin() {
      isPinned = !isPinned;
    }
    
    function handleMouseEnter() {
      isHovered = true;
    }
    
    function handleMouseLeave() {
      isHovered = false;
    }
    
    $: isExpanded = isHovered || isPinned;
  </script>
  
  <aside 
    class="sidebar" 
    class:expanded={isExpanded}
    on:mouseenter={handleMouseEnter}
    on:mouseleave={handleMouseLeave}
  >
    <div class="pin-container">
      <button class="pin-button" class:active={isPinned} on:click={togglePin}>
        {#if isPinned}
          📌
        {:else}
          📍
        {/if}
      </button>
    </div>
    
    <nav class="sidebar-nav">
      {#each views as view}
        <button 
          class="sidebar-button" 
          class:active={activeView === view.id}
          on:click={() => setActiveView(view.id)}
        >
          <span class="sidebar-icon">
            {#if view.icon === 'grid'}
              📊
            {:else if view.icon === 'server'}
              🖥️
            {:else if view.icon === 'settings'}
              ⚙️
            {:else if view.icon === 'file-text'}
              📝
            {/if}
          </span>
          <span class="sidebar-text">{view.name}</span>
        </button>
      {/each}
    </nav>
    
    <div class="sidebar-footer">
      <div class="server-status">
        <div class="status-indicator online"></div>
        <span class="status-text">Servers: 3 Online</span>
      </div>
      
      <div class="memory-usage">
        <div class="progress-bar">
          <div class="progress" style="width: 65%"></div>
        </div>
        <span class="memory-text">Memory: 65%</span>
      </div>
    </div>
  </aside>
  
  <style>
    .sidebar {
      width: var(--sidebar-collapsed-width);
      height: 100%;
      background-color: var(--bg-secondary);
      border-right: 1px solid var(--border-color);
      display: flex;
      flex-direction: column;
      overflow: hidden;
      transition: width var(--transition-speed) ease;
      z-index: 5;
    }
    
    .sidebar.expanded {
      width: var(--sidebar-width);
    }
    
    .pin-container {
      display: flex;
      justify-content: flex-end;
      padding: 0.5rem;
    }
    
    .pin-button {
      padding: 0.25rem;
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: transparent;
      border: none;
      cursor: pointer;
      opacity: 0.6;
      transition: opacity var(--transition-speed) ease;
    }
    
    .pin-button:hover, .pin-button.active {
      opacity: 1;
    }
    
    .sidebar-nav {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
      padding: 0.5rem;
      overflow-y: auto;
    }
    
    .sidebar-button {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem;
      border: none;
      background-color: transparent;
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
      justify-content: flex-start;
      width: 100%;
      text-align: left;
    }
    
    .sidebar-button:hover {
      background-color: var(--bg-hover);
    }
    
    .sidebar-button.active {
      background-color: var(--bg-active);
      color: var(--accent-primary);
    }
    
    .sidebar-icon {
      font-size: 1.2rem;
      min-width: 1.5rem;
      display: flex;
      justify-content: center;
    }
    
    .sidebar-text {
      white-space: nowrap;
      opacity: 0;
      transition: opacity var(--transition-speed) ease;
    }
    
    .expanded .sidebar-text {
      opacity: 1;
    }
    
    .sidebar-footer {
      padding: 1rem 0.5rem;
      border-top: 1px solid var(--border-color);
      font-size: 0.8rem;
      color: var(--text-secondary);
    }
    
    .server-status, .memory-usage {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
      white-space: nowrap;
    }
    
    .status-indicator {
      width: 8px;
      height: 8px;
      border-radius: 50%;
    }
    
    .status-indicator.online {
      background-color: #4caf50;
    }
    
    .progress-bar {
      flex: 1;
      height: 6px;
      background-color: var(--bg-tertiary);
      border-radius: 3px;
      overflow: hidden;
    }
    
    .progress {
      height: 100%;
      background-color: var(--accent-primary);
    }
    
    .status-text, .memory-text {
      opacity: 0;
      transition: opacity var(--transition-speed) ease;
      flex: 1;
    }
    
    .expanded .status-text, .expanded .memory-text {
      opacity: 1;
    }
  </style>