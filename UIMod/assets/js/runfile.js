// static/js/runfile.js
async function initRunfile() {
    try {
        const response = await fetch('/api/v2/runfile/groups');
        const { data, error } = await response.json();
        if (error) {
            showStatus(`Failed to load groups: ${error}`, true);
            return;
        }

        const nav = document.getElementById('runfile-group-nav');
        data.forEach(group => {
            const button = document.createElement('button');
            button.className = 'section-nav-button';
            button.dataset.group = group;
            button.textContent = group;
            button.addEventListener('click', () => showGroup(group));
            nav.appendChild(button);
        });
    } catch (err) {
        showStatus(`Error fetching groups: ${err.message}`, true);
    }
}

async function showGroup(group) {
    try {
        // Hide all sections and prompt
        document.querySelectorAll('.config-section').forEach(section => section.classList.remove('active'));
        document.getElementById('runfile-select-prompt').style.display = 'none';
        document.querySelectorAll('#runfile-group-nav .section-nav-button').forEach(btn => btn.classList.remove('active'));

        // Activate selected group button
        const button = document.querySelector(`#runfile-group-nav .section-nav-button[data-group="${group}"]`);
        button.classList.add('active');

        // Check if section exists, create if not
        let section = document.getElementById(group);
        if (!section) {
            section = document.createElement('div');
            section.id = group;
            section.className = 'config-section';
            section.innerHTML = `<h3 class="section-title">${group} GameServer Settings</h3><div class="channel-grid"></div>`;
            const form = document.getElementById('runfile-config-form');
            const formActions = form.querySelector('.form-actions');
            if (formActions) {
                form.insertBefore(section, formActions);
            } else {
                form.appendChild(section); // Fallback
            }
        }

        // Fetch args
        const response = await fetch(`/api/v2/runfile/args?group=${encodeURIComponent(group)}`);
        const { data, error } = await response.json();
        if (error) {
            showStatus(`Failed to load args: ${error}`, true);
            return;
        }

        // Render args
        const grid = section.querySelector('.channel-grid');
        grid.innerHTML = '';
        data.forEach(arg => {
            const formGroup = document.createElement('div');
            formGroup.className = 'form-group';

            const label = document.createElement('label');
            label.htmlFor = arg.flag;
            label.textContent = arg.ui_label || arg.flag;

            let input;
            if (arg.type === 'bool') {
                input = document.createElement('input');
                input.type = 'checkbox';
                input.id = arg.flag;
                input.dataset.flag = arg.flag;
                input.checked = arg.runtime_value === 'true';
            } else if (arg.type === 'int') {
                input = document.createElement('input');
                input.type = 'number';
                input.id = arg.flag;
                input.dataset.flag = arg.flag;
                input.value = arg.runtime_value || '';
                if (arg.min) input.min = arg.min;
                if (arg.max) input.max = arg.max;
            } else {
                input = document.createElement('input');
                input.type = 'text';
                input.id = arg.flag;
                input.dataset.flag = arg.flag;
                input.value = arg.runtime_value || '';
            }

            if (arg.disabled) input.disabled = true;
            if (arg.required) input.required = true;

            input.addEventListener('change', () => {
                const value = arg.type === 'bool' ? input.checked.toString() : input.value;
                updateArg(arg.flag, value);
            });

            const info = document.createElement('div');
            info.className = 'input-info';
            info.textContent = arg.description || 'No description available';

            formGroup.appendChild(label);
            formGroup.appendChild(input);
            formGroup.appendChild(info);
            grid.appendChild(formGroup);
        });

        section.classList.add('active');
    } catch (err) {
        showStatus(`Error fetching args: ${err.message}`, true);
    }
}

async function updateArg(flag, value) {
    try {
        const response = await fetch('/api/v2/runfile/args/update', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ flag, value })
        });
        const { data, error } = await response.json();
        if (error) {
            showStatus(`Failed to update ${flag}: ${error}`, true);
            return;
        }
        showStatus(`Updated ${flag} to ${value}`, false);
    } catch (err) {
        showStatus(`Error updating ${flag}: ${err.message}`, true);
    }
}

async function saveRunfile() {
    try {
        const response = await fetch('/api/v2/runfile/save', {
            method: 'POST'
        });
        const { data, error } = await response.json();
        if (error) {
            showStatus(`Failed to save runfile: ${error}`, true);
            return;
        }
        showStatus('Runfile saved successfully', false);
    } catch (err) {
        showStatus(`Error saving runfile: ${err.message}`, true);
    }
}

function showStatus(message, isError) {
    const status = document.getElementById('runfile-status');
    status.textContent = message;
    status.style.color = isError ? 'var(--error)' : 'var(--primary)';
    status.style.display = 'block';
    setTimeout(() => {
        status.style.display = 'none';
    }, 3000);
}

document.addEventListener('DOMContentLoaded', initRunfile);