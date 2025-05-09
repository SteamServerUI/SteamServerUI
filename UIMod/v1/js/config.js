// static/js/config.js
function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.tab-button').forEach(button => button.classList.remove('active'));
    document.getElementById(tabId).classList.add('active');
    document.querySelector(`.tab-button[onclick="showTab('${tabId}')"]`).classList.add('active');
}

document.addEventListener('DOMContentLoaded', () => {
    const discordSelect = document.getElementById('isDiscordEnabled');
    if (discordSelect) {
        discordSelect.addEventListener('change', () => {
            if (discordSelect.value === 'false') {
                alert('You are missing out on nice features! Please consider using the Discord Integration.');
                document.querySelector('.discord-container').classList.add('disabled');
            } else {
                document.querySelector('.discord-container').classList.remove('disabled');
            }
        });
        if (discordSelect.value === 'false') {
            document.querySelector('.discord-container').classList.add('disabled');
        }
    }

    document.querySelectorAll('.section-nav-button').forEach(button => {
        button.addEventListener('click', () => {
            document.querySelectorAll('.config-section').forEach(section => section.classList.remove('active'));
            document.getElementById('select-prompt').style.display = 'none';
            document.querySelectorAll('.section-nav-button').forEach(btn => btn.classList.remove('active'));
            const sectionId = button.getAttribute('data-section');
            document.getElementById(sectionId).classList.add('active');
            button.classList.add('active');
        });
    });
});