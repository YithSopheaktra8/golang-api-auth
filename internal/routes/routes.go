package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yithsopheaktra/go-auth-api/internal/controller"
	"github.com/yithsopheaktra/go-auth-api/internal/middleware"
	"github.com/yithsopheaktra/go-auth-api/internal/repository"
	"github.com/yithsopheaktra/go-auth-api/internal/service"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	userRepo := repository.NewUserRepository()

	authService := service.NewAuthService(
		userRepo,
	)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController()

	api := router.Group("/api/v1")

	// auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	users := api.Group("user")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/me", userController.Me)
	}

	return router

}
