package router

import (
	"fileAnalisysService/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.AnalysisHandler) {
	r.POST("/works/:typeWork/reports", h.AnalyzeTypeWork)
	r.GET("/works/:typeWork/reports/last", h.GetLastReport)
}
