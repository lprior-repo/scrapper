# GitHub Codeowners Visualization Tool

A comprehensive GitHub organization scanner and codeowners visualization tool that maps code ownership relationships through an interactive graph interface.

## Features

- **GitHub Organization Scanning**: Scans GitHub organizations, repositories, teams, and users
- **CODEOWNERS Analysis**: Parses and analyzes CODEOWNERS files to extract ownership patterns
- **Interactive Graph Visualization**: React-based interactive graph showing relationships between organizations, repositories, teams, and users
- **Real-time Statistics**: Provides comprehensive statistics about code ownership and coverage
- **REST API**: Complete REST API for programmatic access to all functionality
- **Pure Core Architecture**: Follows Pure Core/Impure Shell architecture for maintainability and testability

## Architecture

- **Backend**: Go with Pure Core/Impure Shell architecture
- **Frontend**: React/TypeScript with Effect-TS for type safety and composability
- **Ingestor**: TypeScript with Effect-TS for functional data processing
- **Database**: Neo4j graph database for storing relationships
- **Package Manager**: Bun for fast TypeScript development
- **API**: RESTful HTTP API with comprehensive endpoints

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Bun (JavaScript runtime and package manager)
- GitHub Personal Access Token (for scanning)
- Task (task runner) - install with `go install github.com/go-task/task/v3/cmd/task@latest`

### 1. Clone and Setup

```bash
git clone <repository-url>
cd scrapper
```

### 2. Setup and Start Services

```bash
# Setup development environment
task setup

# Start full development stack
task dev
```

### 3. Access the Application

- **Web Interface**: http://localhost:3000
- **API**: http://localhost:8081
- **Neo4j Browser**: http://localhost:7474 (username: neo4j, password: password)

### 4. Scan an Organization

```bash
# Using the API (Go backend)
curl -X POST http://localhost:8081/api/scan/kubernetes

# Using the ingestor (TypeScript)
task ingest ORG=kubernetes

# Using the web interface
# Visit http://localhost:3000 and enter organization name
```

## Development Setup

### Using Task (Recommended)

```bash
# Setup everything
task setup

# Start full development stack
task dev

# Start individual services
task start-api      # Go backend only
task start-frontend # React frontend only

# Run tests
task test           # All tests
task test-frontend  # Frontend tests only

# Build everything
task build
```

### Manual Setup

1. **Start database services**:

```bash
docker-compose up -d neo4j
```

2. **Install dependencies**:

```bash
bun install
cd packages/webapp && bun install
```

3. **Run services**:

```bash
# Backend
go run . api

# Frontend
cd packages/webapp && bun run dev
```

### Environment Variables

| Variable         | Description                          | Default                 |
| ---------------- | ------------------------------------ | ----------------------- |
| `GITHUB_TOKEN`   | GitHub Personal Access Token         | Required                |
| `NEO4J_URI`      | Neo4j database URI                   | `bolt://localhost:7687` |
| `NEO4J_USERNAME` | Neo4j username                       | `neo4j`                 |
| `NEO4J_PASSWORD` | Neo4j password                       | `password`              |
| `ENVIRONMENT`    | Environment (development/production) | `development`           |

## API Endpoints

### Organization Endpoints

- `POST /api/scan/{org}` - Scan a GitHub organization
- `GET /api/graph/{org}` - Get graph visualization data
- `GET /api/stats/{org}` - Get organization statistics
- `GET /api/organizations/{org}` - Get organization details

### Repository Endpoints

- `GET /api/repositories/{org}/{repo}` - Get repository details

### Utility Endpoints

- `GET /api/health` - Health check
- `GET /api/version` - Version information

## Available Commands

### Task Commands

```bash
# Development
task dev           # Start full development stack
task setup         # Setup development environment
task clean         # Clean up environment

# Building
task build         # Build everything
task build-api     # Build Go backend
task build-frontend # Build React frontend

# Testing
task test          # Run all tests
task test-frontend # Run frontend tests
task lint          # Run linters

# Services
task start-api     # Start Go API server
task start-frontend # Start React frontend
task ingest ORG=<org> # Run TypeScript ingestor
```

### Direct CLI Commands

```bash
# Go backend
go run . api       # Start API server
go run . scan <org> # Scan organization

# TypeScript packages
cd packages/webapp && bun run dev  # Start frontend
cd packages/ingestor && bun run ingest <org> # Run ingestor
```

## Testing

The project includes comprehensive testing:

```bash
# Run all tests (Go + TypeScript)
task test

# Run Go tests only
go test -v ./...

# Run TypeScript tests only
task test-frontend

# Run specific test categories
go test -v -run "TestGitHub.*"    # GitHub integration tests
go test -v -run "TestGraph.*"     # Graph database tests
go test -v -run "TestPure.*"      # Pure core logic tests
```

## Building

### Using Task

```bash
# Build everything
task build

# Build individual components
task build-api      # Go backend
task build-frontend # React frontend
```

### Manual Build

```bash
# Build Go application
go build -o overseer .

# Build frontend
cd packages/webapp && bun run build
```

## Project Structure

```
├── main.go                    # Go application entry point
├── github_*.go               # GitHub API integration
├── neo4j_*.go                # Graph database operations
├── error_handling.go         # Error handling utilities
├── batch_processing.go       # Batch processing utilities
├── Taskfile.yml              # Task runner configuration
├── docker-compose.yml        # Docker setup
├── package.json              # Root TypeScript dependencies
├── tsconfig.json            # Root TypeScript configuration
├── packages/                 # TypeScript monorepo
│   ├── webapp/              # React frontend with Effect-TS
│   │   ├── src/
│   │   │   ├── App.tsx      # Main application component
│   │   │   ├── components/  # React components
│   │   │   ├── services.ts  # Effect-TS service definitions
│   │   │   └── index.tsx    # Entry point
│   │   ├── package.json     # Frontend dependencies
│   │   └── tsconfig.json    # Frontend TypeScript config
│   └── ingestor/            # TypeScript data ingestor
│       ├── src/
│       │   ├── index.ts     # Ingestor entry point
│       │   └── services.ts  # Effect-TS services
│       └── package.json     # Ingestor dependencies
└── README.md               # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the Pure Core/Impure Shell architecture
4. Write comprehensive tests
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:

- Create an issue on GitHub
- Check existing documentation
- Review the API documentation at `/api/docs`
