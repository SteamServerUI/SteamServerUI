async function sendReset() {
    const game = document.getElementById('gameInput').value;
    
    if (!game) {
        showStatus('Please enter a runfile Identifier.', 'error');
        return;
    }

    try {
        showStatus('Sending reset request...', 'pending');
        
        const response = await fetch('/api/v2/runfile/hardreset', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ game })
        });

        if (response.ok) {
            showStatus('Reset request successful. Game state cleared.', 'success');
        } else {
            const errorData = await response.json().catch(() => ({}));
            const errorMsg = errorData.message || 'Error sending reset request';
            showStatus(`Error: ${errorMsg}`, 'error');
        }
    } catch (error) {
        showStatus(`Network error: ${error.message}`, 'error');
    }
}

function showStatus(message, type) {
    const statusDisplay = document.getElementById('statusDisplay');
    const statusMessage = document.getElementById('statusMessage');
    
    statusMessage.textContent = message;
    statusDisplay.style.display = 'block';
    
    // Remove any existing status classes
    statusDisplay.classList.remove('status-success', 'status-error');
    
    // Add the appropriate status class
    if (type === 'success') {
        statusDisplay.classList.add('status-success');
    } else if (type === 'error') {
        statusDisplay.classList.add('status-error');
    }
}

// Add some keyboard interaction
document.getElementById('gameInput').addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        sendReset();
    }
});