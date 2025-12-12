package handler

import (
	"fileAnalisysService/internal/apllication/manager"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	manager *manager.FileAnalysisManager
}

func NewAnalysisHandler(m *manager.FileAnalysisManager) *AnalysisHandler {
	return &AnalysisHandler{manager: m}
}

// AnalyzeTypeWork godoc
// @Summary Создание/обновление отчета по типу работы
// @Description Проводит анализ всех работ указанного типа на плагиат и сохраняет отчет в S3.
// Если отчет уже существует для данного типа работы, он будет перезаписан.
// @Tags Аналитика
// @Produce json
// @Param typeWork path string true "Тип работы"
// @Success 200 {object} entities.AnalysisReport "Созданный или обновленный отчет по типу работы"
// @Failure 400 {object} map[string]string "error: typeWork не указан"
// @Failure 500 {object} map[string]string "error: внутренняя ошибка при анализе"
// @Router /works/{typeWork}/reports [post]
func (h *AnalysisHandler) AnalyzeTypeWork(c *gin.Context) {
	typeWork := c.Param("typeWork")
	if typeWork == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "typeWork обязательен"})
		return
	}

	report, err := h.manager.AnalyzeTypeWork(typeWork)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetLastReport godoc
// @Summary Получение последнего отчета по типу работы
// @Description Возвращает последний созданный отчет по указанному типу работы.
// @Tags Аналитика
// @Produce json
// @Param typeWork path string true "Тип работы"
// @Success 200 {object} entities.AnalysisReport "Последний отчет по типу работы"
// @Failure 400 {object} map[string]string "error: typeWork не указан"
// @Failure 404 {object} map[string]string "error: отчет не найден"
// @Router /works/{typeWork}/reports/last [get]
func (h *AnalysisHandler) GetLastReport(c *gin.Context) {
	typeWork := c.Param("typeWork")
	if typeWork == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "typeWork обязательен"})
		return
	}

	report, err := h.manager.GetLastReportByTypeWork(typeWork)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
