package main

import (
	"amazing_gateway/internal/adapter/application"
	"amazing_gateway/internal/adapter/handler"
	"amazing_gateway/internal/adapter/proxy"
	"amazing_gateway/internal/adapter/repository"
	"amazing_gateway/internal/auth"
	"amazing_gateway/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	database.InitDB()

	// === Microservice URLs ===
	reviewService := "http://localhost:8082"
	formService := "http://localhost:8081"

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

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
