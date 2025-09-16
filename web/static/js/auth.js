// Authentication utilities
function checkAuth() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (token && user) {
        showAuthenticatedNav();
        return true;
    } else {
        showUnauthenticatedNav();
        return false;
    }
}

function showAuthenticatedNav() {
    const navLinks = document.getElementById('nav-links');
    const authLinks = document.getElementById('auth-links');
    const logoutBtn = document.getElementById('logout-btn');
    const userName = document.getElementById('user-name');
    
    if (navLinks) navLinks.classList.remove('hidden');
    if (authLinks) authLinks.classList.add('hidden');
    if (logoutBtn) logoutBtn.classList.remove('hidden');
    
    if (userName) {
        try {
            const user = JSON.parse(localStorage.getItem('user') || '{}');
        } catch (e) {
            const user = {};
        }
        userName.textContent = user.name || user.email || 'User';
        userName.classList.remove('hidden');
    }
}

function showUnauthenticatedNav() {
    const navLinks = document.getElementById('nav-links');
    const authLinks = document.getElementById('auth-links');
    const logoutBtn = document.getElementById('logout-btn');
    const userName = document.getElementById('user-name');
    
    if (navLinks) navLinks.classList.add('hidden');
    if (authLinks) authLinks.classList.remove('hidden');
    if (logoutBtn) logoutBtn.classList.add('hidden');
    if (userName) userName.classList.add('hidden');
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// Initialize auth state on page load
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }
});

// Redirect to login if not authenticated on protected pages
function requireAuth() {
    if (!checkAuth()) {
        window.location.href = '/login';
    }
}