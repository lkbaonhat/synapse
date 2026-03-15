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

// ImportCSV godoc
// @Summary Import cards from CSV
// @Description Import flashcards into a specific deck from a CSV file
// @Tags Import/Export
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param file formData file true "CSV File (format: front,back)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/import [post]
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

// ExportDeckCSV godoc
// @Summary Export deck to CSV
// @Description Download a deck's cards as a CSV file
// @Tags Import/Export
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 200 {file} file "deck_export.csv"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/export [get]
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

// ExportUserJSON godoc
// @Summary Export all user data to JSON
// @Description Download all folders, decks, cards, and study logs for the authenticated user as a JSON file
// @Tags Import/Export
// @Accept json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {file} file "synapse_export.json"
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /export [get]
func (h *ImportExportHandler) ExportUserJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=synapse_export.json")
	if err := h.ieSvc.ExportUserJSON(c.Request.Context(), middleware.GetUserID(c), c.Writer); err != nil {
		_ = c.Error(err)
	}
}
