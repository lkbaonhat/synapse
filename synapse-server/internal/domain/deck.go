package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Deck holds a collection of flashcards.
type Deck struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	FolderID    *uuid.UUID     `gorm:"type:uuid"`
	Name        string         `gorm:"not null"`
	Description string
	Tags        []Tag          `gorm:"many2many:deck_tags"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
