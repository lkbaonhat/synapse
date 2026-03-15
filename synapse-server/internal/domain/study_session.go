package domain

import (
	"time"

	"github.com/google/uuid"
)

// StudyMode defines the mode of a study session.
type StudyMode string

const (
	StudyModeLearn  StudyMode = "learn"
	StudyModeReview StudyMode = "review"
	StudyModeCram   StudyMode = "cram"
	StudyModeQuiz   StudyMode = "quiz"
)

// StudySession tracks a single study interaction with a deck.
type StudySession struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	DeckID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Mode      StudyMode  `gorm:"not null"`
	StartedAt time.Time  `gorm:"not null"`
	EndedAt   *time.Time
}

// StudyLog records each individual card answer within a session.
type StudyLog struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index"`
	CardID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Rating    int       `gorm:"not null"` // 1=Again 2=Hard 3=Good 4=Easy
	TimeTaken int       `gorm:"not null"` // milliseconds
	LoggedAt  time.Time `gorm:"not null"`
}

// QuizResult summarizes a quiz session.
type QuizResult struct {
	TotalCorrect int                  `json:"totalCorrect"`
	TotalWrong   int                  `json:"totalWrong"`
	WrongAnswers []WrongAnswerSummary `json:"wrongAnswers"`
}

// WrongAnswerSummary represents a single missed question and its correct answer/context.
type WrongAnswerSummary struct {
	CardID      uuid.UUID `json:"cardId"`
	Front       string    `json:"front"`
	CorrectBack string    `json:"correctBack"`
}
