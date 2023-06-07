package routers

import (
	"github.com/NopparootSuree/go-social/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthenRouter(router *gin.Engine, db *gorm.DB) {
	authenHandler := handlers.NewUserHandler(db)
	authen := router.Group("/")
	{
		authen.POST("/login", authenHandler.Login)
		authen.POST("/register", authenHandler.Register)
	}
}
