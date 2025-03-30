# SpurtCMS Logging Guide

This guide explains the improved logging system in SpurtCMS and how to use it effectively in your code.

## Overview

The SpurtCMS logging system has been optimized to provide:

1. Consistent logging across the application
2. Structured logging with contextual information
3. Multiple log levels for better filtering
4. Configurable output formats and destinations
5. HTTP request logging middleware

## Log Levels

The logger supports the following log levels (in order of increasing severity):

1. **Debug**: Detailed information, typically useful only for diagnosing problems
2. **Info**: Confirmation that things are working as expected
3. **Warn**: Indication that something unexpected happened, but the application can continue
4. **Error**: Due to a more serious problem, the application couldn't perform some function
5. **Fatal**: A severe error that causes the application to terminate

## Basic Usage

```go
import "spurt-cms/logger"

// Simple logging
logger.Debug("Debug message")
logger.Info("Information message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal message - will exit the application")

// Logging with context
logger.Info("User logged in", logger.WithField("user_id", "123"))

// Logging with multiple fields
logger.Info("Database query completed", logger.WithFields(map[string]any{
    "query_time_ms": 42,
    "rows_affected": 10,
    "query": "SELECT * FROM users",
}))

// Logging errors
err := someFunction()
if err != nil {
    logger.Error("Failed to execute function", logger.WithError(err))
}
```

## Best Practices

### 1. Use the appropriate log level

- **Debug**: Use for detailed information that is only useful during debugging
- **Info**: Use for general information about system operation
- **Warn**: Use when something unexpected happened but the application can continue
- **Error**: Use when something went wrong and a function couldn't complete
- **Fatal**: Use only when the application cannot continue and must exit

### 2. Include contextual information

Always include relevant context in your log messages:

```go
// Bad
logger.Info("User updated")

// Good
logger.Info("User updated", logger.WithFields(map[string]any{
    "user_id": user.ID,
    "updated_fields": []string{"name", "email"},
    "admin_id": adminUser.ID,
}))
```

### 3. Be consistent with error handling

When handling errors, be consistent:

```go
result, err := someFunction()
if err != nil {
    // Log the error with context
    logger.Error("Failed to execute function", logger.WithFields(map[string]any{
        "function": "someFunction",
        "error": err.Error(),
        "input_params": params,
    }))
    
    // Return or handle the error
    return nil, err
}
```

### 4. Use request IDs in HTTP handlers

In HTTP handlers, include the request ID in logs for traceability:

```go
func MyHandler(c *gin.Context) {
    requestID := logger.GetRequestID(c)
    userID := c.Param("id")
    
    logger.Info("Processing request", logger.WithFields(map[string]any{
        "request_id": requestID,
        "user_id": userID,
        "action": "get_user_details",
    }))
    
    // Process the request...
}
```

### 5. Don't log sensitive information

Never log sensitive information such as:
- Passwords
- API keys
- Personal identifiable information (PII)
- Authentication tokens
- Credit card numbers

### 6. Log at service boundaries

Log at service boundaries (API calls, database queries, external service calls) to help with debugging:

```go
func GetUserFromDatabase(id string) (*User, error) {
    logger.Debug("Fetching user from database", logger.WithField("user_id", id))
    
    // Database query...
    
    if err != nil {
        logger.Error("Database query failed", logger.WithFields(map[string]any{
            "user_id": id,
            "error": err.Error(),
        }))
        return nil, err
    }
    
    logger.Debug("User fetched successfully", logger.WithField("user_id", id))
    return user, nil
}
```

## Configuration

The logger can be configured using environment variables:

- `LOG_LEVEL`: Debug, Info, Warn, Error, Fatal (default: Info)
- `LOG_FORMAT`: Text, JSON (default: Text)
- `LOG_PATH`: Path to log file (default: logs/spurtcms.log)
- `LOG_STDOUT`: true/false - whether to log to stdout (default: true)

## Migrating from Old Code

If you're updating existing code that uses the old logger:

```go
// Old code
errorLog := logger.ErrorLog()
errorLog.Println("This is an error")

// New code
logger.Error("This is an error")

// Old code with formatting
warnLog := logger.WarnLog()
warnLog.Printf("Warning: %s has occurred", "something")

// New code with fields
logger.Warn("Warning: something has occurred", logger.WithField("event", "something"))

// Old code with fmt.Println for errors
if err != nil {
    fmt.Println("Error:", err)
    return err
}

// New code with structured logging
if err != nil {
    logger.Error("Operation failed", logger.WithError(err))
    return err
}
```

## HTTP Request Logging

The logger includes middleware for HTTP request logging:

```go
r := gin.New()
r.Use(logger.RequestLogger())
```

This will automatically log all HTTP requests with:
- Request ID
- Method
- Path
- Status code
- Duration
- Client IP
- User agent

## Conclusion

By following these guidelines, you'll help maintain a consistent and useful logging system throughout the SpurtCMS codebase, making it easier to debug issues and monitor application behavior.
