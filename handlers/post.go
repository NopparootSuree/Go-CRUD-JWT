package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NopparootSuree/go-social/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{
		db: db,
	}
}

type CreatePostResponse struct {
	PostID    uint      `json:"postID" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Body      string    `json:"body" binding:"required"`
	UserID    uint      `json:"userID" binding:"required"`
	Status    string    `json:"status" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreatePostRequest struct {
	Title  string `json:"title" binding:"required,min=6"`
	Body   string `json:"body" binding:"required,min=6"`
	UserID uint   `json:"userID" binding:"required,min=1"`
	Status string `json:"status" binding:"required"`
}

type CreatePostUpdateRequest struct {
	Title  string `json:"title" binding:"required,min=6"`
	Body   string `json:"body" binding:"required,min=6"`
	Status string `json:"status" binding:"required"`
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	var posts []models.Posts
	result := h.db.Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var response []CreatePostResponse

	for _, post := range posts {
		response = append(response, CreatePostResponse{
			PostID:    post.PostID,
			Title:     post.Title,
			Body:      post.Body,
			UserID:    post.UserID,
			Status:    post.Status,
			CreatedAt: post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Posts{
		Title:  req.Title,
		Body:   req.Body,
		UserID: req.UserID,
		Status: req.Status,
	}

	result := h.db.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := CreatePostResponse{
		PostID:    post.PostID,
		Title:     post.Title,
		Body:      post.Body,
		UserID:    post.UserID,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)

}

func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Posts
	result := h.db.First(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := CreatePostResponse{
		PostID:    post.PostID,
		Title:     post.Title,
		Body:      post.Body,
		UserID:    post.UserID,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
	}

	c.JSON(http.StatusOK, response)

}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var req CreatePostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Posts
	result := h.db.First(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	num, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		panic(err)
	}

	post = models.Posts{
		PostID:    uint(num),
		Title:     req.Title,
		Body:      req.Body,
		UserID:    post.UserID,
		Status:    req.Status,
		CreatedAt: post.CreatedAt,
	}

	result = h.db.Save(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := CreatePostResponse{
		PostID:    post.PostID,
		Title:     post.Title,
		Body:      post.Body,
		UserID:    post.UserID,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update record"})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Posts
	result := h.db.Delete(&post, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record is not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Success": "removed record"})
	}
}
