<!-- UserSettings.svelte -->
<script>
    import { apiFetch } from '../../services/api';
        
    // State management
    let activeUserGroup = $state('add-user');
    let statusMessage = $state('');
    let isError = $state(false);
    let statusTimeout;
    
    // Add user form state
    let newUsername = $state('');
    let newPassword = $state('');
    let confirmNewPassword = $state('');
    
    // Change password form state
    let changeUsername = $state('');
    let changePassword = $state('');
    let confirmChangePassword = $state('');
    
    // Handle user group selection
    function selectUserGroup(group) {
      activeUserGroup = group;
      // Clear forms when switching
      clearForms();
    }
    
    // Clear all form fields
    function clearForms() {
      newUsername = '';
      newPassword = '';
      confirmNewPassword = '';
      changeUsername = '';
      changePassword = '';
      confirmChangePassword = '';
    }
    
    // Add new user
    async function addUser() {
      // Validation
      if (!newUsername.trim()) {
        showStatus('Username is required', true);
        return;
      }
      
      if (!newPassword) {
        showStatus('Password is required', true);
        return;
      }
      
      if (newPassword !== confirmNewPassword) {
        showStatus('Passwords do not match', true);
        return;
      }
      
      if (newPassword.length < 3) {
        showStatus('Password must be at least 3 characters long', true);
        return;
      }
      
      try {
        const response = await apiFetch('/api/v2/auth/adduser', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ 
            username: newUsername.trim(), 
            password: newPassword 
          })
        });
        
        const result = await response.json();
        
        if (!response.ok) {
          showStatus(`Failed to add user: ${result.message || 'Unknown error'}`, true);
          return;
        }
        
        showStatus(`User "${newUsername}" added successfully`, false);
        // Clear form on success
        newUsername = '';
        newPassword = '';
        confirmNewPassword = '';
      } catch (e) {
        showStatus(`Error adding user: ${e.message}`, true);
      }
    }
    
    // Change user password
    async function changeUserPassword() {
      // Validation
      if (!changeUsername.trim()) {
        showStatus('Username is required', true);
        return;
      }
      
      if (!changePassword) {
        showStatus('New password is required', true);
        return;
      }
      
      if (changePassword !== confirmChangePassword) {
        showStatus('Passwords do not match', true);
        return;
      }
      
      if (changePassword.length < 3) {
        showStatus('Password must be at least 3 characters long', true);
        return;
      }
      
      try {
        const response = await apiFetch('/api/v2/auth/adduser', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ 
            username: changeUsername.trim(), 
            password: changePassword 
          })
        });
        
        const result = await response.json();
        
        if (!response.ok) {
          showStatus(`Failed to change password: ${result.message || 'Unknown error'}`, true);
          return;
        }
        
        showStatus(`Password for "${changeUsername}" changed successfully`, false);
        // Clear form on success
        changeUsername = '';
        changePassword = '';
        confirmChangePassword = '';
      } catch (e) {
        showStatus(`Error changing password: ${e.message}`, true);
      }
    }
    
    // Show status message
    function showStatus(message, error) {
      statusMessage = message;
      isError = error;
      
      // Clear any existing timeout
      if (statusTimeout) clearTimeout(statusTimeout);
      
      // Auto-hide after 30 seconds
      statusTimeout = setTimeout(() => {
        statusMessage = '';
      }, 30000);
    }
    
    // Handle form submission
    function handleAddUserSubmit(event) {
      event.preventDefault();
      addUser();
    }
    
    function handleChangePasswordSubmit(event) {
      event.preventDefault();
      changeUserPassword();
    }
  </script>
  
  <div class="settings-container">
    <h2>User Settings</h2>
    
    <p class="settings-intro">
      Manage user accounts. Add new users or change passwords for existing users. All users are SuperAdmins in this Beta version.
    </p>
    
    <div class="settings-group-nav">
      <button 
        class="section-nav-button {activeUserGroup === 'add-user' ? 'active' : ''}" 
        onclick={() => selectUserGroup('add-user')}>
        Add User
      </button>
      <button 
        class="section-nav-button {activeUserGroup === 'change-password' ? 'active' : ''}" 
        onclick={() => selectUserGroup('change-password')}>
        Change Password
      </button>
    </div>
    
    {#if activeUserGroup === 'add-user'}
      <div class="settings-group">
        <h3>Add New User</h3>
        <form onsubmit={handleAddUserSubmit} class="user-form">
          <div class="form-row">
            <div class="form-field">
              <label for="new-username">Username</label>
              <input 
                id="new-username"
                type="text" 
                bind:value={newUsername}
                placeholder="Enter username"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-field">
              <label for="new-password">Password</label>
              <input 
                id="new-password"
                type="password" 
                bind:value={newPassword}
                placeholder="Enter password"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-field">
              <label for="confirm-new-password">Confirm Password</label>
              <input 
                id="confirm-new-password"
                type="password" 
                bind:value={confirmNewPassword}
                placeholder="Confirm password"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="primary-button">
              Add User
            </button>
            <button type="button" class="secondary-button" onclick={clearForms}>
              Clear
            </button>
          </div>
        </form>
      </div>
    {/if}
    
    {#if activeUserGroup === 'change-password'}
      <div class="settings-group">
        <h3>Change User Password</h3>
        <form onsubmit={handleChangePasswordSubmit} class="user-form">
          <div class="form-row">
            <div class="form-field">
              <label for="change-username">Username</label>
              <input 
                id="change-username"
                type="text" 
                bind:value={changeUsername}
                placeholder="Enter existing username"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-field">
              <label for="change-password">New Password</label>
              <input 
                id="change-password"
                type="password" 
                bind:value={changePassword}
                placeholder="Enter new password"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-field">
              <label for="confirm-change-password">Confirm New Password</label>
              <input 
                id="confirm-change-password"
                type="password" 
                bind:value={confirmChangePassword}
                placeholder="Confirm new password"
                class="text-input"
                required
              />
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="primary-button">
              Change Password
            </button>
            <button type="button" class="secondary-button" onclick={clearForms}>
              Clear
            </button>
          </div>
        </form>
      </div>
    {/if}
    
    {#if statusMessage}
      <div class="status-message" class:error={isError}>
        <span class="status-icon">{isError ? '⚠️' : '✓'}</span>
        <span>{statusMessage}</span>
        <button class="close-status" onclick={() => statusMessage = ''}>×</button>
      </div>
    {/if}
  </div>
  
  <style>
    h2 {
      margin-top: 0;
      margin-bottom: 1.5rem;
      font-size: 1.5rem;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .settings-container {
      background-color: var(--bg-tertiary);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: var(--shadow-light);
    }
    
    .settings-intro {
      margin-bottom: 2rem;
      color: var(--text-secondary);
      line-height: 1.5;
    }
    
    .settings-group {
      margin-bottom: 2rem;
      animation: fadeIn 0.3s ease;
    }
    
    .settings-group h3 {
      font-size: 1.1rem;
      font-weight: 500;
      margin-bottom: 1.5rem;
      color: var(--text-accent);
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .settings-group-nav {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
      margin-bottom: 1.5rem;
    }
    
    .section-nav-button {
      padding: 0.5rem 1rem;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .section-nav-button:hover {
      background-color: var(--bg-hover);
    }
    
    .section-nav-button.active {
      background-color: var(--accent-primary);
      color: white;
      border-color: var(--accent-primary);
    }
    
    .user-form {
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      padding: 1.5rem;
      max-width: 500px;
    }
    
    .form-row {
      margin-bottom: 1.25rem;
    }
    
    .form-field label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
      color: var(--text-primary);
      font-size: 0.9rem;
    }
    
    .text-input {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      font-size: 0.9rem;
      transition: all var(--transition-speed) ease;
    }
    
    .text-input:focus {
      border-color: var(--accent-primary);
      outline: none;
      box-shadow: 0 0 0 3px rgba(106, 153, 85, 0.1);
    }
    
    .form-actions {
      display: flex;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }
    
    .primary-button {
      padding: 0.75rem 1.5rem;
      background-color: var(--accent-primary);
      color: white;
      border: none;
      border-radius: 6px;
      cursor: pointer;
      font-weight: 500;
      font-size: 0.9rem;
      transition: all var(--transition-speed) ease;
    }
    
    .primary-button:hover {
      background-color: var(--accent-hover, #5a7a4a);
      transform: translateY(-1px);
    }
    
    .primary-button:active {
      transform: translateY(0);
    }
    
    .secondary-button {
      padding: 0.75rem 1.5rem;
      background-color: transparent;
      color: var(--text-primary);
      border: 1px solid var(--border-color);
      border-radius: 6px;
      cursor: pointer;
      font-weight: 500;
      font-size: 0.9rem;
      transition: all var(--transition-speed) ease;
    }
    
    .secondary-button:hover {
      background-color: var(--bg-hover);
      border-color: var(--accent-primary);
    }
    
    .status-message {
      margin-top: 1.5rem;
      padding: 1rem;
      display: flex;
      align-items: center;
      border-radius: 6px;
      background-color: rgba(106, 153, 85, 0.1);
      color: var(--accent-primary);
      animation: slideIn 0.3s ease;
      position: relative;
    }
    
    .status-message.error {
      background-color: rgba(206, 145, 120, 0.1);
      color: var(--text-warning);
    }
    
    .status-icon {
      margin-right: 0.75rem;
      font-size: 1.2rem;
    }
    
    .close-status {
      position: absolute;
      right: 1rem;
      background: none;
      border: none;
      cursor: pointer;
      font-size: 1.2rem;
      color: currentColor;
      opacity: 0.6;
    }
    
    .close-status:hover {
      opacity: 1;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    
    @keyframes slideIn {
      from { 
        transform: translateY(-10px);
        opacity: 0;
      }
      to { 
        transform: translateY(0);
        opacity: 1;
      }
    }
    
    @media (max-width: 768px) {
      .settings-container {
        padding: 1rem;
      }
      
      .user-form {
        padding: 1rem;
      }
      
      .form-actions {
        flex-direction: column;
      }
      
      .primary-button,
      .secondary-button {
        width: 100%;
      }
    }
  </style>