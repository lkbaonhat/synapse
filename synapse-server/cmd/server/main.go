package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/synapse/server/internal/config"
	"github.com/synapse/server/internal/database"
	"github.com/synapse/server/internal/handler"
	"github.com/synapse/server/internal/repository"
	"github.com/synapse/server/internal/router"
	"github.com/synapse/server/internal/service"
)

// @title Synapse API
// @version 1.0
// @description Backend API for the Synapse Study Assistant.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// -- Config --
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// -- Database --
	db, err := database.Open(cfg.DSN(), cfg.Env == "development")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Note: AutoMigrate is commented out because we are using manual SQL 
	// migrations (create_database.sql / tools like golang-migrate).
	// Mixing AutoMigrate with manual SQL migrations often causes constraint naming conflicts.
	/* 
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.Folder{},
		&domain.Tag{},
		&domain.Deck{},
		&domain.Card{},
		&domain.StudySession{},
		&domain.StudyLog{},
	); err != nil {
		slog.Error("auto-migrate failed", "error", err)
		os.Exit(1)
	}
	*/
	// -- Repositories --
	userRepo := repository.NewUserRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	tagRepo := repository.NewTagRepository(db)
	deckRepo := repository.NewDeckRepository(db)
	cardRepo := repository.NewCardRepository(db)
	studyRepo := repository.NewStudyRepository(db)

	// -- Services --
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTAccessTTLMinutes, cfg.JWTRefreshTTLDays)
	deckSvc := service.NewDeckService(folderRepo, deckRepo, tagRepo)
	cardSvc := service.NewCardService(cardRepo, deckRepo)
	studySvc := service.NewStudyService(studyRepo, cardRepo, deckRepo)
	statsSvc := service.NewStatsService(studyRepo, cardRepo)
	ieSvc := service.NewImportExportService(cardRepo, deckRepo, tagRepo)

	// -- Handlers --
	authH := handler.NewAuthHandler(authSvc)
	deckH := handler.NewDeckHandler(deckSvc)
	cardH := handler.NewCardHandler(cardSvc, studySvc, cfg.UploadDir)
	studyH := handler.NewStudyHandler(studySvc)
	statsH := handler.NewStatsHandler(statsSvc)
	ieH := handler.NewImportExportHandler(ieSvc)

	// -- Router --
	r := router.Setup(cfg.JWTSecret, cfg.AllowedOrigins(), authH, deckH, cardH, studyH, statsH, ieH)

	// -- Start server --
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("starting synapse server", "addr", addr, "env", cfg.Env)
	if err := r.Run(addr); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
