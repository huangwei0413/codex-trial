# Repository Guidelines

## Project Structure & Module Organization

This repository follows Go standard project layout with clean architecture:

```
./
├── cmd/api/          # Application entry points
├── internal/         # Private application code
│   ├── handlers/     # HTTP request handlers
│   ├── models/       # Data models and structs
│   └── services/     # Business logic
├── pkg/              # Public library code
│   ├── config/       # Configuration management
│   └── database/     # Database connections
├── deployments/      # Infrastructure manifests
│   └── kubernetes/   # K8s deployment files
├── test/             # Test files
│   ├── unit/         # Unit tests
│   └── integration/  # Integration tests
└── scripts/          # Build and utility scripts
```

## Build, Test, and Development Commands

```bash
# Initialize dependencies
go mod download

# Run the application locally
go run cmd/api/main.go

# Build binary for production
go build -o bin/student-api cmd/api/main.go

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run

# Build Docker image
docker build -f deployments/kubernetes/Dockerfile -t student-api .
```

## Coding Style & Naming Conventions

- Use standard Go formatting (go fmt)
- Package names: lowercase, short, descriptive
- Functions: camelCase with descriptive names
- Constants: UPPER_SNAKE_CASE
- Private functions: start with lowercase letter
- Public functions: start with uppercase letter
- Error handling: always check and handle errors explicitly

Example:
```go
// Good
func (h *StudentHandler) CreateStudent(c *gin.Context) {
    var student models.Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
}
```

## Testing Guidelines

- Use Go built-in testing package
- Test files end with `_test.go`
- Function tests start with `Test`
- Aim for >80% code coverage
- Unit tests: test individual functions in isolation
- Integration tests: test API endpoints with real database

Example:
```go
func TestStudentHandler_CreateStudent(t *testing.T) {
    // Setup test database and handler
    db := setupTestDB(t)
    handler := NewStudentHandler(db)
    
    // Test cases
    tests := []struct {
        name    string
        input   models.Student
        wantErr bool
    }{
        // test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation...
        })
    }
}
```

## Commit & Pull Request Guidelines

**Commit messages**: Follow conventional commits format
```
<type>(<scope>): <description>

Examples:
feat(api): add student deletion endpoint
fix(db): handle database connection errors
docs: update API documentation
```

**Types**: feat, fix, docs, style, refactor, test, chore

**Pull Requirements**:
- Link to related issues in description
- Include tests for new features
- Ensure all tests pass
- Update documentation as needed
- PR title follows commit message format
- Add screenshots for UI changes

## API Design Guidelines

- RESTful principles with proper HTTP methods
- Consistent JSON response format
- Use proper HTTP status codes
- Include health check endpoint (`/health`)
- Version API endpoints (`/api/v1/`)
- Keep the canonical endpoint and payload reference in `docs/api.md` up to date whenever handlers change

Response format example:
```json
{
  "data": {...},
  "error": null,
  "message": "Success"
}
```

## Deployment Guidelines

- Use multi-stage Docker builds
- Include health checks in Kubernetes manifests
- Set appropriate resource limits
- Use environment variables for configuration
- Follow GitFlow branching strategy
