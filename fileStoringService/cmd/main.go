package main

import (
	_ "fileStoringService/docs"
	"fileStoringService/internal/api/handler"
	"fileStoringService/internal/api/router"
	"fileStoringService/internal/application/manager"
	"fileStoringService/internal/config"
	"fileStoringService/internal/repository"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключение к Postgres
	dbHandler := repository.NewPostgresDB(cfg.PostgresDBUrl)
	defer dbHandler.Close()

	// Подключение к MinIO
	s3Storage := repository.NewMinioStorage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
	)

	// Создаем менеджер
	mgr := manager.NewManagerFileStorage(dbHandler, s3Storage)

	// Создаем хэндлер
	fileHandler := handler.NewFileHandler(mgr)

	// Настройка роутов
	r := gin.Default()
	router.RegisterRoutes(r, fileHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Starting server on :%s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
