package routers

import (
	"os"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRouter(router *gin.Engine, db *gorm.DB) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	postHandler := handlers.NewPostHandler(db)
	posts := router.Group("/posts", middlewares.JWTMiddleware(secretKey))
	{
		posts.GET("", postHandler.ListPosts)
		posts.GET("/:id", postHandler.GetPost)
		posts.POST("", postHandler.CreatePost)
		posts.PUT("/:id", postHandler.UpdatePost)
		posts.DELETE("/:id", postHandler.DeletePost)
	}
}
