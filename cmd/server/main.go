package main

import (
	"github.com/yithsopheaktra/go-auth-api/internal/config"
	"github.com/yithsopheaktra/go-auth-api/internal/model"
	"github.com/yithsopheaktra/go-auth-api/internal/routes"
)

func main() {

	config.LoadEnv()

	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&model.User{},
	)

	router := routes.SetupRouter()

	router.Run(":8080")
}