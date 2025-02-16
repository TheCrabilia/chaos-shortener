package shortener

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TheCrabilia/chaos-shortener/internal/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Shortener struct {
	db *db.Queries
}

func generateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:10]
}

func NewShortener(db *db.Queries) *Shortener {
	return &Shortener{
		db: db,
	}
}

// Shorten creates unique identifier for the given URL and saves it in database.
// Returns full redirect URL.
func (s *Shortener) Shorten(ctx context.Context, baseURL, redirectURL string) (string, error) {
	id := generateID()

	if err := s.db.CreateURL(ctx, db.CreateURLParams{ID: id, Url: redirectURL}); err != nil {
		return "", fmt.Errorf("failed to create url in db: %w", err)
	}

	return fmt.Sprintf("%s/r/%s", baseURL, id), nil
}

// RedirectURL gets original URL from database by shortened identifier.
func (s *Shortener) RedirectURL(ctx context.Context, id string) (string, error) {
	url, err := s.db.GetURLById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("failed to get redirect url from db: url with id %s not found", id)
	} else if err != nil {
		return "", fmt.Errorf("failed to get redirect url from db: %w", err)
	}

	return url, nil
}
