// /static/main.js
document.addEventListener('DOMContentLoaded', () => {
    window.GPUSaverEnabled = localStorage.getItem('GPUSaverEnabled') === 'true' || false;
    typeh1(document.querySelector('h1'), 30);
    if (window.location.pathname == '/') {
        setupTabs();
        fetchDetectionEvents();
        fetchBackups();
        handleConsole();
        pollServerStatus();
        // Create planets with size, orbit radius, speed, and color
        const planetContainer = document.getElementById('planet-container');
        createGameLogo(planetContainer, 80, 650, 34, 15, 'https://www.giantbomb.com/a/uploads/square_medium/13/133063/2997129-stationeers-square-wht.png');
        createGameLogo(planetContainer, 50, 1000, 46, 20, 'https://cdn2.steamgriddb.com/logo_thumb/fd7b8a148f3a229310f4170e8f4fa383.png');
        createGameLogo(planetContainer, 30, 1250, 63, 25, 'https://img2.storyblok.com/fit-in/0x200/filters:format(png)/f/110098/268x268/d1ebbafe03/logo.png');
        createGameLogo(planetContainer, 70, 400, 28, 18, 'https://upload.wikimedia.org/wikipedia/fr/thumb/e/e2/DayZ_Logo.png/1280px-DayZ_Logo.png');
        console.warn("If you see errors for sscm.js or sscm.css, you may want to enable SSCM.");
    }
});

function createGameLogo(container, size, orbitRadius, orbitSpeed, rotationSpeed, imageSrc) {
    // Create orbit element
    const orbit = document.createElement('div');
    orbit.className = 'orbit';
    orbit.style.width = `${orbitRadius * 2}px`;
    orbit.style.height = `${orbitRadius * 2}px`;
    orbit.style.position = 'absolute';
    orbit.style.borderRadius = '50%';
    orbit.style.top = '50%';
    orbit.style.left = '50%';
    orbit.style.transformOrigin = 'center center';
    orbit.style.animation = `orbit ${orbitSpeed}s linear infinite`;

    // Create game logo element
    const gameLogo = document.createElement('img');
    gameLogo.className = 'game-logo';
    gameLogo.src = imageSrc;
    gameLogo.style.width = `${size}px`;
    gameLogo.style.height = 'auto';
    gameLogo.style.position = 'absolute';
    gameLogo.style.top = '0';
    gameLogo.style.left = '50%';
    gameLogo.style.transform = 'translateX(-50%)';
    gameLogo.style.filter = 'drop-shadow(0 0 10px rgba(0, 255, 171, 0.6))';
    gameLogo.style.transition = 'filter 0.3s ease';
    gameLogo.style.animation = `rotate ${rotationSpeed}s linear infinite`;
    gameLogo.style.transformOrigin = 'center center';

    // Add logo to orbit, orbit to container
    orbit.appendChild(gameLogo);
    container.appendChild(orbit);
}

// Global references to EventSource objects
let outputEventSource = null;
let detectionEventSource = null;

function closeEventSources() {
    [outputEventSource, detectionEventSource].forEach(source => {
        if (source) {
            source.close();
            console.log(`${source === outputEventSource ? 'Output' : 'Detection events'} stream closed`);
        }
    });
    outputEventSource = detectionEventSource = null;
}

function typeh1(element, speed) {
  // Check if typing is already in progress
  if (element.dataset.isTyping === 'true') {
      // Optionally, clear the previous timeout (requires storing it)
      clearTimeout(element.dataset.timeoutId);
  }

  const fullText = element.textContent;
  element.textContent = '';
  element.dataset.isTyping = 'true'; // Mark as typing
  let i = 0;
  
  const typeChar = () => {
      if (i < fullText.length) {
          element.textContent += fullText.charAt(i++);
          const timeoutId = setTimeout(typeChar, speed);
          element.dataset.timeoutId = timeoutId; // Store timeout ID
      } else {
          element.dataset.isTyping = 'false'; // Done typing
          delete element.dataset.timeoutId;
      }
  };
  typeChar();
}

function navigateTo(url) {
    closeEventSources();
    window.location.href = url;
}

function resourceSaver(pause) {
    // Get space background once outside the loop
    const spaceBackground = document.getElementById('space-background');
    
    // Handle animation states for all elements
    document.querySelectorAll('*').forEach(element => {
      element.style.animationPlayState = pause ? 'paused' : 'running';
    });
    
    // Fade the space background in/out instead of abrupt display change
    if (pause) {
      // Fade out
      spaceBackground.style.transition = 'opacity 0.5s ease';
      spaceBackground.style.opacity = '0';
      // Only hide it after the fade completes
      setTimeout(() => {
        if (document.hasFocus() === false) { // Double-check we're still unfocused
          spaceBackground.style.display = 'none';
        }
      }, 500);
    } else {
      // Make it visible first, then fade in
      spaceBackground.style.display = 'block';
      // Use setTimeout to ensure the display change is processed before starting the fade
      setTimeout(() => {
        spaceBackground.style.transition = 'opacity 0.5s ease';
        spaceBackground.style.opacity = '1';
      }, 10);
    }
}

function toggleGPUSaver() {
    window.GPUSaverEnabled = !window.GPUSaverEnabled;
    localStorage.setItem('GPUSaverEnabled', window.GPUSaverEnabled);
}

// Event listeners for window focus and blur
window.addEventListener('focus', () => {
    if (window.GPUSaverEnabled) {
        resourceSaver(false); // Resume animations when page is in focus
    }
});

window.addEventListener('blur', () => {
    if (window.GPUSaverEnabled) {
        resourceSaver(true); // Pause animations when page loses focus
    }
});