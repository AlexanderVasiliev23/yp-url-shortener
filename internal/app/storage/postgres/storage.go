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
			uuid   uuid primary key,
			token  varchar(255) not null,
			origin varchar(255) not null
		)
	`

	createIndexQuery := `
		create index if not exists short_links_token_index
			on short_links (token);
	`

	tx, err := s.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := s.dbConn.Exec(context.Background(), createTableQuery); err != nil {
		return fmt.Errorf("exec schema creation query: %w", err)
	}
	if _, err := s.dbConn.Exec(context.Background(), createIndexQuery); err != nil {
		return fmt.Errorf("exec schema creation query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) Add(ctx context.Context, token, url string) error {
	shortLink := models.NewShortLink(token, url)

	if err := s.save(ctx, shortLink); err != nil {
		return fmt.Errorf("save short link: %w", err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, token string) (string, error) {
	q := `select origin from short_links where token = $1;`

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
	q := `insert into short_links (uuid, token, origin) values ($1,$2,$3)`

	if _, err := s.dbConn.Exec(ctx, q, shortLink.ID, shortLink.Token, shortLink.Origin); err != nil {
		return fmt.Errorf("exec insert query: %w", err)
	}

	return nil
}
