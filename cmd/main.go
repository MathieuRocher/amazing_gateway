package main

import (
	"amazing_gateway/internal/adapter/application"
	"amazing_gateway/internal/adapter/handler"
	"amazing_gateway/internal/adapter/proxy"
	"amazing_gateway/internal/adapter/repository"
	"amazing_gateway/internal/auth"
	"amazing_gateway/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()

	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	_ = database.DB.AutoMigrate(&repository.User{}, &repository.ClassGroup{})

	reviewPort := os.Getenv("REVIEW_PORT")
	if reviewPort == "" {
		log.Fatal("REVIEW_PORT must be set in environment")
	}

	formPort := os.Getenv("FORM_PORT")
	if formPort == "" {
		log.Fatal("FORM_PORT must be set in environment")
	}

	// === Microservice URLs ===
	reviewService := "http://localhost:" + reviewPort
	formService := "http://localhost:8081" + formPort

	// === Dépendances ===
	userRepo := repository.NewUserRepository()
	userUC := application.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC)
	authHandler := handler.NewAuthHandler(userUC)

	classGroupRepo := repository.NewClassGroupRepository()
	classGroupUC := application.NewClassGroupUseCase(classGroupRepo)
	classGroupHandler := handler.NewClassGroupHandler(classGroupUC)

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
