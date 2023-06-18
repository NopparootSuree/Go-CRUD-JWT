package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/models"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestRegisterUser(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในหน่วยทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	teardownTestDB(db)
	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน CreateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	password, err := utils.HashPassword("password123")
	assert.NoError(t, err)

	// สร้างข้อมูล JSON สำหรับการสร้างผู้ใช้ใหม่
	createUserReq := handlers.CreateUserRequest{
		Username:       "john_doe",
		HashedPassword: password,
		FullName:       "John Doe",
		Email:          "john@example.com",
	}
	createUserJSON, _ := json.Marshal(createUserReq)

	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewReader(createUserJSON))

	var existingUser models.Users
	db.Where("username = ?", createUserReq.Username).Or("email = ?", createUserReq.Email).First(&existingUser)
	if existingUser.ID != 0 {
		assert.Fail(t, "User already exists")
	}

	userHandler.Register(c)

	// ตรวจสอบการสร้างผู้ใช้สำเร็จและรับ JSON กลับจากการเรียกใช้งาน
	assert.Equal(t, http.StatusCreated, w.Code)

	var createUserRes handlers.CreateUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &createUserRes)
	assert.NoError(t, err)

	// ตรวจสอบว่าผู้ใช้ถูกสร้างในฐานข้อมูล
	var user models.Users
	err = db.First(&user, createUserRes.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, createUserReq.Username, user.Username)
	assert.Equal(t, createUserReq.FullName, user.Fullname)
	assert.Equal(t, createUserReq.Email, user.Email)

	// ตรวจสอบการเข้ารหัสพาสเวิร์ดที่ถูกต้อง
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(createUserReq.HashedPassword))
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {
	// เชื่อมต่อฐานข้อมูล
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler พร้อมกำหนดค่าฐานข้อมูล
	userHandler := handlers.NewUserHandler(db)

	// เรียกใช้งานเส้นทางและรับการตอบสนอง
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// สร้างคำขอ HTTP POST ด้วยข้อมูล JSON สำหรับการเข้าสู่ระบบ
	loginReq := handlers.LoginUserRequest{
		Username: "john_doe",
		Password: "password123",
	}
	loginJson, _ := json.Marshal(loginReq)

	c.Request, _ = http.NewRequest("POST", "login", bytes.NewReader(loginJson))
	userHandler.Login(c)

	// ตรวจสอบการตอบสนอง
	assert.Equal(t, http.StatusOK, w.Code)

	// แปลงเนื้อหาของตัวตอบสนองเป็นโครงสร้าง Payload
	var payload struct {
		Payload handlers.Payload `json:"payload"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &payload)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// ตรวจสอบค่าใน Payload
	assert.NotEmpty(t, payload.Payload.Token)
	assert.Equal(t, "admin123", payload.Payload.Username)
	assert.WithinDuration(t, time.Now(), payload.Payload.IssuedAt, time.Second)
	assert.WithinDuration(t, time.Now().Add(time.Hour*3), payload.Payload.ExpiredAt, time.Second)
}
