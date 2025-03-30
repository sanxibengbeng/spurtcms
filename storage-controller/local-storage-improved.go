package storagecontroller

import (
	"encoding/base64"
	"errors"
	"mime/multipart"
	"os"
	"spurt-cms/logger"
	"spurt-cms/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Medias represents a file or directory in the media library
type MediasImproved struct {
	File          bool
	AliaseName    string
	Name          string
	Path          string
	ModTime       time.Time
	TotalSubMedia int
}

// LocalStorageCreationImproved creates the necessary directories for local storage
func LocalStorageCreationImproved() error {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return errors.New("local storage path is not configured")
	}

	// List of directories to create
	directories := []string{
		storagetype.Local,
		storagetype.Local + "/media",
		storagetype.Local + "/entry",
		storagetype.Local + "/member",
		storagetype.Local + "/pages",
		storagetype.Local + "/user",
	}

	// Create each directory if it doesn't exist
	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				logger.Error("Failed to create directory", logger.WithFields(map[string]any{
					"directory": dir,
					"error":     err.Error(),
				}))
				return err
			}
			logger.Info("Created directory", logger.WithField("directory", dir))
		}
	}

	return nil
}

// MediaLocalListImproved lists files and folders in the media directory
func MediaLocalListImproved(search, folderpath string) ([]MediasImproved, []MediasImproved, error) {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return nil, nil, err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return nil, nil, errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/" + folderpath
	entries, err := os.ReadDir(path)
	if err != nil {
		logger.Error("Failed to read directory", logger.WithFields(map[string]any{
			"path":  path,
			"error": err.Error(),
		}))
		return nil, nil, err
	}

	var folders []MediasImproved
	var files []MediasImproved

	for _, e := range entries {
		fileInfo, err := e.Info()
		if err != nil {
			logger.Warn("Failed to get file info", logger.WithFields(map[string]any{
				"file":  e.Name(),
				"error": err.Error(),
			}))
			continue
		}

		// Skip if searching and name doesn't match
		if search != "" && !strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(search)) {
			continue
		}

		media := MediasImproved{
			File:       !fileInfo.IsDir(),
			Name:       fileInfo.Name(),
			AliaseName: fileInfo.Name(),
			Path:       "/" + path,
			ModTime:    fileInfo.ModTime(),
		}

		if fileInfo.IsDir() {
			// Count items in subdirectory
			submedia, err := os.ReadDir(path + "/" + fileInfo.Name())
			if err != nil {
				logger.Warn("Failed to read subdirectory", logger.WithFields(map[string]any{
					"directory": fileInfo.Name(),
					"error":     err.Error(),
				}))
				media.TotalSubMedia = 0
			} else {
				media.TotalSubMedia = len(submedia)
			}
			folders = append(folders, media)
		} else {
			files = append(files, media)
		}
	}

	return folders, files, nil
}

// AddFolderMakeDirImproved creates a new folder in the media directory
func AddFolderMakeDirImproved(name string, folderpath string) error {
	if name == "" {
		return errors.New("folder name is empty, can't create")
	}

	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/" + folderpath + name
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		logger.Error("Failed to create folder", logger.WithFields(map[string]any{
			"path":  path,
			"error": err.Error(),
		}))
		return err
	}

	logger.Info("Created folder", logger.WithField("path", path))
	return nil
}

// UploadImageLocalImproved uploads a file to the local storage
func UploadImageLocalImproved(file multipart.File, fileHeader *multipart.FileHeader, filePath string, c *gin.Context) error {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return errors.New("local storage path is not configured")
	}

	pathEnv := storagetype.Local + "/media/"
	filename := strings.ReplaceAll(fileHeader.Filename, "%", "")
	splitArr := strings.Split(filename, ".")
	ext := splitArr[len(splitArr)-1]
	nameWithoutExt := strings.ReplaceAll(filename, "."+ext, "")

	if len(nameWithoutExt) == 0 {
		logger.Warn("Invalid filename", logger.WithField("filename", filename))
		return errors.New("invalid filename")
	}

	fullPath := pathEnv
	if filePath != "" {
		fullPath += filePath
	}
	fullPath += filename

	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		logger.Error("Failed to save uploaded file", logger.WithFields(map[string]any{
			"path":  fullPath,
			"error": err.Error(),
		}))
		return err
	}

	logger.Info("File uploaded successfully", logger.WithFields(map[string]any{
		"filename": filename,
		"path":     fullPath,
		"size":     fileHeader.Size,
	}))

	return nil
}

// UploadCropImageImproved uploads a base64 encoded image
func UploadCropImageImproved(imageData, imagename string) (string, error) {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return "", err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return "", errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/"
	if imageData == "" {
		logger.Warn("Empty image data")
		return "", errors.New("empty image data")
	}

	_, storagePath, err := ConvertBase64WithNameImproved(imageData, path, imagename)
	if err != nil {
		logger.Error("Failed to convert base64 image", logger.WithFields(map[string]any{
			"imagename": imagename,
			"error":     err.Error(),
		}))
		return "", err
	}

	logger.Info("Image uploaded successfully", logger.WithField("path", storagePath))
	return path, nil
}

// ConvertBase64WithNameImproved converts a base64 encoded image to a file
func ConvertBase64WithNameImproved(imageData string, storagepath string, imagename string) (string, string, error) {
	// Extract the base64 data
	base64data := imageData[strings.IndexByte(imageData, ',')+1:]
	imageName := imagename
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(storagepath, 0755); err != nil {
		logger.Error("Failed to create directory", logger.WithFields(map[string]any{
			"path":  storagepath,
			"error": err.Error(),
		}))
		return "", "", err
	}
	
	storagePath := storagepath + imagename
	decode, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		logger.Error("Failed to decode base64 data", logger.WithError(err))
		return "", "", err
	}
	
	file, err := os.Create(storagePath)
	if err != nil {
		logger.Error("Failed to create file", logger.WithFields(map[string]any{
			"path":  storagePath,
			"error": err.Error(),
		}))
		return "", "", err
	}
	defer file.Close()
	
	if _, err := file.Write(decode); err != nil {
		logger.Error("Failed to write to file", logger.WithFields(map[string]any{
			"path":  storagePath,
			"error": err.Error(),
		}))
		return "", "", err
	}
	
	logger.Info("Base64 image converted and saved", logger.WithField("path", storagePath))
	return imageName, storagePath, nil
}

// DeleteImageFolderImproved deletes a file or folder
func DeleteImageFolderImproved(folderpath, name string) error {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/"
	fullPath := path + folderpath + name
	
	// Update references in the database
	if err := models.RemoveLanguageImagePath("/"+fullPath, TenantId); err != nil {
		logger.Warn("Failed to update image references in database", logger.WithError(err))
		// Continue with deletion even if reference update fails
	}
	
	// Delete the file or directory
	if err := os.RemoveAll(fullPath); err != nil {
		logger.Error("Failed to delete file or directory", logger.WithFields(map[string]any{
			"path":  fullPath,
			"error": err.Error(),
		}))
		return err
	}
	
	logger.Info("File or directory deleted", logger.WithField("path", fullPath))
	return nil
}

// FolderDetailsImproved gets details about a folder
func FolderDetailsImproved(folderpath string) (int, int, []MediasImproved, []MediasImproved, error) {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return 0, 0, nil, nil, err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return 0, 0, nil, nil, errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/" + folderpath
	entries, err := os.ReadDir(path)
	if err != nil {
		logger.Error("Failed to read directory", logger.WithFields(map[string]any{
			"path":  path,
			"error": err.Error(),
		}))
		return 0, 0, nil, nil, err
	}

	var folders []MediasImproved
	var files []MediasImproved
	folderCount := 0
	fileCount := 0

	for _, e := range entries {
		fileInfo, err := e.Info()
		if err != nil {
			logger.Warn("Failed to get file info", logger.WithFields(map[string]any{
				"file":  e.Name(),
				"error": err.Error(),
			}))
			continue
		}

		media := MediasImproved{
			File:       !fileInfo.IsDir(),
			Name:       fileInfo.Name(),
			AliaseName: fileInfo.Name(),
			Path:       path,
			ModTime:    fileInfo.ModTime(),
		}

		if fileInfo.IsDir() {
			// Count items in subdirectory
			submedia, err := os.ReadDir(path + "/" + fileInfo.Name())
			if err != nil {
				logger.Warn("Failed to read subdirectory", logger.WithFields(map[string]any{
					"directory": fileInfo.Name(),
					"error":     err.Error(),
				}))
				media.TotalSubMedia = 0
			} else {
				media.TotalSubMedia = len(submedia)
			}
			folders = append(folders, media)
			folderCount++
		} else {
			files = append(files, media)
			fileCount++
		}
	}

	return folderCount, fileCount, folders, files, nil
}

// UploadCropImageLocalImproved uploads an image from byte data
func UploadCropImageLocalImproved(fileName string, filePath string, imageByte []byte) error {
	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Failed to get selected storage type", logger.WithError(err))
		return err
	}

	if storagetype.Local == "" {
		logger.Warn("Local storage path is not configured")
		return errors.New("local storage path is not configured")
	}

	path := storagetype.Local + "/media/"
	
	// Create directories if they don't exist
	dirPath := path + strings.TrimSuffix(filePath, fileName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		logger.Error("Failed to create directory", logger.WithFields(map[string]any{
			"path":  dirPath,
			"error": err.Error(),
		}))
		return err
	}
	
	fullPath := path + filePath
	
	// Check if directory is writable
	if _, err := os.Stat(dirPath); err != nil {
		logger.Error("Directory access error", logger.WithFields(map[string]any{
			"path":  dirPath,
			"error": err.Error(),
		}))
		return err
	}
	
	// Write the file
	file, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Failed to create file", logger.WithFields(map[string]any{
			"path":  fullPath,
			"error": err.Error(),
		}))
		return err
	}
	defer file.Close()
	
	if _, err := file.Write(imageByte); err != nil {
		logger.Error("Failed to write to file", logger.WithFields(map[string]any{
			"path":  fullPath,
			"error": err.Error(),
		}))
		return err
	}
	
	logger.Info("Image uploaded successfully", logger.WithFields(map[string]any{
		"filename": fileName,
		"path":     fullPath,
		"size":     len(imageByte),
	}))
	
	return nil
}
