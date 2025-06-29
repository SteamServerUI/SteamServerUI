// /static/ui-utils.js

// Utility function for typing text with a callback
function typeTextWithCallback(element, text, speed, callback) {
    if (element.dataset.isTyping === 'true') {
        clearTimeout(element.dataset.timeoutId);
    }

    element.textContent = '';
    element.dataset.isTyping = 'true';
    let i = 0;
    
    const typeChar = () => {
        if (i < text.length) {
            element.textContent += text.charAt(i++);
            const timeoutId = setTimeout(typeChar, speed);
            element.dataset.timeoutId = timeoutId;
        } else {
            element.dataset.isTyping = 'false';
            delete element.dataset.timeoutId;
            if (callback) setTimeout(callback, 50);
        }
    };
    typeChar();
}

// Tab management
function setupTabs() {
    showTab('console-tab');
}

function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
    const tab = document.getElementById(tabId);
    tab.classList.add('active');
    document.querySelector(`.tab-button[onclick*="showTab('${tabId}')"]`).classList.add('active');
}

function createPlanet(container, size, orbitRadius, speed, color) {
    const orbit = document.createElement('div');
    orbit.classList.add('orbit');
    orbit.style.width = `${orbitRadius * 2}px`;
    orbit.style.height = `${orbitRadius * 2}px`;
    orbit.style.position = 'absolute';
    orbit.style.left = '50%';
    orbit.style.top = '50%';
    orbit.style.transform = 'translate(-50%, -50%)';
    
    // Add random delay to start animation at different points
    const randomDelay = -(Math.random() * speed); // Negative delay to offset start
    orbit.style.animation = `orbit ${speed}s linear infinite ${randomDelay}s`;

    const planet = document.createElement('div');
    planet.classList.add('planet');
    planet.style.width = `${size}px`;
    planet.style.height = `${size}px`;
    planet.style.position = 'absolute';
    planet.style.left = '0%';
    planet.style.top = '50%';
    planet.style.backgroundColor = color;
    planet.style.borderRadius = '50%';
    planet.style.boxShadow = `0 0 20px ${color}`;
    
    orbit.appendChild(planet);
    container.appendChild(orbit);
}

function getEventClassName(eventText) {
    const checks = [
        ['Server is ready', 'event-server-ready'],
        ['Server is starting', 'event-server-starting'],
        ['Server error', 'event-server-error'],
        ['Player', 'connecting', 'event-player-connecting'],
        ['Player', 'ready', 'event-player-ready'],
        ['Player', 'disconnected', 'event-player-disconnect'],
        ['World Saved', 'event-world-saved'],
        ['Exception', 'event-exception']
    ];
    
    return checks.find(([text, , condition]) => 
        condition ? eventText.includes(text) && eventText.includes(condition) : eventText.includes(text)
    )?.[1] || '';
}