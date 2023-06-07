package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestListPosts(t *testing.T) {
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

func TestGetPost(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในการเชื่อมต่อกับฐานข้อมูลที่ใช้ในการทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Posts
	err = db.AutoMigrate(&models.Posts{})
	assert.NoError(t, err)

	// สร้าง PostHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	postHandler := handlers.NewPostHandler(db)

	// เพิ่มข้อมูลโพสต์ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	post := models.Posts{
		PostID:    1,
		Title:     "Test Post",
		Body:      "This is a test post",
		UserID:    1,
		Status:    "published",
		CreatedAt: time.Now(),
	}
	err = db.Create(&post).Error
	assert.NoError(t, err)

	// เตรียม HTTP request สำหรับการเรียกใช้งาน GetPost
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/posts/1", nil)

	// เรียกใช้งาน GetPost ผ่าน PostHandler
	postHandler.GetPost(c)

	// ตรวจสอบว่าการดึงข้อมูลโพสต์สำเร็จโดยตรวจสอบสถานะ HTTP response code และแปลง JSON response เป็น CreatePostResponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.CreatePostResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// ตรวจสอบค่าข้อมูลในโพสต์ว่าถูกดึงมาถูกต้องหรือไม่
	assert.Equal(t, post.PostID, response.PostID)
	assert.Equal(t, post.Title, response.Title)
	assert.Equal(t, post.Body, response.Body)
	assert.Equal(t, post.UserID, response.UserID)
	assert.Equal(t, post.Status, response.Status)
	assert.Equal(t, post.CreatedAt, response.CreatedAt)
}

func TestCreatePost(t *testing.T) {
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Posts
	err = db.AutoMigrate(&models.Posts{})
	assert.NoError(t, err)

	// สร้าง PostHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	postHandler := handlers.NewPostHandler(db)

	// เพิ่มข้อมูลโพสต์ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	createPostJson := handlers.CreatePostRequest{
		Title:  "Test Post",
		Body:   "This is a test post",
		UserID: 1,
		Status: "published",
	}

	createPostJSON, _ := json.Marshal(createPostJson)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน CreateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/posts", bytes.NewReader(createPostJSON))
	postHandler.CreatePost(c)

	// ตรวจสอบการสร้างผู้ใช้สำเร็จและรับ JSON กลับจากการเรียกใช้งาน
	assert.Equal(t, http.StatusCreated, w.Code)

	var createPostReq handlers.CreatePostResponse
	err = json.Unmarshal(w.Body.Bytes(), &createPostReq)
	assert.NoError(t, err)

	// ตรวจสอบว่าผู้ใช้ถูกสร้างในฐานข้อมูล
	var post models.Posts
	err = db.First(&post, createPostReq.PostID).Error
	assert.NoError(t, err)
	assert.Equal(t, createPostReq.Title, post.Title)
	assert.Equal(t, createPostReq.Body, post.Body)
	assert.Equal(t, createPostReq.UserID, post.UserID)
	assert.Equal(t, createPostReq.Status, post.Status)
}

func TestUpdatePost(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในการเชื่อมต่อกับฐานข้อมูลที่ใช้ในการทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Posts{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	postHandler := handlers.NewPostHandler(db)
	assert.NoError(t, err)
	// เพิ่มข้อมูลผู้ใช้ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	post := models.Posts{
		Title:  "title123",
		Body:   "body123",
		UserID: 1,
		Status: "status456",
	}
	err = db.Create(&post).Error
	assert.NoError(t, err)

	// เตรียม HTTP request สำหรับการเรียกใช้งาน UpdateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/posts/1", bytes.NewReader([]byte(`{"title": "it title","body": "it body","status": "it status"}`)))
	// // เรียกใช้งาน UpdateUser ผ่าน UserHandler
	postHandler.UpdatePost(c)
	// ตรวจสอบว่าการอัปเดตข้อมูลผู้ใช้สำเร็จโดยตรวจสอบสถานะ HTTP response code และแปลง JSON response เป็น CreateUserResponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.CreatePostResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// ตรวจสอบค่าข้อมูลผู้ใช้หลังการอัปเดตว่าถูกต้องหรือไม่
	assert.Equal(t, uint(1), response.PostID)
	assert.Equal(t, "it title", response.Title)
	assert.Equal(t, "it body", response.Body)
	assert.Equal(t, "it status", response.Status)
}

func TestDeletePost(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในหน่วยทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Posts{})
	assert.NoError(t, err)

	// เตรียมข้อมูลผู้ใช้ในฐานข้อมูลทดสอบ
	post := models.Posts{
		Title:  "title123",
		Body:   "body123",
		UserID: 1,
		Status: "status456",
	}

	err = db.Create(&post).Error
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	postHandler := handlers.NewPostHandler(db)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน DeleteUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: strconv.FormatUint(uint64(post.PostID), 10)})
	postHandler.DeletePost(c)

	// ตรวจสอบการลบผู้ใช้สำเร็จ
	assert.Equal(t, http.StatusOK, w.Code)

	// ตรวจสอบว่าผู้ใช้ถูกลบออกจากฐานข้อมูล
	var deletedPost models.Posts
	err = db.First(&deletedPost, post.PostID).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
