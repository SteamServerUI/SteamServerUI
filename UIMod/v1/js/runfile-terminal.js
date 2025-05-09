async function sendReset() {
    const game = document.getElementById('gameInput').value;
    
    if (!game) {
        showStatus('Please enter a runfile Identifier.', true, 'runfile-init-form');
        return;
    }

    try {
        showStatus('Sending reset request...', false, 'runfile-init-form');
        
        const response = await fetch('/api/v2/runfile/hardreset', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ game })
        });

        if (response.ok) {
            showStatus('Reset request successful. Game state cleared.', false, 'runfile-init-form');
        } else {
            const errorData = await response.json().catch(() => ({}));
            const errorMsg = errorData.message || 'Error sending reset request';
            showStatus(`Error: ${errorMsg}`, true, 'runfile-init-form');
        }
    } catch (error) {
        showStatus(`Network error: ${error.message}`, true, 'runfile-init-form');
    }
}

// Add some keyboard interaction
document.getElementById('gameInput').addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        sendReset();
    }
});