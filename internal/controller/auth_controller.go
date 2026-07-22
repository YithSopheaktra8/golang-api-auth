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

func (a *AuthController) Login(c *gin.Context) {

	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	login, err := a.Service.Login(request)

	if err != nil {

		c.JSON(401, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.SetCookie(
		"refresh_token",
		login.RefreshToken,
		60*60*24*30,
		"/",
		"",
		false, // localhost in production we use true
		true,
	)

	c.JSON(200, gin.H{
		"access_token": login.AccessToken,
	})
}

func (a *AuthController) Refresh(
	c *gin.Context,
) {

	refreshToken, err := c.Cookie(
		"refresh_token",
	)

	if err != nil {

		c.JSON(401, gin.H{
			"error": "missing refresh token",
		})

		return
	}

	newAccess,
		newRefresh,
		err := a.Service.Refresh(
		refreshToken,
	)

	if err != nil {

		c.JSON(401, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.SetCookie(
		"refresh_token",
		newRefresh,
		60*60*24*30,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{

		"access_token": newAccess,
	})

}

func (a *AuthController) Logout(
	c *gin.Context,
) {

	refreshToken, err := c.Cookie(
		"refresh_token",
	)

	if err == nil {

		err = a.Service.Logout(
			refreshToken,
		)

		if err != nil {

			c.JSON(401, gin.H{
				"error": err.Error(),
			})

			return

		}

	}

	// remove cookie

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "logout success",
	})

}
