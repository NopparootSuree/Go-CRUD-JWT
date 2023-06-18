package main

import (
	"log"

	"github.com/NopparootSuree/go-social/routers"
	"github.com/NopparootSuree/go-social/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var secretKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	r := gin.Default()

	db, err := utils.ConnectDatabase()
	if err != nil {
		panic("Failed connect to Database")
	}

	routers.UserRouter(r, db)
	routers.PostRouter(r, db)
	routers.AuthenRouter(r, db)

	r.Use(cors.Default())
	r.Run(":8080")
}
