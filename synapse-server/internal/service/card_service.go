package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/apierror"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/repository"
)

// CardService handles flashcard business logic including content validation.
type CardService interface {
	ListByDeck(ctx context.Context, deckID, userID uuid.UUID, offset, limit int) ([]domain.Card, int64, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Card, error)
	Create(ctx context.Context, card *domain.Card) error
	Update(ctx context.Context, id, userID uuid.UUID, updates *domain.Card) (*domain.Card, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	CountDue(ctx context.Context, deckID uuid.UUID) (int64, error)
}

type cardService struct {
	cardRepo repository.CardRepository
	deckRepo repository.DeckRepository
}

func NewCardService(cardRepo repository.CardRepository, deckRepo repository.DeckRepository) CardService {
	return &cardService{cardRepo: cardRepo, deckRepo: deckRepo}
}

func (s *cardService) ListByDeck(ctx context.Context, deckID, userID uuid.UUID, offset, limit int) ([]domain.Card, int64, error) {
	// Verify deck ownership
	deck, err := s.deckRepo.FindByID(ctx, deckID, userID)
	if err != nil {
		return nil, 0, err
	}
	if deck == nil {
		return nil, 0, apierror.NotFound("deck not found")
	}
	return s.cardRepo.ListByDeck(ctx, deckID, offset, limit)
}

func (s *cardService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Card, error) {
	card, err := s.cardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, apierror.NotFound("card not found")
	}
	// Verify ownership via deck
	deck, err := s.deckRepo.FindByID(ctx, card.DeckID, userID)
	if err != nil {
		return nil, err
	}
	if deck == nil {
		return nil, apierror.Forbidden("access denied")
	}
	return card, nil
}

func (s *cardService) Create(ctx context.Context, card *domain.Card) error {
	if err := validateCardContent(card.Type, card.Content); err != nil {
		return err
	}
	if card.Easiness == 0 {
		card.Easiness = 2.5
	}
	return s.cardRepo.Create(ctx, card)
}

func (s *cardService) Update(ctx context.Context, id, userID uuid.UUID, updates *domain.Card) (*domain.Card, error) {
	card, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if updates.Type != "" {
		card.Type = updates.Type
	}
	if len(updates.Content) > 0 {
		card.Content = updates.Content
	}
	if err := validateCardContent(card.Type, card.Content); err != nil {
		return nil, err
	}
	return card, s.cardRepo.Update(ctx, card)
}

func (s *cardService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	_, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.cardRepo.Delete(ctx, id)
}

func (s *cardService) CountDue(ctx context.Context, deckID uuid.UUID) (int64, error) {
	return s.cardRepo.CountDue(ctx, deckID)
}

// validateCardContent checks that the JSONB content matches the expected schema
// for the given card type.
func validateCardContent(cardType domain.CardType, rawContent []byte) error {
	var generic map[string]interface{}
	if err := json.Unmarshal(rawContent, &generic); err != nil {
		return apierror.BadRequest("card content must be valid JSON")
	}

	switch cardType {
	case domain.CardTypeFlashcard:
		if _, ok := generic["front"]; !ok {
			return apierror.BadRequest("flashcard content must have 'front' field")
		}
		if _, ok := generic["back"]; !ok {
			return apierror.BadRequest("flashcard content must have 'back' field")
		}
	case domain.CardTypeCloze:
		if _, ok := generic["text"]; !ok {
			return apierror.BadRequest("cloze content must have 'text' field")
		}
	case domain.CardTypeFreeResponse:
		if _, ok := generic["prompt"]; !ok {
			return apierror.BadRequest("free_response content must have 'prompt' field")
		}
	default:
		return apierror.BadRequest("unknown card type: " + string(cardType))
	}
	return nil
}
