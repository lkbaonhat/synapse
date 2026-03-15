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

// ListFolders godoc
// @Summary List all folders
// @Description Get a list of all folders for the authenticated user
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Folder
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders [get]
func (h *DeckHandler) ListFolders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	folders, err := h.svc.ListFolders(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, folders)
}

type createFolderRequest struct {
	Name     string     `json:"name"     binding:"required"`
	ParentID *uuid.UUID `json:"parentId"`
}

// CreateFolder godoc
// @Summary Create a new folder
// @Description Create a new folder with an optional parent folder ID
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createFolderRequest true "Folder Information"
// @Success 201 {object} domain.Folder
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders [post]
func (h *DeckHandler) CreateFolder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body createFolderRequest
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

// GetFolder godoc
// @Summary Get a folder by ID
// @Description Get details of a specific folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 200 {object} domain.Folder
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders/{id} [get]
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

type updateFolderRequest struct {
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parentId"`
}

// UpdateFolder godoc
// @Summary Update a folder
// @Description Update a folder's name or parent directory
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Param request body updateFolderRequest true "Folder Information"
// @Success 200 {object} domain.Folder
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders/{id} [put]
func (h *DeckHandler) UpdateFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := middleware.GetUserID(c)
	var body updateFolderRequest
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

// DeleteFolder godoc
// @Summary Delete a folder
// @Description Delete a specific folder
// @Tags Folders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /folders/{id} [delete]
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

// ListDecks godoc
// @Summary List all decks
// @Description Get a paginated list of decks for the authenticated user, optionally filtered by folder or tag
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param folderId query string false "Filter by Folder ID"
// @Param tagId query string false "Filter by Tag ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} domain.Deck
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks [get]
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

type createDeckRequest struct {
	Name        string     `json:"name"        binding:"required"`
	Description string     `json:"description"`
	FolderID    *uuid.UUID `json:"folderId"`
}

// CreateDeck godoc
// @Summary Create a new deck
// @Description Create a new deck with a name, optional description, and optional folder ID
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createDeckRequest true "Deck Information"
// @Success 201 {object} domain.Deck
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks [post]
func (h *DeckHandler) CreateDeck(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body createDeckRequest
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

// GetDeck godoc
// @Summary Get a deck by ID
// @Description Get details of a specific deck
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 200 {object} domain.Deck
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id} [get]
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

// UpdateDeck godoc
// @Summary Update a deck
// @Description Update a deck's information
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param request body domain.Deck true "Deck Information"
// @Success 200 {object} domain.Deck
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id} [put]
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

// DeleteDeck godoc
// @Summary Delete a deck
// @Description Delete a specific deck
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id} [delete]
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

type attachTagsRequest struct {
	TagIDs []uuid.UUID `json:"tagIds" binding:"required"`
}

// AttachTags godoc
// @Summary Attach tags to a deck
// @Description Attach one or more tags to a specific deck
// @Tags Decks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deck ID"
// @Param request body attachTagsRequest true "Tag IDs"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /decks/{id}/tags [post]
func (h *DeckHandler) AttachTags(c *gin.Context) {
	deckID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}
	var body attachTagsRequest
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

// ListTags godoc
// @Summary List all tags
// @Description Get a list of all tags for the authenticated user
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Tag
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags [get]
func (h *DeckHandler) ListTags(c *gin.Context) {
	tags, err := h.svc.ListTags(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, tags)
}

type createTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag for organizing decks
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createTagRequest true "Tag Information"
// @Success 201 {object} domain.Tag
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags [post]
func (h *DeckHandler) CreateTag(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var body createTagRequest
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

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a specific tag
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags/{id} [delete]
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
