package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Antoha2/news/internal/auth/internal/domain/models"
)

// SaveUser saves user to db.
func (r *RepAuth) UserSaver(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.postgres.SaveUser"

	//проверка уникальности пользователя
	var count int
	stmt, err := r.DB.Prepare("SELECT COUNT(id) FROM users WHERE email = $1")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = stmt.QueryRowContext(ctx, email).Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	if count > 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
	}

	//добавление пользователя в базу
	var id int64
	stmt, err = r.DB.Prepare("INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (r *RepAuth) UserProvider(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.User"

	user := new(models.User)
	stmt, err := r.DB.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return *user, nil
}

// App returns app by id.
func (r *RepAuth) App(ctx context.Context, id int) (models.App, error) {
	const op = "repository.postgres.App"

	stmt, err := r.DB.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
