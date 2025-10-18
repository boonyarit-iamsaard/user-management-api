# User Management API

A RESTful API built with Go and Fiber framework for managing users with authentication and authorization features.

## Project Structure

```text
.
├── cmd/                    # Application entry points
│   └── main.go             # Main application entry point
├── internal/               # Private application code
│   ├── domain/             # Domain models and business logic
│   ├── handler/            # HTTP handlers and routing
│   ├── repository/         # Data access layer
│   └── service/            # Business services
├── migrations/             # Database migration files
├── tests/                  # Test files
├── docs/                   # Documentation
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── Makefile                # Build and run commands
└── .gitignore              # Git ignore file
```

## Getting Started

### Prerequisites

- Go 1.25.0 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/boonyarit-iamsaard/user-management-api.git
cd user-management-api
```

1. Download dependencies:

```bash
make deps
```

### Running the Application

1. Run the application:

```bash
make run
```

1. Or run directly:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`

### Available Endpoints

- `GET /` - Welcome message
- `GET /api/v1/health` - Health check endpoint

### Build

To build the application:

```bash
make build
```

This will create a binary in the `bin/` directory.

### Development Commands

- `make run` - Run the application
- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make fmt` - Format code
- `make test` - Run tests
- `make deps` - Download and tidy dependencies
- `make lint` - Run linter

## Technology Stack

- **Go** - Programming language
- **Fiber** - Web framework
- **PostgreSQL** - Database (planned)
- **GORM** - ORM (planned)
- **JWT** - Authentication (planned)
