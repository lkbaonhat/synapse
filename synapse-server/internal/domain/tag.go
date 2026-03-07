package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag is a label that can be attached to decks and cards.
type Tag struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
