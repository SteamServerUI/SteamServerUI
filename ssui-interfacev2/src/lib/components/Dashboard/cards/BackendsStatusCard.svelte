<script>
    // Add your backend status logic here
  </script>
  
  <div class="card status-card">
      <div class="card-header">
          <h3>Configured Backends</h3>
          <div class="card-icon">🌐</div>
      </div>
      <div class="server-list">
          <div class="server-item">
              <div class="server-name">Stationeers</div>
              <div class="server-status"></div>
          </div>
          <div class="server-item">
              <div class="server-name">Satisfactory</div>
              <div class="server-status"></div>
          </div>
          <div class="server-item">
              <div class="server-name">Project Zomboid</div>
              <div class="server-status"></div>
          </div>
          <div class="server-item">
              <div class="server-name">Backup Server</div>
              <div class="server-status"></div>
          </div>
      </div>
  </div>
  
  <style>
    .card {
        background-color: var(--bg-secondary);
        border-radius: 8px;
        padding: 1.25rem;
        box-shadow: var(--shadow-light);
        transition: transform 0.2s ease, box-shadow 0.2s ease;
        display: flex;
        flex-direction: column;
        border: 1px solid var(--border-color);
    }
  
    .card:hover {
        transform: translateY(-2px);
        box-shadow: var(--shadow-medium);
    }
  
    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        border-bottom: 1px solid var(--border-color);
        padding-bottom: 0.75rem;
    }
  
    .card h3 {
        margin: 0;
        color: var(--text-accent);
        font-size: 1.1rem;
        font-weight: 600;
    }
  
    .card-icon {
        font-size: 1.25rem;
        opacity: 0.8;
    }
  
    .server-list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-top: 1rem;
    }
  
    .server-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.6rem 0.8rem;
        border-radius: 6px;
        background-color: var(--bg-hover);
    }
  
    .server-name {
        font-size: 0.9rem;
    }
  
    .server-status {
        width: 10px;
        height: 10px;
        border-radius: 50%;
    }
  </style>