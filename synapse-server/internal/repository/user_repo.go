package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// UserRepository defines DB operations for users and refresh tokens.
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	FindRefreshToken(ctx context.Context, rawToken string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, rawToken string) error
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type userRepo struct{ db *gorm.DB }

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func (r *userRepo) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	token.TokenHash = hashToken(token.TokenHash)
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *userRepo) FindRefreshToken(ctx context.Context, rawToken string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	hash := hashToken(rawToken)
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *userRepo) DeleteRefreshToken(ctx context.Context, rawToken string) error {
	hash := hashToken(rawToken)
	return r.db.WithContext(ctx).Where("token_hash = ?", hash).Delete(&domain.RefreshToken{}).Error
}

func (r *userRepo) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}
