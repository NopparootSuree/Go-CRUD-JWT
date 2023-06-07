package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/models"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestListUsers(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในการเชื่อมต่อกับฐานข้อมูลที่ใช้ในการทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	clocks := clock.NewMock()
	clocks.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	clocks.Add(24 * time.Hour)

	// เพิ่มข้อมูลผู้ใช้ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	user1 := models.Users{
		Username:  "user111",
		Fullname:  "User One",
		Email:     "user1@example.com",
		CreatedAt: clocks.Now(),
	}

	user2 := models.Users{
		Username:  "user222",
		Fullname:  "User two",
		Email:     "user2@example.com",
		CreatedAt: clocks.Now(),
	}
	err = db.Create(&user1).Error
	assert.NoError(t, err)
	err = db.Create(&user2).Error
	assert.NoError(t, err)

	// เตรียม HTTP request สำหรับการเรียกใช้งาน ListUsers
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users", nil)

	// เรียกใช้งาน ListUsers ผ่าน UserHandler
	userHandler.ListUsers(c)

	// ตรวจสอบว่าการดึงข้อมูลผู้ใช้สำเร็จโดยตรวจสอบสถานะ HTTP response code และแปลง JSON response เป็น slice ของ CreateUserResponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response []handlers.CreateUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// ตรวจสอบว่าข้อมูลผู้ใช้ถูกดึงมาทั้งหมด
	assert.Len(t, response, 2)

	// ตรวจสอบค่าข้อมูลในแต่ละผู้ใช้ว่าถูกดึงมาถูกต้องหรือไม่
	assert.Equal(t, user1.Username, response[0].Username)
	assert.Equal(t, user1.Fullname, response[0].FullName)
	assert.Equal(t, user1.Email, response[0].Email)
	assert.Equal(t, user1.CreatedAt, response[0].CreatedAt)

	assert.Equal(t, user2.Username, response[1].Username)
	assert.Equal(t, user2.Fullname, response[1].FullName)
	assert.Equal(t, user2.Email, response[1].Email)
	assert.Equal(t, user2.CreatedAt, response[1].CreatedAt)
}
func TestGetUser(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในการเชื่อมต่อกับฐานข้อมูลที่ใช้ในการทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	clocks := clock.NewMock()
	clocks.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local))
	clocks.Add(24 * time.Hour)

	// เพิ่มข้อมูลผู้ใช้ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	user := models.Users{
		ID:        1,
		Username:  "john_doe",
		Fullname:  "John Doe",
		Email:     "john@example.com",
		CreatedAt: clocks.Now(),
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// เตรียม HTTP request สำหรับการเรียกใช้งาน GetUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/1", nil)

	// เรียกใช้งาน GetUser ผ่าน UserHandler
	userHandler.GetUser(c)

	// ตรวจสอบว่าการดึงข้อมูลผู้ใช้สำเร็จโดยตรวจสอบสถานะ HTTP response code และแปลง JSON response เป็น CreateUserResponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.CreateUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// ตรวจสอบค่าข้อมูลผู้ใช้ว่าถูกดึงมาถูกต้องหรือไม่
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Fullname, response.FullName)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
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

func TestUpdateUser(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในการเชื่อมต่อกับฐานข้อมูลที่ใช้ในการทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	clocks := clock.NewMock()
	clocks.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	clocks.Add(24 * time.Hour)

	password, err := utils.HashPassword("password123")
	assert.NoError(t, err)
	// เพิ่มข้อมูลผู้ใช้ในฐานข้อมูลเพื่อใช้ในการทดสอบ
	user := models.Users{
		Username:       "john_doe",
		HashedPassword: password,
		Fullname:       "John Doe",
		Email:          "john@example.com",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// เตรียม HTTP request สำหรับการเรียกใช้งาน UpdateUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/users/1", bytes.NewReader([]byte(`{"hashedPassword": "test123","fullName": "John Smith"}`)))
	fmt.Println(c.Request.Body)
	// // เรียกใช้งาน UpdateUser ผ่าน UserHandler
	userHandler.UpdateUser(c)

	// ตรวจสอบว่าการอัปเดตข้อมูลผู้ใช้สำเร็จโดยตรวจสอบสถานะ HTTP response code และแปลง JSON response เป็น CreateUserResponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.CreateUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// ตรวจสอบค่าข้อมูลผู้ใช้หลังการอัปเดตว่าถูกต้องหรือไม่
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Fullname, response.FullName)
	assert.Equal(t, user.Email, response.Email)
}

func TestDeleteUser(t *testing.T) {
	// เตรียมฐานข้อมูล MySQL ในหน่วยทดสอบ
	dsn := "root:password@tcp(0.0.0.0:3307)/social?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	// Run migrations สำหรับสร้างตาราง Users
	err = db.AutoMigrate(&models.Users{})
	assert.NoError(t, err)

	// เตรียมข้อมูลผู้ใช้ในฐานข้อมูลทดสอบ
	user := models.Users{
		Username:  "john_doe",
		Fullname:  "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// สร้าง UserHandler โดยใช้ฐานข้อมูลที่เตรียมไว้
	userHandler := handlers.NewUserHandler(db)

	// สร้างเครื่องมือทดสอบ HTTP และเรียกใช้งานฟังก์ชัน DeleteUser
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "id", Value: strconv.FormatUint(uint64(user.ID), 10)})
	userHandler.DeleteUser(c)

	// ตรวจสอบการลบผู้ใช้สำเร็จ
	assert.Equal(t, http.StatusOK, w.Code)

	// ตรวจสอบว่าผู้ใช้ถูกลบออกจากฐานข้อมูล
	var deletedUser models.Users
	err = db.First(&deletedUser, user.ID).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
