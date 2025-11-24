package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description string     `json:"description" gorm:"column:description"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (*TodoItem) TableName() string {
	return "todo_items"
}

func CreateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoItem TodoItem
		//
		if err := c.ShouldBind(&todoItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		if err := db.Create(&todoItem).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"data":   todoItem.Id,
		})
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// Get the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting *sql.DB from GORM: ", err)
	}

	// Set connection pool options
	sqlDB.SetMaxIdleConns(5)            // Max idle connections
	sqlDB.SetMaxOpenConns(10)           // Max open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Max connection lifetime

	log.Println("Connected to database", db)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
		}
	}

	r.Run(":8080")
}
