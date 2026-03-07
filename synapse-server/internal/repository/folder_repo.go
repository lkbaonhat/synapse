package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// FolderRepository defines DB operations for folders.
type FolderRepository interface {
	List(ctx context.Context, userID uuid.UUID) ([]domain.Folder, error)
	FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Folder, error)
	Create(ctx context.Context, folder *domain.Folder) error
	Update(ctx context.Context, folder *domain.Folder) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type folderRepo struct{ db *gorm.DB }

func NewFolderRepository(db *gorm.DB) FolderRepository { return &folderRepo{db: db} }

func (r *folderRepo) List(ctx context.Context, userID uuid.UUID) ([]domain.Folder, error) {
	var folders []domain.Folder
	return folders, r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&folders).Error
}

func (r *folderRepo) FindByID(ctx context.Context, id, userID uuid.UUID) (*domain.Folder, error) {
	var f domain.Folder
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &f, err
}

func (r *folderRepo) Create(ctx context.Context, folder *domain.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *folderRepo) Update(ctx context.Context, folder *domain.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

func (r *folderRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Folder{}).Error
}
