<script>
    import { onMount } from 'svelte';
    
    // Props
    export let activeSidebarTab;
    
    // State management
    let settingsData = [];
    let settingsGroups = [];
    let activeSettingsGroup = '';
    let statusMessage = '';
    let isError = false;
    let statusTimeout;
    
    // Fetch settings data on component mount
    onMount(async () => {
      await fetchSettings();
    });
    
    async function fetchSettings() {
      try {
        const response = await fetch('/api/v2/settings');
        const { data, error } = await response.json();
        
        if (error) {
          showStatus(`Failed to load settings: ${error}`, true);
          return;
        }
        
        settingsData = data;
        // Extract unique groups
        settingsGroups = [...new Set(settingsData.map(s => s.group))];
        
        if (settingsGroups.length > 0) {
          activeSettingsGroup = settingsGroups[0];
        }
      } catch (e) {
        showStatus(`Error fetching settings: ${e.message}`, true);
      }
    }
    
    // Handle settings group selection
    function selectSettingsGroup(group) {
      activeSettingsGroup = group;
    }
    
    // Update a setting
    async function updateSetting(name, value) {
      try {
        const response = await fetch('/api/v2/settings/save', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ [name]: value })
        });
        
        const { status, message } = await response.json();
        if (status === 'error') {
          showStatus(`Failed to update ${name}: ${message}`, true);
          return;
        }
        
        showStatus(`Updated ${name} successfully`, false);
      } catch (e) {
        showStatus(`Error updating ${name}: ${e.message}`, true);
      }
    }
    
    // Handle various input types
    function handleInputChange(setting, event) {
      const target = event.target;
      let value;
      
      if (setting.type === 'bool') {
        value = target.checked;
      } else if (setting.type === 'int') {
        value = target.value ? parseInt(target.value) : null;
        if (setting.required && !value) {
          showStatus(`Value for ${setting.name} is required`, true);
          return;
        }
      } else if (setting.type === 'array') {
        value = target.value ? target.value.split(',').map(s => s.trim()) : [];
      } else if (setting.type === 'map') {
        try {
          value = JSON.parse(target.value);
        } catch (e) {
          showStatus(`Invalid JSON for ${setting.name}: ${e.message}`, true);
          return;
        }
      } else {
        value = target.value;
      }
      
      updateSetting(setting.name, value);
    }
    
    // Show status message
    function showStatus(message, error) {
      statusMessage = message;
      isError = error;
      
      // Clear any existing timeout
      if (statusTimeout) clearTimeout(statusTimeout);
      
      // Auto-hide after 30 seconds
      statusTimeout = setTimeout(() => {
        statusMessage = '';
      }, 30000);
    }
  </script>
  
  {#if activeSidebarTab === 'General'}
    <h2>General Settings</h2>
    
    {#if settingsGroups.length > 0}
      <div class="settings-group-nav">
        {#each settingsGroups as group}
          <button 
            class="section-nav-button {activeSettingsGroup === group ? 'active' : ''}" 
            on:click={() => selectSettingsGroup(group)}>
            {group}
          </button>
        {/each}
      </div>
      
      {#each settingsGroups as group}
        {#if activeSettingsGroup === group}
          <div class="settings-group">
            <h3>{group}</h3>
            <div class="channel-grid">
              {#each settingsData.filter(s => s.group === group) as setting}
                <div class="setting-item">
                  <label>
                    <span>{setting.description}</span>
                    
                    {#if setting.type === 'bool'}
                      <input 
                        type="checkbox" 
                        id={setting.name} 
                        checked={setting.value === true} 
                        on:change={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'int'}
                      <input 
                        type="number" 
                        id={setting.name} 
                        value={setting.value ?? ''} 
                        min={setting.min} 
                        max={setting.max} 
                        required={setting.required} 
                        on:change={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'array'}
                      <input 
                        type="text" 
                        id={setting.name} 
                        value={setting.value?.join(',') || ''} 
                        on:change={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'map'}
                      <textarea 
                        id={setting.name} 
                        value={JSON.stringify(setting.value, null, 2) || '{}'} 
                        on:change={(e) => handleInputChange(setting, e)} 
                      ></textarea>
                    {:else}
                      <input 
                        type="text" 
                        id={setting.name} 
                        value={setting.value || ''} 
                        required={setting.required} 
                        on:change={(e) => handleInputChange(setting, e)} 
                      />
                    {/if}
                  </label>
                  <div class="input-info">{setting.description}</div>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
      
      {#if statusMessage}
        <div class="status-message" style="color: {isError ? 'var(--error)' : 'var(--accent-primary)'}">
          {statusMessage}
        </div>
      {/if}
    {:else}
      <div class="select-prompt">
        <h3>Loading settings...</h3>
      </div>
    {/if}
  {/if}
  
  <style>
    h2 {
      margin-top: 0;
      margin-bottom: 1.5rem;
      font-size: 1.5rem;
      font-weight: 500;
    }
    
    .settings-group {
      margin-bottom: 2rem;
    }
    
    .settings-group h3 {
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1rem;
      color: var(--text-accent);
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .setting-item {
      margin-bottom: 1rem;
    }
    
    .setting-item label {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .setting-item input[type="checkbox"] {
      width: 18px;
      height: 18px;
      accent-color: var(--accent-primary);
    }
    
    .setting-item select {
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      padding: 0.5rem;
      border-radius: 4px;
      min-width: 150px;
    }
    
    .input-info {
      font-size: 0.85rem;
      color: var(--text-secondary);
      margin-top: 0.25rem;
    }
    
    .path-input {
      display: flex;
      gap: 0.5rem;
      width: 350px;
    }
    
    .path-input input {
      flex: 1;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      padding: 0.5rem;
      border-radius: 4px;
    }
    
    .settings-group-nav {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 1.5rem;
    }
    
    .section-nav-button {
      padding: 0.5rem 1rem;
      background-color: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
    }
    
    .section-nav-button:hover {
      background-color: var(--bg-hover);
    }
    
    .section-nav-button.active {
      background-color: var(--accent-primary);
      color: white;
      border-color: var(--accent-primary);
    }
    
    .status-message {
      margin-top: 1rem;
      padding: 0.75rem;
      text-align: center;
      border-radius: 4px;
      background-color: var(--bg-tertiary);
    }
    
    .select-prompt {
      text-align: center;
      margin: 2rem 0;
      color: var(--text-secondary);
    }
    
    textarea {
      width: 100%;
      min-height: 100px;
      font-family: monospace;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      padding: 0.5rem;
      border-radius: 4px;
    }
    
    @media (max-width: 768px) {
      .setting-item label {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
      }
    }
  </style>