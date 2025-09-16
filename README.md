# FinTrack - Cloud-Native Expense Management System

A production-ready microservices-based expense management system built with Go, demonstrating clean architecture, concurrency, and cloud-native practices.

## 🏗️ Architecture

- **User Service** (Port 8081): Authentication and user profile management
- **Expense Service** (Port 8082): CRUD operations for expenses and categories
- **Report Service** (Port 8083): Background report generation with goroutines
- **PostgreSQL**: Primary database for persistence
- **Redis**: Caching and pub/sub for notifications

## 🚀 Tech Stack

- **Backend**: Go 1.21+, Gin, GORM
- **Database**: PostgreSQL, Redis
- **Observability**: Zap logging, Prometheus metrics
- **Containerization**: Docker, Docker Compose
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions

## 📁 Project Structure

```
fintrack/
├── cmd/                    # Application entry points
│   ├── user-service/
│   ├── expense-service/
│   └── report-service/
├── internal/               # Private application code
│   ├── common/            # Shared models and types
│   ├── user/              # User service logic
│   ├── expense/           # Expense service logic
│   └── report/            # Report service logic
├── pkg/                   # Public packages
│   ├── config/            # Configuration management
│   ├── database/          # Database connections
│   └── middleware/        # HTTP middleware
├── k8s/                   # Kubernetes manifests
├── migrations/            # Database migrations
└── .github/workflows/     # CI/CD pipelines
```

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### 1. Clone and Setup
```bash
git clone <repo-url>
cd fintrack
cp .env.example .env
```

### 2. Start with Docker Compose
```bash
# Start all services
docker-compose up -d

# Check service health
curl http://localhost:8081/healthz
curl http://localhost:8082/healthz  
curl http://localhost:8083/healthz
```

### 3. Test the APIs
```bash
# Register a user
curl -X POST http://localhost:8081/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login (save the token)
curl -X POST http://localhost:8081/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Create an expense (use token from login)
curl -X POST http://localhost:8082/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"amount":25.50,"description":"Lunch","category":"Food","date":"2024-01-15"}'
```

## 🔧 Development

### Run Locally
```bash
# Install dependencies
go mod tidy

# Start databases
docker-compose up -d postgres redis

# Run services individually
make run-user      # Port 8081
make run-expense   # Port 8082  
make run-report    # Port 8083
```

### Run Tests
```bash
make test
```

### Build Services
```bash
make build
```

## 📊 API Endpoints

### User Service (Port 8081)
- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/login` - User login
- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics

### Expense Service (Port 8082)
- `GET /api/v1/expenses` - List expenses (with pagination)
- `POST /api/v1/expenses` - Create expense
- `PUT /api/v1/expenses/:id` - Update expense
- `DELETE /api/v1/expenses/:id` - Delete expense
- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics

### Report Service (Port 8083)
- `GET /api/v1/reports/monthly` - Generate monthly report
- `GET /api/v1/reports` - List all reports
- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics

## 🐳 Docker Deployment

### Build Images
```bash
make docker-build
```

### Deploy with Docker Compose
```bash
docker-compose up -d
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes
```bash
# Create namespace and deploy
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/user-service.yaml

# Check deployment status
kubectl get pods -n fintrack
```

## 🔍 Observability

### Logs
Each service uses structured logging with Zap:
```bash
docker-compose logs -f user-service
```

### Metrics
Prometheus metrics available at `/metrics` endpoint:
```bash
curl http://localhost:8081/metrics
```

### Health Checks
```bash
curl http://localhost:8081/healthz
curl http://localhost:8082/healthz
curl http://localhost:8083/healthz
```

## 🧪 Testing

The project includes comprehensive tests:

- **Unit Tests**: Service layer logic
- **Integration Tests**: Database operations
- **API Tests**: HTTP endpoint testing

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...
```

## 🚀 Features Demonstrated

### Backend Development
- ✅ Clean Architecture (cmd/, internal/, pkg/)
- ✅ RESTful API design
- ✅ JWT Authentication
- ✅ Database operations with GORM
- ✅ Input validation and error handling

### Concurrency & Performance
- ✅ Goroutines for background processing
- ✅ Channels for communication
- ✅ Redis pub/sub for notifications
- ✅ Connection pooling

### Cloud-Native Practices
- ✅ Containerized microservices
- ✅ Health checks and graceful shutdown
- ✅ Configuration management
- ✅ Structured logging
- ✅ Metrics collection

### DevOps & Deployment
- ✅ Docker multi-stage builds
- ✅ Docker Compose for local development
- ✅ Kubernetes manifests
- ✅ CI/CD with GitHub Actions
- ✅ Database migrations

## 🔐 Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Input validation and sanitization
- CORS middleware
- Environment-based configuration

## 📈 Scalability Features

- Microservices architecture
- Horizontal scaling with Kubernetes
- Redis caching
- Database connection pooling
- Stateless service design

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

---

**Built with ❤️ for demonstrating production-ready Go microservices**