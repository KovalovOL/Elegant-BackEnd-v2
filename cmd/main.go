package main

import (
	"app/internal/db"
	"app/internal/user"
	"app/internal/auth"
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

	googleOAuth, err := auth.NewGoogleOAuth()
	if err != nil {
		log.Fatal(err)
	}

	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatal(err)
	}
	authService := auth.NewService(googleOAuth, jwtManager, userRepo)
	authHandler := auth.NewHandler(authService)

	router := gin.Default()

	router.GET("/users/:id", userHandler.GetByID)
	router.POST("/users", userHandler.Create)

	router.GET("/auth/google/login", authHandler.LoginGoogle)
	router.GET("/auth/logout", authHandler.Logout)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", authHandler.Me)
	}
	router.Run(":8080")
}