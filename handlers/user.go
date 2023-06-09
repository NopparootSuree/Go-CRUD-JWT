package handlers

import (
	"net/http"
	"time"

	"github.com/NopparootSuree/go-social/models"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

type CreateUserResponse struct {
	ID        uint      `json:"id" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	FullName  string    `json:"fullName" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}

type CreateUserUpdateRequest struct {
	HashedPassword string `json:"hashedPassword" binding:"required,min=6"`
	FullName       string `json:"fullName" binding:"required,min=6"`
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.Users
	result := h.db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var response []CreateUserResponse

	for _, user := range users {
		response = append(response, CreateUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.Fullname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.Users
	result := h.db.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)

}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req CreateUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := utils.HashPassword(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user models.Users
	result := h.db.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	// Prepare the update data
	updatesUser := map[string]interface{}{
		"hashedPassword": hashPassword,
		"fullname":       req.FullName,
	}

	// Update the user's information
	result = h.db.Model(&user).Updates(updatesUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update record"})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.Users
	result := h.db.Delete(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record is not found"})
	} else {
		c.JSON(http.StatusNoContent, gin.H{"Success": "removed record"})
	}
}
