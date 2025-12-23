package http

import (
	"github.com/gin-gonic/gin"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
	// "github.com/EnockYator/saas-photo-listing-platform/backend/internal/http/handlers"
	// "github.com/EnockYator/saas-photo-listing-platform/backend/internal/middleware"
	"gorm.io/gorm"

	_ "github.com/EnockYator/saas-photo-listing-platform/backend/internal/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// // Global middleware
	// r.Use(middleware.CORSMiddleware())
	// r.Use(middleware.Logger())

	// // Auth routes
	// authGroup := r.Group("/auth")
	// {
	// 	authGroup.POST("/login", handlers.LoginHandler(db, cfg))
	// 	authGroup.POST("/register", handlers.RegisterHandler(db, cfg))
	// }

	// // Photo routes (require auth)
	// photoGroup := r.Group("/photos")
	// photoGroup.Use(middleware.AuthMiddleware(cfg))
	// {
	// 	photoGroup.POST("/", handlers.UploadPhotoHandler(db))
	// 	photoGroup.GET("/", handlers.ListPhotosHandler(db))
	// }

	return r
}
