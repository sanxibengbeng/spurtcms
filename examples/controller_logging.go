package examples

import (
	"net/http"
	"spurt-cms/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Example of how to update a controller to use proper logging
// Based on patterns seen in controllers/entriescontroller.go

// Entry represents a content entry
type Entry struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	ChannelID int    `json:"channel_id"`
	AuthorID  int    `json:"author_id"`
}

// GetEntryList retrieves a list of entries
func GetEntryList(c *gin.Context) {
	// Get request parameters
	channelID := c.Query("channel_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	
	// Parse parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logger.Warn("Invalid page parameter", logger.WithField("page", pageStr))
		page = 1
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Warn("Invalid limit parameter", logger.WithField("limit", limitStr))
		limit = 10
	}
	
	// Calculate offset
	offset := (page - 1) * limit
	
	// Log the request with context
	logger.Info("Fetching entry list", logger.WithFields(map[string]any{
		"channel_id": channelID,
		"page": page,
		"limit": limit,
		"offset": offset,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Fetch entries from database
	// This is a placeholder for the actual database query
	entries := []Entry{
		{ID: 1, Title: "Example Entry 1", Status: "published"},
		{ID: 2, Title: "Example Entry 2", Status: "draft"},
	}
	totalCount := 2
	
	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	
	// Log the result
	logger.Debug("Entry list fetched", logger.WithFields(map[string]any{
		"total_entries": totalCount,
		"total_pages": totalPages,
		"current_page": page,
		"entries_count": len(entries),
	}))
	
	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
		"pagination": gin.H{
			"total": totalCount,
			"page": page,
			"limit": limit,
			"pages": totalPages,
		},
	})
}

// GetEntryByID retrieves a single entry by ID
func GetEntryByID(c *gin.Context) {
	// Get entry ID from path parameter
	entryIDStr := c.Param("id")
	entryID, err := strconv.Atoi(entryIDStr)
	if err != nil {
		logger.Error("Invalid entry ID", logger.WithFields(map[string]any{
			"id": entryIDStr,
			"error": err.Error(),
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
		return
	}
	
	// Log the request
	logger.Info("Fetching entry details", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Fetch entry from database
	// This is a placeholder for the actual database query
	entry := Entry{
		ID:        entryID,
		Title:     "Example Entry",
		Content:   "This is the content of the entry.",
		Status:    "published",
		ChannelID: 1,
		AuthorID:  1,
	}
	
	// Check if entry exists
	if entry.ID == 0 {
		logger.Warn("Entry not found", logger.WithFields(map[string]any{
			"entry_id": entryID,
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}
	
	// Log successful retrieval
	logger.Info("Entry details fetched successfully", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"title": entry.Title,
		"status": entry.Status,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Return the response
	c.JSON(http.StatusOK, gin.H{"entry": entry})
}

// CreateEntry creates a new entry
func CreateEntry(c *gin.Context) {
	// Parse request body
	var entry Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		logger.Error("Failed to parse request body", logger.WithFields(map[string]any{
			"error": err.Error(),
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// Validate entry
	if entry.Title == "" {
		logger.Warn("Missing required field", logger.WithFields(map[string]any{
			"field": "title",
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("User ID not found in context", logger.WithField("request_id", logger.GetRequestID(c)))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Set author ID
	entry.AuthorID = userID.(int)
	
	// Log the create operation
	logger.Info("Creating new entry", logger.WithFields(map[string]any{
		"title": entry.Title,
		"channel_id": entry.ChannelID,
		"author_id": entry.AuthorID,
		"status": entry.Status,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Save entry to database
	// This is a placeholder for the actual database operation
	entry.ID = 123 // Simulated ID from database
	
	// Log successful creation
	logger.Info("Entry created successfully", logger.WithFields(map[string]any{
		"entry_id": entry.ID,
		"title": entry.Title,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Return the response
	c.JSON(http.StatusCreated, gin.H{"entry": entry})
}

// UpdateEntry updates an existing entry
func UpdateEntry(c *gin.Context) {
	// Get entry ID from path parameter
	entryIDStr := c.Param("id")
	entryID, err := strconv.Atoi(entryIDStr)
	if err != nil {
		logger.Error("Invalid entry ID", logger.WithFields(map[string]any{
			"id": entryIDStr,
			"error": err.Error(),
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
		return
	}
	
	// Parse request body
	var entry Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		logger.Error("Failed to parse request body", logger.WithFields(map[string]any{
			"error": err.Error(),
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	// Set entry ID
	entry.ID = entryID
	
	// Log the update operation
	logger.Info("Updating entry", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"title": entry.Title,
		"status": entry.Status,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Update entry in database
	// This is a placeholder for the actual database operation
	
	// Log successful update
	logger.Info("Entry updated successfully", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Return the response
	c.JSON(http.StatusOK, gin.H{"entry": entry})
}

// DeleteEntry deletes an entry
func DeleteEntry(c *gin.Context) {
	// Get entry ID from path parameter
	entryIDStr := c.Param("id")
	entryID, err := strconv.Atoi(entryIDStr)
	if err != nil {
		logger.Error("Invalid entry ID", logger.WithFields(map[string]any{
			"id": entryIDStr,
			"error": err.Error(),
			"request_id": logger.GetRequestID(c),
		}))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
		return
	}
	
	// Log the delete operation
	logger.Info("Deleting entry", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Delete entry from database
	// This is a placeholder for the actual database operation
	
	// Log successful deletion
	logger.Info("Entry deleted successfully", logger.WithFields(map[string]any{
		"entry_id": entryID,
		"request_id": logger.GetRequestID(c),
	}))
	
	// Return the response
	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}
