package controllers

import (
	"os"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// EditorController handles the editor page and configuration
func EditorController(c *gin.Context) {
	// Get the base URL from environment
	baseURL := os.Getenv("BASE_URL")
	
	// Create the upload path configuration
	urlpath := map[string]interface{}{
		"path": baseURL + "uploadb64image", 
		"payload": "imagedata", 
		"success": map[string]interface{}{
			"imagepath": "imagepath", 
			"imagename": "imagename",
		},
	}
	
	// Convert to JSON
	configJSON, _ := json.Marshal(urlpath)
	
	// Render the editor page with the configuration
	c.HTML(200, "editor.html", gin.H{
		"title": "SpurtCMS Editor",
		"uploadConfig": string(configJSON),
	})
}
