package main

import (
	"github.com/baigel/lms/main-service/internal/config"
	"github.com/baigel/lms/main-service/internal/handler"
	"github.com/baigel/lms/main-service/internal/middleware"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/internal/service"
	"github.com/baigel/lms/main-service/pkg/database"
	"github.com/baigel/lms/main-service/pkg/logger"

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
	// Initialize logger
	logger.Init()

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Log.Info("Connected to database successfully")

	// Run migrations
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

	// Initialize repositories
	courseRepo := repository.NewCourseRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	lessonRepo := repository.NewLessonRepository(db)

	// Initialize services
	courseSvc := service.NewCourseService(courseRepo)
	chapterSvc := service.NewChapterService(chapterRepo)
	lessonSvc := service.NewLessonService(lessonRepo)

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseSvc)
	chapterHandler := handler.NewChapterHandler(chapterSvc)
	lessonHandler := handler.NewLessonHandler(lessonSvc)

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := router.Group("/api/v1")
	{
		// Courses
		api.POST("/courses", courseHandler.CreateCourse)
		api.GET("/courses", courseHandler.GetAllCourses)
		api.GET("/courses/:id", courseHandler.GetCourseByID)
		api.PUT("/courses/:id", courseHandler.UpdateCourse)
		api.DELETE("/courses/:id", courseHandler.DeleteCourse)

		// Chapters
		api.POST("/chapters", chapterHandler.CreateChapter)
		api.GET("/courses/:course_id/chapters", chapterHandler.GetChaptersByCourseID)
		api.GET("/chapters/:id", chapterHandler.GetChapterByID)
		api.PUT("/chapters/:id", chapterHandler.UpdateChapter)
		api.DELETE("/chapters/:id", chapterHandler.DeleteChapter)

		// Lessons
		api.POST("/lessons", lessonHandler.CreateLesson)
		api.GET("/chapters/:chapter_id/lessons", lessonHandler.GetLessonsByChapterID)
		api.GET("/lessons/:id", lessonHandler.GetLessonByID)
		api.PUT("/lessons/:id", lessonHandler.UpdateLesson)
		api.DELETE("/lessons/:id", lessonHandler.DeleteLesson)
	}

	logger.Log.WithField("port", cfg.AppPort).Info("Starting server...")
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Log.Fatalf("Failed to run server: %v", err)
	}
}
