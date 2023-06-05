package routers

import (
	"github.com/NopparootSuree/go-social/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRouter(router *gin.Engine, db *gorm.DB) {
	postHandler := handlers.NewPostHandler(db)
	posts := router.Group("/posts")
	{
		posts.GET("", postHandler.ListPosts)
		posts.GET("/:id", postHandler.GetPost)
		posts.POST("", postHandler.CreatePost)
		posts.PUT("/:id", postHandler.UpdatePost)
		posts.DELETE("/:id", postHandler.DeletePost)
	}
}
