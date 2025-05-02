<script>
    // No props needed; static logs data
  </script>
  
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
  
  <style>
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
  </style>