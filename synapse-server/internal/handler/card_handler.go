package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/pagination"
	"github.com/synapse/server/internal/service"
)

// CardHandler handles /decks/:id/cards and /cards/:id routes.
type CardHandler struct {
	cardSvc   service.CardService
	studySvc  service.StudyService
	uploadDir string
}

func NewCardHandler(cardSvc service.CardService, studySvc service.StudyService, uploadDir string) *CardHandler {
	return &CardHandler{cardSvc: cardSvc, studySvc: studySvc, uploadDir: uploadDir}
}

// ListCards godoc
// @Summary List cards for a deck
// @Description Get a paginated list of cards for a specific deck
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} domain.Card
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/cards [get]
func (h *CardHandler) ListCards(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	pg := pagination.Parse(c)
	cards, total, err := h.cardSvc.ListByDeck(c.Request.Context(), deckID, middleware.GetUserID(c), pg.Offset, pg.Limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagination.SetTotalCount(c, total)
	c.JSON(http.StatusOK, cards)
}

type createCardRequest struct {
	Type    domain.CardType `json:"type"    binding:"required"`
	Content domain.RawJSON  `json:"content" binding:"required"`
}

// CreateCard godoc
// @Summary Create a new card
// @Description Create a new card in a specific deck
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param request body createCardRequest true "Card Information"
// @Success 201 {object} domain.Card
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/cards [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	var body createCardRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card := &domain.Card{DeckID: deckID, Type: body.Type, Content: body.Content}
	if err := h.cardSvc.Create(c.Request.Context(), card); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, card)
}

// GetCard godoc
// @Summary Get a card by ID
// @Description Get details of a specific card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Card ID"
// @Success 200 {object} domain.Card
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cards/{id} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	card, err := h.cardSvc.GetByID(c.Request.Context(), id, middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, card)
}

// UpdateCard godoc
// @Summary Update a card
// @Description Update a card's information
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Card ID"
// @Param request body domain.Card true "Card Information"
// @Success 200 {object} domain.Card
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cards/{id} [put]
func (h *CardHandler) UpdateCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var updates domain.Card
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := h.cardSvc.Update(c.Request.Context(), id, middleware.GetUserID(c), &updates)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, card)
}

// DeleteCard godoc
// @Summary Delete a card
// @Description Delete a specific card
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Card ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cards/{id} [delete]
func (h *CardHandler) DeleteCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.cardSvc.Delete(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// DueCount godoc
// @Summary Get due count for a deck
// @Description Get the number of cards currently due for review in a specific deck
// @Tags Cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/cards/due [get]
func (h *CardHandler) DueCount(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	count, err := h.cardSvc.CountDue(c.Request.Context(), deckID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// UploadMedia godoc
// @Summary Upload media for a card
// @Description Upload an attachment for a card
// @Tags Cards
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Card ID"
// @Param file formData file true "Media File"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cards/{id}/media [post]
func (h *CardHandler) UploadMedia(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	dir := filepath.Join(h.uploadDir, userID.String(), cardID.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		_ = c.Error(err)
		return
	}
	dst := filepath.Join(dir, header.Filename)
	out, err := os.Create(dst)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"path": dst})
}
