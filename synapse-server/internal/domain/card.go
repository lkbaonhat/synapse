package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CardType identifies the format of a card's content payload.
type CardType string

const (
	CardTypeFlashcard      CardType = "flashcard"
	CardTypeCloze          CardType = "cloze"
	CardTypeFreeResponse   CardType = "free_response"
	CardTypeMultipleChoice CardType = "multiple_choice"
)

// Card is a single learnable item belonging to a Deck.
// Content is a JSONB field whose schema depends on CardType:
//
//	flashcard:     { "front": "...", "back": "..." }
//	cloze:         { "text": "...", "clozeFields": [...] }
//	free_response: { "prompt": "...", "answer": "..." }
type Card struct {
	ID      uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeckID  uuid.UUID      `gorm:"type:uuid;not null;index"`
	Type    CardType       `gorm:"not null"`
	Content datatypes.JSON `gorm:"type:jsonb;not null"`
	Tags    []Tag          `gorm:"many2many:card_tags"`

	// SRS scheduling fields (initialised on first study)
	Interval    int        // days until next review
	Easiness    float64    // SM-2 E-factor (default 2.5)
	Repetitions int        // total successful reviews
	DueAt       *time.Time // nil → card not yet studied (new)

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
