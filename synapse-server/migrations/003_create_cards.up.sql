-- Migration: 003 — Create cards table
-- Up

CREATE TABLE IF NOT EXISTS cards (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deck_id      UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
    type         TEXT NOT NULL,
    content      JSONB NOT NULL DEFAULT '{}',
    interval     INTEGER NOT NULL DEFAULT 0,
    easiness     DOUBLE PRECISION NOT NULL DEFAULT 2.5,
    repetitions  INTEGER NOT NULL DEFAULT 0,
    due_at       TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS card_tags (
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (card_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_cards_deck_due ON cards(deck_id, due_at);
