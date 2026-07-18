package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/yithsopheaktra/go-auth-api/internal/dto/auth"
	"github.com/yithsopheaktra/go-auth-api/internal/service"
)

type AuthController struct {
	Service *service.AuthService
}

func NewAuthController(
	service *service.AuthService,
) *AuthController {

	return &AuthController{
		Service: service,
	}
}

func (a *AuthController) Register(c *gin.Context) {

	var request dto.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(400, gin.H{
			"errors": err.Error(),
		})

		return
	}

	err := a.Service.Register(request)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

func (a *AuthController) Login(
	c *gin.Context,
) {

	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := a.Service.Login(request)

	if err != nil {

		c.JSON(401, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{

		"access_token": token,
	})

}
