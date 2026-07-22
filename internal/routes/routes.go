package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yithsopheaktra/go-auth-api/internal/controller"
	"github.com/yithsopheaktra/go-auth-api/internal/middleware"
	"github.com/yithsopheaktra/go-auth-api/internal/repository"
	"github.com/yithsopheaktra/go-auth-api/internal/service"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	router.Use(cors.New(cors.Config{

		AllowOrigins: []string{
			"http://localhost:5173",
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,
	}))

	userRepo := repository.NewUserRepository()
	refreshTokenRepo := repository.NewRefreshRepository()

	authService := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
	)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController()

	api := router.Group("/api/v1")

	// auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", authController.Logout)
	}

	users := api.Group("/user")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/me", userController.Me)
	}

	return router

}
