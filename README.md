# Product Inventory Management API

A RESTful API service built with Go and Fiber framework for managing product inventory with stock tracking and low-stock monitoring capabilities.

## 📋 Table of Contents

- [Project Description](#project-description)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Running Tests](#running-tests)
- [API Documentation](#api-documentation)
- [Design Choices & Architecture](#design-choices--architecture)
- [Project Structure](#project-structure)
- [Contributing](#contributing)

## 🚀 Project Description

This project is a backend service for managing product inventory operations including:

- Product CRUD operations (Create, Read, Update, Delete)
- Stock increment and decrement operations
- Low stock monitoring with configurable thresholds per product
- Comprehensive error handling and validation
- Clean architecture with domain-driven design principles

Built using modern Go practices with a focus on maintainability, testability, and performance.

## ✨ Features

- **Product Management**: Full CRUD operations for products
- **Stock Operations**: Increment/decrement stock with validation
- **Low Stock Filtering**: Query products below their individual stock thresholds
- **Robust Error Handling**: Custom error types with appropriate HTTP status codes
- **Database Integration**: PostgreSQL with GORM ORM
- **Auto-Migration**: Database schema migration endpoint
- **Comprehensive Testing**: Unit tests with mocks for all business logic
- **Clean Architecture**: Separation of concerns with domain, infrastructure, and transport layers

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 12 or higher
- **Make**: For using the provided Makefile commands
- **Git**: For cloning the repository

Optional (for development):
- **Air**: For hot reloading during development (`go install github.com/air-verse/air@latest`)

## 🛠 Local Setup

### 1. Clone the Repository

```bash
git clone https://github.com/xxthunderblastxx/ase-challenge.git
cd ase-challenge
```

### 2. Environment Configuration

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Configure your `.env` file with the following variables:

```env
# Server Configuration
PORT=8080

# PostgreSQL Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_DB=product_inventory
```

### 3. Install Dependencies

```bash
make tidy
```

### 4. Database Migration

Run the application and trigger migration:

```bash
# Start the server
make dev

# In another terminal, trigger migration
curl http://localhost:8080/api/v1/migrate
```

### 5. Running the Application

#### Development Mode (with hot reload):
> For this you need to have `air` installed.(see Prerequisites)
```bash
make watch
```

#### Development Mode (basic):
```bash
make dev
```

#### Production Mode:
```bash
make build
./bin/aes-challenge-backend
```

The API will be available at `http://localhost:8080`

## 🧪 Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Verbose Output
```bash
go test ./... -v
```

### Run Tests for Specific Package
```bash
# Test domain layer
go test ./internal/domain/product/ -v

# Test handlers
go test ./internal/transport/http/handlers/ -v
```

### Run Tests with Coverage
```bash
go test ./... -cover -v
```

### Available Test Suites

1. **Stock Operations Tests** (`service_stock_test.go`):
   - Increment stock validation and edge cases
   - Decrement stock validation and insufficient stock scenarios
   - Product not found and invalid input handling

2. **Handler Tests** (`product_test.go`):
   - GetAllProducts with and without low-stock filtering
   - Error handling and response validation
   - Query parameter processing

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | Get all products |
| GET | `/products?low-stock=true` | Get all products with low stock |
| GET | `/products/:id` | Get product by ID |
| POST | `/products` | Create new product |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |
| POST | `/products/:id/increment-stock` | Increment product stock |
| POST | `/products/:id/decrement-stock` | Decrement product stock |

#### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/migrate` | Run database migrations |

> For More Deatailed API documentation, run the server and visit: `http://localhost:8080/docs/`

## 🏗 Design Choices & Architecture

### Architecture Pattern

**Clean Architecture** with the following layers:

- **Domain Layer**: Business entities and interfaces
- **Application Layer**: Business logic and use cases
- **Infrastructure Layer**: External dependencies (database, etc.)
- **Transport Layer**: HTTP handlers and routing

### Key Design Decisions

#### 1. **Maximum Reusability Principle**
- **Decision**: Use existing `GetAllProducts()` service method for low-stock filtering
- **Rationale**: Avoid code duplication and maintain consistency
- **Implementation**: Filter products at the handler layer after retrieval

#### 2. **Individual Product Thresholds**
- **Decision**: Each product has its own `LowStockThreshold` field
- **Rationale**: Different products may have different stock requirements
- **Benefits**: Flexible inventory management per product type

#### 3. **Custom Error Handling**
- **Decision**: Implement custom `AppError` types with specific error codes
- **Rationale**: Better error categorization and client-side error handling
- **Examples**: `INSUFFICIENT_STOCK`, `PRODUCT_NOT_FOUND`, `INVALID_INPUT`

#### 4. **Repository Pattern**
- **Decision**: Abstract database operations behind repository interfaces
- **Benefits**: Easy testing with mocks and potential database switching
- **Implementation**: `ProductRepository` interface with PostgreSQL implementation

#### 5. **Stock Operation Validation**
- **Assumptions**:
  - Stock quantities cannot be negative
  - Increment/decrement operations must use positive values
  - Stock operations are atomic (no concurrent modification handling in current version)
- **Edge Cases Handled**:
  - Attempting to decrement more stock than available
  - Invalid product IDs
  - Zero or negative increment/decrement values

#### 6. **HTTP Status Code Strategy**
- `200 OK`: Successful operations
- `201 Created`: Resource creation
- `204 No Content`: Successful deletion
- `400 Bad Request`: Invalid input or missing required fields
- `404 Not Found`: Resource not found
- `409 Conflict`: Insufficient stock operations
- `500 Internal Server Error`: Database or system errors

#### 7. **Testing Strategy**
- **Unit Tests**: Comprehensive coverage for business logic
- **Mocking**: In-memory repository implementations for isolated testing
- **Edge Case Coverage**: All error scenarios and boundary conditions
- **Handler Tests**: HTTP request/response validation

### Assumptions Made

1. **Concurrent Access**: Current implementation doesn't handle concurrent stock modifications
2. **Stock Validation**: Negative stock quantities are not allowed in the system
3. **Database Consistency**: GORM handles basic database consistency requirements
4. **Authentication**: No authentication/authorization implemented (would be added based on requirements)
5. **Pagination**: Not implemented for product listing (could be added for large datasets)


## 📁 Project Structure

```
ase-challenge/
├── cmd/
│   └── main.go                 # Application entry point
├── docs/
│   └── swagger.yaml           # OpenAPI 3.0 specification
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── domain/
│   │   └── product/
│   │       ├── entity.go      # Product entity
│   │       ├── repository.go  # Repository interface
│   │       ├── service.go     # Business logic
│   │       └── *_test.go     # Unit tests
│   ├── infrastructure/
│   │   └── postgres/
│   │       ├── connection.go  # Database connection
│   │       └── product.go     # Repository implementation
│   ├── pkg/
│   │   ├── errors/           # Custom error handling
│   │   ├── model/           # Base models
│   │   └── response/        # Response utilities
│   ├── server/
│   │   └── server.go        # Server setup
│   └── transport/
│       └── http/
│           ├── handlers/    # HTTP handlers
│           └── router/      # Route definitions
├── .air.toml               # Hot reload configuration
├── .env.example           # Environment variables template
├── Makefile              # Build and development commands
├── go.mod               # Go module definition
└── README.md           # This file
```

**Built with ❤️ using Go, Fiber, and PostgreSQL**
