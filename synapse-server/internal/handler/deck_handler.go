package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/pagination"
	"github.com/synapse/server/internal/service"
)

// DeckHandler handles folder, deck, and tag routes.
type DeckHandler struct {
	svc service.DeckService
}

func NewDeckHandler(svc service.DeckService) *DeckHandler { return &DeckHandler{svc: svc} }

// ----- Folders -----

func (h *DeckHandler) ListFolders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	folders, err := h.svc.ListFolders(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, folders)
}

func (h *DeckHandler) CreateFolder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body struct {
		Name     string     `json:"name"     binding:"required"`
		ParentID *uuid.UUID `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	folder := &domain.Folder{UserID: userID, Name: body.Name, ParentID: body.ParentID}
	if err := h.svc.CreateFolder(c.Request.Context(), folder); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, folder)
}

func (h *DeckHandler) GetFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := middleware.GetUserID(c)
	f, err := h.svc.GetFolder(c.Request.Context(), id, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *DeckHandler) UpdateFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := middleware.GetUserID(c)
	var body struct {
		Name     string     `json:"name"`
		ParentID *uuid.UUID `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := h.svc.UpdateFolder(c.Request.Context(), id, userID, body.Name, body.ParentID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *DeckHandler) DeleteFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := middleware.GetUserID(c)
	if err := h.svc.DeleteFolder(c.Request.Context(), id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ----- Decks -----

func (h *DeckHandler) ListDecks(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var folderID, tagID *uuid.UUID
	if fStr := c.Query("folderId"); fStr != "" {
		if id, err := uuid.Parse(fStr); err == nil {
			folderID = &id
		}
	}
	if tStr := c.Query("tagId"); tStr != "" {
		if id, err := uuid.Parse(tStr); err == nil {
			tagID = &id
		}
	}
	decks, total, err := h.svc.ListDecks(c.Request.Context(), userID, folderID, tagID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagination.SetTotalCount(c, total)
	c.JSON(http.StatusOK, decks)
}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body struct {
		Name        string     `json:"name"        binding:"required"`
		Description string     `json:"description"`
		FolderID    *uuid.UUID `json:"folderId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deck := &domain.Deck{
		UserID: userID, Name: body.Name,
		Description: body.Description, FolderID: body.FolderID,
	}
	if err := h.svc.CreateDeck(c.Request.Context(), deck); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, deck)
}

func (h *DeckHandler) GetDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	d, err := h.svc.GetDeck(c.Request.Context(), id, middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body domain.Deck
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.svc.UpdateDeck(c.Request.Context(), id, middleware.GetUserID(c), &body)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteDeck(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DeckHandler) AttachTags(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	var body struct {
		TagIDs []uuid.UUID `json:"tagIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AttachTags(c.Request.Context(), deckID, middleware.GetUserID(c), body.TagIDs); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ----- Tags -----

func (h *DeckHandler) ListTags(c *gin.Context) {
	tags, err := h.svc.ListTags(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, tags)
}

func (h *DeckHandler) CreateTag(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag := &domain.Tag{UserID: userID, Name: body.Name}
	if err := h.svc.CreateTag(c.Request.Context(), tag); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, tag)
}

func (h *DeckHandler) DeleteTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteTag(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
