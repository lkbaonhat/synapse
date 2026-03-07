package srs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/synapse/server/pkg/srs"
)

func scheduleAt(interval int, easiness float64, repetitions int) srs.CardSchedule {
	return srs.CardSchedule{
		Interval:    interval,
		Easiness:    easiness,
		Repetitions: repetitions,
		DueAt:       time.Now(),
	}
}

// --- First-repetition edge cases ---

func TestFirstRepetition_Again(t *testing.T) {
	s := srs.DefaultSchedule()
	next := srs.Compute(s, srs.Again)
	assert.Equal(t, 0, next.Repetitions)
	assert.Equal(t, 1, next.Interval)
	assert.Equal(t, 2.3, next.Easiness) // 2.5 - 0.2
}

func TestFirstRepetition_Hard(t *testing.T) {
	s := srs.DefaultSchedule()
	next := srs.Compute(s, srs.Hard)
	assert.Equal(t, 1, next.Repetitions)
	assert.Equal(t, 1, next.Interval)
	assert.InDelta(t, 2.35, next.Easiness, 0.001) // 2.5 - 0.15
}

func TestFirstRepetition_Good(t *testing.T) {
	s := srs.DefaultSchedule()
	next := srs.Compute(s, srs.Good)
	assert.Equal(t, 1, next.Repetitions)
	assert.Equal(t, 1, next.Interval)
	assert.Equal(t, 2.5, next.Easiness) // unchanged
}

func TestFirstRepetition_Easy(t *testing.T) {
	s := srs.DefaultSchedule()
	next := srs.Compute(s, srs.Easy)
	assert.Equal(t, 1, next.Repetitions)
	assert.Equal(t, 4, next.Interval)
	assert.InDelta(t, 2.5, next.Easiness, 0.001) // clamped at max
}

// --- Multi-repetition branches ---

func TestGoodProgression(t *testing.T) {
	s := srs.DefaultSchedule()
	s = srs.Compute(s, srs.Good) // rep=1, interval=1
	s = srs.Compute(s, srs.Good) // rep=2, interval=6
	s = srs.Compute(s, srs.Good) // rep=3, interval=6*2.5=15
	assert.Equal(t, 3, s.Repetitions)
	assert.Equal(t, 15, s.Interval)
}

func TestEasyBoostsEasiness(t *testing.T) {
	s := scheduleAt(6, 2.5, 2)
	next := srs.Compute(s, srs.Easy)
	assert.Equal(t, 2.5, next.Easiness) // already at max, clamped
}

func TestHardReducesEasiness(t *testing.T) {
	s := scheduleAt(6, 1.4, 2)
	next := srs.Compute(s, srs.Hard)
	assert.InDelta(t, 1.3, next.Easiness, 0.001) // clamped at min
}

// --- E-factor clamping ---

func TestEasinessFloorClamping(t *testing.T) {
	s := scheduleAt(1, 1.31, 0)
	next := srs.Compute(s, srs.Again)
	assert.GreaterOrEqual(t, next.Easiness, 1.3)
}

func TestEasinessCeilingClamping(t *testing.T) {
	s := scheduleAt(10, 2.48, 3)
	next := srs.Compute(s, srs.Easy)
	assert.LessOrEqual(t, next.Easiness, 2.5)
}

// --- DueAt is in the future ---

func TestDueAtIsInFuture(t *testing.T) {
	s := srs.DefaultSchedule()
	next := srs.Compute(s, srs.Good)
	assert.True(t, next.DueAt.After(time.Now().Add(-time.Second)))
}
