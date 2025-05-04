package controllers

import (
	"net/http"
	"simpleAuth/config"
	"simpleAuth/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewAuthController(db *gorm.DB, cfg *config.Config) *AuthController {
	return &AuthController{DB: db, Cfg: cfg}
}

func (a *AuthController) SetupRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	auth.POST("/signin/:id", a.SignInHandler)
}

func (ac *AuthController) SignInHandler(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	tokenPair, err := services.SignIn(ac.DB, ac.Cfg, services.UserInfo{
		UserID:    userID,
		UserIP:    c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
