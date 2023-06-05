package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPostHandler_ListPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// สร้างฐานข้อมูลและเชื่อมต่อ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Posts{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	postHandler := handlers.NewPostHandler(db)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน CreateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// สร้างข้อมูล JSON สำหรับการสร้างผู้ใช้ใหม่
	createPostReq := handlers.CreatePostRequest{
		Title:  "title_test",
		Body:   "unitTest",
		UserID: 1,
		Status: "unit@test.com",
	}

	createPostJSON, _ := json.Marshal(createPostReq)

	c.Request, _ = http.NewRequest("POST", "/posts", bytes.NewReader(createPostJSON))
	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createPostRes handlers.CreatePostResponse
	err = json.Unmarshal(w.Body.Bytes(), &createPostRes)
	assert.NoError(t, err)

	// ตรวจสอบว่าผู้ใช้ถูกสร้างในฐานข้อมูล
	var post models.Posts
	err = db.First(&post, createPostRes.PostID).Error
	assert.NoError(t, err)
	assert.Equal(t, createPostReq.Title, post.Title)
	assert.Equal(t, createPostReq.Body, post.Body)
	assert.Equal(t, createPostReq.UserID, post.UserID)
	assert.Equal(t, createPostReq.Status, post.Status)

}
