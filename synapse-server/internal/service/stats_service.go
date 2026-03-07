package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/repository"
)

// StatsService computes statistics from study logs.
type StatsService interface {
	Overview(ctx context.Context, userID uuid.UUID, totalCards int64) (*domain.StatsOverview, error)
	Activity(ctx context.Context, userID uuid.UUID, days int) ([]domain.DailyActivity, error)
	Forecast(ctx context.Context, userID uuid.UUID, days int) ([]domain.ForecastDay, error)
	DeckStats(ctx context.Context, deckID uuid.UUID) (*domain.DeckStats, error)
}

type statsService struct {
	studyRepo repository.StudyRepository
	cardRepo  repository.CardRepository
}

func NewStatsService(studyRepo repository.StudyRepository, cardRepo repository.CardRepository) StatsService {
	return &statsService{studyRepo: studyRepo, cardRepo: cardRepo}
}

func (s *statsService) Overview(ctx context.Context, userID uuid.UUID, totalCards int64) (*domain.StatsOverview, error) {
	total, err := s.studyRepo.TotalStudied(ctx, userID)
	if err != nil {
		return nil, err
	}
	retention, err := s.studyRepo.RetentionRate(ctx, userID)
	if err != nil {
		return nil, err
	}
	streak, err := s.computeStreak(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &domain.StatsOverview{
		TotalCards:    totalCards,
		RetentionRate: retention,
		CurrentStreak: streak,
		TotalStudied:  total,
	}, nil
}

func (s *statsService) Activity(ctx context.Context, userID uuid.UUID, days int) ([]domain.DailyActivity, error) {
	if days <= 0 {
		days = 30
	}
	return s.studyRepo.DailyActivity(ctx, userID, days)
}

func (s *statsService) Forecast(ctx context.Context, userID uuid.UUID, days int) ([]domain.ForecastDay, error) {
	if days <= 0 {
		days = 7
	}
	return s.studyRepo.Forecast(ctx, userID, days)
}

func (s *statsService) DeckStats(ctx context.Context, deckID uuid.UUID) (*domain.DeckStats, error) {
	newCount, learning, review, mastered, err := s.cardRepo.CountByDeck(ctx, deckID)
	if err != nil {
		return nil, err
	}
	return &domain.DeckStats{
		DeckID:   deckID.String(),
		New:      newCount,
		Learning: learning,
		Review:   review,
		Mastered: mastered,
	}, nil
}

// computeStreak counts consecutive calendar days with at least one log.
func (s *statsService) computeStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	activity, err := s.studyRepo.DailyActivity(ctx, userID, 365)
	if err != nil {
		return 0, err
	}
	if len(activity) == 0 {
		return 0, nil
	}
	// Count backwards from today
	streak := 0
	for i := len(activity) - 1; i >= 0; i-- {
		if activity[i].Count > 0 {
			streak++
		} else {
			break
		}
	}
	return streak, nil
}
