package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	AuthRepository "github.com/Antoha2/news/internal/auth/internal/repository"
	"github.com/Antoha2/news/internal/auth/pkg/jwt"
	"github.com/Antoha2/news/internal/auth/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

func (s *servImpl) RegisterNewUser(ctx context.Context, email string, pass string) (userID int64, err error) {
	const op = "serv.RegisterNewUser"
	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registring user")
	// Генерируем хэш и соль для пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// Сохраняем пользователя в БД
	id, err := s.rep.UserSaver(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *servImpl) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "serv.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("attempting to login user")

	// Достаём пользователя из БД
	user, err := s.rep.UserProvider(ctx, email)
	if err != nil {
		if errors.Is(err, AuthRepository.ErrUserNotFound) {
			s.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		s.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, err := s.rep.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")

	// Создаём токен авторизации
	token, err := jwt.NewToken(user, app, s.rep.TokenTTL)
	if err != nil {
		s.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
