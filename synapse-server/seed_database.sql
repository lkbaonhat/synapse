-- Seed data for synapse_db

-- Clean database to avoid duplicate records on multiple runs
TRUNCATE TABLE users, folders, tags, decks, deck_tags, cards, card_tags, study_sessions, study_logs CASCADE;

-- 1. Create a dummy User
INSERT INTO users (email, password_hash)
VALUES ('demo@example.com', '$2a$10$w3U/d19W677YJ8bXbNVK1uk.1.bXXbT5Zl.h9fBpx/M8K/R6f4WzK');

-- 2. Create Folders
INSERT INTO folders (user_id, name)
SELECT id, 'Programming' FROM users WHERE email = 'demo@example.com' UNION ALL
SELECT id, 'Languages' FROM users WHERE email = 'demo@example.com';

-- 3. Create Tags
INSERT INTO tags (user_id, name)
SELECT id, 'Go' FROM users WHERE email = 'demo@example.com' UNION ALL
SELECT id, 'Backend' FROM users WHERE email = 'demo@example.com' UNION ALL
SELECT id, 'Japanese' FROM users WHERE email = 'demo@example.com';

-- 4. Create Decks
INSERT INTO decks (user_id, folder_id, name, description)
SELECT u.id, f.id, 'Go Fundamentals', 'Basic concepts of Golang'
FROM users u
JOIN folders f ON f.user_id = u.id AND f.name = 'Programming'
WHERE u.email = 'demo@example.com';

INSERT INTO decks (user_id, folder_id, name, description)
SELECT u.id, f.id, 'JLPT N5 Vocab', 'Beginner Japanese Vocabulary'
FROM users u
JOIN folders f ON f.user_id = u.id AND f.name = 'Languages'
WHERE u.email = 'demo@example.com';

-- Link logic for Decks and Tags
INSERT INTO deck_tags (deck_id, tag_id)
SELECT d.id, t.id FROM decks d, tags t WHERE d.name = 'Go Fundamentals' AND t.name IN ('Go', 'Backend');

INSERT INTO deck_tags (deck_id, tag_id)
SELECT d.id, t.id FROM decks d, tags t WHERE d.name = 'JLPT N5 Vocab' AND t.name = 'Japanese';

-- 5. Create Cards
-- Basic cards for "Go Fundamentals"
INSERT INTO cards (deck_id, type, content, interval, easiness, repetitions, due_at)
SELECT id, 'basic', '{"front": "What does goroutine do?", "back": "It is a lightweight thread managed by the Go runtime."}', 0, 2.5, 0, NOW()
FROM decks WHERE name = 'Go Fundamentals' UNION ALL
SELECT id, 'basic', '{"front": "How do you declare a channel in Go?", "back": "ch := make(chan int)"}', 0, 2.5, 0, NOW()
FROM decks WHERE name = 'Go Fundamentals' UNION ALL
SELECT id, 'basic', '{"front": "What keyword is used to defer execution?", "back": "defer"}', 1, 2.6, 1, NOW() + INTERVAL '1 day'
FROM decks WHERE name = 'Go Fundamentals';

-- Basic cards for "JLPT N5 Vocab"
INSERT INTO cards (deck_id, type, content, interval, easiness, repetitions, due_at)
SELECT id, 'basic', '{"front": "猫 (neko)", "back": "Cat"}', 0, 2.5, 0, NOW()
FROM decks WHERE name = 'JLPT N5 Vocab' UNION ALL
SELECT id, 'basic', '{"front": "犬 (inu)", "back": "Dog"}', 0, 2.5, 0, NOW()
FROM decks WHERE name = 'JLPT N5 Vocab';

-- Link logic for Cards and Tags
INSERT INTO card_tags (card_id, tag_id)
SELECT c.id, t.id FROM cards c
JOIN decks d ON c.deck_id = d.id AND d.name = 'Go Fundamentals'
JOIN tags t ON t.name = 'Go'
WHERE c.content->>'front' = 'What does goroutine do?';

INSERT INTO card_tags (card_id, tag_id)
SELECT c.id, t.id FROM cards c
JOIN decks d ON c.deck_id = d.id AND d.name = 'JLPT N5 Vocab'
JOIN tags t ON t.name = 'Japanese'
WHERE c.content->>'front' = '猫 (neko)';
