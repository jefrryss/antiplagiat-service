package handler

import (
	"fileStoringService/internal/application/manager"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	manager *manager.ManagerFileStorage
}

func NewFileHandler(m *manager.ManagerFileStorage) *FileHandler {
	return &FileHandler{manager: m}
}

// UploadFile godoc
// @Summary Загрузить файл
// @Description Загружает файл через менеджер, сохраняет его в MinIO и сохраняет метаданные работы в БД. Возвращает UUID загруженной работы.
// @Tags Файлы
// @Accept multipart/form-data
// @Produce json
// @Param userName formData string true "Имя пользователя, загрузившего файл"
// @Param typeWork formData string true "Тип работы или категории файла"
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string "work_id: UUID загруженной работы"
// @Failure 400 {object} map[string]string "error: неверный запрос или отсутствует файл/параметры"
// @Failure 500 {object} map[string]string "error: ошибка при сохранении файла в MinIO или сохранении работы в БД"
// @Router /upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	userName := c.PostForm("userName")
	typeWork := c.PostForm("typeWork")
	if userName == "" || typeWork == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userName и typeWork обязательны"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось получить файл"})
		return
	}
	defer file.Close()

	workID, err := h.manager.Save(c.Request.Context(), userName, typeWork, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"work_id": workID})
}

// GetFilesList godoc
// @Summary Получить список работ по типу
// @Description Возвращает список всех работ для указанного типа работы, включая информацию о файле (ID, имя, MIME-тип) и данные работы (userName, createdAt, typeWork и т.д.)
// @Tags Файлы
// @Produce json
// @Param typeWork path string true "Тип работы"
// @Success 200 {array} entities.Work "Список работ с полной информацией"
// @Failure 400 {object} map[string]string "error: typeWork не указан"
// @Failure 404 {object} map[string]string "error: работы не найдены"
// @Router /files/list/{typeWork} [get]
func (h *FileHandler) GetFilesList(c *gin.Context) {
	typeWork := c.Param("typeWork")
	if typeWork == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "typeWork обязательен"})
		return
	}

	works, err := h.manager.GetWorksInfoByType(c.Request.Context(), typeWork)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, works)
}

// DownloadFile godoc
// @Summary Скачать файл по UUID
// @Description Отправляет файл из MinIO по UUID работы. В ответе устанавливаются оригинальное имя файла и MIME-тип.
// @Tags Файлы
// @Produce application/octet-stream
// @Param objectName path string true "UUID работы (workID), связанной с файлом"
// @Success 200 {file} binary "Бинарный файл с оригинальным именем и MIME-типом"
// @Failure 400 {object} map[string]string "error: objectName не указан"
// @Failure 404 {object} map[string]string "error: работа или файл не найдены"
// @Failure 500 {object} map[string]string "error: не удалось отправить файл"
// @Router /files/download/{work_id} [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	workID := c.Param("objectName")
	if workID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "objectName обязательен"})
		return
	}

	work, file, err := h.manager.GetWorkWithFile(c.Request.Context(), workID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fileName := work.File.FileName
	if fileName == "" {
		fileName = work.File.ID.String()
	}
	contentType := work.File.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)

	if _, err := io.Copy(c.Writer, file); err != nil {
		c.Error(fmt.Errorf("не удалось отправить файл: %w", err))
		return
	}
}
