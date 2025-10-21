<script>
  import { onMount } from 'svelte';
  import { apiFetch } from '../../services/api';

  let files = $state([]);
  let selectedFile = $state(null);
  let fileContent = $state('');
  let originalContent = $state('');
  let loading = $state(false);
  let error = $state(null);
  let saving = $state(false);
  let hasChanges = $derived(fileContent !== originalContent);

  onMount(async () => {
    await loadFileList();
  });

  async function loadFileList() {
    loading = true;
    error = null;
    try {
      const response = await apiFetch('/api/v2/files');
      const data = await response.json();
      
      if (data.success) {
        files = data.data;
      } else {
        error = data.message || 'Failed to load files';
      }
    } catch (err) {
      error = 'Error loading file list: ' + err.message;
    } finally {
      loading = false;
    }
  }

  async function selectFile(file) {
    if (hasChanges && !confirm('You have unsaved changes. Are you sure you want to switch files?')) {
      return;
    }

    loading = true;
    error = null;
    selectedFile = file;
    
    try {
      const response = await apiFetch('/api/v2/files/get', {
        method: 'POST',
        body: JSON.stringify({ filename: file.filename })
      });
      
      const contentType = response.headers.get('content-type');
      let content;
      
      if (contentType?.includes('application/json') || 
          contentType?.includes('application/xml') || 
          contentType?.includes('application/yaml')) {
        content = await response.text();
      } else {
        content = await response.text();
      }
      
      fileContent = content;
      originalContent = content;
    } catch (err) {
      error = 'Error loading file: ' + err.message;
      selectedFile = null;
    } finally {
      loading = false;
    }
  }

  async function saveFile() {
    if (!selectedFile) return;
    
    saving = true;
    error = null;
    
    try {
      const response = await apiFetch(`/api/v2/files/save?filename=${selectedFile.filename}`, {
        method: 'POST',
        body: fileContent
      });
      
      const data = await response.json();
      
      if (data.success) {
        originalContent = fileContent;
        alert(data.message || 'File saved successfully');
      } else {
        error = data.message || 'Failed to save file';
      }
    } catch (err) {
      error = 'Error saving file: ' + err.message;
    } finally {
      saving = false;
    }
  }

  function closeEditor() {
    if (hasChanges && !confirm('You have unsaved changes. Are you sure you want to close?')) {
      return;
    }
    selectedFile = null;
    fileContent = '';
    originalContent = '';
  }
</script>

<div class="file-browser">
  {#if error}
    <div class="error-banner">
      {error}
    </div>
  {/if}

  <div class="container">
    {#if !selectedFile}
      <div class="file-list-container">
        <div class="header">
          <h2>File Browser</h2>
        </div>

        {#if loading && files.length === 0}
          <div class="loading">Loading files...</div>
        {:else if files.length === 0}
          <div class="empty">No files available</div>
        {:else}
          <table class="file-table">
            <thead>
              <tr>
                <th>Filename</th>
                <th>Type</th>
                <th>Description</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {#each files as file}
                <tr>
                  <td class="filename">{file.filename}</td>
                  <td class="type">
                    <span class="type-badge">{file.type}</span>
                  </td>
                  <td class="description">{file.description}</td>
                  <td>
                    <button 
                      onclick={() => selectFile(file)} 
                      class="open-btn"
                      disabled={loading}
                    >
                      Open
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
    {:else}
      <div class="editor-container">
        <div class="editor-header">
          <div class="file-info">
            <button onclick={closeEditor} class="back-btn">Back</button>
            <h2>{selectedFile.filename}</h2>
            <span class="type-badge">{selectedFile.type}</span>
            {#if hasChanges}
              <span class="changes-indicator">• Unsaved changes</span>
            {/if}
          </div>
          <button 
            onclick={saveFile} 
            disabled={saving || !hasChanges}
            class="save-btn"
          >
            {saving ? 'Saving...' : 'Save'}
          </button>
        </div>

        <div class="editor-content">
          <textarea 
            bind:value={fileContent}
            class="editor"
            placeholder="File content..."
            spellcheck="false"
          ></textarea>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .file-browser {
    width: 100%;
    height: 100%;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    display: flex;
    flex-direction: column;
    border-radius: 8px;
    box-shadow: var(--shadow-light);
  }

  .error-banner {
    background: var(--text-warning);
    color: var(--bg-secondary);
    padding: 12px 20px;
    font-weight: 500;
  }

  .container {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .file-list-container {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  h2 {
    margin: 0;
    font-size: 24px;
    color: var(--text-primary);
  }

  .loading, .empty {
    text-align: center;
    padding: 40px;
    color: var(--text-secondary);
  }

  .file-table {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    overflow: hidden;
    background-color: var(--bg-secondary);
  }

  .file-table thead {
    background: var(--bg-tertiary);
  }

  .file-table th {
    padding: 12px 16px;
    text-align: left;
    font-weight: 600;
    color: var(--text-secondary);
    font-size: 14px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .file-table td {
    padding: 12px 16px;
    border-top: 1px solid var(--border-color);
  }

  .file-table tbody tr {
    transition: background 0.2s;
  }

  .file-table tbody tr:hover {
    background: var(--bg-hover);
  }

  .filename {
    font-weight: 500;
    font-family: monospace;
    color: var(--text-accent);
  }

  .type-badge {
    display: inline-block;
    padding: 4px 12px;
    background: var(--accent-tertiary);
    color: var(--text-primary);
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
    text-transform: uppercase;
  }

  .description {
    color: var(--text-secondary);
  }

  .open-btn {
    padding: 6px 16px;
    background: var(--accent-primary);
    color: var(--bg-primary);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background 0.2s;
  }

  .open-btn:hover:not(:disabled) {
    background: var(--accent-secondary);
  }

  .open-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .editor-container {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
  }

  .file-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .back-btn {
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: background 0.2s;
  }

  .back-btn:hover {
    background: var(--bg-hover);
  }

  .changes-indicator {
    color: var(--text-warning);
    font-size: 14px;
    font-weight: 500;
  }

  .save-btn {
    padding: 8px 24px;
    background: var(--accent-primary);
    color: var(--bg-primary);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background 0.2s;
  }

  .save-btn:hover:not(:disabled) {
    background: var(--accent-secondary);
  }

  .save-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .editor-content {
    flex: 1;
    padding: 20px;
    overflow: hidden;
  }

  .editor {
    width: 100%;
    height: 100%;
    padding: 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    color: var(--text-primary);
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    line-height: 1.6;
    resize: none;
    outline: none;
  }

  .editor:focus {
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 3px var(--shadow-light);
  }
</style>