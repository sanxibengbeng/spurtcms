package storagecontroller

import (
	"os"
	"strings"

	"spurt-cms/models"
)

// GetStorageConfigFromEnv reads storage configuration from environment variables
// instead of the database
func GetStorageConfigFromEnv() (models.TblStorageType, error) {
	storageType := models.TblStorageType{}

	// Get the selected storage type from environment variable
	selectedType := os.Getenv("STORAGE_TYPE")
	if selectedType == "" {
		selectedType = "local" // Default to local storage if not specified
	}

	// Remove any quotes that might be in the env var
	selectedType = strings.Trim(selectedType, "'\"")
	storageType.SelectedType = selectedType

	// Set local storage path if applicable
	if selectedType == "local" {
		localPath := os.Getenv("STORAGE_LOCAL_PATH")
		if localPath == "" {
			localPath = "./storage" // Default path
		}
		localPath = strings.Trim(localPath, "'\"")
		storageType.Local = localPath
	}

	// Set AWS configuration if applicable
	if selectedType == "aws" {
		awsConfig := make(map[string]interface{})
		awsConfig["accessid"] = strings.Trim(os.Getenv("AWS_ACCESS_KEY_ID"), "'\"")
		awsConfig["accesskey"] = strings.Trim(os.Getenv("AWS_SECRET_ACCESS_KEY"), "'\"")
		awsConfig["region"] = strings.Trim(os.Getenv("AWS_DEFAULT_REGION"), "'\"")
		awsConfig["bucketname"] = strings.Trim(os.Getenv("AWS_BUCKET"), "'\"")
		storageType.Aws = awsConfig
	}

	// Set Azure configuration if applicable
	if selectedType == "azure" {
		azureConfig := make(map[string]interface{})
		azureConfig["storageaccount"] = strings.Trim(os.Getenv("AZURE_STORAGE_ACCOUNT"), "'\"")
		azureConfig["accountkey"] = strings.Trim(os.Getenv("AZURE_ACCOUNT_KEY"), "'\"")
		azureConfig["containername"] = strings.Trim(os.Getenv("AZURE_CONTAINER_NAME"), "'\"")
		storageType.Azure = azureConfig
	}

	return storageType, nil
}
