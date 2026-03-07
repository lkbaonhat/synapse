package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/apierror"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/repository"
)

// DeckService handles folder, deck, and tag business logic.
type DeckService interface {
	// Folders
	ListFolders(ctx context.Context, userID uuid.UUID) ([]domain.Folder, error)
	GetFolder(ctx context.Context, id, userID uuid.UUID) (*domain.Folder, error)
	CreateFolder(ctx context.Context, folder *domain.Folder) error
	UpdateFolder(ctx context.Context, id, userID uuid.UUID, name string, parentID *uuid.UUID) (*domain.Folder, error)
	DeleteFolder(ctx context.Context, id, userID uuid.UUID) error
	// Decks
	ListDecks(ctx context.Context, userID uuid.UUID, folderID *uuid.UUID, tagID *uuid.UUID) ([]domain.Deck, int64, error)
	GetDeck(ctx context.Context, id, userID uuid.UUID) (*domain.Deck, error)
	CreateDeck(ctx context.Context, deck *domain.Deck) error
	UpdateDeck(ctx context.Context, id, userID uuid.UUID, updates *domain.Deck) (*domain.Deck, error)
	DeleteDeck(ctx context.Context, id, userID uuid.UUID) error
	AttachTags(ctx context.Context, deckID, userID uuid.UUID, tagIDs []uuid.UUID) error
	// Tags
	ListTags(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error)
	CreateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id, userID uuid.UUID) error
}

type deckService struct {
	folderRepo repository.FolderRepository
	deckRepo   repository.DeckRepository
	tagRepo    repository.TagRepository
}

func NewDeckService(folderRepo repository.FolderRepository, deckRepo repository.DeckRepository, tagRepo repository.TagRepository) DeckService {
	return &deckService{folderRepo: folderRepo, deckRepo: deckRepo, tagRepo: tagRepo}
}

// ----- Folders -----

func (s *deckService) ListFolders(ctx context.Context, userID uuid.UUID) ([]domain.Folder, error) {
	return s.folderRepo.List(ctx, userID)
}

func (s *deckService) GetFolder(ctx context.Context, id, userID uuid.UUID) (*domain.Folder, error) {
	f, err := s.folderRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, apierror.NotFound("folder not found")
	}
	return f, nil
}

func (s *deckService) CreateFolder(ctx context.Context, folder *domain.Folder) error {
	return s.folderRepo.Create(ctx, folder)
}

func (s *deckService) UpdateFolder(ctx context.Context, id, userID uuid.UUID, name string, parentID *uuid.UUID) (*domain.Folder, error) {
	f, err := s.GetFolder(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		f.Name = name
	}
	f.ParentID = parentID
	return f, s.folderRepo.Update(ctx, f)
}

func (s *deckService) DeleteFolder(ctx context.Context, id, userID uuid.UUID) error {
	_, err := s.GetFolder(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.folderRepo.Delete(ctx, id, userID)
}

// ----- Decks -----

func (s *deckService) ListDecks(ctx context.Context, userID uuid.UUID, folderID *uuid.UUID, tagID *uuid.UUID) ([]domain.Deck, int64, error) {
	return s.deckRepo.List(ctx, userID, folderID, tagID)
}

func (s *deckService) GetDeck(ctx context.Context, id, userID uuid.UUID) (*domain.Deck, error) {
	d, err := s.deckRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, apierror.NotFound("deck not found")
	}
	return d, nil
}

func (s *deckService) CreateDeck(ctx context.Context, deck *domain.Deck) error {
	return s.deckRepo.Create(ctx, deck)
}

func (s *deckService) UpdateDeck(ctx context.Context, id, userID uuid.UUID, updates *domain.Deck) (*domain.Deck, error) {
	d, err := s.GetDeck(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if updates.Name != "" {
		d.Name = updates.Name
	}
	d.Description = updates.Description
	d.FolderID = updates.FolderID
	return d, s.deckRepo.Update(ctx, d)
}

func (s *deckService) DeleteDeck(ctx context.Context, id, userID uuid.UUID) error {
	_, err := s.GetDeck(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.deckRepo.Delete(ctx, id, userID)
}

func (s *deckService) AttachTags(ctx context.Context, deckID, userID uuid.UUID, tagIDs []uuid.UUID) error {
	deck, err := s.GetDeck(ctx, deckID, userID)
	if err != nil {
		return err
	}
	var tags []domain.Tag
	for _, tid := range tagIDs {
		t, err := s.tagRepo.FindByID(ctx, tid, userID)
		if err != nil {
			return err
		}
		if t == nil {
			return apierror.NotFound("tag not found: " + tid.String())
		}
		tags = append(tags, *t)
	}
	return s.deckRepo.AddTags(ctx, deck, tags)
}

// ----- Tags -----

func (s *deckService) ListTags(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	return s.tagRepo.List(ctx, userID)
}

func (s *deckService) CreateTag(ctx context.Context, tag *domain.Tag) error {
	return s.tagRepo.Create(ctx, tag)
}

func (s *deckService) DeleteTag(ctx context.Context, id, userID uuid.UUID) error {
	t, err := s.tagRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if t == nil {
		return apierror.NotFound("tag not found")
	}
	return s.tagRepo.Delete(ctx, id, userID)
}
