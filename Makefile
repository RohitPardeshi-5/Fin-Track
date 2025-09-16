.PHONY: build test run-user run-expense run-report docker-build migrate-up migrate-down

# Build all services
build:
	go build -o bin/user-service ./cmd/user-service
	go build -o bin/expense-service ./cmd/expense-service
	go build -o bin/report-service ./cmd/report-service
	go build -o bin/web-frontend ./cmd/web-frontend

# Run tests
test:
	go test -v ./...

# Run individual services
run-user:
	go run ./cmd/user-service

run-expense:
	go run ./cmd/expense-service

run-report:
	go run ./cmd/report-service

run-web:
	go run ./cmd/web-frontend

# Docker operations
docker-build:
	docker build -f Dockerfile.user-service -t fintrack/user-service .
	docker build -f Dockerfile.expense-service -t fintrack/expense-service .
	docker build -f Dockerfile.report-service -t fintrack/report-service .
	docker build -f Dockerfile.web-frontend -t fintrack/web-frontend .

# Database migrations
migrate-up:
	docker-compose exec postgres psql -U user -d fintrack -c "SELECT 1;"

migrate-down:
	docker-compose exec postgres psql -U user -d fintrack -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Development
dev:
	docker-compose up -d postgres redis
	sleep 5
	go run ./cmd/user-service &
	go run ./cmd/expense-service &
	go run ./cmd/report-service &
	go run ./cmd/web-frontend &

# Clean up
clean:
	docker-compose down -v
	rm -rf bin/