package handler

import (
	"net/http"

	"apiGateway/internal/manager"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Service *manager.PlagiarismService
}

func NewReportHandler(service *manager.PlagiarismService) *ReportHandler {
	return &ReportHandler{
		Service: service,
	}
}

// GetReportHandler godoc
// @Summary Получить последний отчёт по типу работы
// @Description Возвращает последний отчёт о плагиате для указанного типа работы
// @Tags Reports
// @Accept json
// @Produce json
// @Param typeWork path string true "Тип работы (лаба, домашка, курс и т.д.)"
// @Success 200 {object} map[string]interface{} "Последний отчёт"
// @Failure 404 {object} map[string]string "Отчёт не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /works/{typeWork}/reports/ [get]
func (h *ReportHandler) GetReportHandler(c *gin.Context) {
	typeWork := c.Param("typeWork")

	report, err := h.Service.GetLatestReport(typeWork)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
