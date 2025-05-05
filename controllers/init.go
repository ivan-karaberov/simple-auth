package controllers

import (
	"simpleAuth/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller interface {
	SetupRoutes(router *gin.Engine)
}

func SetupRoutes(db *gorm.DB, cfg *config.Config, router *gin.Engine) {
	var controllersList []Controller
	controllersList = append(controllersList, NewAuthController(db, cfg))
	controllersList = append(controllersList, NewUserController(db, cfg))

	for _, controller := range controllersList {
		controller.SetupRoutes(router)
	}
}
