.PHONY: build test run clean docker-build docker-run deploy-staging deploy-prod

# Build the application
build:
	go build -o bin/student-api cmd/api/main.go

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run the application locally
run:
	go run cmd/api/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	docker build -f deployments/kubernetes/Dockerfile -t student-api:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 student-api:latest

# Deploy to staging
deploy-staging:
	kubectl apply -f deployments/kubernetes/namespace.yaml
	kubectl apply -f deployments/kubernetes/ -n student-api
	kubectl rollout status deployment/student-api -n student-api

# Deploy to production
deploy-prod:
	kubectl apply -f deployments/kubernetes/namespace.yaml
	kubectl apply -f deployments/kubernetes/ -n student-api
	kubectl rollout status deployment/student-api -n student-api

# Download dependencies
deps:
	go mod download
	go mod tidy
