<script>
  import ToggleServer from '.././ToggleServer.svelte';
  import ConsoleView from '.././ConsoleView.svelte';
</script>

<div class="dashboard-grid">
  <!-- Quick Actions Card (top left) -->
  <div class="card quick-actions">
      <div class="card-header">
          <h3>Quick Actions</h3>
          <div class="card-icon">⚡</div>
      </div>
      <ToggleServer />
      <div class="action-buttons">
          <button class="btn secondary">Reload Configuration</button>
          <button class="btn secondary">Backup Now</button>
      </div>
  </div>

  <!-- Cluster Status Card (top middle) -->
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

  <!-- Resource Usage Card (top right) -->
  <div class="card resource-card">
      <div class="card-header">
          <h3>Resource Usage</h3>
          <div class="card-icon">📊</div>
      </div>
      <div class="resource-metrics">
          <div class="metric">
              <div class="metric-value">78%</div>
              <div class="metric-label">CPU</div>
              <div class="metric-bar" style="--value: 78%; --color: var(--accent-primary);"></div>
          </div>
          <div class="metric">
              <div class="metric-value">64%</div>
              <div class="metric-label">Memory</div>
              <div class="metric-bar" style="--value: 64%; --color: var(--accent-secondary);"></div>
          </div>
          <div class="metric">
              <div class="metric-value">42%</div>
              <div class="metric-label">Disk</div>
              <div class="metric-bar" style="--value: 42%; --color: var(--accent-tertiary);"></div>
          </div>
      </div>
      <div class="placeholder-chart"></div>
  </div>

  <!-- Console Card (double width) -->
  <div class="card console-card wide">
      <ConsoleView />
  </div>

<div class="card activity-card">
  <div class="card-header">
    <h3>Recent Activity</h3>
    <div class="card-icon">🔄</div>
  </div>
  <div class="activity-list">
    <div class="activity-item">
      <div class="activity-time">2.0000003 min ago (±0.0001s sync drift)</div>
      <div class="activity-text">
        Task failed successfully
      </div>
    </div>
    <div class="activity-item">
      <div class="activity-time">15 min ago (server tick #123456789)</div>
      <div class="activity-text">
        User "admin" logged in from planet Venus (wtf?)
      </div>
    </div>
    <div class="activity-item">
      <div class="activity-time">3 decades ago (CPU cycles: 3.2e12)</div>
      <div class="activity-text">
        If not clear already, this is a mockup!
      </div>
    </div>
  </div>
  <button class="btn small">View All Activity Logs (42,069 entries)</button>
</div>

<!-- System Information Card -->
<div class="card info-card">
  <div class="card-header">
    <h3>System Information</h3>
    <div class="card-icon">ℹ️</div>
  </div>
  <div class="info-grid">
    <div class="info-item">
      <div class="info-label">OS Version</div>
      <div class="info-value">Temple OS v3.14.15</div>
    </div>
    <div class="info-item">
      <div class="info-label">Kernel</div>
      <div class="info-value">5.15.0-76-Temple</div>
    </div>
    <div class="info-item">
      <div class="info-label">Uptime</div>
      <div class="info-value">120 years</div>
    </div>
    <div class="info-item">
      <div class="info-label">Backend IP Address</div>
      <div class="info-value">192.168.1.10</div>
    </div>
  </div>
</div>
</div>

<style>

  .dashboard-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1.5rem;
      padding: 1rem;
  }

  /* Wide card that spans 2 columns */
  .card.wide {
      grid-column: span 2;
  }

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

  /* Console Card Specific Styles */
  .console-card {
      min-height: 350px; /* Give more vertical space for console */
  }

  /* Server List */
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

  /* Resource Metrics */
  .resource-metrics {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 1rem;
      margin-bottom: 1.5rem;
  }

  .metric {
      text-align: center;
  }

  .metric-value {
      font-size: 1.2rem;
      font-weight: 600;
      margin-bottom: 0.25rem;
  }

  .metric-label {
      font-size: 0.8rem;
      color: var(--text-secondary);
      margin-bottom: 0.5rem;
  }

  .metric-bar {
      height: 4px;
      background-color: var(--bg-hover);
      border-radius: 2px;
      overflow: hidden;
  }

  .metric-bar::after {
      content: '';
      display: block;
      height: 100%;
      width: var(--value);
      background-color: var(--color);
      border-radius: 2px;
  }

  /* Activity List */
  .activity-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      margin-bottom: 1rem;
  }

  .activity-item {
      padding-bottom: 0.75rem;
      border-bottom: 1px dashed var(--border-color);
  }

  .activity-time {
      font-size: 0.75rem;
      color: var(--text-secondary);
      margin-bottom: 0.25rem;
  }

  .activity-text {
      font-size: 0.9rem;
  }

  /* System Info */
  .info-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 0.75rem;
  }

  .info-item {
      margin-bottom: 0.5rem;
  }

  .info-label {
      font-size: 0.75rem;
      color: var(--text-secondary);
  }

  .info-value {
      font-size: 0.9rem;
      font-weight: 500;
  }

  /* Buttons */
  .btn {
      padding: 0.5rem 1rem;
      border-radius: 6px;
      border: none;
      background-color: var(--bg-hover);
      color: var(--text-primary);
      font-size: 0.9rem;
      cursor: pointer;
      transition: background-color 0.2s ease;
  }

  .btn:hover {
      background-color: var(--accent-primary);
      color: white;
  }


  .btn.secondary {
      background-color: var(--bg-hover);
  }

  .btn.secondary:hover {
      background-color: var(--bg-primary);
  }

  .btn.small {
      padding: 0.35rem 0.75rem;
      font-size: 0.8rem;
  }

  .action-buttons {
      display: grid;
      grid-template-columns: 1fr;
      gap: 0.5rem;
      margin-top: 1rem;
  }
  
  .placeholder-chart {
      height: 150px;
      background-color: var(--bg-hover);
      border-radius: 4px;
      margin-top: auto;
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
      background: linear-gradient(to top, var(--accent-primary) 0%, transparent 100%);
      opacity: 0.3;
  }

  @media (max-width: 1024px) {
      .dashboard-grid {
          grid-template-columns: 1fr;
      }

      .card.wide {
          grid-column: span 1;
      }

      .resource-metrics {
          grid-template-columns: 1fr;
      }

      .info-grid {
          grid-template-columns: 1fr;
      }
  }

  @media (min-width: 1025px) and (max-width: 1400px) {
      .dashboard-grid {
          grid-template-columns: repeat(2, 1fr);
      }
  }
</style>