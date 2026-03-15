package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/apierror"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/repository"
	"github.com/synapse/server/pkg/srs"
)

const defaultStudyBatch = 20

// StudyService manages study sessions and card scheduling.
type StudyService interface {
	StartSession(ctx context.Context, userID, deckID uuid.UUID, mode domain.StudyMode) (*domain.StudySession, []domain.Card, error)
	NextCards(ctx context.Context, sessionID, userID uuid.UUID) ([]domain.Card, error)
	Answer(ctx context.Context, sessionID, userID, cardID uuid.UUID, rating int, timeTaken int) error
	EndSession(ctx context.Context, sessionID, userID uuid.UUID) error
	GetQuizResult(ctx context.Context, sessionID, userID uuid.UUID) (*domain.QuizResult, error)
}

type studyService struct {
	studyRepo repository.StudyRepository
	cardRepo  repository.CardRepository
	deckRepo  repository.DeckRepository
}

func NewStudyService(studyRepo repository.StudyRepository, cardRepo repository.CardRepository, deckRepo repository.DeckRepository) StudyService {
	return &studyService{studyRepo: studyRepo, cardRepo: cardRepo, deckRepo: deckRepo}
}

func (s *studyService) StartSession(ctx context.Context, userID, deckID uuid.UUID, mode domain.StudyMode) (*domain.StudySession, []domain.Card, error) {
	// Verify deck ownership
	deck, err := s.deckRepo.FindByID(ctx, deckID, userID)
	if err != nil {
		return nil, nil, err
	}
	if deck == nil {
		return nil, nil, apierror.NotFound("deck not found")
	}

	session := &domain.StudySession{
		UserID:    userID,
		DeckID:    deckID,
		Mode:      mode,
		StartedAt: time.Now().UTC(),
	}
	if err := s.studyRepo.CreateSession(ctx, session); err != nil {
		return nil, nil, err
	}

	cards, err := s.fetchCards(ctx, deckID, mode, defaultStudyBatch)
	if err != nil {
		return nil, nil, err
	}
	return session, cards, nil
}

func (s *studyService) NextCards(ctx context.Context, sessionID, userID uuid.UUID) ([]domain.Card, error) {
	session, err := s.studyRepo.FindSession(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, apierror.NotFound("session not found")
	}
	return s.fetchCards(ctx, session.DeckID, session.Mode, defaultStudyBatch)
}

func (s *studyService) Answer(ctx context.Context, sessionID, userID, cardID uuid.UUID, rating, timeTaken int) error {
	session, err := s.studyRepo.FindSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return apierror.NotFound("session not found")
	}

	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return err
	}
	if card == nil {
		return apierror.NotFound("card not found")
	}

	// Log the answer
	log := &domain.StudyLog{
		SessionID: sessionID,
		CardID:    cardID,
		Rating:    rating,
		TimeTaken: timeTaken,
		LoggedAt:  time.Now().UTC(),
	}
	if err := s.studyRepo.CreateLog(ctx, log); err != nil {
		return err
	}

	// Update SRS schedule (skip for cram mode)
	if session.Mode != domain.StudyModeCram {
		current := srs.CardSchedule{
			Interval:    card.Interval,
			Easiness:    card.Easiness,
			Repetitions: card.Repetitions,
		}
		if card.DueAt != nil {
			current.DueAt = *card.DueAt
		}
		next := srs.Compute(current, srs.DifficultyRating(rating))
		card.Interval = next.Interval
		card.Easiness = next.Easiness
		card.Repetitions = next.Repetitions
		card.DueAt = &next.DueAt
		if err := s.cardRepo.Update(ctx, card); err != nil {
			return err
		}
	}
	return nil
}

func (s *studyService) EndSession(ctx context.Context, sessionID, userID uuid.UUID) error {
	session, err := s.studyRepo.FindSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return apierror.NotFound("session not found")
	}
	now := time.Now().UTC()
	return s.studyRepo.EndSession(ctx, sessionID, now)
}

func (s *studyService) fetchCards(ctx context.Context, deckID uuid.UUID, mode domain.StudyMode, limit int) ([]domain.Card, error) {
	switch mode {
	case domain.StudyModeLearn:
		return s.cardRepo.FindNewCards(ctx, deckID, limit)
	case domain.StudyModeReview:
		return s.cardRepo.FindDueCards(ctx, deckID, limit)
	case domain.StudyModeCram:
		return s.cardRepo.FindAllCards(ctx, deckID, limit)
	default:
		return nil, apierror.BadRequest("unknown study mode: " + string(mode))
	}
}

func (s *studyService) GetQuizResult(ctx context.Context, sessionID, userID uuid.UUID) (*domain.QuizResult, error) {
	session, err := s.studyRepo.FindSession(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, apierror.NotFound("session not found")
	}

	logs, err := s.studyRepo.FindLogsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	result := &domain.QuizResult{
		TotalCorrect: 0,
		TotalWrong:   0,
		WrongAnswers: make([]domain.WrongAnswerSummary, 0),
	}

	for _, log := range logs {
		if log.Rating >= 3 {
			result.TotalCorrect++
		} else {
			result.TotalWrong++
			card, err := s.cardRepo.FindByID(ctx, log.CardID)
			if err != nil || card == nil {
				continue
			}

			var content map[string]interface{}
			_ = json.Unmarshal(card.Content, &content)

			front := ""
			back := ""

			switch card.Type {
			case domain.CardTypeFlashcard:
				if f, ok := content["front"].(string); ok {
					front = f
				}
				if b, ok := content["back"].(string); ok {
					back = b
				}
			case domain.CardTypeCloze:
				if t, ok := content["text"].(string); ok {
					front = t
				}
			case domain.CardTypeFreeResponse:
				if p, ok := content["prompt"].(string); ok {
					front = p
				}
				if a, ok := content["answer"].(string); ok {
					back = a
				}
			case domain.CardTypeMultipleChoice:
				if p, ok := content["prompt"].(string); ok {
					front = p
				}
				if opts, ok := content["options"].([]interface{}); ok {
					if idxFloat, ok := content["correctIndex"].(float64); ok {
						idx := int(idxFloat)
						if idx >= 0 && idx < len(opts) {
							if optStr, ok := opts[idx].(string); ok {
								back = optStr
							}
						}
					}
				}
			}

			result.WrongAnswers = append(result.WrongAnswers, domain.WrongAnswerSummary{
				CardID:      card.ID,
				Front:       front,
				CorrectBack: back,
			})
		}
	}

	return result, nil
}
