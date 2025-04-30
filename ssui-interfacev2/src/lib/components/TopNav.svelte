<script>
    import { onMount } from 'svelte';
    
    export let views = [];
    export let activeView = 'dashboard';
    export let setActiveView;
    
    let currentTime = new Date();
    
    // Update the time every minute
    onMount(() => {
      const interval = setInterval(() => {
        currentTime = new Date();
      }, 60000);
      
      return () => {
        clearInterval(interval);
      };
    });
    
    $: formattedTime = currentTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  </script>
  
  <nav class="top-nav">
    <div class="logo">
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

    <div class="backend-dropdown">
      <div class="dropdown-header">
        <h3>Select Backend</h3>
      </div>
    </div>
    
    <div class="user-area">
      <span class="time">{formattedTime}</span>
      <button class="user-button">
        <span class="user-icon">👤</span>
      </button>
    </div>
  </nav>
  
  <style>
    .top-nav {
      height: var(--top-nav-height);
      background-color: var(--bg-secondary);
      display: flex;
      align-items: center;
      padding: 0 1rem;
      box-shadow: var(--shadow-light);
      z-index: 10;
    }
    
    .logo {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-weight: 600;
      font-size: 1.1rem;
      margin-right: 2rem;
    }
    
    .logo-icon {
      font-size: 1.3rem;
    }
    
    .nav-buttons {
      display: flex;
      gap: 0.5rem;
      flex: 1;
    }
    
    .user-area {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    
    .time {
      color: var(--text-secondary);
      font-size: 0.9rem;
    }
    
    .user-button {
      padding: 0.25rem;
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 50%;
    }

    .backend-dropdown {
      position: absolute;
      top: 100%;
      left: 0;
      right: 0;
      background: var(--bg-secondary);
      border-radius: 8px;
      box-shadow: var(--shadow-medium);
      margin-top: 0.5rem;
      z-index: 10;
      overflow: hidden;
      animation: slideDown 0.2s ease-out;
    }
    
  </style>