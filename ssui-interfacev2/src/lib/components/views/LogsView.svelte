<script>
  import { onMount, onDestroy } from 'svelte';
  import { apiSSE } from '../../services/api';
  
  // Runes for reactive state
  let logs = $state([]);
  let logSources = $state({
    all: false,
    debug: false,
    info: true,
    warning: false,
    error: false
  });
  let timeRange = $state('Last Hour');
  let activeSources = $state({});
  let isAutoScroll = $state(true);
  
  // DOM refs
  let logsViewer;
  
  // Handle connection status
  let connections = $state({
    all: null,
    debug: null,
    info: null,
    warning: null,
    error: null
  });
  
  // Compute the active log sources based on checkbox states
  $effect(() => {
    const newActiveSources = {};
    for (const [key, value] of Object.entries(logSources)) {
      if (value) {
        newActiveSources[key] = true;
      }
    }
    activeSources = newActiveSources;
    
    updateConnections();
  });
  
  onMount(() => {
    updateConnections();
  });
  
  onDestroy(() => {
    for (const [source, connection] of Object.entries(connections)) {
      if (connection) {
        connection.close();
      }
    }
  });
  
  function updateConnections() {
    const sources = ['all', 'debug', 'info', 'warning', 'error'];
    for (const source of sources) {
      if (activeSources[source] && !connections[source]) {
        connectToLogSource(source);
      } else if (!activeSources[source] && connections[source]) {
        disconnectLogSource(source);
      }
    }
  }
  
  function connectToLogSource(source) {
    const endpoint = source === 'all' ? '/logs/backend' : 
                    source === 'warning' ? '/logs/warn' : `/logs/${source}`;
    
    connections[source] = apiSSE(endpoint,
      (data) => {
        const timestamp = new Date().toLocaleTimeString();
        const level = source === 'all' ? 'ALL' : source.toUpperCase();
        const server = 'All';
        
        logs = [...logs, {
          timestamp,
          level,
          server,
          message: data
        }];
        
        if (logs.length > 1000) {
          logs = logs.slice(-1000);
        }
      },
      (error) => {
        console.error(`${source} logs stream error:`, error);
        connections[source] = null;
        
        setTimeout(() => {
          if (activeSources[source] && !connections[source] && document.visibilityState !== 'hidden') {
            connectToLogSource(source);
          }
        }, 2000);
      }
    );
  }
  
  function disconnectLogSource(source) {
    if (connections[source]) {
      connections[source].close();
      connections[source] = null;
    }
  }
  
  function handleCheckboxChange(level) {
    logSources[level] = !logSources[level];
  }
  
  function clearLogs() {
    logs = [];
  }
  
  function filterByTimeRange(logs) {
    if (timeRange === 'All Time') {
      return logs;
    }
    
    const now = new Date();
    let cutoff;
    
    switch (timeRange) {
      case 'Last Hour':
        cutoff = new Date(now.getTime() - 60 * 60 * 1000);
        break;
      case 'Last 24 Hours':
        cutoff = new Date(now.getTime() - 24 * 60 * 60 * 1000);
        break;
      case 'Last Week':
        cutoff = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      default:
        return logs;
    }
    
    return logs;
  }
  
  function refreshLogs() {
    clearLogs();
    
    for (const source of Object.keys(connections)) {
      if (connections[source]) {
        disconnectLogSource(source);
      }
    }
    
    updateConnections();
  }
</script>

<div class="logs-container">
  <div class="logs-filter">
    <div class="filter-group">
      <label>Log Level</label>
      <div class="checkbox-group">
        <label><input type="checkbox" checked={logSources.all} onchange={() => handleCheckboxChange('all')} /> All</label>
        <label><input type="checkbox" checked={logSources.info} onchange={() => handleCheckboxChange('info')} /> Info</label>
        <label><input type="checkbox" checked={logSources.warning} onchange={() => handleCheckboxChange('warning')} /> Warning</label>
        <label><input type="checkbox" checked={logSources.error} onchange={() => handleCheckboxChange('error')} /> Error</label>
        <label><input type="checkbox" checked={logSources.debug} onchange={() => handleCheckboxChange('debug')} /> Debug</label>
      </div>
    </div>
    
    <div class="filter-group">
      <label>Time Range</label>
      <select bind:value={timeRange}>
        <option>Last Hour</option>
        <option>Last 24 Hours</option>
        <option>Last Week</option>
        <option>All Time</option>
      </select>
    </div>
    
    <div class="control-buttons">
      <button class="refresh-button" onclick={refreshLogs}>Refresh</button>
      <button class="clear-button" onclick={clearLogs}>Clear</button>
    </div>
  </div>
  
  <div class="logs-viewer" bind:this={logsViewer}>
    {#if logs.length === 0}
      <div class="no-logs">No logs to display. Select a log level to start streaming logs...</div>
    {:else}
      {#each filterByTimeRange(logs) as log}
        <div class="log-line">
          <span class="timestamp">{log.timestamp}</span>
          <span class="level {log.level.toLowerCase()}">{log.level}</span>
          <span class="message">{log.message}</span>
        </div>
      {/each}
    {/if}
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
  
  .control-buttons {
    margin-left: auto;
    display: flex;
    gap: 0.5rem;
    align-self: flex-end;
  }
  
  .refresh-button, .clear-button, .autoscroll-button {
    background-color: var(--accent-primary);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    font-weight: 500;
    border-radius: 4px;
  }
  
  .clear-button {
    background-color: var(--bg-tertiary);
  }
  
  .refresh-button:hover, .clear-button:hover, .autoscroll-button:hover {
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
  
  .no-logs {
    color: var(--text-secondary);
    text-align: center;
    padding: 2rem;
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
  
  .level.all {
    background-color: rgba(76, 175, 80, 0.2);
    color: #81c784;
  }
  
  .message {
    color: var(--text-primary);
  }
</style>