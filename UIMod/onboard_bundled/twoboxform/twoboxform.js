document.addEventListener('DOMContentLoaded', () => {
    // Initialize components
    initializeFormHandlers();
});
/**
 * Initialize form submission and skip button handlers
 */
function initializeFormHandlers() {
    const form = document.getElementById('two-box-form');
    const skipBtn = document.getElementById('skip-btn');

    if (form) {
        form.addEventListener('submit', handleFormSubmission);
    }

    if (skipBtn) {
        skipBtn.addEventListener('click', handleSkipButton);
    }

    // Handle finalize button (event delegation for dynamically added elements)
    document.addEventListener('click', handleFinalizeButton);
}

/**
 * Handle form submission based on current step
 */
async function handleFormSubmission(e) {
    e.preventDefault();
    
    const formData = getFormData();
    const { step, nextStep } = formData;

    // Handle navigation-only steps
    if (isNavigationOnlyStep(step)) {
        navigateToStep(nextStep);
        return;
    }

    try {
        showPreloader();
        const success = await submitFormData(formData);
        
        if (success) {
            handleSuccessfulSubmission(step, nextStep);
        }
    } catch (error) {
        hidePreloader();
        console.error('Form submission error:', error);
        showNotification('Something went wrong!', 'error');
    }
}

/**
 * Extract form data and determine next step
 */
function getFormData() {
    const stepEl = document.getElementById('step');
    const modeEl = document.getElementById('mode');
    const configFieldEl = document.getElementById('config-field');
    const nextStepEl = document.getElementById('next-step');
    const primaryFieldEl = document.getElementById('primary-field');
    const secondaryFieldEl = document.getElementById('secondary-field');

    return {
        step: stepEl ? stepEl.value : '',
        mode: modeEl ? modeEl.value : '',
        configField: configFieldEl ? configFieldEl.value : '',
        nextStep: nextStepEl ? nextStepEl.value : '',
        primaryValue: primaryFieldEl ? primaryFieldEl.value : '',
        secondaryValue: secondaryFieldEl ? secondaryFieldEl.value : ''
    };
}

/**
 * Check if step requires only navigation (no API call)
 */
function isNavigationOnlyStep(step) {
    return ['welcome', 'beta_warning', 'finalize'].includes(step);
}

/**
 * Submit form data to appropriate API endpoint
 */
async function submitFormData(formData) {
    const { step, configField, primaryValue, secondaryValue } = formData;
    
    let url, body;

    if (step === 'admin_account') {
        // Admin account setup
        url = '/api/v2/auth/setup/register';
        body = JSON.stringify({
            username: primaryValue,
            password: secondaryValue
        });
    } else if (configField) {
        // Configuration setting
        url = '/api/v2/settings/save';
        body = JSON.stringify({
            [configField]: isBooleanField(configField) ? 
                convertToBoolean(primaryValue) : 
                primaryValue
        });
    } else {
        // Login
        url = '/auth/login';
        body = JSON.stringify({
            username: primaryValue,
            password: secondaryValue
        });
    }

    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: body
    });

    const data = await response.json();

    if (!response.ok) {
        hidePreloader();
        showNotification(data.error || 'Action failed!', 'error');
        return false;
    }

    return { success: true, data };
}

/**
 * Handle successful form submission
 */
function handleSuccessfulSubmission(step, nextStep) {
    hidePreloader();
    
    if (step === 'admin_account') {
        showNotification('Admin account saved!', 'success');
        setTimeout(() => navigateToStep(nextStep), 800);
    } else if (document.getElementById('mode').value === 'login') {
        showNotification('Login Successful!', 'success');
        setTimeout(async () => {
            await preloadNextPage();
            hidePreloader();
            window.location.href = '/';
        }, 600);
    } else {
        showNotification('Config saved!', 'success');
        setTimeout(() => navigateToStep(nextStep), 800);
    }
}

/**
 * Handle skip button clicks
 */
function handleSkipButton() {
    const step = document.getElementById('step').value;
    const nextStep = document.getElementById('next-step').value;
    
    if (step === 'welcome') {
        window.location.href = '/';
        return;
    }
    
    if (step === 'finalize') {
        showNotification('Setup completed, Auth disabled!', 'success');
        setTimeout(() => window.location.href = '/', 1000);
        return;
    }
    
    navigateToStep(nextStep);
}

/**
 * Handle finalize button clicks (event delegation)
 */
async function handleFinalizeButton(e) {
    if (!e.target || e.target.id !== 'finalize-btn') return;

    try {
        showPreloader();
        const response = await fetch('/api/v2/auth/setup/finalize', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        
        const data = await response.json();
        
        if (response.ok) {
            hidePreloader();
            showNotification(`${data.message}\n${data.restart_hint}`, 'success');
            setTimeout(() => window.location.href = '/login', 2000);
        } else {
            hidePreloader();
            showNotification(data.error || 'Finalize failed!', 'error');
        }
    } catch (error) {
        hidePreloader();
        console.error('Finalize error:', error);
        showNotification('Error finalizing setup!', 'error');
    }
}

/**
 * Navigate to a specific setup step
 */
function navigateToStep(step) {
    window.location.href = `/setup?step=${step}`;
}

/**
 * Check if a config field should be treated as boolean
 */
function isBooleanField(fieldName) {
    const booleanFields = [
        'IsDiscordEnabled', 
        'UPNPEnabled', 
        'ServerVisible', 
        'UseSteamP2P', 
        'IsSSCMEnabled'
    ];
    return booleanFields.includes(fieldName);
}

/**
 * Convert string input to boolean for config fields
 */
function convertToBoolean(value) {
    if (typeof value !== 'string') return false;
    
    const normalizedValue = value.trim().toLowerCase();
    return ['yes', 'true', '1'].includes(normalizedValue);
}

/**
 * Show notification message
 */
function showNotification(message, type = 'error') {
    const existingNotification = document.querySelector('.notification');
    if (existingNotification) existingNotification.remove();

    const notification = document.createElement('div');
    notification.classList.add('notification', type);
    notification.textContent = message;
    
    document.body.appendChild(notification);
    notification.offsetHeight; // Force reflow
    notification.classList.add('show');

    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => notification.remove(), 500);
    }, 3000);
}

/**
 * Show preloader
 */
function showPreloader() {
    const preloader = document.getElementById('preloader');
    if (preloader) preloader.classList.add('show');
}

/**
 * Hide preloader
 */
function hidePreloader() {
    const preloader = document.getElementById('preloader');
    if (preloader) preloader.classList.remove('show');
}

/**
 * Preload next page resources
 */
async function preloadNextPage() {
    try {
        const response = await fetch('/static/favicon.ico', { 
            method: 'HEAD', 
            cache: 'force-cache' 
        });
        return response.ok;
    } catch (error) {
        console.error('Preload failed:', error);
        return false;
    }
}