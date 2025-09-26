# Structured Error Handling System

This package provides a comprehensive error handling system for the ASE Challenge application with proper HTTP status code mappings, structured error responses, and consistent error formats across the entire application.

## Overview

The error handling system consists of three main components:

1. **Custom Error Types** (`errors.go`) - Predefined error types with specific codes and HTTP status mappings
2. **Error Handler Middleware** (`handler.go`) - Automatic error handling and response formatting
3. **Enhanced Response Structure** (`../response/response.go`) - Structured response formats

## Error Types and Codes

### Validation Errors (400 Bad Request)
- `VALIDATION_ERROR` - General validation failure
- `INVALID_INPUT` - Invalid input format or type
- `MISSING_REQUIRED_DATA` - Required field is missing
- `INVALID_FORMAT` - Field format is incorrect

### Not Found Errors (404 Not Found)
- `NOT_FOUND` - General resource not found
- `PRODUCT_NOT_FOUND` - Specific product not found
- `USER_NOT_FOUND` - Specific user not found

### Business Logic Errors (422 Unprocessable Entity / 409 Conflict)
- `BUSINESS_LOGIC_ERROR` - General business logic violation
- `INSUFFICIENT_STOCK` - Not enough stock available (409 Conflict)
- `DUPLICATE_ENTRY` - Resource already exists (409 Conflict)

### Database Errors (500 Internal Server Error / 503 Service Unavailable)
- `DATABASE_ERROR` - General database operation failure
- `CONNECTION_ERROR` - Database connection issues (503)
- `MIGRATION_ERROR` - Database migration failure

### Internal Server Errors (500 Internal Server Error)
- `INTERNAL_SERVER_ERROR` - General internal error
- `UNKNOWN_ERROR` - Unexpected error

### Authentication/Authorization Errors (401/403)
- `UNAUTHORIZED` - Authentication required (401)
- `FORBIDDEN` - Access denied (403)
- `TOKEN_EXPIRED` - JWT token expired (401)

## Usage in Services

### Creating Custom Errors

```go
package service

import (
    apperrors "github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
)

func (s *service) CreateProduct(product *Product) error {
    // Validation errors
    if product.Name == "" {
        return apperrors.NewMissingRequiredDataError("name")
    }
    
    if product.StockQuantiy < 0 {
        return apperrors.NewInvalidInputError("stock quantity cannot be negative")
    }
    
    // Business logic errors
    if existingProduct, _ := s.repo.GetByName(product.Name); existingProduct != nil {
        return apperrors.NewDuplicateEntryError("name", product.Name)
    }
    
    // Database errors
    if err := s.repo.Create(product); err != nil {
        return apperrors.NewDatabaseError("failed to create product: " + err.Error())
    }
    
    return nil
}

func (s *service) DecrementStock(id string, quantity int) error {
    product, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return apperrors.NewProductNotFoundError(id)
        }
        return apperrors.NewDatabaseError("failed to retrieve product")
    }
    
    // Business logic error with details
    if product.StockQuantiy < quantity {
        return apperrors.NewInsufficientStockError(product.StockQuantiy, quantity)
    }
    
    return nil
}
```

## Usage in Handlers

### Using Error Helper Functions

```go
package handlers

import (
    "github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
)

func (h *ProductHandler) CreateProduct() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var product Product
        
        // Parse request body
        if err := c.BodyParser(&product); err != nil {
            return errors.HandleError(c, errors.NewInvalidInputError("invalid request body format"))
        }
        
        // Call service
        if err := h.service.CreateProduct(&product); err != nil {
            return errors.HandleError(c, err) // Automatically handles all error types
        }
        
        // Success response
        return errors.HandleCreatedSuccess(c, product)
    }
}

func (h *ProductHandler) GetProductByID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        id := c.Params("id")
        if id == "" {
            return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
        }
        
        product, err := h.service.GetProductByID(id)
        if err != nil {
            return errors.HandleError(c, err)
        }
        
        return errors.HandleSuccess(c, product)
    }
}
```

## Response Formats

### Success Response Format

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    "id": "uuid",
    "name": "Product Name",
    "stock_quantity": 10
  }
}
```

### Error Response Format

```json
{
  "success": false,
  "code": "PRODUCT_NOT_FOUND",
  "message": "Product with ID abc-123 not found"
}
```

### Error Response with Details

```json
{
  "success": false,
  "code": "INSUFFICIENT_STOCK",
  "message": "Insufficient stock. Available: 5, Required: 10",
  "details": {
    "available": 5,
    "required": 10
  }
}
```

### Validation Error Response

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": {
    "errors": [
      {
        "field": "name",
        "message": "Name is required",
        "value": ""
      },
      {
        "field": "stock_quantity",
        "message": "Stock quantity must be non-negative",
        "value": -5
      }
    ]
  }
}
```

## Middleware Setup

The error handler middleware is automatically configured in the server setup:

```go
// internal/server/server.go
app := fiber.New(fiber.Config{
    ErrorHandler: errors.ErrorHandler(),
})
```

This middleware automatically:
- Catches all unhandled errors
- Maps custom AppErrors to appropriate HTTP status codes
- Formats error responses consistently
- Logs errors for debugging

## Helper Functions

### Error Creation Functions

- `NewValidationError(message)` - 400 Bad Request
- `NewInvalidInputError(message)` - 400 Bad Request  
- `NewMissingRequiredDataError(field)` - 400 Bad Request
- `NewProductNotFoundError(id)` - 404 Not Found
- `NewInsufficientStockError(available, required)` - 409 Conflict
- `NewDatabaseError(message)` - 500 Internal Server Error
- `NewConnectionError(message)` - 503 Service Unavailable

### Response Helper Functions

- `HandleError(c, err)` - Handle any error type
- `HandleSuccess(c, data, statusCode?)` - Success response
- `HandleCreatedSuccess(c, data)` - 201 Created response
- `HandleNoContent(c)` - 204 No Content response

## Best Practices

1. **Use Specific Error Types**: Always use the most specific error type available
2. **Provide Context**: Include relevant details in error messages
3. **Handle at Service Layer**: Create custom errors in the service layer where business logic resides
4. **Use Helper Functions**: Use the provided helper functions in handlers for consistency
5. **Don't Expose Internal Details**: Avoid exposing internal implementation details in error messages
6. **Log for Debugging**: The middleware automatically logs errors for debugging purposes

## Error Code Mapping

| Error Code | HTTP Status | Use Case |
|------------|-------------|----------|
| `VALIDATION_ERROR` | 400 | Input validation failures |
| `INVALID_INPUT` | 400 | Malformed request data |
| `MISSING_REQUIRED_DATA` | 400 | Required fields missing |
| `PRODUCT_NOT_FOUND` | 404 | Product doesn't exist |
| `INSUFFICIENT_STOCK` | 409 | Not enough inventory |
| `DUPLICATE_ENTRY` | 409 | Resource already exists |
| `DATABASE_ERROR` | 500 | Database operation failed |
| `CONNECTION_ERROR` | 503 | Database connection issues |
| `INTERNAL_SERVER_ERROR` | 500 | Unexpected server errors |

This structured approach ensures consistent error handling across the entire application while providing clear, actionable error messages to API consumers.