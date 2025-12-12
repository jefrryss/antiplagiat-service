package handler

import (
	"apiGateway/internal/manager"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	Service *manager.PlagiarismService
}

func NewFileHandler(service *manager.PlagiarismService) *FileHandler {
	return &FileHandler{
		Service: service,
	}
}

// UploadFileHandler godoc
// @Summary Загрузить файл и создать облако слов
// @Description Загружает файл студента, создаёт отчёт о плагиате и возвращает ссылку на облако слов
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param userName formData string true "Имя пользователя"
// @Param typeWork formData string true "Тип работы (лаба, домашка и т.д.)"
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string "Ссылка на облако слов"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /files/upload [post]
func (h *FileHandler) UploadFileHandler(c *gin.Context) {
	userName := c.PostForm("userName")
	typeWork := c.PostForm("typeWork")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer file.Close()

	cloudURL, err := h.Service.UploadFileAndGetWordCloud(userName, typeWork, fileHeader.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wordCloudURL": cloudURL})
}
