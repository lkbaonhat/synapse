package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	domain "github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/service"
)

// StatsHandler handles /stats/* and /decks/:id/stats routes.
type StatsHandler struct {
	statsSvc service.StatsService
}

func NewStatsHandler(statsSvc service.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: statsSvc}
}

// Overview godoc
// @Summary Get learning overview stats
// @Description Get overall learning statistics for the current user
// @Tags Stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param totalCards query int false "Total cards in library"
// @Success 200 {object} domain.StatsOverview
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/overview [get]
func (h *StatsHandler) Overview(c *gin.Context) {
	userID := middleware.GetUserID(c)
	// Total cards is passed as a query param from the client (pre-fetched)
	total, _ := strconv.ParseInt(c.DefaultQuery("totalCards", "0"), 10, 64)
	overview, err := h.statsSvc.Overview(c.Request.Context(), userID, total)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, overview)
}

// Activity godoc
// @Summary Get learning activity
// @Description Get the learning activity over a specified number of days
// @Tags Stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days" default(30)
// @Success 200 {array} domain.DailyActivity
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/activity [get]
func (h *StatsHandler) Activity(c *gin.Context) {
	userID := middleware.GetUserID(c)
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	activity, err := h.statsSvc.Activity(c.Request.Context(), userID, days)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, activity)
}

// Forecast godoc
// @Summary Get review forecast
// @Description Get a forecast of upcoming reviews for the next few days
// @Tags Stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days" default(7)
// @Success 200 {array} domain.ForecastDay
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stats/forecast [get]
func (h *StatsHandler) Forecast(c *gin.Context) {
	userID := middleware.GetUserID(c)
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	forecast, err := h.statsSvc.Forecast(c.Request.Context(), userID, days)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, forecast)
}

// DeckStats godoc
// @Summary Get deck statistics
// @Description Get learning statistics for a specific deck
// @Tags Stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 200 {object} domain.DeckStats
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/stats [get]
func (h *StatsHandler) DeckStats(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	stats, err := h.statsSvc.DeckStats(c.Request.Context(), deckID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
