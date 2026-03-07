package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
