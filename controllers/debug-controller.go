package controllers

import (
	"fmt"
	"os"
	"spurt-cms/models"
	"spurt-cms/storage-controller"

	"github.com/gin-gonic/gin"
)

// DebugStorageController provides diagnostic information about the storage configuration
func DebugStorageController(c *gin.Context) {
	var result = make(map[string]interface{})
	
	// Get storage configuration
	storageType, err := models.GetStorageValue(TenantId)
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to get storage configuration: %v", err)
		c.JSON(200, result)
		return
	}
	
	// Basic storage info
	result["storage_type"] = storageType.SelectedType
	
	// Check local storage
	if storageType.Local != "" {
		localInfo := make(map[string]interface{})
		localInfo["path"] = storageType.Local
		
		// Check if the directory exists
		localPath := storageType.Local
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			localInfo["directory_exists"] = false
			localInfo["error"] = fmt.Sprintf("Local storage directory '%s' does not exist", localPath)
		} else {
			localInfo["directory_exists"] = true
			
			// Check if media directory exists
			mediaPath := localPath + "/media"
			if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
				localInfo["media_directory_exists"] = false
			} else {
				localInfo["media_directory_exists"] = true
			}
			
			// Check if entries directory exists
			entriesPath := mediaPath + "/entries"
			if _, err := os.Stat(entriesPath); os.IsNotExist(err) {
				localInfo["entries_directory_exists"] = false
			} else {
				localInfo["entries_directory_exists"] = true
			}
			
			// Check write permissions
			testFile := localPath + "/test_write_permission.txt"
			file, err := os.Create(testFile)
			if err != nil {
				localInfo["write_permission"] = false
				localInfo["write_error"] = err.Error()
			} else {
				file.Close()
				os.Remove(testFile)
				localInfo["write_permission"] = true
			}
		}
		
		result["local_storage"] = localInfo
	}
	
	// Check AWS S3
	if storageType.Aws != nil {
		awsInfo := make(map[string]interface{})
		
		// Check environment variables
		awsInfo["aws_access_key_id_set"] = os.Getenv("AWS_ACCESS_KEY_ID") != ""
		awsInfo["aws_secret_access_key_set"] = os.Getenv("AWS_SECRET_ACCESS_KEY") != ""
		awsInfo["aws_default_region"] = os.Getenv("AWS_DEFAULT_REGION")
		awsInfo["aws_bucket"] = os.Getenv("AWS_BUCKET")
		
		result["aws_s3"] = awsInfo
	}
	
	// Check environment variables
	envInfo := make(map[string]interface{})
	envInfo["base_url"] = os.Getenv("BASE_URL")
	envInfo["tenant_id"] = os.Getenv("TENANT_ID")
	result["environment"] = envInfo
	
	// Generate recommendations
	var recommendations []string
	if storageType.SelectedType == "" {
		recommendations = append(recommendations, "Set a storage type in the database (local, aws, or azure)")
	}
	
	if storageType.SelectedType == "local" && storageType.Local == "" {
		recommendations = append(recommendations, "Configure the local storage path in the database")
	}
	
	if storageType.SelectedType == "aws" {
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || 
		   os.Getenv("AWS_DEFAULT_REGION") == "" || os.Getenv("AWS_BUCKET") == "" {
			recommendations = append(recommendations, "Set all required AWS environment variables")
		}
	}
	
	result["recommendations"] = recommendations
	
	// Test image upload
	if c.Query("test") == "true" {
		// Create a simple test image
		testImage := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
		
		// Try to upload it
		imageName, imagePath, imageByte, err := ConvertBase64toByte(testImage, "test")
		if err != nil {
			result["test_upload"] = map[string]interface{}{
				"success": false,
				"stage": "convert_base64",
				"error": err.Error(),
			}
		} else {
			// Use the storage handler to upload the image
			uerr := storagecontroller.UploadImage(imageName, imagePath, imageByte)
			if uerr != nil {
				result["test_upload"] = map[string]interface{}{
					"success": false,
					"stage": "upload_image",
					"error": uerr.Error(),
				}
			} else {
				result["test_upload"] = map[string]interface{}{
					"success": true,
					"image_name": imageName,
					"image_path": imagePath,
				}
			}
		}
	}
	
	c.JSON(200, result)
}
