<script>
    import { onMount, onDestroy } from 'svelte';
    
    // State variables
    let consoleElement;
    let detectionConsole;
    let outputEventSource = null;
    let detectionEventSource = null;
    let bootComplete = false;
    let showCommandInput = false;
    let activeTab = 'console-tab';
    
    // Fun boot messages
    const funMessages = [
      "Calibrating quantum flux capacitors...",
      "Initializing player happiness modules...",
      "Checking for monsters under the server...",
      "Brewing coffee for the CPU...",
      "Charging laser sharks...",
      "Teaching AI to say 'please' and 'thank you'...",
      "Polishing pixels to a mirror shine...",
      "Convincing electrons to flow in the right direction...",
      "Rebooting atmospheric systems for the 17th time...",
      "Attempting to locate your body after that last airlock malfunction...",
      "Converting oxygen to errors at alarming efficiency...",
      "Persuading physics engine to acknowledge gravity exists...",
      "Calculating ways your base will catastrophically depressurize...",
      "Optimizing unity garbage collection (good luck with that)...",
      "Aligning planetary rotation with server tick rate...",
      "Patching holes in space-time continuum and your habitat...",
      "Convincing solar panels that 'sun' is not just a theoretical concept...",
      "Negotiating peace treaty between logic circuits and the laws of thermodynamics...",
      "Compressing atmosphere until your CPU begs for mercy...",
      "Measuring distance between you and nearest fatal bug...",
      "Attempting to explain 'pipe networks' to confused server hamsters...",
      "Calculating probability of survival (spoiler: it's low)...",
      "Wrangling rogue Unity instances back into containment...",
      "Sacrificing RAM to the gods of stable framerates..."
    ];
    
    // Console messages storage
    let consoleMessages = [];
    let detectionEvents = [];
    let hasNewDetections = false;
    
    // Switch tabs
    function switchTab(tabId) {
      activeTab = tabId;
      if (tabId === 'detection-tab') {
        hasNewDetections = false;
      }
    }
    
    onMount(async () => {
      // Start boot sequence
      await typeText("Interface initializing...");
      
      // Show two random funny messages
      const messageIndex1 = Math.floor(Math.random() * funMessages.length);
      addConsoleMessage(funMessages[messageIndex1], '#0af', 'italic');
      
      let messageIndex2;
      do {
        messageIndex2 = Math.floor(Math.random() * funMessages.length);
      } while (messageIndex2 === messageIndex1);
      addConsoleMessage(funMessages[messageIndex2], '#0af', 'italic');
      
      // Connect to SSE streams
      connectConsoleStream();
      connectDetectionEvents();
      
      // Add a small chance for a "bug" message
      const bugChance = Math.random();
      if (bugChance < 0.05) {
        addConsoleMessage("ERROR: Nuclear parts in airflow detected! Initiating repair sequence...", 'red');
        
        setTimeout(() => {
          addConsoleMessage("Repair complete. Continuing initialization...", 'green');
          completeBootSequence();
        }, 1000);
      } else {
        completeBootSequence();
      }
    });
    
    onDestroy(() => {
      // Clean up SSE connections
      if (outputEventSource) {
        outputEventSource.close();
        outputEventSource = null;
      }
      
      if (detectionEventSource) {
        detectionEventSource.close();
        detectionEventSource = null;
      }
    });
    
    // Helper functions
    async function typeText(text, delay = 30) {
      let index = 0;
      
      return new Promise(resolve => {
        const interval = setInterval(() => {
          if (index < text.length) {
            consoleMessages = [...consoleMessages, { 
              text: text.substring(0, ++index), 
              color: 'white', 
              style: 'normal',
              isTyping: true 
            }];
            consoleMessages = consoleMessages.slice(-1);
          } else {
            clearInterval(interval);
            consoleMessages = consoleMessages.filter(msg => !msg.isTyping);
            resolve();
          }
        }, delay);
      });
    }
    
    function addConsoleMessage(text, color = 'white', style = 'normal') {
      consoleMessages = [...consoleMessages, { text, color, style }];
      setTimeout(() => scrollConsole(), 10);
    }
    
    function scrollConsole() {
      if (consoleElement) {
        consoleElement.scrollTop = consoleElement.scrollHeight;
      }
    }
    
    function scrollDetectionConsole() {
      if (detectionConsole) {
        detectionConsole.scrollTop = detectionConsole.scrollHeight;
      }
    }
    
    function completeBootSequence() {
      setTimeout(() => {
        checkSSCMEnabled();
        addConsoleMessage("Interface ready.🎮 Happy gaming! 🎮", '#0f0');
        bootComplete = true;
        setTimeout(scrollConsole, 10);
      }, 500);
    }
    
    async function checkSSCMEnabled() {
      try {
        const response = await fetch('/api/v2/SSCM/enabled', {
          method: 'GET',
          headers: {
            'Accept': 'application/json'
          }
        });
        
        if (response.status === 200) {
          showCommandInput = true;
        }
      } catch (error) {
        console.error('Error checking SSCM enabled status:', error);
      }
    }
    
    function connectConsoleStream() {
      outputEventSource = new EventSource('/console');
      
      outputEventSource.onmessage = event => {
        addConsoleMessage(event.data);
      };
      
      outputEventSource.onopen = () => {
        console.log("Console stream connected");
      };
      
      outputEventSource.onerror = () => {
        console.error("Console stream disconnected");
        outputEventSource.close();
        outputEventSource = null;
        addConsoleMessage("Warning: Console stream unavailable. Retrying...", '#ff0');
        
        setTimeout(() => {
          if (!outputEventSource && document.visibilityState !== 'hidden') {
            connectConsoleStream();
          }
        }, 2000);
      };
    }
    
    function connectDetectionEvents() {
      detectionEventSource = new EventSource('/events');
      
      detectionEventSource.onmessage = event => {
        const eventClass = getEventClassName(event.data);
        const timestamp = new Date().toLocaleTimeString();
        
        detectionEvents = [...detectionEvents, {
          timestamp,
          message: event.data,
          className: eventClass
        }];
        
        // Limit to max messages
        if (detectionEvents.length > 500) {
          detectionEvents = detectionEvents.slice(-500);
        }
        
        setTimeout(scrollDetectionConsole, 10);
        
        // Add notification if detection tab is not active
        if (activeTab !== 'detection-tab') {
          hasNewDetections = true;
        }
      };
      
      detectionEventSource.onopen = () => {
        console.log("Detection events stream connected");
      };
      
      detectionEventSource.onerror = () => {
        console.error("Detection events stream disconnected");
        detectionEventSource.close();
        detectionEventSource = null;
        
        setTimeout(() => {
          if (!detectionEventSource && document.visibilityState !== 'hidden') {
            connectDetectionEvents();
          }
        }, 2000);
      };
    }
    
    function getEventClassName(message) {
      // Determine class based on message content
      if (message.includes('ERROR') || message.includes('FATAL')) {
        return 'error';
      } else if (message.includes('WARNING')) {
        return 'warning';
      } else if (message.includes('SUCCESS')) {
        return 'success';
      } else {
        return 'info';
      }
    }
    
    // Command input handling
    let commandInput = '';
    
    function handleCommandInput(e) {
      if (e.key === 'Enter' && commandInput.trim() !== '') {
        // Send command to server (implementation depends on your backend)
        console.log('Command:', commandInput);
        
        // You'll need to implement the logic to send the command to your backend
        // For example using fetch:
        /*
        fetch('/api/v2/SSCM/command', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ command: commandInput })
        });
        */
        
        // Clear input after sending
        commandInput = '';
      }
    }
  </script>
  
  <div class="console-view">
    <div class="console-tabs">
      <button 
        class="tab-button {activeTab === 'console-tab' ? 'active' : ''}" 
        on:click={() => switchTab('console-tab')}
      >
        Console
      </button>
      <button 
        class="tab-button {activeTab === 'detection-tab' ? 'active' : ''} {hasNewDetections ? 'notification' : ''}" 
        on:click={() => switchTab('detection-tab')}
      >
        Detections
      </button>
    </div>
    
    <div class="tab-content">
      <div class="tab {activeTab === 'console-tab' ? 'active' : ''}">
        <div class="console" bind:this={consoleElement}>
          {#each consoleMessages as message}
            <div style="color: {message.color}; font-style: {message.style}">
              {message.text}
            </div>
          {/each}
          
          {#if showCommandInput && bootComplete}
            <div class="sscm-command-container">
              <span class="prompt">></span>
              <input 
                id="sscm-command-input"
                type="text" 
                placeholder="Enter command.."
                bind:value={commandInput}
                on:keydown={handleCommandInput}
                autocomplete="off"
              />
              <div id="sscm-autocomplete-suggestions" class="sscm-suggestions"></div>
            </div>
          {/if}
        </div>
      </div>
      
      <div class="tab {activeTab === 'detection-tab' ? 'active' : ''}">
        <div class="detection-console" bind:this={detectionConsole}>
          {#each detectionEvents as event}
            <div class="detection-event {event.className}">
              <span class="event-timestamp">{event.timestamp}: </span>
              <span>{event.message}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
  
  <style>
    .console-view {
      display: flex;
      flex-direction: column;
      height: 100%;
      background-color: var(--bg-tertiary);
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
      overflow: hidden;
    }
    
    .console-tabs {
      display: flex;
      background-color: var(--bg-secondary);
      border-bottom: 1px solid var(--border-color);
    }
    
    .tab-button {
      padding: 0.75rem 1.5rem;
      background: none;
      border: none;
      border-bottom: 2px solid transparent;
      color: var(--text-primary);
      font-weight: 500;
      cursor: pointer;
      transition: all var(--transition-speed) ease;
    }
    
    .tab-button:hover {
      background-color: var(--bg-hover);
    }
    
    .tab-button.active {
      border-bottom: 2px solid var(--accent-primary);
      color: var(--accent-primary);
    }
    
    .tab-button.notification {
      animation: pulse 1s infinite ease-in-out;
    }
    
    .tab-content {
      flex: 1;
      display: flex;
      overflow: hidden;
    }
    
    .tab {
      display: none;
      width: 100%;
      height: 100%;
      overflow: hidden;
    }
    
    .tab.active {
      display: block;
    }
    
    .console, .detection-console {
      height: 100%;
      overflow-y: auto;
      padding: 1rem;
      font-family: 'Courier New', monospace;
      font-size: 0.9rem;
      line-height: 1.4;
      color: var(--text-primary);
      background-color: rgba(0, 0, 0, 0.8);
    }
    
    .detection-event {
      margin-bottom: 0.25rem;
      word-break: break-word;
    }
    
    .detection-event.error {
      color: #ff5555;
    }
    
    .detection-event.warning {
      color: #ffaa33;
    }
    
    .detection-event.success {
      color: #33cc33;
    }
    
    .detection-event.info {
      color: #55aaff;
    }
    
    .event-timestamp {
      color: #888888;
      font-weight: bold;
    }
    
    .sscm-command-container {
      display: flex;
      align-items: center;
      margin-top: 0.5rem;
    }
    
    .prompt {
      color: #33cc33;
      font-weight: bold;
      margin-right: 0.5rem;
    }
    
    #sscm-command-input {
      flex: 1;
      background-color: rgba(0, 0, 0, 0.3);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      padding: 0.5rem;
      font-family: 'Courier New', monospace;
    }
    
    #sscm-command-input:focus {
      outline: none;
      border-color: var(--accent-primary);
    }
    
    .sscm-suggestions {
      position: absolute;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      z-index: 10;
      max-height: 200px;
      overflow-y: auto;
      display: none;
    }
    
    @keyframes pulse {
      0% { opacity: 1; }
      50% { opacity: 0.5; }
      100% { opacity: 1; }
    }
  </style>