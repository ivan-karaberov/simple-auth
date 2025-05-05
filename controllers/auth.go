package controllers

import (
	"net/http"
	"simpleAuth/config"
	"simpleAuth/errors"
	"simpleAuth/middleware"
	"simpleAuth/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	auth.POST("/refresh", a.RefreshTokenHandler)
	auth.POST("/signout", middleware.AuthMiddleware(a.DB, a.Cfg), a.SignOutHandler)
}

func (ac *AuthController) SignInHandler(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		errors.APIError(c, errors.ErrBadRequestBody)
		return
	}

	tokenPair, err := services.SignIn(ac.DB, ac.Cfg, services.UserInfo{
		UserID:    userID,
		UserIP:    c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})
	if err != nil {
		logrus.WithError(err).Error("Failed signin")
		errors.APIError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

func (ac *AuthController) RefreshTokenHandler(c *gin.Context) {
	var tokenPair services.TokenPair
	if err := c.ShouldBindJSON(&tokenPair); err != nil {
		errors.APIError(c, errors.ErrBadRequestBody)
		return
	}

	userIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	newTokenPair, err := services.RefreshToken(ac.DB, ac.Cfg, &tokenPair, userIP, userAgent)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh token")
		errors.APIError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, newTokenPair)
}

func (ac *AuthController) SignOutHandler(c *gin.Context) {
	sessionID := c.Value("sessionID")
	if sessionID == nil {
		logrus.Error("Failed signout sessionID is empty")
		errors.APIError(c, errors.ErrInternalServer)
		return
	}

	err := services.SignOut(ac.DB, sessionID.(string))
	if err != nil {
		logrus.WithError(err).Error("Failed signout")
		errors.APIError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sign out success"})
}
