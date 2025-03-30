package logger

// This file contains examples of how to use the logger in different scenarios.
// It's meant to be used as a reference and is not part of the actual implementation.

/*
Example 1: Basic logging
-----------------------

import "spurt-cms/logger"

func someFunction() {
    // Simple logging at different levels
    logger.Debug("Debug message - only shown when debug level is enabled")
    logger.Info("Information message - general information about system operation")
    logger.Warn("Warning message - something might be wrong but operation continues")
    logger.Error("Error message - something went wrong")
    
    // Fatal will log and then exit the application with status code 1
    // logger.Fatal("Fatal message - application cannot continue")
}

Example 2: Logging with context
------------------------------

import "spurt-cms/logger"

func processUser(userID string) {
    // Add context to your logs with fields
    logger.Debug("Processing user", logger.WithField("user_id", userID))
    
    // Multiple fields
    logger.Info("User details retrieved", logger.WithFields(map[string]any{
        "user_id": userID,
        "status": "active",
        "last_login": "2023-01-01T12:00:00Z",
    }))
}

Example 3: Error handling
------------------------

import (
    "errors"
    "spurt-cms/logger"
)

func doSomething() error {
    err := someOperation()
    if err != nil {
        // Log the error with context and return it
        logger.Error("Operation failed", logger.WithError(err))
        return err
    }
    
    logger.Info("Operation succeeded")
    return nil
}

Example 4: Using in HTTP handlers
-------------------------------

import (
    "spurt-cms/logger"
    "github.com/gin-gonic/gin"
)

func UserHandler(c *gin.Context) {
    userID := c.Param("id")
    requestID := logger.GetRequestID(c)
    
    logger.Debug("Processing user request", logger.WithFields(map[string]any{
        "user_id": userID,
        "request_id": requestID,
    }))
    
    // Do something...
    
    if err != nil {
        logger.Error("Failed to process user", logger.WithFields(map[string]any{
            "user_id": userID,
            "request_id": requestID,
            "error": err.Error(),
        }))
        c.JSON(500, gin.H{"error": "Failed to process user"})
        return
    }
    
    logger.Info("User request completed", logger.WithFields(map[string]any{
        "user_id": userID,
        "request_id": requestID,
    }))
    
    c.JSON(200, gin.H{"status": "success"})
}

Example 5: Migrating from old logger
----------------------------------

// Old code using the previous logger
import (
    "log"
    "spurt-cms/logger"
)

func oldFunction() {
    // Old way
    errorLog := logger.ErrorLog()
    errorLog.Println("This is an error")
    
    // New way
    logger.Error("This is an error")
    
    // Old way with formatting
    warnLog := logger.WarnLog()
    warnLog.Printf("Warning: %s has occurred", "something")
    
    // New way with fields
    logger.Warn("Warning: something has occurred", logger.WithField("event", "something"))
}
*/
