package examples

import (
	"bytes"
	"mime/multipart"
	"os"
	"spurt-cms/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

// Example of how to update the AWS S3 storage controller to use proper logging
// Based on the existing storage-controller/aws-s3-storage.go

// UploadFileToS3 uploads a file to AWS S3
func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, filePath string) (string, error) {
	// Get AWS credentials from environment
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_DEFAULT_REGION")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	
	// Log the operation with context
	logger.Debug("Preparing to upload file to S3", logger.WithFields(map[string]any{
		"filename": fileHeader.Filename,
		"size": fileHeader.Size,
		"path": filePath,
		"bucket": bucketName,
		"region": region,
	}))
	
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"",
		),
	})
	
	if err != nil {
		logger.Error("Failed to create AWS session", logger.WithError(err))
		return "", err
	}
	
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	
	// Read file content
	fileBytes := make([]byte, fileHeader.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		logger.Error("Failed to read file content", logger.WithError(err))
		return "", err
	}
	
	// Reset file pointer
	file.Seek(0, 0)
	
	// Set up the upload parameters
	uploadParams := &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filePath + fileHeader.Filename),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	}
	
	// Upload the file to S3
	result, err := uploader.Upload(uploadParams)
	if err != nil {
		logger.Error("Failed to upload file to S3", logger.WithFields(map[string]any{
			"filename": fileHeader.Filename,
			"error": err.Error(),
		}))
		return "", err
	}
	
	// Log successful upload
	logger.Info("File uploaded to S3 successfully", logger.WithFields(map[string]any{
		"filename": fileHeader.Filename,
		"location": result.Location,
		"size": fileHeader.Size,
	}))
	
	return result.Location, nil
}

// DeleteFileFromS3 deletes a file from AWS S3
func DeleteFileFromS3(filePath string) error {
	// Get AWS credentials from environment
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_DEFAULT_REGION")
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	
	// Log the operation
	logger.Debug("Preparing to delete file from S3", logger.WithFields(map[string]any{
		"path": filePath,
		"bucket": bucketName,
	}))
	
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"",
		),
	})
	
	if err != nil {
		logger.Error("Failed to create AWS session", logger.WithError(err))
		return err
	}
	
	// Create S3 service client
	svc := s3.New(sess)
	
	// Set up the delete parameters
	deleteParams := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	}
	
	// Delete the file from S3
	_, err = svc.DeleteObject(deleteParams)
	if err != nil {
		logger.Error("Failed to delete file from S3", logger.WithFields(map[string]any{
			"path": filePath,
			"error": err.Error(),
		}))
		return err
	}
	
	// Log successful deletion
	logger.Info("File deleted from S3 successfully", logger.WithField("path", filePath))
	
	return nil
}

// Example of how to update a controller that uses storage
func UploadProfileImage(c *gin.Context) {
	// Get the file from the request
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		logger.Error("Failed to get file from request", logger.WithError(err))
		c.JSON(400, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()
	
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("User ID not found in context")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Log the upload attempt
	logger.Info("Profile image upload requested", logger.WithFields(map[string]any{
		"user_id": userID,
		"filename": fileHeader.Filename,
		"size": fileHeader.Size,
	}))
	
	// Upload the file to S3
	filePath := "users/" + userID.(string) + "/profile/"
	fileURL, err := UploadFileToS3(file, fileHeader, filePath)
	if err != nil {
		logger.Error("Failed to upload profile image", logger.WithFields(map[string]any{
			"user_id": userID,
			"error": err.Error(),
		}))
		c.JSON(500, gin.H{"error": "Failed to upload profile image"})
		return
	}
	
	// Update user profile in database
	// ...
	
	// Log successful upload
	logger.Info("Profile image uploaded successfully", logger.WithFields(map[string]any{
		"user_id": userID,
		"file_url": fileURL,
	}))
	
	c.JSON(200, gin.H{"url": fileURL})
}
