package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// CardRepository defines DB operations for cards.
type CardRepository interface {
	ListByDeck(ctx context.Context, deckID uuid.UUID, offset, limit int) ([]domain.Card, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Card, error)
	Create(ctx context.Context, card *domain.Card) error
	Update(ctx context.Context, card *domain.Card) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountDue(ctx context.Context, deckID uuid.UUID) (int64, error)
	// Study helpers
	FindNewCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error)
	FindDueCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error)
	FindAllCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error)
	// Stats helpers
	CountByDeck(ctx context.Context, deckID uuid.UUID) (new, learning, review, mastered int64, err error)
}

type cardRepo struct{ db *gorm.DB }

func NewCardRepository(db *gorm.DB) CardRepository { return &cardRepo{db: db} }

func (r *cardRepo) ListByDeck(ctx context.Context, deckID uuid.UUID, offset, limit int) ([]domain.Card, int64, error) {
	var cards []domain.Card
	var total int64
	q := r.db.WithContext(ctx).Preload("Tags").Where("deck_id = ?", deckID)
	if err := q.Model(&domain.Card{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return cards, total, q.Offset(offset).Limit(limit).Find(&cards).Error
}

func (r *cardRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Card, error) {
	var c domain.Card
	err := r.db.WithContext(ctx).Preload("Tags").First(&c, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *cardRepo) Create(ctx context.Context, card *domain.Card) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *cardRepo) Update(ctx context.Context, card *domain.Card) error {
	return r.db.WithContext(ctx).Save(card).Error
}

func (r *cardRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Card{}, "id = ?", id).Error
}

func (r *cardRepo) CountDue(ctx context.Context, deckID uuid.UUID) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&domain.Card{}).
		Where("deck_id = ? AND due_at IS NOT NULL AND due_at <= ?", deckID, time.Now()).
		Count(&count).Error
}

func (r *cardRepo) FindNewCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error) {
	var cards []domain.Card
	return cards, r.db.WithContext(ctx).
		Where("deck_id = ? AND due_at IS NULL", deckID).
		Limit(limit).Find(&cards).Error
}

func (r *cardRepo) FindDueCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error) {
	var cards []domain.Card
	return cards, r.db.WithContext(ctx).
		Where("deck_id = ? AND due_at IS NOT NULL AND due_at <= ?", deckID, time.Now()).
		Order("due_at ASC").Limit(limit).Find(&cards).Error
}

func (r *cardRepo) FindAllCards(ctx context.Context, deckID uuid.UUID, limit int) ([]domain.Card, error) {
	var cards []domain.Card
	return cards, r.db.WithContext(ctx).
		Where("deck_id = ?", deckID).
		Limit(limit).Find(&cards).Error
}

func (r *cardRepo) CountByDeck(ctx context.Context, deckID uuid.UUID) (new, learning, review, mastered int64, err error) {
	base := r.db.WithContext(ctx).Model(&domain.Card{}).Where("deck_id = ?", deckID)
	base.Where("due_at IS NULL").Count(&new)
	base.Where("due_at IS NOT NULL AND repetitions < 3").Count(&learning)
	base.Where("repetitions >= 3 AND interval < 21").Count(&review)
	base.Where("interval >= 21").Count(&mastered)
	return
}
