package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Folder organises decks in a tree structure.
type Folder struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name      string         `gorm:"not null"`
	ParentID  *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
