<script>
    // No props needed; using static data for now
  </script>
  
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
  
  <style>
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
  
    @media (max-width: 768px) {
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