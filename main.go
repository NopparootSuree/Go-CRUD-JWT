package main

import (
	"log"

	"github.com/NopparootSuree/go-social/routers"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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
