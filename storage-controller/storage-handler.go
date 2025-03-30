package storagecontroller

import (
	"fmt"
	"spurt-cms/logger"
)

// UploadImage handles file uploads based on the selected storage type
func UploadImage(fileName string, filePath string, imageByte []byte) error {
	// Get the selected storage type
	storageType, err := GetSelectedType()
	if err != nil {
		return fmt.Errorf("failed to get storage type: %v", err)
	}

	// Log the selected storage type for debugging
	logger.Info(fmt.Sprintf("Using storage type: %s\n", storageType.SelectedType))

	// Use the appropriate upload function based on the selected storage type
	switch storageType.SelectedType {
	case "local":
		// For local storage, use the local storage implementation
		logger.Info("Using local storage implementation")
		return UploadCropImageLocal(fileName, filePath, imageByte)
	case "aws":
		// For AWS S3 storage
		logger.Info("Using AWS S3 storage implementation")
		return UploadCropImageS3(fileName, filePath, imageByte)
	case "azure":
		// For Azure storage (if implemented)
		return fmt.Errorf("azure storage not implemented yet")
	default:
		// Default to local storage if type is not recognized
		logger.Info(fmt.Sprintf("Unknown storage type '%s', defaulting to local storage\n", storageType.SelectedType))
		return UploadCropImageLocal(fileName, filePath, imageByte)
	}
}
