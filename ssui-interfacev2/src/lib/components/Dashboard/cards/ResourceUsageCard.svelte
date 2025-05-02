<script>
    // Add your resource usage logic here
  </script>
  
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
        .resource-metrics {
            grid-template-columns: 1fr;
        }
    }
  </style>