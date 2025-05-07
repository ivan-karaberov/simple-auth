package main

import (
	"context"
	"simpleAuth/config"
	"simpleAuth/controllers"
	"simpleAuth/models"

	_ "simpleAuth/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @Description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	cfg := config.LoadConfig(ctx, ".env", "certs/jwt-private.pem", "certs/jwt-public.pem")

	gin.SetMode(gin.ReleaseMode)

	db := models.NewDBConnection(cfg)

	router := gin.Default()

	controllers.SetupRoutes(db, cfg, router)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":3000")
}
