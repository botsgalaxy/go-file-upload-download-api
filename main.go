package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/sqlite"
)

const (
	storagePath     = "./uploads"
	maxFileSize     = 300 << 20 // 300 MB
	expirationDelta = 30 * time.Minute
)

var db *sql.DB

func main() {
	// Initialize database
	initDB()

	// Ensure storage directory exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware to handle CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	router.POST("/upload", handleUpload)
	router.GET("/download/:file_id", handleDownload)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "file-storage.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create file table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			path TEXT,
			expiration DATETIME
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Clean up expired files
	go func() {
		for {
			cleanupExpiredFiles()
			time.Sleep(time.Minute * 1) // Check every minute
		}
	}()
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file size
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit"})
		return
	}

	// Generate a unique name for the file
	fileName := generateUniqueFileName(file.Filename)

	// Save the file to disk
	filePath := filepath.Join(storagePath, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate download URL
	fileID := addFileToDB(fileName)

	downloadURL := fmt.Sprintf("/download/%d", fileID)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"url": downloadURL}})
}

func handleDownload(c *gin.Context) {
	fileID := c.Param("file_id")

	// Fetch file path and name from database
	var (
		filePath string
		fileName string
	)
	err := db.QueryRow("SELECT path, name FROM files WHERE id = ?", fileID).Scan(&filePath, &fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file with original name
	c.FileAttachment(filePath, fileName)
}

func generateUniqueFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	base := fileName[:len(fileName)-len(ext)]
	randomDigits, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("%s_%04d%s", base, randomDigits, ext)
}

func addFileToDB(fileName string) int64 {
	expiration := time.Now().Add(expirationDelta)
	result, err := db.Exec("INSERT INTO files (name, path, expiration) VALUES (?, ?, ?)", fileName, filepath.Join(storagePath, fileName), expiration)
	if err != nil {
		log.Printf("Failed to add file to database: %v", err)
		return -1
	}

	fileID, _ := result.LastInsertId()
	return fileID
}

func cleanupExpiredFiles() {
	// Delete expired files from database
	_, err := db.Exec("DELETE FROM files WHERE expiration <= ?", time.Now())
	if err != nil {
		log.Printf("Failed to delete expired files from database: %v", err)
	}

	// Remove expired files from disk
	fileRows, err := db.Query("SELECT path FROM files WHERE expiration <= ?", time.Now())
	if err != nil {
		log.Printf("Failed to query expired files: %v", err)
		return
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var filePath string
		if err := fileRows.Scan(&filePath); err != nil {
			log.Printf("Failed to scan file path: %v", err)
			continue
		}
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove file: %v", err)
		}
	}
}
