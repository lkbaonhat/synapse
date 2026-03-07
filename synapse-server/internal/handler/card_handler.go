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

func (h *CardHandler) CreateCard(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	var body struct {
		Type    domain.CardType `json:"type"    binding:"required"`
		Content domain.RawJSON  `json:"content" binding:"required"`
	}
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

// UploadMedia saves an attachment for a card.
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
