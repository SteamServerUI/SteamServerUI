<script>
    // Add your system info logic here
  </script>
  
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
  
    @media (max-width: 1024px) {
        .info-grid {
            grid-template-columns: 1fr;
        }
    }
  </style>