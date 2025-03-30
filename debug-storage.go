// This file is for debugging storage configuration
// It's not meant to be built with the main application
// To use it, run: go run debug-storage.go

// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"spurt-cms/logger"
)

// This is a debugging tool to check your storage configuration
func debugStorage() {
	// Initialize database connection
	// Note: This function is commented out as it's causing build conflicts
	// err := models.InitDB()
	// if err != nil {
	//     log.Fatalf("Failed to initialize database: %v", err)
	// }

	// // Get storage configuration
	// storageType, err := models.GetStorageValue(1)
	// if err != nil {
	//     log.Fatalf("Failed to get storage configuration: %v", err)
	// }

	logger.Info("=== Storage Configuration ===")
	// logger.Info(fmt.Sprintf("Selected Type: %s\n", storageType.SelectedType))
	
	// Example debug code
	logger.Info("This is a debug utility for storage configuration")
	logger.Info("To use it properly, uncomment the code in debug-storage.go")
	
	// Check environment variables
	logger.Info("\n=== Environment Variables ===")
	logger.Info(fmt.Sprintf("BASE_URL: %s\n", os.Getenv("BASE_URL")))
	logger.Info(fmt.Sprintf("TENANT_ID: %s\n", os.Getenv("TENANT_ID")))
	
	// AWS environment variables
	logger.Info("\nAWS Environment Variables:")
	logger.Info(fmt.Sprintf("  AWS_ACCESS_KEY_ID: %s\n", maskString(os.Getenv("AWS_ACCESS_KEY_ID"))))
	logger.Info(fmt.Sprintf("  AWS_SECRET_ACCESS_KEY: %s\n", maskString(os.Getenv("AWS_SECRET_ACCESS_KEY"))))
	logger.Info(fmt.Sprintf("  AWS_DEFAULT_REGION: %s\n", os.Getenv("AWS_DEFAULT_REGION")))
	logger.Info(fmt.Sprintf("  AWS_BUCKET: %s\n", os.Getenv("AWS_BUCKET")))
}

// Mask sensitive information
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****"
}
