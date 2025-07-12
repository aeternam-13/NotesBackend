package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"notes_api/models"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Note{})
}

func main() {

	initDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:36527", "http://localhost:36527"}, // Flutter dev server URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return true // TEMP: Accept any origin during dev
		},
	}))

	r.POST("/notes", createNote)
	r.GET("/notes", getNotes)
	r.GET("/notes/:id", getNote)
	r.PUT("/notes/:id", updateNote)
	r.DELETE("/notes/:id", deleteNote)

	r.Run(":8080")
}

func createNote(c *gin.Context) {
	fmt.Println("Received POST /notes")
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if note.ID == -1 {
		note.ID = 0
	}

	db.Create(&note)
	c.JSON(http.StatusCreated, note)
}

func getNotes(c *gin.Context) {
	var notes []models.Note
	db.Find(&notes)
	c.JSON(http.StatusOK, notes)
}

func getNote(c *gin.Context) {
	var note models.Note
	if err := db.First(&note, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, note)
}

func updateNote(c *gin.Context) {
	var note models.Note
	if err := db.First(&note, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var input models.Note
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.Title = input.Title
	note.Content = input.Content
	note.Color = input.Color

	db.Save(&note)

	c.JSON(http.StatusOK, note)
}

func deleteNote(c *gin.Context) {
	db.Delete(&models.Note{}, c.Param("id"))
	c.Status(http.StatusNoContent)
}
