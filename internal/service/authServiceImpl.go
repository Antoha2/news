package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Antoha2/news/internal/repository"
	"github.com/Antoha2/news/pkg/jwt"
	"github.com/Antoha2/news/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

func (s *authServImpl) RegisterNewUser(ctx context.Context, email string, pass string) (userID int64, err error) {
	const op = "serv.RegisterNewUser"
	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registring user")
	// generate a hash and salt for the password
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// Saving the user in the database
	id, err := s.rep.UserSaver(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *authServImpl) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "serv.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("attempting to login user")

	// Getting the user from the database
	user, err := s.rep.UserProvider(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		s.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Checking the correctness of the received password
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	// Create an authorization token
	token, err := jwt.NewToken(s.cfg, user)
	if err != nil {
		s.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
