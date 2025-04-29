<!-- BackendSettings.svelte - Component for managing backend connections -->
<script>
  import { onMount } from 'svelte';
  import { backendConfig, setBackend, setActiveBackend, initializeApiService, apiFetch } from '../../services/api';
  
  let currentConfig;
  let newBackendId = '';
  let newBackendUrl = '';
  let activeBackend = '';
  let backends = [];
  
  // Subscribe to the backend config store
  const unsubscribe = backendConfig.subscribe(value => {
    currentConfig = value;
    activeBackend = value.active;
    backends = Object.keys(value.backends);
  });
  
  onMount(() => {
    // Initialize the API service
    initializeApiService();
    
    // Cleanup subscription when component is destroyed
    return () => {
      unsubscribe();
    };
  });
  
  // Add a new backend
  function addBackend() {
    if (!newBackendId || !newBackendUrl) return;
    
    setBackend(newBackendId, newBackendUrl);
    
    // Clear input fields
    newBackendId = '';
    newBackendUrl = '';
  }
  
  // Handle active backend change
  function changeActiveBackend() {
    setActiveBackend(activeBackend);
  }
  
  // Remove a backend
  function removeBackend(id) {
    if (id === 'default') {
      alert('Cannot remove the default backend');
      return;
    }
    
    backendConfig.update(config => {
      // If removing the active backend, switch to default
      if (config.active === id) {
        config.active = 'default';
        activeBackend = 'default';
      }
      
      // Remove the backend
      delete config.backends[id];
      return config;
    });
  }
  
  // Test connection to a backend
  async function testConnection(id) {
    const backendUrl = currentConfig.backends[id].url;
    try {
      const response = await apiFetch(`/api/v2/server/status`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      if (response.ok) {
        alert(`Connection to ${id} successful!`);
      } else {
        alert(`Failed to connect to ${id}. Status: ${response.status}`);
      }
    } catch (error) {
      alert(`Error connecting to ${id}: ${error.message}`);
    }
  }
</script>

<div class="backend-settings">
  <div class="section-header">
    <h2>Backend Connections</h2>
  </div>
  
  <div class="current-backends">
    <div class="backend-controls">
      <h3>Configured Backends</h3>
      <div class="search-box">
        <input type="text" placeholder="Search backends..." />
      </div>
    </div>
    
    <div class="backend-list">
      {#each backends as backendId}
        <div class="backend-item">
          <div class="backend-icon">{backendId === activeBackend ? '🟢' : '⚪'}</div>
          <div class="backend-info">
            <h3>{backendId}</h3>
            <p>{currentConfig.backends[backendId].url}</p>
          </div>
          <div class="backend-status">
            <label class="status-toggle">
              <input 
                type="radio" 
                name="activeBackend" 
                value={backendId} 
                bind:group={activeBackend} 
                on:change={changeActiveBackend}
              />
              <span class="status-label">{backendId === activeBackend ? 'Active' : 'Set Active'}</span>
            </label>
          </div>
          <div class="backend-actions">
            <button class="action-button" on:click={() => testConnection(backendId)}>Test</button>
            {#if backendId !== 'default'}
              <button class="action-button danger" on:click={() => removeBackend(backendId)}>Remove</button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
  
  <div class="add-backend">
    <div class="section-header">
      <h3>Add New Backend</h3>
    </div>
    
    <div class="form-container">
      <div class="form-group">
        <label for="backend-id">Backend ID</label>
        <input 
          type="text" 
          id="backend-id" 
          bind:value={newBackendId} 
          placeholder="e.g., production"
        />
      </div>
      
      <div class="form-group">
        <label for="backend-url">Backend URL</label>
        <input 
          type="text" 
          id="backend-url" 
          bind:value={newBackendUrl} 
          placeholder="e.g., https://api.example.com"
        />
      </div>
      
      <button class="primary-button" on:click={addBackend}>+ Add Backend</button>
    </div>
  </div>
</div>

<style>
  .backend-settings {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1rem;
  }
  
  .section-header {
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--border-color);
  }
  
  .section-header h2, .section-header h3 {
    margin: 0;
    color: var(--text-primary);
  }
  
  .backend-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .search-box {
    flex: 0 0 200px;
  }
  
  .search-box input {
    width: 100%;
    padding: 0.5rem;
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: 4px;
  }
  
  .backend-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  
  .backend-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    background-color: var(--bg-secondary);
    border-radius: 4px;
    padding: 1rem;
    box-shadow: var(--shadow-light);
  }
  
  .backend-icon {
    font-size: 1.5rem;
  }
  
  .backend-info {
    flex: 1;
  }
  
  .backend-info h3 {
    margin: 0 0 0.25rem 0;
    font-size: 1.1rem;
  }
  
  .backend-info p {
    margin: 0;
    color: var(--text-secondary);
    font-size: 0.9rem;
  }
  
  .backend-status {
    display: flex;
    align-items: center;
  }
  
  .status-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .status-label {
    color: var(--text-secondary);
    font-size: 0.9rem;
  }
  
  .backend-actions {
    display: flex;
    gap: 0.5rem;
  }
  
  .action-button {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .action-button:hover {
    background-color: var(--bg-hover);
  }
  
  .action-button.danger {
    color: var(--danger);
  }
  
  .action-button.danger:hover {
    background-color: var(--danger-bg);
  }
  
  .form-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 4px;
    box-shadow: var(--shadow-light);
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .form-group label {
    color: var(--text-primary);
    font-size: 0.9rem;
    font-weight: 500;
  }
  
  .form-group input {
    padding: 0.5rem;
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: 4px;
  }
  
  .form-group input:focus {
    border-color: var(--accent-primary);
    outline: none;
  }
  
  .primary-button {
    background-color: var(--accent-primary);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    font-weight: 500;
    border-radius: 4px;
    cursor: pointer;
    align-self: flex-start;
    transition: background-color 0.2s ease;
  }
  
  .primary-button:hover {
    background-color: var(--accent-secondary);
  }
  
  @media (max-width: 768px) {
    .backend-item {
      flex-wrap: wrap;
    }
    
    .backend-actions {
      width: 100%;
      margin-top: 0.5rem;
      justify-content: flex-end;
    }
    
    .backend-status {
      margin-left: auto;
    }
  }
</style>