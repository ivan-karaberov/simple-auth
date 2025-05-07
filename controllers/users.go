package controllers

import (
	"net/http"
	"simpleAuth/config"
	"simpleAuth/errors"
	"simpleAuth/middleware"
	"simpleAuth/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserController struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewUserController(db *gorm.DB, cfg *config.Config) *UserController {
	return &UserController{DB: db, Cfg: cfg}
}

func (u *UserController) SetupRoutes(router *gin.Engine) {
	user := router.Group("/users")

	user.GET("/me", middleware.AuthMiddleware(u.DB, u.Cfg), u.UserDetailHandler)
}

// @Summary Get current user info
// @Description Gets current user info by token authorization
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /users/me [get]
func (u *UserController) UserDetailHandler(c *gin.Context) {
	userID := c.Value("userID")
	if userID != nil {
		c.JSON(http.StatusOK, models.UserResponse{UserID: userID.(string)})
		return
	}
	logrus.Error("Failed get user detail, userID is empty")
	errors.APIError(c, errors.ErrInternalServer)
}
