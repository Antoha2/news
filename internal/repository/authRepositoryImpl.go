package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Antoha2/news/internal/lib/models"
	"github.com/pkg/errors"
)

// SaveUser saves user to db.
func (r *RepImpl) UserSaver(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.postgres.SaveUser"

	//checking user uniqueness
	var count int
	stmt := r.DB.QueryRowContext(ctx, "SELECT COUNT(id) FROM users WHERE email = $1", email)

	if err := stmt.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
	}

	//adding a user to the database
	var id int64
	stmt = r.DB.QueryRowContext(ctx, "INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id", email, passHash)
	if err := stmt.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (r *RepImpl) UserProvider(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.User"

	user := new(models.User)
	stmt := r.DB.QueryRowContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1", email)

	err := stmt.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return *user, nil
}
