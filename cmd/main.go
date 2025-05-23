package main

import (
	"amazing_gateway/internal/adapter/application"
	"amazing_gateway/internal/adapter/handler"
	"amazing_gateway/internal/adapter/proxy"
	"amazing_gateway/internal/adapter/repository"
	"amazing_gateway/internal/auth"
	"amazing_gateway/internal/infrastructure/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	database.InitDB()
	_ = database.DB.AutoMigrate(&repository.User{}, &repository.ClassGroup{})

	_ = godotenv.Load()

	reviewHost := os.Getenv("REVIEW_HOST")
	reviewPort := os.Getenv("REVIEW_PORT")
	formHost := os.Getenv("FORM_HOST")
	formPort := os.Getenv("FORM_PORT")

	if reviewHost == "" || reviewPort == "" || formHost == "" || formPort == "" {
		log.Fatal("❌ REVIEW_HOST, REVIEW_PORT, FORM_HOST et FORM_PORT doivent être définis")
	}

	reviewService := fmt.Sprintf("http://%s:%s", reviewHost, reviewPort)
	formService := fmt.Sprintf("http://%s:%s", formHost, formPort)

	// === Dépendances ===
	userRepo := repository.NewUserRepository()
	userUC := application.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC)
	authHandler := handler.NewAuthHandler(userUC)

	classGroupRepo := repository.NewClassGroupRepository()
	classGroupUC := application.NewClassGroupUseCase(classGroupRepo)
	classGroupHandler := handler.NewClassGroupHandler(classGroupUC)

	// apgeoerbo
	// === Public routes ===
	public := router.Group("/api")
	{
		public.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello world"})
		})
		authHandler.RegisterRoutes(public)
		userHandler.RegisterPublicRoutes(public)
	}

	// === Protected routes (with JWT middleware) ===
	private := router.Group("/api")
	private.Use(auth.JWTMiddleware())
	{
		userHandler.RegisterProtectedRoutes(private)
		classGroupHandler.RegisterRoutes(private)

		// Proxy config
		proxyRoutes := map[string]string{
			"/reviews":            reviewService,
			"/form-answers":       reviewService,
			"/forms":              formService,
			"/form-questions":     formService,
			"/courses":            formService,
			"/course-assignments": formService,
		}

		// Pour proxy propre
		for route, service := range proxyRoutes {
			private.Any(route, proxy.ReverseProxy(service))               // ex: /api/reviews
			private.Any(route+"/*proxyPath", proxy.ReverseProxy(service)) // ex: /api/reviews/42
		}

	}

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		log.Fatal("GATEWAY_PORT must be set in environment")
	}
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
