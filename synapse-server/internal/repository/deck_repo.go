package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// DeckRepository defines DB operations for decks.
type DeckRepository interface {
	List(ctx context.Context, userID uuid.UUID, folderID *uuid.UUID, tagID *uuid.UUID) ([]domain.Deck, int64, error)
	FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Deck, error)
	Create(ctx context.Context, deck *domain.Deck) error
	Update(ctx context.Context, deck *domain.Deck) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	AddTags(ctx context.Context, deck *domain.Deck, tags []domain.Tag) error
}

type deckRepo struct{ db *gorm.DB }

func NewDeckRepository(db *gorm.DB) DeckRepository { return &deckRepo{db: db} }

func (r *deckRepo) List(ctx context.Context, userID uuid.UUID, folderID *uuid.UUID, tagID *uuid.UUID) ([]domain.Deck, int64, error) {
	q := r.db.WithContext(ctx).Preload("Tags").Where("user_id = ?", userID)
	if folderID != nil {
		q = q.Where("folder_id = ?", folderID)
	}
	if tagID != nil {
		q = q.Joins("JOIN deck_tags ON deck_tags.deck_id = decks.id").
			Where("deck_tags.tag_id = ?", tagID)
	}
	var decks []domain.Deck
	var total int64
	if err := q.Model(&domain.Deck{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return decks, total, q.Find(&decks).Error
}

func (r *deckRepo) FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Deck, error) {
	var d domain.Deck
	err := r.db.WithContext(ctx).Preload("Tags").
		Where("id = ? AND user_id = ?", id, userID).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &d, err
}

func (r *deckRepo) Create(ctx context.Context, deck *domain.Deck) error {
	return r.db.WithContext(ctx).Create(deck).Error
}

func (r *deckRepo) Update(ctx context.Context, deck *domain.Deck) error {
	return r.db.WithContext(ctx).Save(deck).Error
}

func (r *deckRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Deck{}).Error
}

func (r *deckRepo) AddTags(ctx context.Context, deck *domain.Deck, tags []domain.Tag) error {
	return r.db.WithContext(ctx).Model(deck).Association("Tags").Replace(tags)
}
