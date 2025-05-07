package controllers

import (
	"net/http"
	"simpleAuth/config"
	"simpleAuth/errors"
	"simpleAuth/middleware"
	"simpleAuth/models"
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

// @Summary User Sign In
// @Description Signs in a user and returns a token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} services.TokenPair
// @Failure 400 {object} errors.ErrorResponse "Bad Request body"
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/signin/{id} [post]
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

// @Summary Refreshes the access and refresh tokens
// @Description Refreshes the access and refresh tokens using the provided token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param tokenPair body services.TokenPair true "Token pair containing refresh token"
// @Success 200 {object} services.TokenPair "New token pair"
// @Failure 400 {object} errors.ErrorResponse "Bad request body"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
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

// @Summary Signs out the user
// @Description Signs out the user by invalidating the session ID
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SignOutResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /auth/signout [post]
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

	c.JSON(http.StatusOK, models.SignOutResponse{Message: "Sign out success"})
}
