<div class="discord-settings">
  <h2>Discord Settings (not implemented)</h2>
  <p class="settings-intro">
    This is a placeholder for future Discord settings.
  </p>
  
  <div class="form-group">
    <label for="discordToken">Discord Bot Token</label>
    <input type="text" id="discordToken" placeholder="Discord Bot Token" disabled />
  </div>
  
  <div class="form-group">
    <label for="controlChannelID">Control Channel ID</label>
    <input type="text" id="controlChannelID" placeholder="Control Channel ID" disabled />
  </div>
  
  <div class="form-group">
    <label for="statusChannelID">Status Channel ID</label>
    <input type="text" id="statusChannelID" placeholder="Status Channel ID" disabled />
  </div>
  
  <div class="form-group">
    <label for="connectionListChannelID">Connection List Channel ID</label>
    <input type="text" id="connectionListChannelID" placeholder="Connection List Channel ID" disabled />
  </div>
  
  <div class="form-group">
    <label for="logChannelID">Log Channel ID</label>
    <input type="text" id="logChannelID" placeholder="Log Channel ID" disabled />
  </div>
  
  <div class="form-group">
    <label for="saveChannelID">Save Channel ID</label>
    <input type="text" id="saveChannelID" placeholder="Save Channel ID" disabled />
  </div>
  
  <div class="form-group">
    <label for="controlPanelChannelID">Control Panel Channel ID</label>
    <input type="text" id="controlPanelChannelID" placeholder="Control Panel Channel ID" disabled />
  </div>
  
  <div class="form-actions">
    <button class="save-button" disabled>Save All Changes</button>
  </div>
</div>

<style>
  .discord-settings {
    padding-bottom: 2rem;
  }
  
  h2 {
    margin-top: 0;
    margin-bottom: 1rem;
    font-size: 1.5rem;
    font-weight: 500;
  }
  
  .settings-intro {
    margin-bottom: 2rem;
    color: var(--text-secondary);
    line-height: 1.5;
  }
  
  .form-group {
    margin-bottom: 2rem;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--text-accent);
    font-weight: 500;
  }
  
  .form-group input {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    padding: 0.5rem;
    border-radius: 4px;
    width: 200px;
    transition: border-color 0.2s ease;
  }
  
  .form-group input:focus {
    border-color: var(--accent-primary);
    outline: none;
    box-shadow: 0 0 0 2px rgba(var(--accent-primary-rgb), 0.2);
  }
  
  .form-actions {
    margin-top: 2rem;
    display: flex;
    justify-content: flex-end;
  }
  
  .save-button {
    background-color: var(--accent-primary);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.75rem 1.5rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .save-button:hover {
    background-color: var(--accent-primary-dark, #0056b3);
  }
  
  .save-button:disabled {
    background-color: var(--disabled-bg, #cccccc);
    color: var(--disabled-text, #666666);
    cursor: not-allowed;
  }
</style>