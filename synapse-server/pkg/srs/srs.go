package srs

import (
	"math"
	"time"
)

// DifficultyRating maps to Anki-style user ratings.
type DifficultyRating int

const (
	Again DifficultyRating = 1 // complete blackout — reset
	Hard  DifficultyRating = 2 // correct but difficult
	Good  DifficultyRating = 3 // correct with some hesitation
	Easy  DifficultyRating = 4 // perfect recall
)

const (
	minEasiness     = 1.3
	maxEasiness     = 2.5
	defaultEasiness = 2.5
)

// CardSchedule holds the current SRS state for a card.
type CardSchedule struct {
	Interval    int     // days until next review
	Easiness    float64 // SM-2 E-factor
	Repetitions int     // total successful repetitions
	DueAt       time.Time
}

// DefaultSchedule returns a fresh schedule for a card that has never been studied.
func DefaultSchedule() CardSchedule {
	return CardSchedule{
		Interval:    0,
		Easiness:    defaultEasiness,
		Repetitions: 0,
	}
}

// Compute returns the next CardSchedule after applying the SM-2 algorithm for
// the given rating. It is a pure function with no DB access.
func Compute(current CardSchedule, rating DifficultyRating) CardSchedule {
	next := current

	switch rating {
	case Again:
		// Full reset — card goes back to the beginning.
		next.Repetitions = 0
		next.Interval = 1
		next.Easiness = clamp(current.Easiness-0.2, minEasiness, maxEasiness)

	case Hard:
		// Penalty — slow the interval growth, reduce E-factor.
		if current.Repetitions == 0 {
			next.Interval = 1
		} else {
			next.Interval = int(math.Round(float64(current.Interval) * 1.2))
		}
		next.Easiness = clamp(current.Easiness-0.15, minEasiness, maxEasiness)
		next.Repetitions = current.Repetitions + 1

	case Good:
		// Standard SM-2 schedule.
		next.Repetitions = current.Repetitions + 1
		switch current.Repetitions {
		case 0:
			next.Interval = 1
		case 1:
			next.Interval = 6
		default:
			next.Interval = int(math.Round(float64(current.Interval) * current.Easiness))
		}
		// E-factor unchanged for Good.

	case Easy:
		// Accelerated scheduling — big interval boost.
		next.Repetitions = current.Repetitions + 1
		switch current.Repetitions {
		case 0:
			next.Interval = 4
		case 1:
			next.Interval = 10
		default:
			next.Interval = int(math.Round(float64(current.Interval) * current.Easiness * 1.3))
		}
		next.Easiness = clamp(current.Easiness+0.15, minEasiness, maxEasiness)
	}

	next.DueAt = time.Now().UTC().AddDate(0, 0, next.Interval)
	return next
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
