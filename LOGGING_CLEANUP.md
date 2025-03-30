# SpurtCMS Logging Cleanup Guide

This document outlines the plan to remove all direct `fmt.Print` statements from the SpurtCMS codebase and replace them with proper structured logging.

## Why Remove fmt.Print Statements?

Using `fmt.Print` statements for logging has several disadvantages:

1. **No log levels**: Cannot filter logs by severity
2. **No structured data**: Cannot easily parse or analyze logs
3. **No control over output**: Always goes to stdout
4. **No timestamps or context**: Missing important metadata
5. **No consistency**: Different formats across the codebase

## Replacement Strategy

All `fmt.Print` statements should be replaced with the appropriate logger method from our improved logging system:

| Current | Replacement | When to Use |
|---------|------------|------------|
| `fmt.Println(err)` | `logger.Error("Operation failed", logger.WithError(err))` | For error messages |
| `fmt.Println("Info message")` | `logger.Info("Info message")` | For informational messages |
| `fmt.Printf("Value: %v", val)` | `logger.Info("Message", logger.WithField("value", val))` | For variable values |
| `fmt.Println("Warning")` | `logger.Warn("Warning message")` | For warning messages |
| `fmt.Printf("Debug: %v", val)` | `logger.Debug("Debug message", logger.WithField("value", val))` | For debug information |

## Files Requiring Updates

The following files contain `fmt.Print` statements that need to be replaced:

### Middleware
- `middleware/authmiddleware.go`
- `middleware/csrfmiddleware.go`

### Storage Controllers
- `storage-controller/aws-s3-storage.go`
- `storage-controller/local-storage.go`
- `storage-controller/storage-handler.go`

### Controllers
- `controllers/authcontroller.go`
- `controllers/categorycontroller.go`
- `controllers/channelscontroller.go`
- `controllers/channelsettingcontroller.go`
- `controllers/common.go`
- `controllers/contentaccesscontrol.go`
- `controllers/datacontroller.go`
- `controllers/emailtemplates.go`
- `controllers/entriescontroller.go`
- `controllers/general-settings.go`
- `controllers/graphql-settings.go`
- `controllers/languagecontroller.go`
- `controllers/media_settingscontroller.go`
- `controllers/membercontroller.go`
- `controllers/membergroupcontroller.go`
- `controllers/menucontroller.go`
- `controllers/rolespermissioncontroller.go`
- `controllers/settingscontroller.go`
- `controllers/templates.go`
- `controllers/usercontroller.go`

### GraphQL
- `graphql/controller/channels.go`
- `graphql/controller/common.go`
- `graphql/logger/logger.go`
- `graphql/model/member.go`

### Models
- `models/emailtemplates.go`
- `models/language.go`

### Config and Migration
- `config/dbconfig.go`
- `migration/migration.go`

### Debug and Tests
- `debug-storage.go`
- `main_test.go`

## Example Replacements

### Error Logging

```go
// Before
if err != nil {
    fmt.Println(err)
    return err
}

// After
if err != nil {
    logger.Error("Operation failed", logger.WithError(err))
    return err
}
```

### Informational Logging

```go
// Before
fmt.Println("Processing user:", username)

// After
logger.Info("Processing user", logger.WithField("username", username))
```

### Debug Logging

```go
// Before
fmt.Printf("Query params: limit=%d, offset=%d\n", limit, offset)

// After
logger.Debug("Query parameters", logger.WithFields(map[string]any{
    "limit": limit,
    "offset": offset,
}))
```

### Warning Logging

```go
// Before
fmt.Println("csrf token mismatch")

// After
logger.Warn("CSRF token mismatch")
```

## Implementation Plan

1. **Update the logger package first**
   - Ensure the improved logger is in place
   - Add backward compatibility for existing code

2. **Update core components**
   - Start with middleware and storage controllers
   - These are critical components that affect the entire application

3. **Update controllers**
   - Focus on one controller at a time
   - Test thoroughly after each controller is updated

4. **Update models and other components**
   - Replace remaining fmt.Print statements
   - Ensure consistent logging patterns

5. **Update tests and debug tools**
   - Replace fmt.Print statements in tests
   - Update debug tools to use the logger

## Testing

After updating each file:

1. Run the application and verify it starts correctly
2. Test the functionality related to the updated file
3. Check the logs to ensure they contain the expected information
4. Verify that no information is lost compared to the previous implementation

## Conclusion

By replacing all `fmt.Print` statements with proper structured logging, we will:

1. Improve the quality and consistency of logs
2. Make it easier to debug issues
3. Enable better monitoring and alerting
4. Follow best practices for production applications
