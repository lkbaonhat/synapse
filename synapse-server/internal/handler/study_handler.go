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

type startSessionRequest struct {
	DeckID uuid.UUID       `json:"deckId" binding:"required"`
	Mode   domain.StudyMode `json:"mode"   binding:"required"`
}

// StartSession godoc
// @Summary Start a study session
// @Description Start a new study session for a deck in a specific mode (review, cram, quiz)
// @Tags Study
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body startSessionRequest true "Session Parameters"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /study/sessions [post]
func (h *StudyHandler) StartSession(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body startSessionRequest
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

// NextCards godoc
// @Summary Get next due cards for session
// @Description Fetch the next batch of due cards for an active session
// @Tags Study
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {array} domain.Card
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /study/sessions/{id}/cards/next [get]
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

type answerRequest struct {
	CardID    uuid.UUID `json:"cardId"    binding:"required"`
	Rating    int       `json:"rating"    binding:"required,min=1,max=4"`
	TimeTaken int       `json:"timeTaken" binding:"required,min=0"`
}

// Answer godoc
// @Summary Submit an answer
// @Description Submit an answer rating for a card in an active study session
// @Tags Study
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Param request body answerRequest true "Answer Details"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /study/sessions/{id}/answer [post]
func (h *StudyHandler) Answer(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	var body answerRequest
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

// EndSession godoc
// @Summary End a study session
// @Description End an active study session and save progress
// @Tags Study
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /study/sessions/{id}/end [post]
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

// GetQuizResult godoc
// @Summary Get quiz results
// @Description Retrieve the scorecard and wrong answer details for a completed quiz session
// @Tags Study
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID"
// @Success 200 {object} domain.QuizResult
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /study/sessions/{id}/results [get]
func (h *StudyHandler) GetQuizResult(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	
	result, err := h.studySvc.GetQuizResult(c.Request.Context(), sessionID, middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	
	c.JSON(http.StatusOK, result)
}
