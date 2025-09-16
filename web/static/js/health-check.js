// Service health check utility
async function checkServiceHealth() {
    const services = [
        { name: 'User Service', url: 'http://localhost:8081/healthz' },
        { name: 'Expense Service', url: 'http://localhost:8082/healthz' },
        { name: 'Report Service', url: 'http://localhost:8083/healthz' }
    ];
    
    const results = [];
    
    for (const service of services) {
        try {
            const response = await fetch(service.url, { 
                method: 'GET',
                timeout: 5000 
            });
            results.push({
                ...service,
                status: response.ok ? 'healthy' : 'unhealthy',
                statusCode: response.status
            });
        } catch (error) {
            results.push({
                ...service,
                status: 'unreachable',
                error: error.message
            });
        }
    }
    
    return results;
}

// Display service status
function displayServiceStatus(results) {
    const statusDiv = document.getElementById('service-status');
    if (!statusDiv) return;
    
    const html = results.map(service => {
        const statusClass = service.status === 'healthy' ? 'text-green-600' : 'text-red-600';
        return `
            <div class="flex justify-between items-center p-2 border-b">
                <span>${service.name}</span>
                <span class="${statusClass}">${service.status}</span>
            </div>
        `;
    }).join('');
    
    statusDiv.innerHTML = html;
}

// Auto-check services on page load
document.addEventListener('DOMContentLoaded', async () => {
    try {
        const results = await checkServiceHealth();
        displayServiceStatus(results);
    } catch (error) {
        console.error('Health check failed:', error);
    }
});