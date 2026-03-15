package domain

// Stats response DTOs — computed from study_logs, no dedicated table.

// StatsOverview is the payload for GET /api/v1/stats/overview.
// @Description Overall statistics for a user's library and study patterns
type StatsOverview struct {
	TotalCards    int64   `json:"totalCards"`
	RetentionRate float64 `json:"retentionRate"` // percentage
	CurrentStreak int     `json:"currentStreak"` // calendar days
	TotalStudied  int64   `json:"totalStudied"`
}

// DailyActivity is one data-point for the heatmap chart.
// @Description A data-point representing learning activity for a specific day
type DailyActivity struct {
	Date  string `json:"date"`  // "YYYY-MM-DD"
	Count int    `json:"count"` // cards reviewed that day
}

// ForecastDay represents the projected review workload for a future date.
// @Description Projected review workload for a future date
type ForecastDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// DeckStats is the per-deck breakdown of card mastery tiers.
// @Description Breakdown of card mastery tiers for a specific deck
type DeckStats struct {
	DeckID   string `json:"deckId"`
	New      int64  `json:"new"`
	Learning int64  `json:"learning"`
	Review   int64  `json:"review"`
	Mastered int64  `json:"mastered"`
}
