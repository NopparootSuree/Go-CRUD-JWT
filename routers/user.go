package routers

import (
	"github.com/NopparootSuree/go-social/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRouter(router *gin.Engine, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)
	users := router.Group("/users")
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
