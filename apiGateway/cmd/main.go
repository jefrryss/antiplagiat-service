package main

import (
	_ "apiGateway/docs"
	"apiGateway/internal/api/handler"
	"apiGateway/internal/api/router"
	"apiGateway/internal/config"
	"apiGateway/internal/manager"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Создаем менеджеры для работы с микросервисами
	fileMgr := manager.NewFileManager(cfg.FileStorageURL)
	reportMgr := manager.NewReportManager(cfg.AnalysisService)
	wordCloudMgr := manager.NewWordCloudManager()

	// Создаем сервис, объединяющий логику всех менеджеров
	plagService := manager.NewPlagiarismService(fileMgr, reportMgr, wordCloudMgr)

	// Создаем хэндлеры с использованием сервиса
	fileHandler := handler.NewFileHandler(plagService)
	reportHandler := handler.NewReportHandler(plagService)

	// Настройка роутера Gin
	r := gin.Default()
	router.RegisterRoutes(r, fileHandler, reportHandler)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Starting API Gateway on port :%s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run API Gateway: %v", err)
	}
}
