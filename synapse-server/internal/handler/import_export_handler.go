package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/service"
)

// ImportExportHandler handles import/export routes.
type ImportExportHandler struct {
	ieSvc service.ImportExportService
}

func NewImportExportHandler(ieSvc service.ImportExportService) *ImportExportHandler {
	return &ImportExportHandler{ieSvc: ieSvc}
}

func (h *ImportExportHandler) ImportCSV(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	result, err := h.ieSvc.ImportCSV(c.Request.Context(), deckID, file)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ImportExportHandler) ExportDeckCSV(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=deck_export.csv")
	if err := h.ieSvc.ExportDeckCSV(c.Request.Context(), deckID, middleware.GetUserID(c), c.Writer); err != nil {
		_ = c.Error(err)
	}
}

func (h *ImportExportHandler) ExportUserJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=synapse_export.json")
	if err := h.ieSvc.ExportUserJSON(c.Request.Context(), middleware.GetUserID(c), c.Writer); err != nil {
		_ = c.Error(err)
	}
}
