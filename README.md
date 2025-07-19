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
- **Frontend**: React/TypeScript with Bun package manager
- **Database**: Neo4j graph database for storing relationships
- **API**: RESTful HTTP API with comprehensive endpoints

## Quick Start

### Prerequisites

- Docker and Docker Compose
- GitHub Personal Access Token (for scanning)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd scrapper
```

### 2. Start Services

```bash
# Start all services (Neo4j, Redis, API, UI)
docker-compose up -d

# Wait for services to be healthy
docker-compose ps
```

### 3. Access the Application

- **Web Interface**: http://localhost:3000
- **API**: http://localhost:8081
- **Neo4j Browser**: http://localhost:7474 (username: neo4j, password: password)

### 4. Scan an Organization

```bash
# Set your GitHub token
export GITHUB_TOKEN="your_github_token_here"

# Scan an organization
./overseer scan microsoft

# Or using Docker
docker exec overseer-app ./overseer scan microsoft
```

## Development Setup

### Local Development

1. **Start only the database services**:
```bash
docker-compose -f docker-compose.dev.yml up -d
```

2. **Run the backend locally**:
```bash
export GITHUB_TOKEN="your_token"
export NEO4J_URI="bolt://localhost:7687"
export NEO4J_USERNAME="neo4j"
export NEO4J_PASSWORD="password"

go run . api
```

3. **Run the frontend locally**:
```bash
cd ui
bun install
bun run dev
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GITHUB_TOKEN` | GitHub Personal Access Token | Required |
| `NEO4J_URI` | Neo4j database URI | `bolt://localhost:7687` |
| `NEO4J_USERNAME` | Neo4j username | `neo4j` |
| `NEO4J_PASSWORD` | Neo4j password | `password` |
| `ENVIRONMENT` | Environment (development/production) | `development` |

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

## CLI Commands

```bash
# Show help
./overseer help

# Start API server
./overseer api

# Scan an organization
./overseer scan <organization>

# Clean up processes
./overseer cleanup
```

## Testing

The project includes comprehensive testing:

```bash
# Run all tests
go test -v ./...

# Run specific test categories
go test -v -run "TestGitHub.*"    # GitHub integration tests
go test -v -run "TestGraph.*"     # Graph database tests
go test -v -run "TestPure.*"      # Pure core logic tests

# Run frontend tests
cd ui
bun test
```

## Building

### Local Build

```bash
# Build Go application
go build -o overseer .

# Build frontend
cd ui
bun run build
```

### Docker Build

```bash
# Build all services
docker-compose build

# Build specific service
docker-compose build app
```

## Project Structure

```
├── main.go                 # Application entry point
├── config.go              # Configuration management
├── github_*.go            # GitHub API integration
├── graph_*.go             # Graph database operations
├── http_server.go         # HTTP server and API endpoints
├── migrations.go          # Database migrations
├── error_handling.go      # Error handling utilities
├── batch_processing.go    # Batch processing utilities
├── docker-compose.yml     # Production Docker setup
├── docker-compose.dev.yml # Development Docker setup
├── Dockerfile            # Application Docker image
├── ui/                   # React frontend
│   ├── src/
│   │   ├── App.tsx       # Main application component
│   │   ├── components/   # React components
│   │   └── main.tsx      # Entry point
│   ├── Dockerfile        # UI Docker image
│   └── package.json      # UI dependencies
└── README.md            # This file
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
