package main

import (
	_ "fileAnalisysService/docs"
	"fileAnalisysService/internal/antiplagiat"
	"fileAnalisysService/internal/api/handler"
	"fileAnalisysService/internal/api/router"
	"fileAnalisysService/internal/apllication/manager"
	"fileAnalisysService/internal/config"
	"fileAnalisysService/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Загружаем конфиг из env
	cfg := config.LoadConfig()

	// Инициализация S3 для хранения отчетов
	reportStorage := repository.NewMinioReportStorage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
	)

	// HTTP клиент для FileStorageService
	fsClient := repository.NewHttpFileStorageClient(cfg.FileStorageURL)

	// AntiPlagiarismEngine
	engine := antiplagiat.NewBitwiseEngine()

	// Менеджер
	manager := manager.NewFileAnalysisManager(fsClient, engine, reportStorage)

	// Handler
	analysisHandler := handler.NewAnalysisHandler(manager)

	// Router
	r := gin.Default()
	router.RegisterRoutes(r, analysisHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	log.Printf("Starting server on :%s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
