package main

import (
	"github.com/NopparootSuree/go-social/routers"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := utils.ConnectDatabase()
	if err != nil {
		panic("Failed connect to Database")
	}

	r := gin.Default()
	routers.UserRouter(r, db)
	routers.PostRouter(r, db)

	r.Run(":8080")
}
