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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ClearUsersTable() error {
	var db *gorm.DB
	err := db.Delete(&models.Users{}).Error
	if err != nil {
		return err
	}
	return nil
}

func TestCreateUser(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในหน่วยทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน CreateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// สร้างข้อมูล JSON สำหรับการสร้างผู้ใช้ใหม่
	createUserReq := handlers.CreateUserRequest{
		Username:       "john_doe",
		HashedPassword: "password123",
		FullName:       "John Doe",
		Email:          "john@example.com",
	}
	createUserJSON, _ := json.Marshal(createUserReq)

	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewReader(createUserJSON))
	userHandler.CreateUser(c)

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
