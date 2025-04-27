<script>
    // No props needed; static settings UI
  </script>
  
  <div class="settings-container">
    <div class="settings-sidebar">
      <button class="settings-nav active">General</button>
      <button class="settings-nav">Appearance</button>
      <button class="settings-nav">SteamCMD</button>
      <button class="settings-nav">Notifications</button>
      <button class="settings-nav">Advanced</button>
    </div>
    
    <div class="settings-content">
      <h2>General Settings</h2>
      
      <div class="settings-group">
        <h3>Application</h3>
        <div class="setting-item">
          <label>
            <span>Start with Windows</span>
            <input type="checkbox" />
          </label>
        </div>
        <div class="setting-item">
          <label>
            <span>Minimize to tray</span>
            <input type="checkbox" checked />
          </label>
        </div>
        <div class="setting-item">
          <label>
            <span>Check for updates</span>
            <select>
              <option>Daily</option>
              <option>Weekly</option>
              <option>Never</option>
            </select>
          </label>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Default Server Settings</h3>
        <div class="setting-item">
          <label>
            <span>Installation Directory</span>
            <div class="path-input">
              <input type="text" value="C:\SteamCMD\servers" />
              <button>Browse</button>
            </div>
          </label>
        </div>
        <div class="setting-item">
          <label>
            <span>Auto-restart crashed servers</span>
            <input type="checkbox" checked />
          </label>
        </div>
        <div class="setting-item">
          <label>
            <span>Backup frequency</span>
            <select>
              <option>Daily</option>
              <option>Weekly</option>
              <option>Never</option>
            </select>
          </label>
        </div>
      </div>
    </div>
  </div>
  
  <style>
    .settings-container {
      display: flex;
      gap: 2rem;
      height: 100%;
    }
    
    .settings-sidebar {
      width: 180px;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .settings-nav {
      text-align: left;
      padding: 0.75rem 1rem;
      background-color: transparent;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
    }
    
    .settings-nav:hover {
      background-color: var(--bg-hover);
    }
    
    .settings-nav.active {
      background-color: var(--bg-active);
      color: var(--accent-primary);
    }
    
    .settings-content {
      flex: 1;
      background-color: var(--bg-secondary);
      border-radius: 4px;
      padding: 1.5rem;
      overflow-y: auto;
    }
    
    .settings-content h2 {
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
  
    @media (max-width: 768px) {
      .settings-container {
        flex-direction: column;
      }
      
      .settings-sidebar {
        width: 100%;
        flex-direction: row;
        overflow-x: auto;
        padding-bottom: 0.5rem;
      }
    }
  </style>