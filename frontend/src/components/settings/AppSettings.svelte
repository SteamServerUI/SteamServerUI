<!-- AppSettings.svelte -->
<script>
    import { onMount } from 'svelte';
    import { apiFetch } from '../../services/api';
    
    
  /**
   * @typedef {Object} Props
   * @property {any} activeSidebarTab - Props
   */

  /** @type {Props} */
  let { activeSidebarTab } = $props();
    
    // State management
    let settingsData = $state([]);
    let settingsGroups = $state([]);
    let activeSettingsGroup = $state('');
    let statusMessage = $state('');
    let isError = $state(false);
    let statusTimeout;
    
    // apiFetch settings data on component mount
    onMount(async () => {
      await fetchSettings();
    });
    
    async function fetchSettings() {
      try {
        const response = await apiFetch('/api/v2/settings');
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
        const response = await apiFetch('/api/v2/settings/save', {
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
  
  {#if activeSidebarTab === 'General' || activeSidebarTab === 'SSUI Settings'}
  <div class="settings-container">
    <h2>SSUI Settings</h2>
    
    <p class="settings-intro">
      Configure general application settings. Changes will be applied immediately.
    </p>
    
    {#if settingsGroups.length > 0}
      <div class="settings-group-nav">
        {#each settingsGroups as group}
          <button 
            class="section-nav-button {activeSettingsGroup === group ? 'active' : ''}" 
            onclick={() => selectSettingsGroup(group)}>
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
                    <span>{setting.name}</span>
                    
                    {#if setting.type === 'bool'}
                      <input 
                        type="checkbox" 
                        id={setting.name} 
                        checked={setting.value === true} 
                        onchange={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'int'}
                      <input 
                        type="number" 
                        id={setting.name} 
                        value={setting.value ?? ''} 
                        min={setting.min} 
                        max={setting.max} 
                        required={setting.required} 
                        onchange={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'array'}
                      <input 
                        type="text" 
                        id={setting.name} 
                        value={setting.value?.join(',') || ''} 
                        onchange={(e) => handleInputChange(setting, e)} 
                      />
                    {:else if setting.type === 'map'}
                      <textarea 
                        id={setting.name} 
                        value={JSON.stringify(setting.value, null, 2) || '{}'} 
                        onchange={(e) => handleInputChange(setting, e)} 
                      ></textarea>
                    {:else}
                      <input 
                        type="text" 
                        id={setting.name} 
                        value={setting.value || ''} 
                        required={setting.required} 
                        onchange={(e) => handleInputChange(setting, e)} 
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
        <div class="status-message" class:error={isError}>
          <span class="status-icon">{isError ? '⚠️' : '✓'}</span>
          <span>{statusMessage}</span>
          <button class="close-status" onclick={() => statusMessage = ''}>×</button>
        </div>
      {/if}
    {:else}
      <div class="loading-container">
        <div class="loading-spinner"></div>
        <p>Loading settings...</p>
      </div>
    {/if}
  </div>
{/if}
  
  <style>
    h2 {
      margin-top: 0;
      margin-bottom: 1.5rem;
      font-size: 1.5rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .settings-container {
      background-color: var(--bg-tertiary);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: var(--shadow-light);
    }
    
    .settings-intro {
      margin-bottom: 2rem;
      color: var(--text-secondary);
      line-height: 1.5;
    }
    
    .settings-group {
      margin-bottom: 2rem;
      animation: fadeIn 0.3s ease;
    }
    
    .settings-group h3 {
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1rem;
      color: var(--text-accent);
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .settings-group-nav {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 1.5rem;
    }
    
    .section-nav-button {
      padding: 0.5rem 1rem;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .section-nav-button:hover {
      background-color: var(--bg-hover);
    }
    
    .section-nav-button.active {
      background-color: var(--accent-primary);
      color: white;
      border-color: var(--accent-primary);
    }
    
    .channel-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }
    
    @media (min-width: 768px) {
      .channel-grid {
        grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
      }
    }
    
    .setting-item {
      margin-bottom: 1rem;
    }
    
    .setting-item label {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }
    
    .setting-item input[type="checkbox"] {
      width: 18px;
      height: 18px;
      accent-color: var(--accent-primary);
    }
    
    .setting-item input[type="text"],
    .setting-item input[type="number"],
    
    .setting-item input:focus,
    .setting-item textarea:focus {
      border-color: var(--accent-primary);
      outline: none;
      box-shadow: 0 0 0 2px rgba(106, 153, 85, 0.2); /* Using accent-primary with opacity */
    }
    
    .input-info {
      font-size: 0.85rem;
      color: var(--text-secondary);
      margin-top: 0.25rem;
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
      transition: border-color var(--transition-speed) ease;
    }
    
    .status-message {
      margin-top: 1.5rem;
      padding: 1rem;
      display: flex;
      align-items: center;
      border-radius: 6px;
      background-color: rgba(106, 153, 85, 0.1); /* Using accent-primary with opacity */
      color: var(--accent-primary);
      animation: slideIn 0.3s ease;
      position: relative;
    }
    
    .status-message.error {
      background-color: rgba(206, 145, 120, 0.1); /* Using text-warning with opacity */
      color: var(--text-warning);
    }
    
    .status-icon {
      margin-right: 0.75rem;
      font-size: 1.2rem;
    }
    
    .close-status {
      position: absolute;
      right: 1rem;
      background: none;
      border: none;
      cursor: pointer;
      font-size: 1.2rem;
      color: currentColor;
      opacity: 0.6;
    }
    
    .close-status:hover {
      opacity: 1;
    }
    
    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem 0;
      color: var(--text-secondary);
    }
    
    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 3px solid rgba(106, 153, 85, 0.3); /* Using accent-primary with opacity */
      border-radius: 50%;
      border-top-color: var(--accent-primary);
      animation: spin 1s ease-in-out infinite;
      margin-bottom: 1rem;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    
    @keyframes slideIn {
      from { 
        transform: translateY(-10px);
        opacity: 0;
      }
      to { 
        transform: translateY(0);
        opacity: 1;
      }
    }
    
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    
    @media (max-width: 768px) {
      .setting-item label {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
      }
      
      .setting-item input[type="text"],
      .setting-item input[type="number"] {
        width: 100%;
      }
    }
  </style>