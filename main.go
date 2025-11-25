package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
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

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id"`
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	Status      string `json:"status" gorm:"column:status"`
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title"`
	Description *string `json:"description" gorm:"column:description"`
	Status      string  `json:"status" gorm:"column:status"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }

func CreateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoItem TodoItemCreation
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

func GetDetailItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoItem TodoItem

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		if err := db.Where("id=?", id).First(&todoItem).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   todoItem,
			"data2":  "Thời gian phản hồi mà bạn thấy (200 OK, 14 ms) là thời gian mà Go server trả về cho client khi thực hiện một yêu cầu HTTP. Nếu bạn muốn so sánh thời gian phản hồi giữa Go và PHP, thời gian phản hồi sẽ phụ thuộc vào nhiều yếu tố, bao gồm:\n\nCấu hình Server:\n\nNếu Go và PHP chạy trên cùng một máy chủ với cùng một cấu hình, khả năng xử lý của mỗi ngôn ngữ sẽ có sự khác biệt. Go thường có hiệu suất tốt hơn nhờ vào khả năng xử lý đồng thời tốt hơn và việc biên dịch sang mã máy giúp tối ưu hóa tốc độ.\n\nPHP, mặc dù rất phổ biến trong các ứng dụng web, nhưng thường có overhead do phải chạy trong môi trường máy chủ (như Apache hoặc Nginx với PHP-FPM).\n\nTốc độ xử lý:\n\nGo: Go được thiết kế để có tốc độ rất nhanh, đặc biệt là trong các tác vụ xử lý đồng thời và kết nối với cơ sở dữ liệu. Nó có thể xử lý hàng nghìn yêu cầu HTTP đồng thời với ít tài nguyên hơn nhờ vào goroutines.\n\nPHP: PHP không được tối ưu cho các tác vụ đồng thời, mặc dù có thể chạy đồng thời với các công cụ như PHP-FPM (FastCGI Process Manager) hoặc sử dụng worker threads. Tuy nhiên, mỗi yêu cầu PHP thường phải tạo ra một quá trình mới, điều này gây overhead.\n\nMôi trường PHP:\n\nTrong PHP, thời gian phản hồi có thể sẽ cao hơn một chút, đặc biệt nếu PHP không được tối ưu. Một ứng dụng PHP đơn giản với Apache hoặc Nginx + PHP-FPM có thể có thời gian phản hồi từ khoảng 20ms đến 100ms cho các tác vụ cơ bản.\n\nNếu bạn sử dụng các framework PHP như Laravel, Symfony hoặc Zend, thời gian phản hồi có thể sẽ chậm hơn do overhead của các framework này.\n\nTóm lại:\n\nGo thường có hiệu suất vượt trội trong việc xử lý yêu cầu HTTP, và thời gian phản hồi có thể chỉ khoảng 14ms, như bạn đã thấy.\n\nPHP có thể mất khoảng 20ms đến 100ms hoặc lâu hơn, tùy thuộc vào cách cấu hình và tối ưu server, cũng như việc sử dụng các framework.\n\nTuy nhiên, mức độ chênh lệch này không quá lớn đối với các ứng dụng nhỏ hoặc trung bình. Sự khác biệt sẽ rõ rệt hơn khi hệ thống của bạn cần xử lý một lượng lớn yêu cầu đồng thời, vì Go có lợi thế lớn trong khả năng xử lý đồng thời.",
		})
	}
}

func UpdateItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todoItem TodoItemUpdate

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		if err := c.ShouldBind(&todoItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		if err := db.Where("id=?", id).Updates(&todoItem).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   todoItem,
		})
	}
}

func DeleteItem(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		if err := db.Where("id=?", id).Delete(nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
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
			items.GET("/:id", GetDetailItem(db))
			items.PUT("/:id", UpdateItem(db))
		}
	}

	r.Run(":8080")
}
