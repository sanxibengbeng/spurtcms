package storagecontroller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"spurt-cms/logger"
	"spurt-cms/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrorLog *log.Logger

var WarnLog *log.Logger

func init() {

	ErrorLog = logger.ErrorLog()

	WarnLog = logger.WarnLog()

}

// var category categories.Category

// var channelAuth channels.Channel

type Medias struct {
	File          bool
	AliaseName    string
	Name          string
	Path          string
	ModTime       time.Time
	TotalSubMedia int
}

/*local folder creation*/
func LocalStorageCreation() {

	storagetype, err := GetSelectedType()

	if err != nil {
		//need to handle
		logger.Error("Error occurred", logger.WithError(err))
	}

	if storagetype.Local != "" {

		if _, folerr := os.Stat(storagetype.Local); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local, os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

		if _, folerr := os.Stat(storagetype.Local + "/media"); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local+"/media", os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

		if _, folerr := os.Stat(storagetype.Local + "/entry"); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local+"/entry", os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

		if _, folerr := os.Stat(storagetype.Local + "/member"); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local+"/member", os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

		if _, folerr := os.Stat(storagetype.Local + "/pages"); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local+"/pages", os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

		if _, folerr := os.Stat(storagetype.Local + "/user"); os.IsNotExist(folerr) {
			if err := os.Mkdir(storagetype.Local+"/user", os.ModePerm); err != nil {
				WarnLog.Println(err)
			}
		}

	}

}

func MediaLocalList(search, folderpath string) ([]Medias, []Medias, error) {

	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Error occurred", logger.WithError(err))
	}

	var Path string
	if storagetype.Local != "" {
		Path = storagetype.Local + "/media/" + folderpath
	}

	entries, err := os.ReadDir(Path)
	if err != nil {
		log.Println(err)
		return []Medias{}, []Medias{}, err
	}

	var (
		Folder []Medias
		File   []Medias
	)

	if search != "" {

		for _, e := range entries {

			fileInfo, _ := e.Info()

			if strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(search)) {

				var med Medias
				med.File = fileInfo.IsDir()
				med.Name = fileInfo.Name()
				med.AliaseName = fileInfo.Name()
				med.Path = "/" + Path
				med.ModTime = fileInfo.ModTime()

				if fileInfo.IsDir() {
					submedia, err := os.ReadDir(Path + fileInfo.Name())
					if err != nil {
						log.Println(err)
					}

					med.TotalSubMedia = len(submedia)
					Folder = append(Folder, med)

				} else {
					File = append(File, med)
				}
			}

		}

	} else {

		for _, e := range entries {

			var med Medias
			fileInfo, _ := e.Info()
			med.File = fileInfo.IsDir()
			med.Name = fileInfo.Name()
			med.AliaseName = fileInfo.Name()
			med.Path = "/" + Path
			med.ModTime = fileInfo.ModTime()

			if fileInfo.IsDir() {
				submedia, err := os.ReadDir(Path + fileInfo.Name())
				if err != nil {
					log.Println(err)
				}

				med.TotalSubMedia = len(submedia)
				Folder = append(Folder, med)

			} else {
				File = append(File, med)
			}

		}
	}

	return Folder, File, nil
}

/*Make Add folder for Media library */
func AddFolderMakeDir(name string, folderpath string) error {

	storagetype, err := GetSelectedType()
	if err != nil {
		logger.Error("Error occurred", logger.WithError(err))
	}

	if name != "" {
		Path := storagetype.Local + "/media/"
		if err := os.Mkdir(Path+folderpath+name, os.ModePerm); err != nil {
			return err
		}
		return nil

	}
	return errors.New("foldername is empty can't create")
}

/**/
func UploadImageLocal(file multipart.File, fileHeader *multipart.FileHeader, filePath string, c *gin.Context) error {

	storagetype, serr := GetSelectedType()
	if serr != nil {
		logger.Info(fmt.Sprintf("%v", serr))
	}

	pathEnv := storagetype.Local + "/media/"
	filename := strings.ReplaceAll(fileHeader.Filename, "%", "")
	splitArr := strings.Split(filename, ".")
	ext := splitArr[len(splitArr)-1]
	nameWithoutExt := strings.ReplaceAll(filename, "."+ext, "")

	if len(nameWithoutExt) == 0 {
		return nil
	}

	var err error

	// You can customize the file storage path and name as per your requirement.
	if filePath != "" {
		err = c.SaveUploadedFile(fileHeader, pathEnv+filePath+filename)
	} else {
		err = c.SaveUploadedFile(fileHeader, pathEnv+filename)
	}

	return err
}

func UploadCropImage(imageData, imagename string) (path string) {

	storagetype, serr := GetSelectedType()
	if serr != nil {
		logger.Info(fmt.Sprintf("%v", serr))
	}

	Path := storagetype.Local + "/media/"
	if imageData != "" {
		_, _, err := ConvertBase64WithName(imageData, Path, imagename)
		if err != nil {
			return
		}

	}

	return Path
}

func ConvertBase64WithName(imageData string, storagepath string, imagename string) (imgname string, path string, err error) {

	// extEndIndex := strings.Index(imageData, ";base64,")
	base64data := imageData[strings.IndexByte(imageData, ',')+1:]
	// var ext = imageData[11:extEndIndex]
	// rand_num := strconv.Itoa(int(time.Now().Unix()))
	// imageName := "IMG-" + rand_num + "." + ext
	imageName := imagename
	os.MkdirAll(storagepath, 0755)
	storagePath := storagepath + imagename
	decode, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {
		logger.Error("Error occurred", logger.WithError(err))
	}
	file, err := os.Create(storagePath)
	if err != nil {
		logger.Error("Error occurred", logger.WithError(err))
	}
	if _, err := file.Write(decode); err != nil {
		logger.Error("Error occurred", logger.WithError(err))
	}

	return imageName, storagePath, err
}

/*Delete*/
func DeleteImageFolder(folderpath, name string) error {

	// sp.Authority = &AUTH

	// channelAuth.Authority = AUTH

	// category.Authority = AUTH

	storagetype, serr := GetSelectedType()

	if serr != nil {

		logger.Info(fmt.Sprintf("%v", serr))

	}

	Path := storagetype.Local + "/media/"

	// sp.RemoveSpaceImage("/" + Path + folderpath + val)

	// channelAuth.RemoveEntriesCoverImage("/" + Path + folderpath + name)

	// category.UpdateImagePath("/" + Path + folderpath + name)

	models.RemoveLanguageImagePath("/"+Path+folderpath+name, TenantId)
	err := os.RemoveAll(Path + folderpath + name)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func FolderDetails(folderpath string) (folderCount int, fileCount int, Folder []Medias, File []Medias, err error) {

	storagetype, serr := GetSelectedType()
	if serr != nil {
		logger.Info(fmt.Sprintf("%v", serr))
	}

	Path := storagetype.Local + "/media/"
	entries, err := os.ReadDir(Path + folderpath)
	if err != nil {
		logger.Info(fmt.Sprintf("FolderDetails Error : %s/n", err))
	}

	for _, e := range entries {

		var med Medias
		fileInfo, _ := e.Info()
		med.File = fileInfo.IsDir()
		med.Name = fileInfo.Name()
		med.AliaseName = fileInfo.Name()
		med.Path = Path + folderpath
		med.ModTime = fileInfo.ModTime()

		if fileInfo.IsDir() {
			submedia, err := os.ReadDir(Path + folderpath + "/" + fileInfo.Name())
			if err != nil {
				log.Println(err)
			}

			med.TotalSubMedia = len(submedia)
			Folder = append(Folder, med)
			folderCount++

		} else {

			File = append(File, med)
			fileCount++

		}

	}

	return folderCount, fileCount, Folder, File, nil
}
func UploadCropImageLocal(fileName string, filePath string, imageByte []byte) error {
	storagetype, serr := GetSelectedType()
	if serr != nil {
		return fmt.Errorf("failed to get storage type: %v", serr)
	}

	if storagetype.Local == "" {
		return fmt.Errorf("local storage path is not configured in database")
	}

	Path := storagetype.Local + "/media/"
	
	// Create directories if they don't exist
	dirPath := Path + strings.TrimSuffix(filePath, fileName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
	}
	
	fullPath := Path + filePath
	
	// Check if directory is writable
	if _, err := os.Stat(dirPath); err != nil {
		return fmt.Errorf("directory access error: %v", err)
	}
	
	// Write the file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fullPath, err)
	}
	defer file.Close()
	
	if _, err := file.Write(imageByte); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", fullPath, err)
	}
	
	return nil
}
