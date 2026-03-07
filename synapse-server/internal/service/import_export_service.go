package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/synapse/server/internal/apierror"
	"github.com/synapse/server/internal/domain"
	"github.com/synapse/server/internal/repository"
)

// ImportResult reports the outcome of a bulk import.
type ImportResult struct {
	Imported int      `json:"imported"`
	Errors   []string `json:"errors"`
}

// ImportExportService handles CSV import and CSV/JSON export.
type ImportExportService interface {
	ImportCSV(ctx context.Context, deckID uuid.UUID, file multipart.File) (*ImportResult, error)
	ExportDeckCSV(ctx context.Context, deckID, userID uuid.UUID, w io.Writer) error
	ExportUserJSON(ctx context.Context, userID uuid.UUID, w io.Writer) error
}

type importExportService struct {
	cardRepo repository.CardRepository
	deckRepo repository.DeckRepository
	tagRepo  repository.TagRepository
}

func NewImportExportService(cardRepo repository.CardRepository, deckRepo repository.DeckRepository, tagRepo repository.TagRepository) ImportExportService {
	return &importExportService{cardRepo: cardRepo, deckRepo: deckRepo, tagRepo: tagRepo}
}

func (s *importExportService) ImportCSV(ctx context.Context, deckID uuid.UUID, file multipart.File) (*ImportResult, error) {
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, apierror.BadRequest("cannot read CSV headers")
	}

	colIndex := make(map[string]int)
	for i, h := range headers {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	result := &ImportResult{}
	row := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: read error: %v", row, err))
			row++
			continue
		}

		cardType := domain.CardType(strings.TrimSpace(safeGet(record, colIndex, "type")))
		front := strings.TrimSpace(safeGet(record, colIndex, "front"))
		back := strings.TrimSpace(safeGet(record, colIndex, "back"))

		if cardType == "" || front == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: missing required fields (type, front)", row))
			row++
			continue
		}

		content, err := buildContent(cardType, front, back)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", row, err))
			row++
			continue
		}

		card := &domain.Card{
			DeckID:   deckID,
			Type:     cardType,
			Content:  content,
			Easiness: 2.5,
		}
		if err := s.cardRepo.Create(ctx, card); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: db error: %v", row, err))
		} else {
			result.Imported++
		}
		row++
	}
	return result, nil
}

func (s *importExportService) ExportDeckCSV(ctx context.Context, deckID, userID uuid.UUID, w io.Writer) error {
	_, err := s.deckRepo.FindByID(ctx, deckID, userID)
	if err != nil {
		return err
	}
	cards, _, err := s.cardRepo.ListByDeck(ctx, deckID, 0, 10000)
	if err != nil {
		return err
	}

	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"type", "front", "back", "interval", "easiness", "repetitions", "due_at"})
	for _, c := range cards {
		var payload map[string]interface{}
		_ = json.Unmarshal(c.Content, &payload)
		front, _ := payload["front"].(string)
		back, _ := payload["back"].(string)
		dueAt := ""
		if c.DueAt != nil {
			dueAt = c.DueAt.Format(time.RFC3339)
		}
		_ = cw.Write([]string{
			string(c.Type), front, back,
			fmt.Sprint(c.Interval), fmt.Sprint(c.Easiness), fmt.Sprint(c.Repetitions), dueAt,
		})
	}
	cw.Flush()
	return cw.Error()
}

func (s *importExportService) ExportUserJSON(ctx context.Context, userID uuid.UUID, w io.Writer) error {
	decks, _, err := s.deckRepo.List(ctx, userID, nil, nil)
	if err != nil {
		return err
	}
	type export struct {
		Decks     []domain.Deck            `json:"decks"`
		Cards     map[string][]domain.Card `json:"cards"`
		ExportedAt time.Time               `json:"exportedAt"`
	}
	ex := export{Decks: decks, Cards: make(map[string][]domain.Card), ExportedAt: time.Now().UTC()}
	for _, d := range decks {
		cards, _, _ := s.cardRepo.ListByDeck(ctx, d.ID, 0, 100000)
		ex.Cards[d.ID.String()] = cards
	}
	return json.NewEncoder(w).Encode(ex)
}

func safeGet(record []string, index map[string]int, key string) string {
	i, ok := index[key]
	if !ok || i >= len(record) {
		return ""
	}
	return record[i]
}

func buildContent(cardType domain.CardType, front, back string) ([]byte, error) {
	switch cardType {
	case domain.CardTypeFlashcard:
		return json.Marshal(map[string]string{"front": front, "back": back})
	case domain.CardTypeCloze:
		return json.Marshal(map[string]string{"text": front})
	case domain.CardTypeFreeResponse:
		return json.Marshal(map[string]string{"prompt": front, "answer": back})
	default:
		return nil, apierror.BadRequest("unknown card type: " + string(cardType))
	}
}
