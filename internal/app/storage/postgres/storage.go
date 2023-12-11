package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	dbConn *pgx.Conn
}

func New(ctx context.Context, dbConn *pgx.Conn) (*Storage, error) {
	s := &Storage{
		dbConn: dbConn,
	}

	if err := s.createSchema(ctx); err != nil {
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	return s, nil
}

func (s *Storage) createSchema(ctx context.Context) error {
	createTableQuery := `
		create table if not exists short_links
		(
			id       uuid primary key,
			token    varchar(255) not null,
			original varchar(255) not null
		)
	`

	createIndexQuery := `
		create index if not exists short_links_token_index
			on short_links (token);
	`

	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("exec schema creation query: %w", err)
	}
	if _, err := tx.Exec(ctx, createIndexQuery); err != nil {
		return fmt.Errorf("exec schema creation query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) Add(ctx context.Context, token, url string) error {
	shortLink := models.NewShortLink(token, url)

	if err := s.SaveBatch(ctx, []*models.ShortLink{shortLink}); err != nil {
		return fmt.Errorf("save short link: %w", err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, token string) (string, error) {
	q := `select original from short_links where token = $1;`

	var link string

	if err := s.dbConn.QueryRow(ctx, q, token).Scan(&link); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrNotFound
		}

		return "", fmt.Errorf("query link by token: %w", err)
	}

	return link, nil
}

func (s *Storage) save(ctx context.Context, shortLink *models.ShortLink) error {
	q := `insert into short_links (id, token, original) values ($1,$2,$3)`

	if _, err := s.dbConn.Exec(ctx, q, shortLink.ID, shortLink.Token, shortLink.Original); err != nil {
		return fmt.Errorf("exec insert query: %w", err)
	}

	return nil
}

func (s *Storage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	if len(shortLinks) == 0 {
		return nil
	}

	var entries [][]any
	columns := []string{"id", "token", "original"}
	tableName := "short_links"

	for _, shortLink := range shortLinks {
		entries = append(entries, []any{shortLink.ID, shortLink.Token, shortLink.Original})
	}

	_, err := s.dbConn.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(entries),
	)

	if err != nil {
		return fmt.Errorf("copying into %s table: %w", tableName, err)
	}

	return nil
}
