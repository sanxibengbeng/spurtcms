# SpurtCMS Logger

This package provides a structured logging solution for SpurtCMS with multiple log levels, formats, and output destinations.

## Features

- Multiple log levels: Debug, Info, Warn, Error, Fatal
- Structured logging with fields
- Text and JSON output formats
- Multiple output destinations (file, stdout)
- Request logging middleware for HTTP requests
- Backward compatibility with existing code

## Usage

### Basic Logging

```go
package main

import "spurt-cms/logger"

func main() {
    // Simple logging
    logger.Debug("This is a debug message")
    logger.Info("This is an info message")
    logger.Warn("This is a warning message")
    logger.Error("This is an error message")
    
    // Logging with fields
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
}
```

### HTTP Request Logging

```go
package main

import (
    "spurt-cms/logger"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()
    
    // Add the request logger middleware
    r.Use(logger.RequestLogger())
    
    r.GET("/ping", func(c *gin.Context) {
        // Get the request ID from the context
        requestID := logger.GetRequestID(c)
        
        // Log with request ID
        logger.Info("Processing ping request", logger.WithField("request_id", requestID))
        
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    r.Run(":8080")
}
```

### Configuration

The logger can be configured using environment variables:

- `LOG_LEVEL`: Debug, Info, Warn, Error, Fatal (default: Info)
- `LOG_FORMAT`: Text, JSON (default: Text)
- `LOG_PATH`: Path to log file (default: logs/spurtcms.log)
- `LOG_STDOUT`: true/false - whether to log to stdout (default: true)

## Backward Compatibility

For backward compatibility with existing code, the following functions are provided:

```go
// Get standard loggers
errorLogger := logger.ErrorLog()
warnLogger := logger.WarnLog()
infoLogger := logger.InfoLog()

// Use them like standard loggers
errorLogger.Println("This is an error")
```
