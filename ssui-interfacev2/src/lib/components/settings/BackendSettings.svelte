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
        const response = await apiFetch(`${backendUrl}/api/v2/server/status`, {
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
    <h2>Backend Connections</h2>
    
    <div class="current-backends">
      <h3>Configured Backends</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>URL</th>
            <th>Active</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each backends as backendId}
            <tr>
              <td>{backendId}</td>
              <td>{currentConfig.backends[backendId].url}</td>
              <td>
                <input 
                  type="radio" 
                  name="activeBackend" 
                  value={backendId} 
                  bind:group={activeBackend} 
                  on:change={changeActiveBackend}
                />
              </td>
              <td>
                <button on:click={() => testConnection(backendId)}>Test</button>
                {#if backendId !== 'default'}
                  <button on:click={() => removeBackend(backendId)}>Remove</button>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    
    <div class="add-backend">
      <h3>Add New Backend</h3>
      <div class="form-group">
        <label for="backend-id">Backend ID:</label>
        <input 
          type="text" 
          id="backend-id" 
          bind:value={newBackendId} 
          placeholder="e.g., production"
        />
      </div>
      
      <div class="form-group">
        <label for="backend-url">Backend URL:</label>
        <input 
          type="text" 
          id="backend-url" 
          bind:value={newBackendUrl} 
          placeholder="e.g., https://api.example.com"
        />
      </div>
      
      <button on:click={addBackend}>Add Backend</button>
    </div>
  </div>
  
  <style>
    .backend-settings {
      padding: 1rem;
    }
    
    .form-group {
      margin-bottom: 1rem;
    }
    
    .form-group label {
      display: block;
      margin-bottom: 0.5rem;
    }
    
    input[type="text"] {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid #ccc;
      border-radius: 4px;
    }
    
    button {
      padding: 0.5rem 1rem;
      background-color: #4a5568;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      margin-right: 0.5rem;
    }
    
    button:hover {
      background-color: #2d3748;
    }
    
    table {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 1rem;
    }
    
    table th, table td {
      padding: 0.5rem;
      border: 1px solid #e2e8f0;
      text-align: left;
    }
    
    table th {
      background-color: #f7fafc;
    }
  </style>