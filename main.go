package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Greeting struct {
	gorm.Model
	ID       int    `gorm:"primary_key"`
	Greeting string `gorm:"type:varchar(100);"`
}

var db *gorm.DB

func init() {
	godotenv.Load()
}

func main() {
	var err error

	dbUser, dbPassword, dbHost, dbName, dbPort := getDbDetails()
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword))
	if err != nil {
		fmt.Println(dbUser, dbPassword, dbHost, dbName, dbPort)
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Greeting{})

	r := gin.Default()

	r.GET("/greeting/:id", GetGreeting)
	r.POST("/greeting", CreateGreeting)
	r.PUT("/greeting/:id", UpdateGreeting)
	r.DELETE("/greeting/:id", DeleteGreeting)

	r.Run()
}

func getDbDetails() (string, string, string, string, string) {
	return os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")
}

func GetGreeting(c *gin.Context) {
	var greeting Greeting
	id := c.Param("id")

	if err := db.First(&greeting, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, greeting)
}

func CreateGreeting(c *gin.Context) {
	var greeting Greeting
	if err := c.ShouldBindJSON(&greeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&greeting)

	c.JSON(http.StatusOK, greeting)
}

func UpdateGreeting(c *gin.Context) {
	var greeting Greeting
	id := c.Param("id")

	if err := db.First(&greeting, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	if err := c.ShouldBindJSON(&greeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&greeting)

	c.JSON(http.StatusOK, greeting)
}

func DeleteGreeting(c *gin.Context) {
	var greeting Greeting
	id := c.Param("id")

	if err := db.First(&greeting, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&greeting)

	c.JSON(http.StatusOK, gin.H{"data": "Record deleted!"})
}
