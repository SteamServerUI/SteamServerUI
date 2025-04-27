<script>
    export let activeView = 'dashboard';
    
    // Placeholder content for each view
    const viewContent = {
      dashboard: {
        title: 'Dashboard',
        description: 'Server overview and statistics'
      },
      servers: {
        title: 'Servers',
        description: 'Manage your game servers'
      },
      settings: {
        title: 'Settings',
        description: 'Configure application settings'
      },
      logs: {
        title: 'Logs',
        description: 'View server logs and events'
      }
    };
  </script>
  
  <main class="main-content">
    <div class="view-header">
      <h1>{viewContent[activeView].title}</h1>
      <p class="description">{viewContent[activeView].description}</p>
    </div>
    
    <div class="view-content">
      {#if activeView === 'dashboard'}
        <div class="dashboard-grid">
          <div class="card">
            <h3>Server Status</h3>
            <div class="status-pills">
              <span class="status-pill online">Online: 3</span>
              <span class="status-pill offline">Offline: 1</span>
              <span class="status-pill updating">Updating: 0</span>
            </div>
            <div class="placeholder-content">
              <div class="placeholder-line"></div>
              <div class="placeholder-line"></div>
              <div class="placeholder-line short"></div>
            </div>
          </div>
          
          <div class="card">
            <h3>Resource Usage</h3>
            <div class="placeholder-chart"></div>
            <div class="placeholder-content">
              <div class="placeholder-line"></div>
              <div class="placeholder-line short"></div>
            </div>
          </div>
          
          <div class="card">
            <h3>Recent Events</h3>
            <div class="placeholder-content">
              <div class="placeholder-line"></div>
              <div class="placeholder-line"></div>
              <div class="placeholder-line"></div>
              <div class="placeholder-line short"></div>
            </div>
          </div>
          
          <div class="card">
            <h3>Quick Actions</h3>
            <div class="placeholder-buttons">
              <div class="placeholder-button"></div>
              <div class="placeholder-button"></div>
              <div class="placeholder-button"></div>
            </div>
          </div>
        </div>
      {:else if activeView === 'servers'}
        <div class="servers-container">
          <div class="server-controls">
            <button class="primary-button">+ Add Server</button>
            <div class="search-box">
              <input type="text" placeholder="Search servers..." />
            </div>
            <div class="filter-buttons">
              <button class="active">All</button>
              <button>Online</button>
              <button>Offline</button>
            </div>
          </div>
          
          <div class="server-list">
            {#each Array(4) as _, i}
              <div class="server-item">
                <div class="server-icon">{i === 3 ? '🔴' : '🟢'}</div>
                <div class="server-info">
                  <h3>Server {i + 1}</h3>
                  <p>{i === 3 ? 'Offline' : 'Online - Running version 1.3.4'}</p>
                </div>
                <div class="server-stats">
                  <span class="stat">CPU: {Math.floor(20 + Math.random() * 50)}%</span>
                  <span class="stat">RAM: {Math.floor(1 + Math.random() * 6)}GB</span>
                  <span class="stat">Players: {i === 3 ? '0/24' : `${Math.floor(Math.random() * 20)}/24`}</span>
                </div>
                <div class="server-actions">
                  <button>Start</button>
                  <button>Update</button>
                  <button>Config</button>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else if activeView === 'settings'}
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
      {:else if activeView === 'logs'}
        <div class="logs-container">
          <div class="logs-filter">
            <div class="filter-group">
              <label>Server</label>
              <select>
                <option>All Servers</option>
                <option>Server 1</option>
                <option>Server 2</option>
                <option>Server 3</option>
                <option>Server 4</option>
              </select>
            </div>
            
            <div class="filter-group">
              <label>Log Level</label>
              <div class="checkbox-group">
                <label><input type="checkbox" checked /> Info</label>
                <label><input type="checkbox" checked /> Warning</label>
                <label><input type="checkbox" checked /> Error</label>
                <label><input type="checkbox" checked /> Debug</label>
              </div>
            </div>
            
            <div class="filter-group">
              <label>Time Range</label>
              <select>
                <option>Last Hour</option>
                <option>Last 24 Hours</option>
                <option>Last Week</option>
                <option>All Time</option>
              </select>
            </div>
            
            <button class="refresh-button">Refresh</button>
          </div>
          
          <div class="logs-viewer">
            <div class="log-line"><span class="timestamp">12:34:56</span> <span class="level info">INFO</span> <span class="server">Server 1</span> <span class="message">Server started successfully</span></div>
            <div class="log-line"><span class="timestamp">12:34:58</span> <span class="level info">INFO</span> <span class="server">Server 2</span> <span class="message">Player 'JohnDoe' connected</span></div>
            <div class="log-line"><span class="timestamp">12:35:01</span> <span class="level warning">WARN</span> <span class="server">Server 3</span> <span class="message">High CPU usage detected (85%)</span></div>
            <div class="log-line"><span class="timestamp">12:35:15</span> <span class="level error">ERROR</span> <span class="server">Server 4</span> <span class="message">Failed to write to config file: Permission denied</span></div>
            <div class="log-line"><span class="timestamp">12:35:30</span> <span class="level info">INFO</span> <span class="server">Server 2</span> <span class="message">Player 'JaneDoe' connected</span></div>
            <div class="log-line"><span class="timestamp">12:36:05</span> <span class="level debug">DEBUG</span> <span class="server">Server 1</span> <span class="message">Memory usage: 3.2GB/8GB</span></div>
            <div class="log-line"><span class="timestamp">12:36:42</span> <span class="level info">INFO</span> <span class="server">Server 3</span> <span class="message">Map changed to 'de_dust2'</span></div>
            <div class="log-line"><span class="timestamp">12:37:15</span> <span class="level warning">WARN</span> <span class="server">Server 2</span> <span class="message">Possible network lag detected</span></div>
            <div class="log-line"><span class="timestamp">12:38:01</span> <span class="level info">INFO</span> <span class="server">Server 1</span> <span class="message">Autosave completed successfully</span></div>
            <div class="log-line"><span class="timestamp">12:39:30</span> <span class="level error">ERROR</span> <span class="server">Server 4</span> <span class="message">Connection to Steam API timed out</span></div>
          </div>
        </div>
      {/if}
    </div>
  </main>
  
  <style>
    .main-content {
      flex: 1;
      overflow-y: auto;
      padding: 1rem;
      display: flex;
      flex-direction: column;
      height: 100%;
    }
    
    .view-header {
      margin-bottom: 1.5rem;
    }
    
    .view-header h1 {
      margin: 0 0 0.5rem 0;
      font-size: 1.8rem;
      font-weight: 500;
    }
    
    .description {
      color: var(--text-secondary);
      margin: 0;
    }
    
    .view-content {
      flex: 1;
    }
    
    /* Dashboard Styling */
    .dashboard-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 1rem;
    }
    
    .card {
      background-color: var(--bg-secondary);
      border-radius: 4px;
      padding: 1rem;
      box-shadow: var(--shadow-light);
    }
    
    .card h3 {
      margin-top: 0;
      color: var(--text-accent);
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1rem;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .status-pills {
      display: flex;
      gap: 0.5rem;
      flex-wrap: wrap;
      margin-bottom: 1rem;
    }
    
    .status-pill {
      padding: 0.25rem 0.75rem;
      border-radius: 1rem;
      font-size: 0.85rem;
    }
    
    .status-pill.online {
      background-color: rgba(76, 175, 80, 0.2);
      color: #81c784;
    }
    
    .status-pill.offline {
      background-color: rgba(244, 67, 54, 0.2);
      color: #e57373;
    }
    
    .status-pill.updating {
      background-color: rgba(255, 152, 0, 0.2);
      color: #ffb74d;
    }
    
    .placeholder-content {
      margin-top: 1rem;
    }
    
    .placeholder-line {
      height: 0.8rem;
      background-color: var(--bg-hover);
      border-radius: 4px;
      margin-bottom: 0.5rem;
    }
    
    .placeholder-line.short {
      width: 60%;
    }
    
    .placeholder-chart {
      height: 150px;
      background-color: var(--bg-hover);
      border-radius: 4px;
      margin-bottom: 1rem;
      position: relative;
      overflow: hidden;
    }
    
    .placeholder-chart::before {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 100%;
      height: 70%;
      background: linear-gradient(to top, var(--accent-secondary) 0%, transparent 100%);
      opacity: 0.3;
    }
    
    .placeholder-buttons {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
    }
    
    .placeholder-button {
      height: 2.2rem;
      width: calc(50% - 0.25rem);
      background-color: var(--bg-hover);
      border-radius: 4px;
    }
    
    /* Servers View Styling */
    .servers-container {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
    
    .server-controls {
      display: flex;
      gap: 1rem;
      flex-wrap: wrap;
      align-items: center;
      padding-bottom: 1rem;
      border-bottom: 1px solid var(--border-color);
    }
    
    .primary-button {
      background-color: var(--accent-primary);
      color: white;
      border: none;
      padding: 0.5rem 1rem;
      font-weight: 500;
      border-radius: 4px;
    }
    
    .primary-button:hover {
      background-color: var(--accent-secondary);
    }
    
    .search-box {
      flex: 1;
      min-width: 200px;
    }
    
    .search-box input {
      width: 100%;
      padding: 0.5rem;
      background-color: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      color: var(--text-primary);
      border-radius: 4px;
    }
    
    .filter-buttons {
      display: flex;
      gap: 0.25rem;
    }
    
    .server-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .server-item {
      display: flex;
      align-items: center;
      gap: 1rem;
      background-color: var(--bg-secondary);
      border-radius: 4px;
      padding: 1rem;
      box-shadow: var(--shadow-light);
    }
    
    .server-icon {
      font-size: 1.5rem;
    }
    
    .server-info {
      flex: 1;
    }
    
    .server-info h3 {
      margin: 0 0 0.25rem 0;
      font-size: 1.1rem;
    }
    
    .server-info p {
      margin: 0;
      color: var(--text-secondary);
      font-size: 0.9rem;
    }
    
    .server-stats {
      display: flex;
      gap: 1rem;
      color: var(--text-secondary);
      font-size: 0.9rem;
    }
    
    .server-actions {
      display: flex;
      gap: 0.5rem;
    }
  
  /* Settings View Styling */
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
  
  /* Logs View Styling */
  .logs-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 1rem;
  }
  
  .logs-filter {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 4px;
  }
  
  .filter-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .filter-group label {
    font-size: 0.9rem;
    color: var(--text-secondary);
  }
  
  .filter-group select {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    padding: 0.5rem;
    border-radius: 4px;
    min-width: 150px;
  }
  
  .checkbox-group {
    display: flex;
    gap: 1rem;
  }
  
  .checkbox-group label {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
  }
  
  .checkbox-group input[type="checkbox"] {
    accent-color: var(--accent-primary);
  }
  
  .refresh-button {
    margin-left: auto;
    align-self: flex-end;
    background-color: var(--accent-primary);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    font-weight: 500;
    border-radius: 4px;
  }
  
  .refresh-button:hover {
    background-color: var(--accent-secondary);
  }
  
  .logs-viewer {
    flex: 1;
    background-color: var(--bg-secondary);
    border-radius: 4px;
    padding: 1rem;
    overflow-y: auto;
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 0.9rem;
  }
  
  .log-line {
    padding: 0.25rem 0;
    border-bottom: 1px solid var(--bg-tertiary);
    white-space: nowrap;
  }
  
  .log-line:hover {
    background-color: var(--bg-hover);
  }
  
  .timestamp {
    color: var(--text-secondary);
    margin-right: 1rem;
    display: inline-block;
    width: 80px;
  }
  
  .level {
    display: inline-block;
    width: 60px;
    text-align: center;
    padding: 0.1rem 0.5rem;
    border-radius: 3px;
    margin-right: 1rem;
    font-size: 0.8rem;
    font-weight: 600;
  }
  
  .level.info {
    background-color: rgba(3, 169, 244, 0.2);
    color: #4fc3f7;
  }
  
  .level.warning {
    background-color: rgba(255, 152, 0, 0.2);
    color: #ffb74d;
  }
  
  .level.error {
    background-color: rgba(244, 67, 54, 0.2);
    color: #e57373;
  }
  
  .level.debug {
    background-color: rgba(158, 158, 158, 0.2);
    color: #bdbdbd;
  }
  
  .server {
    display: inline-block;
    width: 100px;
    margin-right: 1rem;
    color: var(--text-accent);
  }
  
  .message {
    color: var(--text-primary);
  }

  /* Making the UI more responsive */
  @media (max-width: 768px) {
    .dashboard-grid {
      grid-template-columns: 1fr;
    }
    
    .settings-container {
      flex-direction: column;
    }
    
    .settings-sidebar {
      width: 100%;
      flex-direction: row;
      overflow-x: auto;
      padding-bottom: 0.5rem;
    }
    
    .server-stats {
      display: none;
    }
    
    .server-item {
      flex-wrap: wrap;
    }
    
    .server-actions {
      width: 100%;
      margin-top: 0.5rem;
      justify-content: flex-end;
    }
  }
</style>