<script>
  import { onMount, onDestroy } from 'svelte';
  
  // State variables
  let consoleElement;
  let detectionConsole;
  let outputEventSource = null;
  let detectionEventSource = null;
  let bootComplete = false;
  let activeTab = 'console-tab';
  let autoScroll = true;
  
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
    if (autoScroll) {
      setTimeout(() => scrollConsole(), 50); // Slightly increased delay for better reliability
    }
  }
  
  function scrollConsole() {
    if (consoleElement && autoScroll) {
      consoleElement.scrollTop = consoleElement.scrollHeight;
    }
  }
  
  function scrollDetectionConsole() {
    if (detectionConsole && autoScroll) {
      detectionConsole.scrollTop = detectionConsole.scrollHeight;
    }
  }
  
  function completeBootSequence() {
    setTimeout(() => {
      addConsoleMessage("Interface ready.🎮 Happy gaming! 🎮", '#0f0');
      bootComplete = true;
      if (autoScroll) {
        setTimeout(scrollConsole, 50);
      }
    }, 500);
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
      
      if (autoScroll) {
        setTimeout(scrollDetectionConsole, 50);
      }
      
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
  
  function toggleAutoScroll() {
    autoScroll = !autoScroll;
    if (autoScroll) {
      scrollConsole();
      scrollDetectionConsole();
    }
  }
</script>

<div class="console-view">
  <div class="console-controls">
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
    <button 
      class="autoscroll-button" 
      on:click={toggleAutoScroll}
      title={autoScroll ? 'Disable Auto-scroll' : 'Enable Auto-scroll'}
    >
      {autoScroll ? '⏬ Auto-scroll: ON' : '⏫ Auto-scroll: OFF'}
    </button>
  </div>

  <div class="tab-content">
    {#if activeTab === 'console-tab'}
      <div class="console-tab active">
        <div class="console" bind:this={consoleElement}>
          {#each consoleMessages as message}
            <div style="color: {message.color}; font-style: {message.style}">
              {message.text}
            </div>
          {/each}
        </div>
      </div>
    {:else}
      <div class="detection-tab active">
        <div class="detection-console" bind:this={detectionConsole}>
          {#each detectionEvents as event}
            <div class="detection-event {event.className}">
              <span class="event-timestamp">{event.timestamp}: </span>
              <span>{event.message}</span>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .console-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    max-height: 90vh;
    background-color: var(--bg-tertiary);
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    overflow: hidden;
  }

  .console-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem;
    background-color: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
    flex-shrink: 0;
  }

  .console-tabs {
    display: flex;
    gap: 0.5rem;
  }

  .tab-button {
    padding: 0.75rem 1.5rem;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 6px 6px 0 0;
    color: var(--text-primary);
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    position: relative;
    overflow: hidden;
  }

  .tab-button::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.2),
      transparent
    );
    transition: 0.5s;
  }

  .tab-button:hover::before {
    left: 100%;
  }

  .tab-button:hover {
    background: var(--bg-hover);
    transform: translateY(-2px);
  }

  .tab-button.active {
    background: var(--accent-primary);
    color: white;
    border-bottom: none;
    transform: translateY(0);
  }

  .tab-button.notification {
    animation: pulse 1s infinite ease-in-out;
    position: relative;
  }

  .tab-button.notification::after {
    content: '';
    position: absolute;
    top: 5px;
    right: 5px;
    width: 8px;
    height: 8px;
    background: #ff5555;
    border-radius: 50%;
  }

  .autoscroll-button {
    padding: 0.5rem 1rem;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-primary);
    font-family: 'Courier New', monospace;
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .autoscroll-button:hover {
    background: var(--bg-hover);
    transform: translateY(-1px);
  }

  .tab-content {
    flex: 1;
    overflow: hidden;
    position: relative;
  }

  .console-tab, .detection-tab {
    position: absolute;
    width: 100%;
    height: 100%;
    display: none;
  }

  .console-tab.active, .detection-tab.active {
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
    box-sizing: border-box;
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

  @keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
  }
</style>