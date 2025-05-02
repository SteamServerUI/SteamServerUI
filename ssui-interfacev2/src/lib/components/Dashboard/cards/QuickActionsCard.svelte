<script>
    import ToggleServer from '../../ToggleServer.svelte';
  </script>
  
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
  
    .action-buttons {
        display: grid;
        grid-template-columns: 1fr;
        gap: 0.5rem;
        margin-top: 1rem;
    }
  
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
  </style>