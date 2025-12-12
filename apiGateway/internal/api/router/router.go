package router

import (
	"apiGateway/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, fileHandler *handler.FileHandler, reportHandler *handler.ReportHandler) {
	r.POST("/files/upload", fileHandler.UploadFileHandler)
	r.GET("/works/:typeWork/reports/", reportHandler.GetReportHandler)
}
