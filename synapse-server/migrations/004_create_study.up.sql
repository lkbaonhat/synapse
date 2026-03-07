-- Migration: 004 — Create study sessions and logs
-- Up

CREATE TABLE IF NOT EXISTS study_sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    deck_id    UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
    mode       TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS study_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES study_sessions(id) ON DELETE CASCADE,
    card_id    UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    rating     INTEGER NOT NULL,
    time_taken INTEGER NOT NULL,
    logged_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_study_logs_session ON study_logs(session_id);
CREATE INDEX IF NOT EXISTS idx_study_logs_logged_at ON study_logs(logged_at);
