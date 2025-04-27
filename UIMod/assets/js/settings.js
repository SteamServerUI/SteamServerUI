async function initSettings() {
    try {
        const response = await fetch('/api/v2/settings');
        const { data: settings, error } = await response.json();
        if (error) {
            showStatus(`Failed to load settings: ${error}`, true, 'server-config-form');
            return;
        }

        // Group settings by category
        const groups = [...new Set(settings.map(s => s.group))];
        const nav = document.getElementById('server-nav');
        nav.innerHTML = '';
        groups.forEach(group => {
            const button = document.createElement('button');
            button.className = 'section-nav-button';
            button.dataset.section = group;
            button.textContent = group;
            button.addEventListener('click', () => showSection(group));
            nav.appendChild(button);
        });

        // Render settings for each group
        const form = document.getElementById('server-config-form');
        groups.forEach(group => {
            const section = document.createElement('div');
            section.id = group;
            section.className = 'config-section';
            section.innerHTML = `<h3 class="section-title">${group}</h3><div class="channel-grid"></div>`;
            const grid = section.querySelector('.channel-grid');

            settings.filter(s => s.group === group).forEach(setting => {
                const div = document.createElement('div');
                div.className = 'form-group';

                const label = document.createElement('label');
                label.htmlFor = setting.name;
                label.textContent = setting.description;

                let input;
                if (setting.type === 'bool') {
                    input = document.createElement('input');
                    input.type = 'checkbox';
                    input.id = setting.name;
                    input.dataset.name = setting.name;
                    input.checked = setting.value === true;
                    input.addEventListener('change', () => updateSetting(setting.name, input.checked));
                } else if (setting.type === 'int') {
                    input = document.createElement('input');
                    input.type = 'number';
                    input.id = setting.name;
                    input.dataset.name = setting.name;
                    input.value = setting.value ?? '';
                    if (setting.min !== undefined) input.min = setting.min;
                    if (setting.max !== undefined) input.max = setting.max;
                    if (setting.required) input.required = true;
                    input.addEventListener('change', () => {
                        const value = input.value ? parseInt(input.value) : null;
                        if (input.required && !value) {
                            showStatus(`Value for ${setting.name} is required`, true, 'server-config-form');
                            return;
                        }
                        updateSetting(setting.name, value);
                    });
                } else if (setting.type === 'array') {
                    input = document.createElement('input');
                    input.type = 'text';
                    input.id = setting.name;
                    input.dataset.name = setting.name;
                    input.value = setting.value.join(',') || '';
                    input.addEventListener('change', () => {
                        const value = input.value ? input.value.split(',').map(s => s.trim()) : [];
                        updateSetting(setting.name, value);
                    });
                } else if (setting.type === 'map') {
                    input = document.createElement('textarea');
                    input.id = setting.name;
                    input.dataset.name = setting.name;
                    input.value = JSON.stringify(setting.value, null, 2) || '{}';
                    input.addEventListener('change', () => {
                        try {
                            const value = JSON.parse(input.value);
                            updateSetting(setting.name, value);
                        } catch (e) {
                            showStatus(`Invalid JSON for ${setting.name}: ${e.message}`, true, 'server-config-form');
                        }
                    });
                } else {
                    input = document.createElement('input');
                    input.type = 'text';
                    input.id = setting.name;
                    input.dataset.name = setting.name;
                    input.value = setting.value || '';
                    if (setting.required) input.required = true;
                    input.addEventListener('change', () => updateSetting(setting.name, input.value));
                }

                const info = document.createElement('div');
                info.className = 'input-info';
                info.textContent = setting.description;

                div.appendChild(label);
                div.appendChild(input);
                div.appendChild(info);
                grid.appendChild(div);
            });

            const actions = form.querySelector('.form-actions');
            form.insertBefore(section, actions);
        });

        // Show the first section by default
        if (groups.length > 0) {
            showSection(groups[0]);
        } else {
            document.getElementById('select-prompt').style.display = 'block';
        }
    } catch (e) {
        showStatus(`Error fetching settings: ${e.message}`, true, 'server-config-form');
    }
}

async function showSection(group) {
    document.querySelectorAll('.config-section').forEach(section => section.classList.remove('active'));
    document.getElementById('select-prompt').style.display = 'none';
    document.querySelectorAll('#server-nav .section-nav-button').forEach(button => button.classList.remove('active'));

    const button = document.querySelector(`#server-nav .section-nav-button[data-section="${group}"]`);
    button.classList.add('active');

    const section = document.getElementById(group);
    section.classList.add('active');
}

async function updateSetting(name, value) {
    try {
        const response = await fetch('/api/v2/settings/save', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ [name]: value })
        });
        const { status, message } = await response.json();
        if (status === 'error') {
            showStatus(`Failed to update ${name}: ${message}`, true, 'server-config-form');
            return;
        }
        showStatus(`Updated ${name} successfully`, false, 'server-config-form');
    } catch (e) {
        showStatus(`Error updating ${name}: ${e.message}`, true, 'server-config-form');
    }
}

function showStatus(message, isError, formId) {
    if (!formId) {
        console.error('formId is required for showStatus');
        return;
    }

    // Map formId to statusDivId
    const statusDivId = 
        formId === 'server-config-form' ? 'server-status' :
        formId === 'runfile-init-form' ? 'runfile-init-status' :
        'runfile-status';
    
    let status = document.getElementById(statusDivId);
    if (!status) {
        status = document.createElement('div');
        status.id = statusDivId;
        status.style.textAlign = 'center';
        const form = document.getElementById(formId);
        if (!form) {
            console.error(`Form with ID ${formId} not found`);
            return;
        }
        form.appendChild(status);
    }

    status.textContent = message;
    status.style.color = isError ? 'var(--error)' : 'var(--primary)';
    status.style.display = 'block';

    setTimeout(() => {
        status.style.display = 'none';
    }, 30000);
}

document.addEventListener('DOMContentLoaded', initSettings);