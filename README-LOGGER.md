# SpurtCMS Logger Implementation

## Overview

This document describes the logger implementation in SpurtCMS. We've replaced all `fmt.Println` statements with structured logging using the `logger` package, which provides different log levels and formatting options.

## Logger Features

- **Log Levels**: Debug, Info, Warn, Error, Fatal
- **Structured Logging**: Support for adding fields to log messages
- **Multiple Output Formats**: Text and JSON formats
- **File and Console Output**: Logs can be written to files and/or stdout

## Usage Examples

### Basic Logging

```go
// Info level logging
logger.Info("User logged in successfully")

// Error level logging
logger.Error("Failed to connect to database")

// Debug level logging
logger.Debug("Processing request parameters")

// Warning level logging
logger.Warn("API rate limit approaching threshold")

// Fatal level logging (will exit the application)
logger.Fatal("Critical configuration error")
```

### Structured Logging with Fields

```go
// Log with a single field
logger.Info("User logged in", logger.WithField("user_id", userId))

// Log with multiple fields
logger.Info("Request processed", logger.WithFields(map[string]any{
    "duration_ms": duration,
    "status_code": statusCode,
    "path": requestPath,
}))

// Log with error information
if err != nil {
    logger.Error("Operation failed", logger.WithError(err))
}
```

## Configuration

The logger can be configured through environment variables:

- `LOG_LEVEL`: Sets the minimum log level (DEBUG, INFO, WARN, ERROR, FATAL)
- `LOG_FORMAT`: Sets the output format (TEXT, JSON)
- `LOG_PATH`: Sets the log file path
- `LOG_STDOUT`: Controls whether logs are also written to stdout (true/false)

## Best Practices

1. Use appropriate log levels:
   - `Debug`: Detailed information useful during development
   - `Info`: General operational information
   - `Warn`: Warning conditions that don't affect normal operation
   - `Error`: Error conditions that affect specific operations
   - `Fatal`: Critical errors that prevent the application from running

2. Include relevant context in log messages:
   - Use structured fields for machine-readable data
   - Make message text human-readable and descriptive

3. Log at the right points:
   - Log at entry/exit points of important functions
   - Log errors where they occur
   - Log significant state changes

4. Don't log sensitive information:
   - Avoid logging passwords, tokens, or personal data
   - Mask sensitive fields when necessary

## Implementation Details

The logger implementation is in the `logger` package and provides both global functions and a logger instance that can be customized. The implementation supports backward compatibility with the standard Go logger.
