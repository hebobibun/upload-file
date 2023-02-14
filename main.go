package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/uploaded", "./uploaded/")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	r.POST("/upload", func(c *gin.Context) {
		// call the UploadFile function
		err := UploadFile(c)
		if err != nil {
			// display an error message if the file upload failed
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			// display a success message if the file upload succeeded
			c.JSON(http.StatusOK, gin.H{
				"message": "File uploaded successfully!",
			})
		}
	})

	r.Run(":8080")
}

func UploadFile(c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return err
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid file type",
		})
		return err
	}

	// Create directory if it doesn't exist
	dir := "./uploaded/"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return err
	}

	// get original file name
	filename := file.Filename

	// replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// generate timestamp
	t := time.Now().Format("2006-01-02_15-04-05")

	// Generate a unique filename
	newName := filename[:len(filename)-len(ext)] + "-" + t + ext

	// Save the file to disk
	err = c.SaveUploadedFile(file, filepath.Join(dir, newName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": newName,
		"message":  "sukses menambahkan file",
	})

	return nil
}
