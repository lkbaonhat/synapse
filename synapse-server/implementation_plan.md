# Synapse Backend — Go Implementation Plan

The Synapse server is a **greenfield Go REST API** that powers the Synapse study-assistant client (Vue 3 / TanStack Query). It implements Spaced Repetition (SRS/SM-2), multiple card formats, deck/folder/tag management, statistics, and import/export utilities.

---

## Architecture Decisions

| Concern | Choice | Rationale |
|---|---|---|
| Language | **Go 1.22+** | Performance, strong typing, excellent concurrency |
| Web framework | **Gin** | Minimal, battle-tested, great middleware ecosystem |
| ORM / DB layer | **GORM + PostgreSQL** | Postgres for JSONB support (card content), mature ORM |
| Auth | **JWT (access + refresh tokens)** | Stateless, matches Vue axios interceptor pattern |
| Validation | **go-playground/validator** | Struct-tag driven, integrates with Gin |
| Config | **godotenv + viper** | 12-factor app env management |
| Migrations | **golang-migrate** | SQL-first raw migrations, version controlled |
| File storage | **local disk → configurable S3** | Multimedia attachments (images, audio) |
| Testing | **testify + httptest** | Unit + integration tests per package |
| Containerization | **Docker + docker-compose** | Postgres + server orchestration |

---

## Directory Structure

```
synapse-server/
├── cmd/
│   └── server/
│       └── main.go               # Entry point
├── internal/
│   ├── config/                   # Viper config loader
│   ├── database/                 # DB connection, migrations runner
│   ├── middleware/               # Auth, CORS, logger, error handler
│   ├── domain/                   # Pure domain types & interfaces
│   │   ├── user.go
│   │   ├── folder.go
│   │   ├── deck.go
│   │   ├── card.go
│   │   ├── tag.go
│   │   ├── study_session.go
│   │   └── stats.go
│   ├── repository/               # DB access layer (GORM)
│   │   ├── user_repo.go
│   │   ├── folder_repo.go
│   │   ├── deck_repo.go
│   │   ├── card_repo.go
│   │   ├── tag_repo.go
│   │   └── study_repo.go
│   ├── service/                  # Business logic
│   │   ├── auth_service.go
│   │   ├── deck_service.go
│   │   ├── card_service.go
│   │   ├── study_service.go      # SRS engine lives here
│   │   ├── stats_service.go
│   │   └── import_export_service.go
│   ├── handler/                  # Gin HTTP handlers (thin)
│   │   ├── auth_handler.go
│   │   ├── deck_handler.go
│   │   ├── card_handler.go
│   │   ├── study_handler.go
│   │   ├── stats_handler.go
│   │   └── import_export_handler.go
│   └── router/
│       └── router.go             # Route registration
├── migrations/                   # SQL migration files
├── pkg/
│   └── srs/                      # SM-2 algorithm (pure functions, zero deps)
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## Implementation Phases

### Test-Driven Development (TDD) Requirement

All feature development globally must follow a strict Test-Driven Development (TDD) cycle:
1. **Red**: Write a failing test for the next piece of functionality (e.g., a unit test for a pure function in `pkg/srs` or an `httptest` integration test for a new Gin handler).
2. **Green**: Write the minimum amount of Go code necessary to make the test pass.
3. **Refactor**: Clean up the code while ensuring the tests continue to pass, adhering to Go idioms and project architecture.

Tests are not an afterthought; they are the primary driver of the implementation. For HTTP handlers, the route and expected JSON response should be tested via `httptest` before the underlying service logic is fully fleshed out.

---

### Phase 1 — Project Scaffold & Infrastructure
**Goal:** A runnable Go server with DB connection, config, and health endpoint.

#### Tasks
- Initialize Go module (`go mod init`)
- Add dependencies: Gin, GORM, postgres driver, godotenv, viper, golang-migrate, jwt-go, validator, testify
- `cmd/server/main.go` — wires up config → DB → router → server
- `internal/config/` — load `.env` / viper
- `internal/database/` — GORM open + `AutoMigrate` call + ping
- `internal/router/router.go` — register `GET /health`
- `Dockerfile` + `docker-compose.yml` (Postgres + server)
- `migrations/` — first migration: `users` table

#### Deliverables
- `GET /api/health` → `{ "status": "ok" }`
- Docker Compose brings up Postgres + Go server

---

### Phase 2 — Authentication (JWT)
**Goal:** Secure register / login / token refresh, matching the client's axios interceptor contract.

#### Domain Models
```go
type User struct {
  ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
  Email        string    `gorm:"uniqueIndex;not null"`
  PasswordHash string    `gorm:"not null"`
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
```

#### API Endpoints
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/auth/register` | Create user, return tokens |
| `POST` | `/api/auth/login` | Verify password, return tokens |
| `POST` | `/api/auth/refresh` | Swap refresh token for new access token |
| `POST` | `/api/auth/logout` | Invalidate refresh token |

#### Key Implementation Points
- **Access token**: short-lived JWT (15 min), signed with `HS256`
- **Refresh token**: long-lived (7 days), stored in DB (`refresh_tokens` table)
- **Middleware** `AuthRequired` — extracts & validates bearer token, injects `userID` into Gin context
- Passwords hashed with `bcrypt` (cost = 12)
- Validation: email format + password min length via struct tags

---

### Phase 3 — Library Management (Folders, Decks, Tags)
**Goal:** Full CRUD for the content organisation layer.

#### Domain Models

```go
type Folder struct {
  ID       uuid.UUID
  UserID   uuid.UUID
  Name     string
  ParentID *uuid.UUID // nullable — supports nesting
}

type Deck struct {
  ID          uuid.UUID
  UserID      uuid.UUID
  FolderID    *uuid.UUID
  Name        string
  Description string
  Tags        []Tag `gorm:"many2many:deck_tags"`
}

type Tag struct {
  ID     uuid.UUID
  UserID uuid.UUID
  Name   string
}
```

#### API Endpoints
| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/api/folders` | List / create folders |
| `GET/PUT/DELETE` | `/api/folders/:id` | Get / update / delete folder |
| `GET/POST` | `/api/decks` | List / create decks (filter by folder, tag) |
| `GET/PUT/DELETE` | `/api/decks/:id` | Get / update / delete deck |
| `GET/POST` | `/api/tags` | List / create tags |
| `DELETE` | `/api/tags/:id` | Delete tag |
| `POST` | `/api/decks/:id/tags` | Attach tags to deck |

#### Notes
- All list endpoints scoped to `userID` from JWT — no data leakage
- Soft delete (`DeletedAt`) via GORM so stats history is preserved

---

### Phase 4 — Card Management (Multiple Formats)
**Goal:** CRUD for flashcards supporting Flashcard, Cloze, and Free Response formats.

#### Domain Model
```go
type CardType string

const (
  CardTypeFlashcard    CardType = "flashcard"
  CardTypeCloze        CardType = "cloze"
  CardTypeFreeResponse CardType = "free_response"
)

type Card struct {
  ID       uuid.UUID
  DeckID   uuid.UUID
  Type     CardType
  // JSONB field stores format-specific payload:
  // flashcard: { front, back }
  // cloze:     { text, clozeFields[] }
  // free_response: { prompt, answer }
  Content  datatypes.JSON `gorm:"type:jsonb"`
  Tags     []Tag `gorm:"many2many:card_tags"`
  // SRS scheduling fields (initialised on first study)
  Interval    int       // days
  Easiness    float64   // SM-2 E-factor (default 2.5)
  Repetitions int
  DueAt       *time.Time
  CreatedAt   time.Time
  UpdatedAt   time.Time
}
```

#### API Endpoints
| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/api/decks/:id/cards` | List / create cards in a deck |
| `GET/PUT/DELETE` | `/api/cards/:id` | Get / update / delete card |
| `POST` | `/api/cards/:id/media` | Upload image/audio attachment |

#### Notes
- Media uploads saved to `/uploads/{userID}/{cardID}/` (local for MVP, swappable to S3)
- Content validated per `CardType` — a cloze card without cloze markers returns `400`

---

### Phase 5 — Study Engine (SRS Core)
**Goal:** The science heart of the app — SM-2 scheduling, study session tracking.

#### `pkg/srs/` — Pure Algorithm Package
```go
// DifficultyRating maps to Anki-style ratings
type DifficultyRating int

const (
  Again DifficultyRating = 1 // complete blackout
  Hard  DifficultyRating = 2
  Good  DifficultyRating = 3
  Easy  DifficultyRating = 4
)

type CardSchedule struct {
  Interval    int
  Easiness    float64
  Repetitions int
  DueAt       time.Time
}

// Compute returns the next schedule after a given rating.
// Pure function — no DB access, easily unit-tested.
func Compute(current CardSchedule, rating DifficultyRating) CardSchedule
```

**SM-2 rules implemented:**
- `Again` → reset interval to 1 day, repetitions = 0
- `Hard` → interval × 1.2, E-factor -= 0.15
- `Good` → standard SM-2 formula
- `Easy` → interval × E-factor × 1.3, E-factor += 0.15
- E-factor clamped to `[1.3, 2.5]`

#### Study Session Domain
```go
type StudySession struct {
  ID        uuid.UUID
  UserID    uuid.UUID
  DeckID    uuid.UUID
  Mode      StudyMode // "learn" | "review" | "cram"
  StartedAt time.Time
  EndedAt   *time.Time
}

type StudyLog struct {
  ID        uuid.UUID
  SessionID uuid.UUID
  CardID    uuid.UUID
  Rating    DifficultyRating
  TimeTaken int // ms
  LoggedAt  time.Time
}
```

#### API Endpoints
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/study/sessions` | Start a session (returns session ID + first batch of cards) |
| `GET` | `/api/study/sessions/:id/next` | Get next card(s) for session |
| `POST` | `/api/study/sessions/:id/answer` | Submit rating for a card → updates SRS schedule |
| `POST` | `/api/study/sessions/:id/end` | End session, persist summary |
| `GET` | `/api/decks/:id/due-count` | How many cards are due today (for dashboard badge) |

#### Mode Logic
- **Learn**: fetch cards where `DueAt IS NULL` (never studied), batch of 20
- **Review**: fetch cards where `DueAt <= now()`, ordered by due date
- **Cram**: fetch all cards in deck regardless of schedule — does **not** update `DueAt`

---

### Phase 6 — Statistics & Progress
**Goal:** Power the stats dashboard and mastery/streak features.

#### API Endpoints
| Method | Path | Description |
|---|---|---|
| `GET` | `/api/stats/overview` | Total cards, retention rate, current streak |
| `GET` | `/api/stats/activity` | Daily counts for last N days (heatmap data) |
| `GET` | `/api/stats/forecast` | Review workload per day for next 7 days |
| `GET` | `/api/decks/:id/stats` | Per-deck breakdown (new / learning / review / mastered) |

#### Streak Logic
- A "streak day" = at least 1 card reviewed that calendar day (user's local timezone via header)
- Computed from `StudyLog`, grouped by date — no separate table needed for MVP

#### Mastery Tiers (per card)
| Tier | Condition |
|---|---|
| New | `DueAt IS NULL` |
| Learning | `repetitions < 3` |
| Review | `repetitions >= 3 AND interval < 21` |
| Mastered | `interval >= 21` |

---

### Phase 7 — Import / Export
**Goal:** Bulk upload from CSV/Excel; full data export for backups.

#### Import
- `POST /api/decks/:id/import` — multipart form upload
- Accepts `.csv` with columns: `type, front, back, tags`
- Uses `encoding/csv` (stdlib) for CSV; `github.com/xuri/excelize` for `.xlsx`
- Validates each row, returns partial success report (`{ imported: N, errors: [...] }`)

#### Export
- `GET /api/decks/:id/export?format=csv` — streams CSV
- `GET /api/user/export` — full JSON export of all user data (GDPR-friendly)

---

### Phase 8 — Cross-Cutting Concerns & Hardening
**Goal:** Production-ready API quality.

#### Tasks
- **Global error handler middleware** — maps domain errors to HTTP codes; never leaks stack traces
- **Request logging** — structured JSON logs (Zap or slog)
- **Rate limiting** — `golang.org/x/time/rate` per IP (login endpoint especially)
- **CORS** — whitelist client origin from config
- **Pagination** — all list endpoints: `?page=&limit=` with `X-Total-Count` header
- **OpenAPI / Swagger** — `swaggo/swag` annotations on handlers, generate `docs/`
- **Integration tests** — `httptest` + real Postgres (test containers or dedicated test DB)
- **CI** — GitHub Actions: `go vet`, `staticcheck`, `go test ./...`

---

## Database Schema Summary (ERD)

```
users ──< decks ──< cards
  │           │
  │          deck_tags >── tags
  │
  ├──< folders (self-referential)
  ├──< refresh_tokens
  ├──< study_sessions ──< study_logs
  └──< tags
```

---

## API Versioning

All routes prefixed with `/api/v1/`. The router group pattern in Gin:
```go
v1 := r.Group("/api/v1")
v1.Use(middleware.AuthRequired())
{
  v1.GET("/decks", handler.ListDecks)
  // ...
}
```

---

## Verification Plan

### Automated Tests

#### Unit Tests — SRS Algorithm
```bash
cd synapse-server
go test ./pkg/srs/... -v
```
Tests cover all four rating branches, E-factor clamping, and edge cases (first repetition).

#### Unit Tests — Service Layer
```bash
go test ./internal/service/... -v
```
Services tested with mock repositories (using `testify/mock`).

#### Integration Tests — HTTP Handlers
```bash
# Requires Docker Compose to be running
docker compose up -d db
go test ./internal/handler/... -v -tags=integration
```
Uses `net/http/httptest` against a real test database. Tests cover:
- Auth flow (register → login → access protected route → refresh → logout)
- Full deck + card CRUD lifecycle
- Study session: start → answer (all 4 ratings) → verify card schedule updated → end
- Import CSV, verify cards created

#### Build & Static Analysis
```bash
go build ./...
go vet ./...
```

### Manual Verification
1. Run `docker compose up` — confirm `GET http://localhost:8080/api/v1/health` returns `200 { "status": "ok" }`.
2. Use the Bruno / Postman collection (to be committed in `docs/api/`) to walk through auth flow.
3. Start the Vue client (`cd synapse-client && pnpm dev`) and confirm Login / Register screens can authenticate against the running server.
