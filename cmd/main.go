package main

import (
	"app/internal/db"
	"app/internal/user"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
    	log.Printf("Error loading .env file: %v. Proceeding without it.", err)
    }

	conn := db.GetPostgresConnection()

	userRepo := user.NewRepository(conn)
	userSerc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSerc)

	router := gin.Default()
	router.GET("/users/:id", userHandler.GetByID)
	router.POST("/users", userHandler.Create)

	router.Run(":8080")
}