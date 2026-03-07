package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/domain"
	"gorm.io/gorm"
)

// StudyRepository defines DB operations for study sessions and logs.
type StudyRepository interface {
	CreateSession(ctx context.Context, s *domain.StudySession) error
	FindSession(ctx context.Context, id, userID uuid.UUID) (*domain.StudySession, error)
	EndSession(ctx context.Context, id uuid.UUID, endedAt time.Time) error
	CreateLog(ctx context.Context, log *domain.StudyLog) error
	// Stats queries
	DailyActivity(ctx context.Context, userID uuid.UUID, days int) ([]domain.DailyActivity, error)
	TotalStudied(ctx context.Context, userID uuid.UUID) (int64, error)
	RetentionRate(ctx context.Context, userID uuid.UUID) (float64, error)
	Forecast(ctx context.Context, userID uuid.UUID, days int) ([]domain.ForecastDay, error)
}

type studyRepo struct{ db *gorm.DB }

func NewStudyRepository(db *gorm.DB) StudyRepository { return &studyRepo{db: db} }

func (r *studyRepo) CreateSession(ctx context.Context, s *domain.StudySession) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *studyRepo) FindSession(ctx context.Context, id, userID uuid.UUID) (*domain.StudySession, error) {
	var s domain.StudySession
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}

func (r *studyRepo) EndSession(ctx context.Context, id uuid.UUID, endedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.StudySession{}).Where("id = ?", id).
		Update("ended_at", endedAt).Error
}

func (r *studyRepo) CreateLog(ctx context.Context, log *domain.StudyLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *studyRepo) DailyActivity(ctx context.Context, userID uuid.UUID, days int) ([]domain.DailyActivity, error) {
	since := time.Now().UTC().AddDate(0, 0, -days)
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT DATE(sl.logged_at) AS date, COUNT(*) AS count
		FROM study_logs sl
		JOIN study_sessions ss ON ss.id = sl.session_id
		WHERE ss.user_id = ? AND sl.logged_at >= ?
		GROUP BY DATE(sl.logged_at)
		ORDER BY date ASC`, userID, since).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DailyActivity
	for rows.Next() {
		var da domain.DailyActivity
		if err := rows.Scan(&da.Date, &da.Count); err != nil {
			return nil, err
		}
		result = append(result, da)
	}
	return result, nil
}

func (r *studyRepo) TotalStudied(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM study_logs sl
		JOIN study_sessions ss ON ss.id = sl.session_id
		WHERE ss.user_id = ?`, userID).Scan(&count).Error
	return count, err
}

func (r *studyRepo) RetentionRate(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total, correct int64
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM study_logs sl
		JOIN study_sessions ss ON ss.id = sl.session_id
		WHERE ss.user_id = ?`, userID).Scan(&total)
	r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM study_logs sl
		JOIN study_sessions ss ON ss.id = sl.session_id
		WHERE ss.user_id = ? AND sl.rating >= 3`, userID).Scan(&correct)
	if total == 0 {
		return 0, nil
	}
	return float64(correct) / float64(total) * 100, nil
}

func (r *studyRepo) Forecast(ctx context.Context, userID uuid.UUID, days int) ([]domain.ForecastDay, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT DATE(c.due_at) AS date, COUNT(*) AS count
		FROM cards c
		JOIN decks d ON d.id = c.deck_id
		WHERE d.user_id = ? AND c.due_at BETWEEN NOW() AND NOW() + INTERVAL '? days'
		GROUP BY DATE(c.due_at)
		ORDER BY date ASC`, userID, days).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ForecastDay
	for rows.Next() {
		var fd domain.ForecastDay
		if err := rows.Scan(&fd.Date, &fd.Count); err != nil {
			return nil, err
		}
		result = append(result, fd)
	}
	return result, nil
}
