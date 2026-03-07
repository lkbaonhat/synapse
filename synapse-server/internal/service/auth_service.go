package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/synapse/server/internal/apierror"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/middleware"
	"github.com/synapse/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles registration, login, token refresh, and logout.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*TokenPair, error)
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Refresh(ctx context.Context, rawRefreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

// TokenPair groups the access and refresh tokens returned to the client.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type authService struct {
	userRepo         repository.UserRepository
	jwtSecret        string
	accessTTLMinutes int
	refreshTTLDays   int
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, accessTTL, refreshTTL int) AuthService {
	return &authService{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		accessTTLMinutes: accessTTL,
		refreshTTLDays:   refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*TokenPair, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apierror.Conflict("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("bcrypt: %w", err)
	}

	user := &domain.User{Email: email, PasswordHash: string(hash)}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apierror.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, apierror.Unauthorized("invalid credentials")
		}
		return nil, err
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) Refresh(ctx context.Context, rawRefreshToken string) (*TokenPair, error) {
	stored, err := s.userRepo.FindRefreshToken(ctx, rawRefreshToken)
	if err != nil {
		return nil, err
	}
	if stored == nil || time.Now().After(stored.ExpiresAt) {
		return nil, apierror.Unauthorized("invalid or expired refresh token")
	}

	// Rotate: delete old, issue new
	if err := s.userRepo.DeleteRefreshToken(ctx, rawRefreshToken); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, stored.UserID)
}

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	return s.userRepo.DeleteRefreshToken(ctx, rawRefreshToken)
}

// issueTokens generates and persists a fresh access + refresh token pair.
func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	// Access token
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	// Refresh token — random 32-byte opaque token
	rawRefresh, err := randomBase64(32)
	if err != nil {
		return nil, err
	}

	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: rawRefresh, // repo will hash it
		ExpiresAt: time.Now().Add(time.Duration(s.refreshTTLDays) * 24 * time.Hour),
	}
	if err := s.userRepo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: rawRefresh}, nil
}

func (s *authService) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.accessTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func randomBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
