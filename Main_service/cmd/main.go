package main

import (
	"context"

	"github.com/baigel/lms/main-service/internal/config"
	"github.com/baigel/lms/main-service/internal/handler"
	"github.com/baigel/lms/main-service/internal/middleware"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/database"
	"github.com/baigel/lms/main-service/pkg/keycloak"
	"github.com/baigel/lms/main-service/pkg/logger"
	"github.com/baigel/lms/main-service/pkg/storage"

	_ "github.com/baigel/lms/main-service/docs" // swagger docs
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // postgres driver for goose
	"github.com/pressly/goose/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title LMS Main Service API
// @version 1.0
// @description API Server for LMS Main Service
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger.Init()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Log.Info("Connected to database successfully")

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatalf("Failed to get sql.DB: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Log.Fatalf("Failed to set goose dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		logger.Log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Log.Info("Migrations applied successfully")

	courseRepo := repository.NewCourseRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	lessonRepo := repository.NewLessonRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	courseSvc := service.NewCourseService(courseRepo)
	chapterSvc := service.NewChapterService(chapterRepo)
	lessonSvc := service.NewLessonService(lessonRepo)

	fileStorage, err := storage.NewMinIOStorage(context.Background(), storage.MinIOConfig{
		Endpoint:  cfg.MinIOEndpoint,
		AccessKey: cfg.MinIOAccessKey,
		SecretKey: cfg.MinIOSecretKey,
		Bucket:    cfg.MinIOBucket,
		UseSSL:    cfg.MinIOUseSSL,
	})
	if err != nil {
		logger.Log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	attachmentSvc := service.NewAttachmentService(attachmentRepo, lessonRepo, fileStorage)

	courseHandler := handler.NewCourseHandler(courseSvc)
	chapterHandler := handler.NewChapterHandler(chapterSvc)
	lessonHandler := handler.NewLessonHandler(lessonSvc)
	attachmentHandler := handler.NewAttachmentHandler(attachmentSvc)

	kc := keycloak.New(
		cfg.KeycloakURL,
		cfg.KeycloakClientID,
		cfg.KeycloakClientSecret,
		cfg.KeycloakRealm,
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// Публичные маршруты — чтение без авторизации
	public := api.Group("")
	{
		public.GET("/courses", courseHandler.GetAllCourses)
		public.GET("/courses/:id", courseHandler.GetCourseByID)
		public.GET("/courses/:course_id/chapters", chapterHandler.GetChaptersByCourseID)
		public.GET("/chapters/:id", chapterHandler.GetChapterByID)
		public.GET("/chapters/:chapter_id/lessons", lessonHandler.GetLessonsByChapterID)
		public.GET("/lessons/:id", lessonHandler.GetLessonByID)
	}

	// Защищённые маршруты — только admin
	admin := api.Group("", middleware.RequireAdmin(kc))
	{
		admin.POST("/courses", courseHandler.CreateCourse)
		admin.PUT("/courses/:id", courseHandler.UpdateCourse)
		admin.DELETE("/courses/:id", courseHandler.DeleteCourse)

		admin.POST("/chapters", chapterHandler.CreateChapter)
		admin.PUT("/chapters/:id", chapterHandler.UpdateChapter)
		admin.DELETE("/chapters/:id", chapterHandler.DeleteChapter)

		admin.POST("/lessons", lessonHandler.CreateLesson)
		admin.PUT("/lessons/:id", lessonHandler.UpdateLesson)
		admin.DELETE("/lessons/:id", lessonHandler.DeleteLesson)
	}

	teacherFiles := api.Group("", middleware.RequireAdminOrTeacher(kc))
	{
		teacherFiles.POST("/upload", attachmentHandler.Upload)
	}

	authenticated := api.Group("", middleware.RequireAuth(kc))
	{
		authenticated.GET("/download", attachmentHandler.Download)
	}

	logger.Log.WithField("port", cfg.AppPort).Info("Starting server...")
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Log.Fatalf("Failed to run server: %v", err)
	}
}
