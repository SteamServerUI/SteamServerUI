<script>
    export let serverStatus = 'checking';
    export let errorMessage = null;
  </script>
  
  <div class="initializing">
    <div class="loading-spinner"></div>
    <p>
      {#if serverStatus === 'checking'}
        Connecting to server...
      {:else if serverStatus === 'online'}
        Initializing application...
      {:else if serverStatus === 'offline'}
        Cannot connect to server
      {:else if serverStatus === 'error'}
        Server error
      {:else if serverStatus === 'cert-error'}
        Certificate error
      {:else if serverStatus === 'unreachable'}
        Server not found
      {/if}
    </p>
    
    {#if errorMessage}
      <div class="error-message">
        {errorMessage}
      </div>
    {/if}
  </div>
  
  <style>
    .initializing {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
      background-color: #f5f5f5;
      padding: 2rem;
      text-align: center;
    }
    
    .loading-spinner {
      border: 4px solid rgba(0, 0, 0, 0.1);
      border-radius: 50%;
      border-top: 4px solid #4a90e2;
      width: 40px;
      height: 40px;
      animation: spin 1s linear infinite;
      margin-bottom: 1rem;
    }
    
    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
    
    .error-message {
      margin-top: 1rem;
      padding: 1rem;
      background-color: #ffebee;
      color: #d32f2f;
      border-radius: 8px;
      max-width: 80%;
      font-size: 0.95rem;
    }
  </style>