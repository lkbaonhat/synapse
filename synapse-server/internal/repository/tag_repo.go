package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// TagRepository defines DB operations for tags.
type TagRepository interface {
	List(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error)
	FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Tag, error)
	FindOrCreate(ctx context.Context, userID uuid.UUID, name string) (*domain.Tag, error)
	Create(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type tagRepo struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) TagRepository { return &tagRepo{db: db} }

func (r *tagRepo) List(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	var tags []domain.Tag
	return tags, r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tags).Error
}

func (r *tagRepo) FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Tag, error) {
	var t domain.Tag
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *tagRepo) FindOrCreate(ctx context.Context, userID uuid.UUID, name string) (*domain.Tag, error) {
	var t domain.Tag
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", userID, name).
		First(&t).Error
	if err == nil {
		return &t, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	t = domain.Tag{UserID: userID, Name: name}
	return &t, r.db.WithContext(ctx).Create(&t).Error
}

func (r *tagRepo) Create(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Tag{}).Error
}
