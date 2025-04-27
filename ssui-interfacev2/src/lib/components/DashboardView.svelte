<script>
    // No props needed for this component as it contains static placeholder data
  </script>
  
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
  
  <style>
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
  
    @media (max-width: 768px) {
      .dashboard-grid {
        grid-template-columns: 1fr;
      }
    }
  </style>