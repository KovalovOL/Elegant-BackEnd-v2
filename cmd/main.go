package main

import (
	"app/internal/auth"
	"app/internal/auth/google"
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

	conn := db.GetPostgresDB()
	pool := db.GetPool()

	userRepo := user.NewRepository(conn)
	userSerc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSerc)

	googleOAuth, err := google.NewGoogleOAuth()
	if err != nil {
		log.Fatal(err)
	}

	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatal(err)
	}
	authRepo := auth.NewRepository(conn)
	authService := google.NewService(googleOAuth, jwtManager, userRepo, authRepo, pool)
	authHandler := google.NewHandler(authService)

	router := gin.Default()

	router.GET("/users/:id", userHandler.GetByID)
	router.POST("/users", userHandler.Create)

	router.GET("/auth/google/login", authHandler.LoginGoogle)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)
	router.POST("/auth/logout", authHandler.Logout)
	router.POST("/auth/refresh", authHandler.Refresh)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(jwtManager))
	{
		protected.GET("/auth/me", authHandler.Me)
	}
	router.Run(":8080")
}
