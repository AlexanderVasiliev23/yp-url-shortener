package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage missing godoc.
type Storage struct {
	dbConn *pgxpool.Pool
}

// New missing godoc.
func New(ctx context.Context, dbConn *pgxpool.Pool) (*Storage, error) {
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
			original varchar(255) not null,
			user_id  bigint not null,
			deleted_at timestamp null
		)
	`

	createTokenIndexQuery := `
		create index if not exists short_links_token_index
			on short_links (token);
	`

	createOriginalURLUniqueIndexQuery := `
		create unique index if not exists short_links_original_unique_index 
			on short_links (original);
	`

	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("createTableQuery: %w", err)
	}
	if _, err := tx.Exec(ctx, createTokenIndexQuery); err != nil {
		return fmt.Errorf("createTokenIndexQuery: %w", err)
	}
	if _, err := tx.Exec(ctx, createOriginalURLUniqueIndexQuery); err != nil {
		return fmt.Errorf("createOriginalURLUniqueIndexQuery: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Add missing godoc.
func (s *Storage) Add(ctx context.Context, shortLink *models.ShortLink) error {
	if err := s.save(ctx, shortLink); err != nil {
		return fmt.Errorf("save short link: %w", err)
	}

	return nil
}

// Get missing godoc.
func (s *Storage) Get(ctx context.Context, token string) (*models.ShortLink, error) {
	q := `select id, token, original, user_id, deleted_at from short_links where token = $1;`

	link := new(models.ShortLink)

	if err := s.dbConn.QueryRow(ctx, q, token).Scan(
		&link.ID,
		&link.Token,
		&link.Original,
		&link.UserID,
		&link.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("query link by token: %w", err)
	}

	return link, nil
}

// GetTokenByURL missing godoc.
func (s *Storage) GetTokenByURL(ctx context.Context, url string) (string, error) {
	var token string

	if err := s.dbConn.QueryRow(ctx, "select token from short_links where original = $1", url).Scan(&token); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrNotFound
		}

		return "", fmt.Errorf("GetTokenByURL: %w", err)
	}

	return token, nil
}

// SaveBatch missing godoc.
func (s *Storage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	if len(shortLinks) == 0 {
		return nil
	}

	var entries [][]any
	columns := []string{"id", "token", "original", "user_id"}
	tableName := "short_links"

	for _, shortLink := range shortLinks {
		entries = append(entries, []any{shortLink.ID, shortLink.Token, shortLink.Original, shortLink.UserID})
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

// FindByUserID missing godoc.
func (s *Storage) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	q := `select id, token, original, user_id from short_links where user_id = $1 and deleted_at is null`

	rows, err := s.dbConn.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("select short links by user: %w", err)
	}
	defer rows.Close()

	var result []*models.ShortLink
	for rows.Next() {
		var model models.ShortLink
		if err := rows.Scan(
			&model.ID,
			&model.Token,
			&model.Original,
			&model.UserID,
		); err != nil {
			return nil, fmt.Errorf("scan row to struct: %w", err)
		}

		result = append(result, &model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find by user id: rows err: %w", err)
	}

	return result, nil
}

// DeleteByTokens missing godoc.
func (s *Storage) DeleteByTokens(ctx context.Context, tokens []string) error {
	q := `update short_links set deleted_at = now() where token = ANY ($1)`

	if _, err := s.dbConn.Exec(ctx, q, tokens); err != nil {
		return fmt.Errorf("exec delete tokens query: %w", err)
	}

	return nil
}

// FilterOnlyThisUserTokens missing godoc.
func (s *Storage) FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error) {
	q := `
		select token
		from short_links
		where user_id = $1 and token = any($2);
	`

	rows, err := s.dbConn.Query(ctx, q, userID, tokens)
	if err != nil {
		return nil, fmt.Errorf("select tokens by user and tokens: %w", err)
	}

	res := make([]string, 0, len(tokens))
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}
		res = append(res, token)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("filter only this user tokens: rows err: %w", err)
	}

	return res, nil
}

func (s *Storage) save(ctx context.Context, shortLink *models.ShortLink) error {
	q := `insert into short_links (id, token, original, user_id) values ($1,$2,$3,$4)`

	if _, err := s.dbConn.Exec(ctx, q, shortLink.ID, shortLink.Token, shortLink.Original, shortLink.UserID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrAlreadyExists
		}
		return fmt.Errorf("exec insert query: %w", err)
	}

	return nil
}
