package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	dbpool, err := pgxpool.New(ctx, storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	return &Storage{db: dbpool}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id`,
		email,
		passHash,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User
	err := s.db.QueryRow(
		ctx,
		`SELECT id, email, pass_hash FROM users WHERE email=$1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"
	var app models.App

	err := s.db.QueryRow(
		ctx,
		`SELECT id, name, secret FROM apps WHERE id = $1`,
		id,
	).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"
	var isAdmin bool

	err := s.db.QueryRow(ctx, `SELECT is_admin FROM users WHERE id = $1`, userID).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}
