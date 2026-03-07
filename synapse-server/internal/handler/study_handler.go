package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/service"
)

// StudyHandler handles /study/sessions/* routes.
type StudyHandler struct {
	studySvc service.StudyService
}

func NewStudyHandler(studySvc service.StudyService) *StudyHandler {
	return &StudyHandler{studySvc: studySvc}
}

func (h *StudyHandler) StartSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body struct {
		DeckID uuid.UUID       `json:"deckId" binding:"required"`
		Mode   domain.StudyMode `json:"mode"   binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, cards, err := h.studySvc.StartSession(c.Request.Context(), userID, body.DeckID, body.Mode)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"session": session, "cards": cards})
}

func (h *StudyHandler) NextCards(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	cards, err := h.studySvc.NextCards(c.Request.Context(), sessionID, middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, cards)
}

func (h *StudyHandler) Answer(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	var body struct {
		CardID    uuid.UUID `json:"cardId"    binding:"required"`
		Rating    int       `json:"rating"    binding:"required,min=1,max=4"`
		TimeTaken int       `json:"timeTaken" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.studySvc.Answer(c.Request.Context(), sessionID, middleware.GetUserID(c), body.CardID, body.Rating, body.TimeTaken); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *StudyHandler) EndSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.studySvc.EndSession(c.Request.Context(), sessionID, middleware.GetUserID(c)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
