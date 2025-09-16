// API configuration - reads from .env file
const API_BASE_URL = {
    USER: `http://localhost:${window.ENV?.USER_SERVICE_PORT || '8001'}/api/v1`,
    EXPENSE: `http://localhost:${window.ENV?.EXPENSE_SERVICE_PORT || '8002'}/api/v1`,
    REPORT: `http://localhost:${window.ENV?.REPORT_SERVICE_PORT || '8003'}/api/v1`
};

// Fallback for direct port access
if (!window.ENV) {
    API_BASE_URL.USER = 'http://localhost:8001/api/v1';
    API_BASE_URL.EXPENSE = 'http://localhost:8002/api/v1';
    API_BASE_URL.REPORT = 'http://localhost:8003/api/v1';
}

console.log('API Configuration:', API_BASE_URL);

// API utility functions
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
    };
}

async function apiRequest(url, options = {}) {
    const config = {
        headers: getAuthHeaders(),
        ...options
    };
    
    try {
        const response = await fetch(url, config);
        
        if (response.status === 401) {
            // Token expired or invalid
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
            return null;
        }
        
        return response;
    } catch (error) {
        console.error('API request failed:', error);
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            throw new Error('Service unavailable. Please check if the server is running.');
        }
        throw new Error('Network error. Please try again.');
    }
}

// User API functions
const userAPI = {
    async register(userData) {
        return apiRequest(`${API_BASE_URL.USER}/users/register`, {
            method: 'POST',
            body: JSON.stringify(userData)
        });
    },
    
    async login(credentials) {
        return apiRequest(`${API_BASE_URL.USER}/users/login`, {
            method: 'POST',
            body: JSON.stringify(credentials)
        });
    }
};

// Expense API functions
const expenseAPI = {
    async getExpenses(params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const url = `${API_BASE_URL.EXPENSE}/expenses${queryString ? '?' + queryString : ''}`;
        return apiRequest(url);
    },
    
    async createExpense(expenseData) {
        return apiRequest(`${API_BASE_URL.EXPENSE}/expenses`, {
            method: 'POST',
            body: JSON.stringify(expenseData)
        });
    },
    
    async updateExpense(id, expenseData) {
        return apiRequest(`${API_BASE_URL.EXPENSE}/expenses/${id}`, {
            method: 'PUT',
            body: JSON.stringify(expenseData)
        });
    },
    
    async deleteExpense(id) {
        return apiRequest(`${API_BASE_URL.EXPENSE}/expenses/${id}`, {
            method: 'DELETE'
        });
    }
};

// Report API functions
const reportAPI = {
    async getReports() {
        return apiRequest(`${API_BASE_URL.REPORT}/reports`);
    },
    
    async generateMonthlyReport() {
        return apiRequest(`${API_BASE_URL.REPORT}/reports/monthly`);
    }
};

// Health check function
async function checkServiceHealth() {
    const services = [
        { name: 'User Service', url: 'http://localhost:8081/healthz' },
        { name: 'Expense Service', url: 'http://localhost:8082/healthz' },
        { name: 'Report Service', url: 'http://localhost:8083/healthz' }
    ];
    
    const results = await Promise.allSettled(
        services.map(async service => {
            try {
                const response = await fetch(service.url);
                return { ...service, status: response.ok ? 'healthy' : 'unhealthy' };
            } catch (error) {
                return { ...service, status: 'unreachable' };
            }
        })
    );
    
    return results.map(result => result.value);
}