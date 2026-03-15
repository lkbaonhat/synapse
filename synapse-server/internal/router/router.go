package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/synapse/server/docs" // imported for swagger docs
	"github.com/synapse/server/internal/handler"
	"github.com/synapse/server/internal/middleware"
)

// Setup registers all routes and returns the Gin engine.
func Setup(
	jwtSecret string,
	allowedOrigins []string,
	authH *handler.AuthHandler,
	deckH *handler.DeckHandler,
	cardH *handler.CardHandler,
	studyH *handler.StudyHandler,
	statsH *handler.StatsHandler,
	ieH *handler.ImportExportHandler,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (unauthenticated)
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (unauthenticated, with rate limiting)
	auth := r.Group("/api/v1/auth")
	auth.Use(middleware.RateLimit(5, 10))
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", authH.Logout)
	}

	// Protected routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthRequired(jwtSecret))
	{
		// Folders
		v1.GET("/folders", deckH.ListFolders)
		v1.POST("/folders", deckH.CreateFolder)
		v1.GET("/folders/:id", deckH.GetFolder)
		v1.PUT("/folders/:id", deckH.UpdateFolder)
		v1.DELETE("/folders/:id", deckH.DeleteFolder)

		// Decks
		v1.GET("/decks", deckH.ListDecks)
		v1.POST("/decks", deckH.CreateDeck)
		v1.GET("/decks/:id", deckH.GetDeck)
		v1.PUT("/decks/:id", deckH.UpdateDeck)
		v1.DELETE("/decks/:id", deckH.DeleteDeck)
		v1.POST("/decks/:id/tags", deckH.AttachTags)
		v1.GET("/decks/:id/stats", statsH.DeckStats)
		v1.GET("/decks/:id/due-count", cardH.DueCount)

		// Deck cards
		v1.GET("/decks/:id/cards", cardH.ListCards)
		v1.POST("/decks/:id/cards", cardH.CreateCard)
		v1.POST("/decks/:id/import", ieH.ImportCSV)
		v1.GET("/decks/:id/export", ieH.ExportDeckCSV)

		// Cards
		v1.GET("/cards/:id", cardH.GetCard)
		v1.PUT("/cards/:id", cardH.UpdateCard)
		v1.DELETE("/cards/:id", cardH.DeleteCard)
		v1.POST("/cards/:id/media", cardH.UploadMedia)

		// Tags
		v1.GET("/tags", deckH.ListTags)
		v1.POST("/tags", deckH.CreateTag)
		v1.DELETE("/tags/:id", deckH.DeleteTag)

		// Study sessions
		v1.POST("/study/sessions", studyH.StartSession)
		v1.GET("/study/sessions/:id/next", studyH.NextCards)
		v1.POST("/study/sessions/:id/answer", studyH.Answer)
		v1.POST("/study/sessions/:id/end", studyH.EndSession)
		v1.GET("/study/sessions/:id/results", studyH.GetQuizResult)

		// Statistics
		v1.GET("/stats/overview", statsH.Overview)
		v1.GET("/stats/activity", statsH.Activity)
		v1.GET("/stats/forecast", statsH.Forecast)

		// Data export
		v1.GET("/user/export", ieH.ExportUserJSON)
	}

	return r
}
