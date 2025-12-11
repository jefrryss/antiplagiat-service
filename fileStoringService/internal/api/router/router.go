package router

import (
	"fileStoringService/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.FileHandler) {
	r.POST("/upload", h.UploadFile)
	r.GET("/files/list/:typeWork", h.GetFilesList)
	r.GET("/files/download/:objectName", h.DownloadFile)
}
