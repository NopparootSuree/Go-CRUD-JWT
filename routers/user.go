package routers

import (
	"os"

	"github.com/NopparootSuree/go-social/handlers"
	"github.com/NopparootSuree/go-social/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRouter(router *gin.Engine, db *gorm.DB) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	userHandler := handlers.NewUserHandler(db)
	users := router.Group("/users", middlewares.JWTMiddleware(secretKey))
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
