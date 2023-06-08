package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/NopparootSuree/go-social/models"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// จัดการ payload
type Payload struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// จัดการ req ของ login
type LoginUserRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}

// จัดการ req ของ register
type CreateUserRequest struct {
	Username       string `json:"username" binding:"required,min=6"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
	FullName       string `json:"fullName" binding:"required,min=6"`
	Email          string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := utils.HashPassword(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Users
	h.db.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&existingUser)
	if existingUser.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	user := models.Users{
		Username:       req.Username,
		HashedPassword: hashPassword,
		Fullname:       req.FullName,
		Email:          req.Email,
	}

	created := h.db.Create(&user)
	if created.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": created.Error.Error()})
		return
	}

	response := CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)

}

func (h *UserHandler) Login(c *gin.Context) {
	var loginReq LoginUserRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//check user Exists
	var user models.Users
	result := h.db.Where("username = ?", loginReq.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record is not found"})
		return
	}

	err := utils.ComparePasswords(user.HashedPassword, loginReq.Password)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login failed"})
		return
	}

	//สร้าง secretKey
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	//ปรับแต่ง key
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	}

	// สร้าง Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เข้ารหัส Token เป็นสตริง JWT โดยใช้คีย์สั่งลับ
	token, errSecretKey := jwtToken.SignedString(secretKey)
	if errSecretKey != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errSecretKey.Error()})
		return
	}

	payload := &Payload{
		Token:     token,
		Username:  user.Username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 3),
	}

	c.JSON(http.StatusOK, gin.H{"payload": payload})
}
